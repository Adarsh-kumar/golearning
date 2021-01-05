package main

import (
	"context"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"syscall"
	"time"
)

// walkFiles starts a goroutine to walk the directory and send the
// path of each regular file on the string channel.

func walkFiles(ctx context.Context, wgglobal *sync.WaitGroup, root string) (<-chan string, <-chan error) {
	paths := make(chan string, 100)
	errc := make(chan error, 10)

	// Add a value to waitgroup because we are to launch a new goroutine
	// defer wgglobal.Done() marks done for this gorotine
	// Now waiting caller can know that

	wgglobal.Add(1)

	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		defer wgglobal.Done()
		defer close(errc)

		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !info.Mode().IsRegular() {
				//fmt.Println("here i am")
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				fmt.Println("walk cencelled1")
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  []byte
	err  error
}

// digester reads path names from paths, calcculate hash of that file
// and sends on a channel.
func digester(ctx context.Context, wgglobal *sync.WaitGroup, paths <-chan string) <-chan result {
	// create a result channel and close once function completes
	resultchan := make(chan result, 100)

	go func() {
		defer close(resultchan)
		defer wgglobal.Done()
		for path := range paths { // HLpaths

			f, err := os.Open(path)
			defer f.Close()

			h := md5.New()

			if err != nil {
				//log.Fatal(err)
				resultchan <- result{path, h.Sum(nil), err}
			} else {

				if _, err := io.Copy(h, f); err != nil {
					resultchan <- result{path, h.Sum(nil), err}
				}
				select {
				case resultchan <- result{path, h.Sum(nil), err}:
				case <-ctx.Done():
					{
						fmt.Println("walk cencelled2")
						return
					}
				}
			}
		}

	}()
	return resultchan
}

func merge(ctx context.Context, wgglobal *sync.WaitGroup, resultchans []<-chan result) <-chan result {
	var wg sync.WaitGroup
	out := make(chan result, 100)

	output := func(ctx context.Context, c <-chan result) {
		defer wg.Done()
		defer wgglobal.Done()
		for n := range c {
			//fmt.Println(n.path, " ", n.sum)
			select {
			case out <- n:
			case <-ctx.Done():
				{
					fmt.Println("walk cencelled3")
					return
				}
			}
		}

	}

	// configure the wait group and start one goroutine for each input channel
	wg.Add(len(resultchans))
	wgglobal.Add(len(resultchans))

	for _, c := range resultchans {
		go output(ctx, c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.
	go func() {
		wg.Wait()
		//	fmt.Println("closing the out channle")
		close(out)
	}()

	return out
}

// Base function which is called by calculateMd5
// launches a number of goroutines for the task and
// returns a map of file path and md5 checksum

func MD5All(ctx context.Context, wgglobal *sync.WaitGroup, root string, maxgr int) (map[string][]byte, error) {

	paths, errc := walkFiles(ctx, wgglobal, root)

	// create a goroutine which will be listening to the
	// errc , new go channel so that this function's execution  is not blocked

	go func() {
		for err := range errc {
			if err != nil {
				fmt.Println("got some error currently ignoring", err)
			}
		}

	}()

	// Create a list of channels of length maxgr
	// so that each goroutine launched to calculate checksum will push results in it's seprate channel
	resultchans := make([]<-chan result, maxgr)
	wgglobal.Add(maxgr)

	for i := 0; i < maxgr; i++ {
		resultchans[i] = digester(ctx, wgglobal, paths)
	}

	// merge the results calculated by each goroutine in a single channel
	//done <- struct{}{}

	ans := merge(ctx, wgglobal, resultchans)
	// store the results in a map
	m := make(map[string][]byte)
	//m2 := make(map[string]error)

	for r := range ans {
		if r.err != nil {
			fmt.Println("Error occurred for file, ignoring from final result ", r.path)
			fmt.Println(r.err)

		} else {
			m[r.path] = r.sum
		}
	}
	// Check whether the Walk failed.
	/*	if err := <-errc; err != nil {
		// if we want to handle errors gracefully (i.e. if one file from directory can't be open or some error)
		// we are ignoring the error and will show in final result ,
		// if you want to  terminate uncomment the return line.
		fmt.Println("inside error")
		_ = err
		return nil, err
	}*/

	return m, nil

}

func calculateMd5(path string, maxgr int, memdebug bool, printresult bool) {

	// intialise a done channle and pass to downstream functions
	// so that if upsteram function expliciitly cencels the task
	// downstream goroutines listening to this done channel can close their
	// channels and memory leak can be avoided.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// waitgroup to make sure that every downstream goroutine finishes gracefully in case of cencellation
	wgglobal := &sync.WaitGroup{}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func(ctx context.Context, cancel context.CancelFunc) {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		//push value in done channel so that downsteram goroutines can gracefully exit
		cancel()

		//os.Exit(0)
	}(ctx, cancel)

	// Calculate the MD5 sum of all files under the specified directory,

	var mem runtime.MemStats

	// if memdebug flag is true print the momory statistics before and after calling the program
	// Note - Heapprofile is a sampling based and memory statics gives the full idea of every allocation
	// and deallocation. ref -

	if memdebug {
		fmt.Println("memory baseline...")
		runtime.ReadMemStats(&mem)
		log.Println(mem.Alloc)
		log.Println(mem.TotalAlloc)
		log.Println(mem.HeapAlloc)
		log.Println(mem.HeapSys)
	}

	m, err := MD5All(ctx, wgglobal, path, maxgr)

	if memdebug {
		fmt.Println("memory comparison...")
		runtime.ReadMemStats(&mem)
		log.Println(mem.Alloc)
		log.Println(mem.TotalAlloc)
		log.Println(mem.HeapAlloc)
		log.Println(mem.HeapSys)
	}

	_ = m
	_ = err

	if printresult {
		var paths []string
		if err != nil {
			fmt.Println(err)
			//return m, paths
		}

		for path := range m {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			fmt.Printf("%x  %s\n", m[path], path)
		}
	}

	wgglobal.Wait()

}

func normal(root string) error {
	start := time.Now()
	t := time.Now()
	maxt := t.Sub(start)

	errc := make(chan error, 1)
	errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		start = time.Now()
		data, err := ioutil.ReadFile(path)
		_ = md5.Sum(data)
		t = time.Now()

		if maxt < t.Sub(start) {
			maxt = t.Sub(start)
		}
		//fmt.Println("Time taken in this element ", t.Sub(start))
		return nil
	})
	//fmt.Println("Max sequential time ", maxt)
	_ = start
	_ = maxt
	_ = t
	return nil
}

func main() {

	// Initialise the flags
	grPtr := flag.Int("gr", 16, "Number of goroutines to be launced for calculating hash")

	printresult := flag.Bool("printresult", false, "If you want to print the results")

	memdebug := flag.Bool("memdebug", false, "prints the memory allocation deallocation profile")

	// parse the flags

	flag.Parse()

	start := time.Now()
	var path string = "C://Users//Administrator//Downloads//checksum_data"

	// call the md5 calculator
	calculateMd5(path, *grPtr, *memdebug, *printresult)
	fmt.Println("time taken ", time.Since(start))

}
