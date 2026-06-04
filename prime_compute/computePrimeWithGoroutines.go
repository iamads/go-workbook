package main

import (
	"math"
	"sync"
)

type PrimeRegistry struct {
	isPrimeArr []bool
	mu         sync.RWMutex
	wg         sync.WaitGroup
}

func NewPrimeRegistry(n int) *PrimeRegistry {
	npr := PrimeRegistry{isPrimeArr: make([]bool, n+1), mu: sync.RWMutex{},
		wg: sync.WaitGroup{}}
	return &npr
}

func isReallyPrime(n int) bool {
	for j := 2; j <= int(math.Sqrt(float64(n))); j++ {
		if n%j == 0 {
			return false
		}
	}
	return true
}

func markMultipleNotPrimeConcurrent(pr *PrimeRegistry, i int) {
	defer pr.wg.Done()

	cur := i * i
	for cur < len(pr.isPrimeArr) {
		pr.mu.Lock()
		pr.isPrimeArr[cur] = false
		pr.mu.Unlock()
		cur += i
	}
}

func computePrimesWithGoroutines(n int) []int {

	pr := NewPrimeRegistry(n)

	for i := 0; i < len(pr.isPrimeArr); i++ {
		if i == 0 || i == 1 {
			pr.isPrimeArr[i] = false
		} else {
			pr.isPrimeArr[i] = true
		}
	}

	i := 2 // firstPrime

	for i <= n {
		pr.mu.RLock()
		isPrime := pr.isPrimeArr[i]
		pr.mu.RUnlock()
		if isPrime {
			if isReallyPrime(i) {
				pr.wg.Add(1)
				go markMultipleNotPrimeConcurrent(pr, i)
			}
		}
		i++
	}

	pr.wg.Wait()

	res := []int{}
	for i, v := range pr.isPrimeArr {
		if v {
			res = append(res, i)
		}
	}

	return res
}
