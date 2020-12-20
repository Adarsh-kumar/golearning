package main

import (
	"crypto/md5"
	//	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	//         "sort"
	"sync"
	"time"
)

// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel.

func walkFiles(root string) (<-chan string, <-chan error) {
	paths := make(chan string, 100)
	errc := make(chan error, 1)
	go func() { // HL
		// Close the paths channel after Walk returns.
		defer close(paths) // HL
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			paths <- path
			return nil
		})
	}()
	return paths, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// digester reads path names from paths and sends digest
func digester(paths <-chan string) <-chan result {
	// create a result channel and close once function completes
	resultchan := make(chan result, 100)
	//  defer close(resultchan)
	go func() {
		for path := range paths { // HLpaths
			data, err := ioutil.ReadFile(path)
			resultchan <- result{path, md5.Sum(data), err}

		}
		close(resultchan)
	}()
	return resultchan
}

func merge(resultchans []<-chan result) <-chan result {
	var wg sync.WaitGroup
	out := make(chan result, 100)

	output := func(c <-chan result) {
		for n := range c {
			//fmt.Println(n.path, " ", n.sum)
			out <- n
		}
		wg.Done()
	}

	// configure the wait group and start one goroutine for each input channel
	wg.Add(len(resultchans))
	for _, c := range resultchans {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func prepareResult(ans <-chan result) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	for r := range ans {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	return m, nil
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.

func MD5All(root string, maxgr int) <-chan result {

	paths, errc := walkFiles(root)

	// Create a list of channels that will be store the result of each digester
	resultchans := make([]<-chan result, maxgr)
	//fmt.Println("here is the grs", maxgr)
	//defer close(resultchans)
	//	var wg sync.WaitGroup
	//	numDigesters := maxgr
	//	wg.Add(numDigesters)
	for i := range resultchans {
		resultchans[i] = digester(paths)
	}

	// wait for all goroutines to complete
	/*	go func() {
		wg.Wait()
	}()*/

	ans := merge(resultchans)
	for c := range ans {
		_ = c
		//fmt.Println(c.path, c.sum)
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		fmt.Println("inside error", err)
		return nil
	}

	//m,err:= prepareResult(ans)
	return ans

}

func MD5Allsync(root string, maxgr int) <-chan result {

	paths, errc := walkFiles(root)

	// Create a list of channels that will be store the result of eaach digester
	resultchans := make([]<-chan result, maxgr)
	//defer close(resultchans)
	var wg sync.WaitGroup
	numDigesters := maxgr
	wg.Add(numDigesters)
	for i := range resultchans {
		resultchans[i] = digester(paths)
	}

	// wait for all goroutines to complete
	go func() {
		wg.Wait()
	}()

	ans := merge(resultchans)
	for c := range ans {
		_ = c
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		fmt.Println("inside error")
		return nil
	}

	//m,err:= prepareResult(ans)
	return ans

}

func calculateMd5(path string, maxgr int) <-chan result {
	// Calculate the MD5 sum of all files under the specified directory,
	// then print the results sorted by path name.
	// ans := MD5All(path, maxgr)
	ans := MD5All(path, maxgr)

	return ans
	/*	var paths []string
			if err != nil {
				fmt.Println(err)
				return m, paths
			}

			for path := range m {
				paths = append(paths, path)
			}
			sort.Strings(paths)
			for _, path := range paths {
				fmt.Printf("%x  %s\n", m[path], path)
			}
		        return m, paths
	*/
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

func init() {

	os.Setenv("GODEBUG", "gctrace=1")

}

func main() {

	//fmt.Println("Al")
	fmt.Println("GODEBUG", os.Getenv("GODEBUG"))

	start := time.Now()
	var path string = "C://Users//Administrator//Downloads//checksum_data"
	_ = calculateMd5(path, 16)
	fmt.Println("time taken ", time.Since(start))
	/*for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}*/
}
