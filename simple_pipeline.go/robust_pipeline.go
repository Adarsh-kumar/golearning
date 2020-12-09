package main

import (
	"fmt"
	"sync"
	"time"
)

func prepare(done <-chan struct{}, arr []int) <-chan int {
	// make a channle to push the numbers from array in a channel
	in := make(chan int, 20)

	go func() {
		for _, n := range arr {

			in <- n
		}
		close(in)
	}()

	return in
}

func fib(n int) int {
	if n < 3 {
		return n
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func process(done chan struct{}, in <-chan int) chan int {
	out := make(chan int, 20)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- fib(n):
			case <-done:
				return
			}
		}
	}()

	//close(out)

	return out
}

func merge(done chan struct{}, channels []chan (int)) <-chan (int) {
	// we need a wait group so that we can wait till every goroutine finishes
	var wg sync.WaitGroup
	// make a channel to push the outputs got from the worker goroutines
	out := make(chan int, 20)

	// output function definition
	output := func(in <-chan int) {
		for i := range in {
			select {
			case out <- i:
			case <-done:
				return
			}
		}
		// mark a done in wg because one goroutine finished
		wg.Done()
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}

func calculate(n int) <-chan int {
	// list of availble integers

	done := make(chan struct{})
	defer close(done)

	arr := make([]int, 18)

	for i := range arr {
		arr[i] = i
	}

	// put these a channel so that i/o can be done concurrently

	in := prepare(done, arr)

	chans := make([]chan int, n)
	//fmt.Println("number to fan in")
	//fmt.Println(n)
	for i := range chans {
		//chans[i] = make(chan int, 4)
		chans[i] = process(done, in)
	}

	// distribute to multiple workers

	// merge the results

	out := merge(done, chans)

	/*for i := range out {
		fmt.Println(i)
	} */
	return out

}

func main() {
	start := time.Now()
	calculate(2)
	fmt.Println(time.Since(start))
}
