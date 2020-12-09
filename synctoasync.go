package main
import ( "fmt"
  //"io/ioutil"
	"net/http"
)


type data struct
{
  Body []byte
  Error error
}

//takes url.string as input and writess to a channel once done
func getdata(url string) <- chan(bool) {
  c := make(chan bool)

  go func() {
		resp, _:= http.Get(url)
		//body, _:= ioutil.ReadAll(resp.Body)

	  fmt.Println(string(resp.StatusCode))
     c <- true
	}()

  return c
}


func main() {
    done := make(chan bool)
    url:="http://google.com"

    go func() {
  		resp, _:= http.Get(url)
  		//body, _:= ioutil.ReadAll(resp.Body)

  	  fmt.Println((resp.StatusCode))
       done <- true
  	}()


    go func() {
  		resp, _:= http.Get(url)
  		//body, _:= ioutil.ReadAll(resp.Body)

  	  fmt.Println((resp.StatusCode))
       done <- true
  	}()


    go func() {
  		resp, _:= http.Get(url)
  		//body, _:= ioutil.ReadAll(resp.Body)

  	  fmt.Println((resp.StatusCode))
       done <- true
  	}()

	// do many other things

	<-done
  <-done
  <-done
}
