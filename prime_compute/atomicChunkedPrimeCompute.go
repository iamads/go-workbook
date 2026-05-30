package main

import (
	// "runtime"
	// "fmt"
	"sync"
)

// the problem we are facing is performing isReallyPrime
// If we do linearly it is easy to findPrime
// but when we do it concurrently we are not sure when
// a particular value will be reached
//
// I want to use the logic where if I have a prime value p
// then I can find all non-prime numbers with values [2, p)
//
// I will chunk primes and keep on increasing the length

type PrimeRegistryWithLastPrime struct {
	isPrimeArr   []bool
	wg           sync.WaitGroup
	lastPrime    int
	computedTill int
}

func markMultipleNotPrimeAtomicLatPrime(apr *PrimeRegistryWithLastPrime, i int) {
	// fmt.Println("Malking multiples for ", i)
	cur := i * i
	for cur < len(apr.isPrimeArr) {
		apr.isPrimeArr[cur] = false
		cur += i
	}
}

func (pr *PrimeRegistryWithLastPrime) thisORn(n int) int {
	if n < len(pr.isPrimeArr) {
		return n
	}
	return len(pr.isPrimeArr) - 1
}

func (pr *PrimeRegistryWithLastPrime) sweep() {
	start := pr.lastPrime
	// ch := make(chan struct{}, runtime.NumCPU())
	if start*start < len(pr.isPrimeArr) {
		for i := start; i < pr.thisORn(start*start); i++ {
			if pr.isPrimeArr[i] && i*i < len(pr.isPrimeArr) {
				// ch <- struct{}{}
				pr.wg.Add(1)
				go func(pr *PrimeRegistryWithLastPrime, i int) {

					defer func() {
						// <-ch
						pr.wg.Done()
					}()
					markMultipleNotPrimeAtomicLatPrime(pr, i)
				}(pr, i)
			}
		}
	}

	pr.wg.Wait()

	pr.computedTill = start*start - 1
	for i := pr.thisORn(start*start - 1); i > start; i-- {
		if pr.isPrimeArr[i] {
			pr.lastPrime = i
			break
		}
	}
}

func NewPrimeRegistryWithLastPrime(n int) *PrimeRegistryWithLastPrime {
	pr := PrimeRegistryWithLastPrime{isPrimeArr: make([]bool, n+1), wg: sync.WaitGroup{}}
	pr.lastPrime = 2
	return &pr
}

func computePrimeChunkedAndAtomic(n int) []int {

	apr := NewPrimeRegistryWithLastPrime(n)

	for i := 0; i < len(apr.isPrimeArr); i++ {
		if i == 0 || i == 1 {
			apr.isPrimeArr[i] = false
		} else {
			apr.isPrimeArr[i] = true
		}
	}

	for int(apr.computedTill) < n {
		// fmt.Println(apr.computedTill, apr.lastPrime)
		apr.sweep()
	}

	res := []int{}

	for i := 0; i < len(apr.isPrimeArr); i++ {
		if apr.isPrimeArr[i] {
			res = append(res, i)
		}
	}
	return res
}
