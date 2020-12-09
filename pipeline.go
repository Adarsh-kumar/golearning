package main
import ( "fmt"
)



func generate(nums ...int) <- chan int {
done := make(chan int);

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

func main() {
    // Set up the pipeline.
    c := generate(2, 3)
    out := produce(c)

    // Consume the output.
    fmt.Println(<-out) // 4
    fmt.Println(<-out) // 9
}
