package main

import "testing"

func benchmarkhelper(gr int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculate(gr)
	}
}

func BenchmarkCalculate1(b *testing.B) {

	benchmarkhelper(1, b)
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkCalculate4(b *testing.B) {

	benchmarkhelper(4, b)
}

func BenchmarkCalculate8(b *testing.B) {

	benchmarkhelper(8, b)
}

func BenchmarkCalculate11(b *testing.B) {

	benchmarkhelper(16, b)
}
