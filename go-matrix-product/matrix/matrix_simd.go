//go:build goexperiment.simd

package matrix

import (
	"fmt"
	"simd"
)

func SIMDProduct(a Matrix, b Matrix) Matrix {
	size := a.size
	if size != b.size {
		panic(fmt.Sprintf("size mismatch: %d != %d", size, b.size))
	}

	if simd.Emulated() {
		// fallback to SimpleProduct if simd is emulated
		return SimpleProduct(a, b)
	}

	result := Matrix{
		data: make([]float32, size*size),
		size: size,
	}

	// how many float32 can we put inside the largest vector
	lanes := simd.VectorBitSize() / 32
	for row := range size {
		resultRow := result.data[row*size : (row+1)*size]
		for k := range size {
			rowElement := a.data[row*size+k]
			aVec := simd.BroadcastFloat32s(rowElement)
			bRow := b.data[k*size : (k+1)*size]

			col := 0
			for ; col+lanes <= size; col += lanes {
				accumulator := simd.LoadFloat32s(resultRow[col:]) // loads lanes float32 values
				bVec := simd.LoadFloat32s(bRow[col:])
				aVec.MulAdd(bVec, accumulator).Store(resultRow[col:])
			}

			// handle the tail separately with partial
			// can be a manual loop instead
			if col < size {
				accumulator, _ := simd.LoadFloat32sPart(resultRow[col:])
				bVec, _ := simd.LoadFloat32sPart(bRow[col:])
				aVec.MulAdd(bVec, accumulator).StorePart(resultRow[col:])
			}
		}
	}

	return result
}
