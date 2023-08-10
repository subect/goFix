package test

import "testing"

func BenchmarkMakeSliceWithoutAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeSliceWithoutAlloc()
	}
}

func BenchmarkMakeSliceWithPreAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeSliceWithPreAlloc()
	}
}
