package main

import (
	"reflect"
	"strconv"
	"testing"
)

func TestComputePrimeReference_KnownValues(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{name: "below two", n: 1, want: []int{}},
		{name: "two", n: 2, want: []int{2}},
		{name: "ten", n: 10, want: []int{2, 3, 5, 7}},
		{name: "twenty", n: 20, want: []int{2, 3, 5, 7, 11, 13, 17, 19}},
		{name: "thirty", n: 30, want: []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computePrimeReference(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("computePrimeReference(%d) mismatch:\nwant %v\ngot  %v", tt.n, tt.want, got)
			}
		})
	}
}

func TestChunkedImplementationMatchesReference(t *testing.T) {
	tests := []int{2, 3, 10, 20, 50, 100, 250, 500, 1000000}

	for _, n := range tests {
		t.Run("n="+strconv.Itoa(n), func(t *testing.T) {
			want := computePrimeReference(n)
			got := computePrimeChunked(n)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("chunked implementation mismatch for n=%d:\nwant %v\ngot  %v", n, want, got)
			}
		})
	}
}

func TestAtomicChunkedImplementationMatchesReference(t *testing.T) {
	tests := []int{2, 3, 10, 20, 50, 100, 250, 500, 1000000}

	for _, n := range tests {
		t.Run("n="+strconv.Itoa(n), func(t *testing.T) {
			want := computePrimeReference(n)
			got := computePrimeChunkedAtomic(n)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("atomic chunked implementation mismatch for n=%d:\nwant %v\ngot  %v", n, want, got)
			}
		})
	}
}

func TestReferenceAndSingleMatch(t *testing.T) {
	tests := []int{2, 3, 10, 20, 50, 100, 250, 500, 1000000}

	for _, n := range tests {
		t.Run("n="+strconv.Itoa(n), func(t *testing.T) {
			want := computePrimeReference(n)
			got := computePrimeSingle(n)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("single implementation mismatch for n=%d:\nwant %v\ngot  %v", n, want, got)
			}
		})
	}
}
