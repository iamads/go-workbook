package main

import (
	"fmt"
	"testing"
)

var testArr []int = []int{10, 20, 100, 1000, 100000, 1000000}

func BenchmarkComputePrimeSingle(b *testing.B) {
	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_prime_single to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimeSingle(n)
			}
		})
	}
}

func BenchmarkComputePrimeWithGoroutines(b *testing.B) {

	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_primes_with_goroutines to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimesWithGoroutines(n)
			}
		})
	}
}
func BenchmarkComputePrimeWithGoroutinesBoundedConcurrency(b *testing.B) {

	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_primes_with_goroutines_bounded_concurrency to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimesWithGoroutinesBoundedConcurrency(n)
			}
		})
	}
}

func BenchmarkComputePrimeGoroutinesAtomic(b *testing.B) {

	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_prime_goroutines_atomic to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimeGoroutinesAtomic(n)
			}
		})
	}
}
func BenchmarkComputePrimeGoroutinesAtomicBounded(b *testing.B) {

	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_prime_goroutines_atomic_bounded to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimeGoroutinesAtomicBounded(n)
			}
		})
	}
}

func BenchmarkComputePrimeChunked(b *testing.B) {

	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_prime_chunked to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimeChunked(n)
			}
		})
	}
}

func BenchmarkComputePrimeChunkedAtomic(b *testing.B) {

	for _, n := range testArr {
		b.Run(fmt.Sprintf("Using compute_prime_chunked_atomic to find primes till n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = computePrimeChunkedAtomic(n)
			}
		})
	}
}
