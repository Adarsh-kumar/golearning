package main
import ( "fmt"
  "sync"
  "time"
)

func fib(n int) int {
        if n < 3 {
                return n
        } else {
                return fib(n-1) + fib(n-2)
        }
}


func generate(nums []int) <- chan int {
done := make(chan int,10);

go func(){
for _,i := range nums{
  done <- i
  }
  close(done)
}()

  return done
}

func produce(done <- chan int) <- chan int {
out:= make(chan int,10);

go func(){
  for n := range done {
    out <- fib(n)
  }
  close(out)
}()

return out
}

func merge(cs []<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int,10)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }

    // configure the wait group and start one goroutine for each input channel
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func makeRange(min, max int) []int {
    a := make([]int, max-min+1)
    for i := range a {
        a[i] = min + i
    }
    return a
}

func normal(arr []int){


   for i :=range arr{
   fmt.Println(fib(arr[i]))
}
}
func calculate(numberofgr int ) {

    // Set up the pipeline.
    arr:= makeRange(40,50)
//    start := time.Now()
//    normal(arr)
//    fmt.Println("took ", time.Since(start))
    start := time.Now()
    c := generate(arr)
    
    channels :=make([]<-chan int,numberofgr)

    for i:= range channels{
        channels[i]= produce(c)
}

    // Consume the merged output from c1 and c2.
    ans := make([]int, 11)
    var index int= 0
    for n := range merge(channels) {
        //fmt.Println(n) 
        ans[index]=n
        index++  
        // 4 then 9, or 9 then 4
    }
//    elapsed := time.Since(start)
//    fmt.Println("took ", elapsed)

}

func main(){
fmt.Println("go")
}

