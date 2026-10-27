//go:build goexperiment.simd

package matrix

import (
	"simd"
	"testing"
)

func benchmarkSIMDProduct(b *testing.B, size int) {
	b.Helper()

	matrixA := NewRandomMatrix(size, seed)
	matrixB := NewRandomMatrix(size, seed+1)

	b.ReportMetric(float64(simd.VectorBitSize()), "vec-bits")
	b.ReportAllocs()
	for b.Loop() {
		_ = SIMDProduct(matrixA, matrixB)
	}
}

func BenchmarkSIMDProduct256(b *testing.B) {
	const size = 256
	benchmarkSIMDProduct(b, size)
}

func BenchmarkSIMDProduct512(b *testing.B) {
	const size = 512
	benchmarkSIMDProduct(b, size)
}

func BenchmarkSIMDProduct1024(b *testing.B) {
	const size = 1024
	benchmarkSIMDProduct(b, size)
}

func BenchmarkSIMDProduct2048(b *testing.B) {
	const size = 2048
	benchmarkSIMDProduct(b, size)
}

func BenchmarkSIMDProduct4096(b *testing.B) {
	const size = 4096
	benchmarkSIMDProduct(b, size)
}
