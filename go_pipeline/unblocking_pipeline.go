package main
import ( "fmt"
  "sync"
  "time"
)



func generate(nums ...int) <- chan int {
done := make(chan int,len(nums));

go func(){
for _,i := range nums{
  done <- i
  }
  close(done)
}()

  return done
}

func produce(done <- chan int) <- chan int {
out:= make(chan int);

go func(){
  for n := range done {
    out <- n*n
  }
  close(out)
}()

return out
}

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int,1)

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

func main() {

    // Set up the pipeline.
    //arr:= makeRange(1,100)
    start := time.Now()
    c := generate(1,2,3,4,5,6,7,8,9,12,13,14,15)
    out1 :=produce(c)
    out2 :=produce(c)

    // Consume the merged output from c1 and c2.
    for n := range merge(out1, out2) {
        fmt.Println(n) // 4 then 9, or 9 then 4
    }
    elapsed := time.Since(start)
    fmt.Println("took ", elapsed)

}
