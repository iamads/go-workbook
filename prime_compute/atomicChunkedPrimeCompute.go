package main

import (
	"sync"
	"sync/atomic"
)

type AtomicChunkedPrimeRegistry struct {
	isPrimeArr   []atomic.Bool
	wg           sync.WaitGroup
	lastPrime    int
	computedTill int
}

func NewAtomicChunkedPrimeRegistry(n int) *AtomicChunkedPrimeRegistry {
	pr := AtomicChunkedPrimeRegistry{isPrimeArr: make([]atomic.Bool, n+1)}
	pr.lastPrime = 2
	return &pr
}

func markMultipleNotPrimeAtomicChunked(pr *AtomicChunkedPrimeRegistry, i int) {
	cur := i * i
	for cur < len(pr.isPrimeArr) {
		pr.isPrimeArr[cur].Store(false)
		cur += i
	}
}

func (pr *AtomicChunkedPrimeRegistry) thisORn(n int) int {
	if n < len(pr.isPrimeArr) {
		return n
	}
	return len(pr.isPrimeArr) - 1
}

func (pr *AtomicChunkedPrimeRegistry) sweep() {
	start := pr.lastPrime
	if start*start < len(pr.isPrimeArr) {
		for i := start; i < pr.thisORn(start*start); i++ {
			if pr.isPrimeArr[i].Load() && i*i < len(pr.isPrimeArr) {
				pr.wg.Add(1)
				go func(pr *AtomicChunkedPrimeRegistry, i int) {
					defer pr.wg.Done()
					markMultipleNotPrimeAtomicChunked(pr, i)
				}(pr, i)
			}
		}
	}

	pr.wg.Wait()

	pr.computedTill = start*start - 1
	for i := pr.thisORn(start*start - 1); i > start; i-- {
		if pr.isPrimeArr[i].Load() {
			pr.lastPrime = i
			break
		}
	}
}

func computePrimeChunkedAtomic(n int) []int {
	apr := NewAtomicChunkedPrimeRegistry(n)

	for i := 0; i < len(apr.isPrimeArr); i++ {
		apr.isPrimeArr[i].Store(i != 0 && i != 1)
	}

	for apr.computedTill < n {
		apr.sweep()
	}

	res := []int{}
	for i := 0; i < len(apr.isPrimeArr); i++ {
		if apr.isPrimeArr[i].Load() {
			res = append(res, i)
		}
	}
	return res
}
