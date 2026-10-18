package matrix

import "testing"

const seed int64 = 123

func benchmarkSimpleProduct(b *testing.B, size int) {
	b.Helper()

	matrixA := NewRandomMatrix(size, seed)
	matrixB := NewRandomMatrix(size, seed+1)

	b.ReportAllocs()
	for b.Loop() {
		_ = SimpleProduct(matrixA, matrixB)
	}
}

func BenchmarkSimpleProduct256(b *testing.B) {
	const size = 256
	benchmarkSimpleProduct(b, size)
}

func BenchmarkSimpleProduct512(b *testing.B) {
	const size = 512
	benchmarkSimpleProduct(b, size)
}

func BenchmarkSimpleProduct1024(b *testing.B) {
	const size = 1024
	benchmarkSimpleProduct(b, size)
}

func BenchmarkSimpleProduct2048(b *testing.B) {
	const size = 2048
	benchmarkSimpleProduct(b, size)
}

//func BenchmarkSimpleProduct4096(b *testing.B) {
//	const size = 4096
//	benchmarkSimpleProduct(b, size)
//}
