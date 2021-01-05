# golearning

Golang concurrency practice

I took the example problem from Sameer Ajmani's blog Go concurrency pipeline and cencellation, and wrote a fan-in-out pipeline,benchmark tests and  improved the scalability to calculate the checksum of a files concurrently. Some details are described below.

go to  folder- go_pipeline/file_pipeline

bounded_file_checksum.go - code file
file_pipeline_test.go - Benchmark test file 

change the path variable to a local folder(having some mb of data) of your machine.

run -
go build bounded_file_checksum.go
./bounded_file_checksum.exe

What I am able to do till now 
------------------------------------------
1. Minimum sync block of goroutines due to channel syncronisation.
2. Graceful error handling and cencelling handling (SIGTERM) to avoid any memory memory leak.
3. Written beanchmark tests.
4. Collect profiling data ( cpuprofile,blocking profile, goroutine trace etc attached in go_pipeline/file_pipeline folder ) and understand every espect of excution. 
5. Avoided reading full file in memory and doing operation at once which was causing too frequent GC, and used the io.Writer interface of MD5 and process chunk by chunk.

Results of benchmark test
----------------------------------------------------------------
go test -bench=BenchmarkCalculate -benchtime=30s

goos: windows
goarch: amd64
| BenchmarkFunction | | total calls | | time taken per call |
| --------------- | --------------- | --------------- |
| BenchmarkCalculate1-16 |            |   5   |     | 7059610060 ns/op |
| BenchmarkCalculate4-16 |             |  18  |     | 1840480306 ns/op |
BenchmarkCalculate8-16                37         952633289 ns/op
BenchmarkCalculate16-16               67         536108457 ns/op
BenchmarkCalculate32-16               64         553310138 ns/op
BenchmarkCalculate64-16               63         568016006 ns/op
PASS

where suffix 1-16 means that one goroutine was launched on 16 represents the numberr of cores on my machine.

The second column means how many time program was executed in 30 seconds and last columns says time per execution.

Note that using one goroutine we are doing the task in around 7 seconds while with 16 goroutines we are close to one 500 millisecond.

Further Work
-------------------------------------------
Make this frameworrk more generic.


