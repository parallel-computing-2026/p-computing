//go:build goexperiment.simd && amd64

package matrix

import (
	"reflect"
	"testing"
)

func TestSIMDAMD64Product(t *testing.T) {
	tests := []struct {
		name string
		a    Matrix
		b    Matrix
		want Matrix
	}{
		{
			name: "valid 1x1",
			a: Matrix{
				data: []float32{1},
				size: 1,
			},
			b: Matrix{
				data: []float32{2},
				size: 1,
			},
			want: Matrix{
				data: []float32{2},
				size: 1,
			},
		},
		{
			name: "valid 2x2",
			a: Matrix{
				data: []float32{1, 2, 3, 4},
				size: 2,
			},
			b: Matrix{
				data: []float32{5, 6, 7, 8},
				size: 2,
			},
			want: Matrix{
				data: []float32{19, 22, 43, 50},
				size: 2,
			},
		},
		{
			name: "square 3x3",
			a: Matrix{
				data: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 3,
			},
			b: Matrix{
				data: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 3,
			},
			want: Matrix{
				data: []float32{30, 36, 42, 66, 81, 96, 102, 126, 150},
				size: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SIMDAMD64Product(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SIMDProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
