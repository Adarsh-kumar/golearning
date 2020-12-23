package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

func osio(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x", h.Sum(nil))
}

func ioutilio(path string) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x", md5.Sum(data))
}

func main() {

	var path = "C://Users//Administrator//Downloads//checksum_data//books9//ai(18).pdf"

	var mem runtime.MemStats

	fmt.Println("memory baseline...")

	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)

	ioutilio(path)

	fmt.Println("memory after funtion call...")

	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)

}
