package main

import (
	"sync"
	"sync/atomic"
)

type PrimeRegistryAtomic struct {
	isPrimeArr []atomic.Bool
	wg         sync.WaitGroup
}

func NewAtomicPrimeRegistry(n int) *PrimeRegistryAtomic {
	apr := PrimeRegistryAtomic{isPrimeArr: make([]atomic.Bool, n+1), wg: sync.WaitGroup{}}
	return &apr
}

func markMultipleNotPrimeAtomic(apr *PrimeRegistryAtomic, i int) {
	cur := i * i
	for cur < len(apr.isPrimeArr) {
		apr.isPrimeArr[cur].Store(false)
		cur += i
	}
}

func computePrimeGoroutinesAtomic(n int) []int {
	apr := NewAtomicPrimeRegistry(n)

	for i := 0; i < len(apr.isPrimeArr); i++ {
		if i == 0 || i == 1 {
			apr.isPrimeArr[i].Store(false)
		} else {
			apr.isPrimeArr[i].Store(true)
		}
	}

	i := 2 // first Prime

	for i <= n {
		if apr.isPrimeArr[i].Load() && isReallyPrime(i) {
			go markMultipleNotPrimeAtomic(apr, i)
		}
		i++
	}

	res := []int{}
	for i, v := range apr.isPrimeArr {
		if v.Load() {
			res = append(res, i)
		}
	}

	return res
}
