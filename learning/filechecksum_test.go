package main

import (
	"testing"
)

var path string = "C://Users//Administrator//adarsh//golearning"

func benchmarkCalculate(path string, gr int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateMd5(path, gr)
	}
}

func benchmarkCalculateBuffered(path string, gr int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateMd5Buffered(path, gr)
	}
}

func BenchmarkCalculate1(b *testing.B)         { benchmarkCalculate(path, 1, b) }
func BenchmarkCalculateBuffered1(b *testing.B) { benchmarkCalculateBuffered(path, 1, b) }
func BenchmarkCalculate4(b *testing.B)         { benchmarkCalculate(path, 4, b) }
func BenchmarkCalculateBuffered4(b *testing.B) { benchmarkCalculateBuffered(path, 4, b) }
func BenchmarkCalculate8(b *testing.B)         { benchmarkCalculate(path, 8, b) }

func BenchmarkCalculateBuffered8(b *testing.B)  { benchmarkCalculateBuffered(path, 8, b) }
func BenchmarkCalculate16(b *testing.B)         { benchmarkCalculate(path, 16, b) }
func BenchmarkCalculateBuffered16(b *testing.B) { benchmarkCalculateBuffered(path, 16, b) }

/*
func BenchmarkCalculate32(b *testing.B) { benchmarkCalculate(path, 32, b) }

func BenchmarkCalculate64(b *testing.B) { benchmarkCalculate(path, 64, b) }

func BenchmarkCalculate128(b *testing.B) { benchmarkCalculate(path, 128, b) }
*/
