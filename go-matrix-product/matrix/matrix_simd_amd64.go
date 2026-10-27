//go:build goexperiment.simd && amd64

package matrix

import (
	"fmt"
	"simd/archsimd"
)

// SIMDAMD64Product multiplies a and b using AVX2 SIMD.
//
// It panics if a.size is not equal to b.size.
func SIMDAMD64Product(a Matrix, b Matrix) Matrix {
	size := a.size
	if size != b.size {
		panic(fmt.Sprintf("size mismatch: %d != %d", size, b.size))
	}

	result := Matrix{
		data: make([]float32, size*size),
		size: size,
	}

	// AVX2 256 bits aka 8 float32 values
	const lanes = 8
	for row := range size {
		resultRow := result.data[row*size : (row+1)*size]
		for k := range size {
			rowElement := a.data[row*size+k]
			aVec := archsimd.BroadcastFloat32x8(rowElement)
			bRow := b.data[k*size : (k+1)*size]

			col := 0
			for ; col+lanes <= size; col += lanes {
				accumulator := archsimd.LoadFloat32x8(resultRow[col:]) // loads 8 float32 values
				bVec := archsimd.LoadFloat32x8(bRow[col:])
				aVec.MulAdd(bVec, accumulator).Store(resultRow[col:])
			}

			// handle the tail separately with partial
			// can be a manual loop instead
			if col < size {
				accumulator, _ := archsimd.LoadFloat32x8Part(resultRow[col:])
				bVec, _ := archsimd.LoadFloat32x8Part(bRow[col:])
				aVec.MulAdd(bVec, accumulator).StorePart(resultRow[col:])
			}
		}
	}

	return result
}
