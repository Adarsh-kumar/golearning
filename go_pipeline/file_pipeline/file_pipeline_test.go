package main

import "testing"

var path = "C://Users//Administrator//Downloads//checksum_data"

func benchmarkhelper(gr int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateMd5(path, gr)
	}
}

func BenchmarkCalculate1(b *testing.B) {

	benchmarkhelper(1, b)
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normal(path)
	}
}

func BenchmarkCalculate4(b *testing.B) {

	benchmarkhelper(4, b)
}

func BenchmarkCalculate8(b *testing.B) {

	benchmarkhelper(8, b)
}

func BenchmarkCalculate16(b *testing.B) {

	benchmarkhelper(16, b)
}
