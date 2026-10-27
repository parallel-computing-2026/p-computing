package matrix

import (
	"fmt"
	"math/rand"
)

// Matrix is a square matrix of float values.
type Matrix struct {
	data []float32
	size int
}

// NewRandomMatrix generates a random matrix of size nXn
// with values in range of [0; 10).
func NewRandomMatrix(n int, seed int64) Matrix {
	r := rand.New(rand.NewSource(seed))

	square := n * n
	matrix := Matrix{
		data: make([]float32, square),
		size: n,
	}

	data := matrix.data
	for i := range square {
		data[i] = r.Float32() * 10 // values in range [0; 10)
	}

	return matrix
}

// SimpleProduct multiplies a and b row by column 0(n^3).
//
// It panics if a.size is not equal to b.size.
func SimpleProduct(a Matrix, b Matrix) Matrix {
	size := a.size
	if size != b.size {
		panic(fmt.Sprintf("size mismatch: %d != %d", size, b.size))
	}

	result := Matrix{
		data: make([]float32, size*size),
		size: size,
	}

	for row := range size {
		for col := range size {
			sum := float32(0)
			for k := range size {
				sum += a.data[row*size+k] * b.data[k*size+col]
			}

			result.data[row*size+col] = sum
		}
	}

	return result
}
