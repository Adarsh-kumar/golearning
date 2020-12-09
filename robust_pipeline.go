package main
import ( "fmt"
  "sync"
)


func prepare(done <- chan struct{},arr []int) <- chan int {
// make a channle to push the numbers from array in a channel
in:= make(chan int);

go func(){
  for _,n:=range arr {

  in <- n
}
close(in);
}()

return in
}

func process(done chan struct{},in <- chan int) <- chan int {
  out:= make(chan int)

  go func(){
    defer close(out)
  for n:=range in{
    select {
    case out <- n*n:
    case <-done:
    return
    }
}
close(out)
}()

  //close(out)

  return out
}

func merge(done chan struct{},channels ...<-chan(int)) <- chan(int){
  // we need a wait group so that we can wait till every goroutine finishes
  var wg sync.WaitGroup
  // make a channel to push the outputs got from the worker goroutines
  out:=make(chan int)

  // output function definition
  output:= func(in <- chan int){
    for i:= range in {
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


func main(){
  // list of availble integers

  done := make(chan struct{})
  defer close(done)

  arr:=[]int{1,2,3,4,5,6,76,7}

  // put these a channel so that i/o can be done concurrently

  in:= prepare(done, arr)

  // distribute to multiple workers

  w1:= process(done,in)
  w2:= process(done,in)

  // merge the results

  out:= merge(done,w1,w2)

  for n := range out{
  fmt.Println(n)
}
}
