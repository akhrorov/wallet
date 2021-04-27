package sum

import "testing"

func BenchmarkRegular(b *testing.B) {
	want := int64(2000)
	for i := 0; i < b.N; i++ {
		result := Regular()
		if result != want {
			b.Fatalf("Invalid result, got %v, want %v", result, want)
		}
	}
}

func BenchmarkConcurrently(b *testing.B) {
	want := int64(2000)
	for i := 0; i < b.N; i++ {
		result := Concurrently()
		if result != want {
			b.Fatalf("Invalid result, got %v, want %v", result, want)
		}
	}
}