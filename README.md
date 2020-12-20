# golearning

Golang concurrency practice

I took the example problem from Sameer Ajmani's blog Go concurrency pipeline and cencellation, and wrote a fan-in-out pipeline,benchmark tests and  improved the scalability to calculate the checksum of a files concurrently. Some details are described below.

go to  folder- go_pipeline/file_pipeline

bounded_file_checksum.go - code file
file_pipeline_test.go - Benchmark test file 

What I am able to do till now 
------------------------------------------
1. Minimum sync block of goroutines due to channel syncronisation.
2. Minimize lock contention using seprate channels to push result of each goroutine w.r.t shared channel.
3. Graceful error handling.
4. Tuning garbage collection parameter to decraese tool frequent garbage collection (every five millisecond), which improved scalability.
5. Collect profiling data ( cpuprofile,blocking profile, goroutine trace etc attached in go_pipeline/file_pipeline folder ) and understand every espect of excution. 

Results of benchmark test
------------------------------------------

goos: windows
goarch: amd64
BenchmarkCalculate1-16                 8        8003875975 ns/op
BenchmarkNormal-16                     8        7939249962 ns/op
BenchmarkCalculate4-16                28        2317645121 ns/op
BenchmarkCalculate8-16                50        1310940174 ns/op
BenchmarkCalculate16-16               75        1129502856 ns/op
PASS

where suffix 1-16 means that one goroutine was launched on 16 represents the numberr of cores on my machine.

Note that using one goroutine we are doing the task in around 8 seconds while with 16 goroutines we are close to one second.

Further improvement 
-------------------------------------------

I checked the cpu profile and it shows around 20 % of the total time taken is used in sync block and syscall blocking. I am trying to understand and remove that overhead.

