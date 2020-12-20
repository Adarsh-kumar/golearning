# golearning

Golang concurrency practice

I took the example problem from Sameer Ajmani's blog Go concurrency pipeline and cencellation, and wrote a fan-in-out pipeline to calculate the checksum of a files concurrently.

What I am able to do till now 
------------------------------------------
1. Minimum sync block of goroutines due to channel syncronisation.
2. Graceful error handling.
3. Tuning garbage collection parameter to decraese tool frequent garbage collection (every five millisecond), which improved scalability.
4. Collect profiling data ( cpuprofile,blocking profile, goroutine trace etc attached in go_pipeline/file_pipeline folder ) and understand every espect of excution. 

Results of benchmark test
------------------------------------------

