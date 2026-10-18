package matrix

import (
	"fmt"
	"math/rand"
)

type Matrix struct {
	data []float32
	size int
}

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
