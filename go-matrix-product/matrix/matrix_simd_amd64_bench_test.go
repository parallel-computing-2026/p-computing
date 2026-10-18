//go:build goexperiment.simd && amd64

package matrix

import (
	"testing"
)

func benchmarkSIMDAMD64Product(b *testing.B, size int) {
	b.Helper()

	matrixA := NewRandomMatrix(size, seed)
	matrixB := NewRandomMatrix(size, seed+1)

	b.ReportMetric(256, "vec-bits")
	b.ReportAllocs()
	for b.Loop() {
		_ = SIMDAMD64Product(matrixA, matrixB)
	}
}

func BenchmarkSIMDAMD64Product256(b *testing.B) {
	const size = 256
	benchmarkSIMDAMD64Product(b, size)
}

func BenchmarkSIMDAMD64Product512(b *testing.B) {
	const size = 512
	benchmarkSIMDAMD64Product(b, size)
}

func BenchmarkSIMDAMD64Product1024(b *testing.B) {
	const size = 1024
	benchmarkSIMDAMD64Product(b, size)
}

func BenchmarkSIMDAMD64Product2048(b *testing.B) {
	const size = 2048
	benchmarkSIMDAMD64Product(b, size)
}

func BenchmarkSIMDAMD64Product4096(b *testing.B) {
	const size = 4096
	benchmarkSIMDAMD64Product(b, size)
}
