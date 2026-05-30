package main

import (
	"runtime"
	"sync"
)

func computePrimeGoroutinesAtomicBounded(n int) []int {
	apr := NewAtomicPrimeRegistry(n)
	ch := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup

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
			ch <- struct{}{}
			wg.Add(1)
			go func(apr *PrimeRegistryAtomic, i int) {
				defer func() {
					<-ch
					wg.Done()
				}()
				markMultipleNotPrimeAtomic(apr, i)
			}(apr, i)
		}
		i++
	}

	wg.Wait()

	res := []int{}
	for i := 0; i < len(apr.isPrimeArr); i++ {
		if apr.isPrimeArr[i].Load() {
			res = append(res, i)
		}
	}

	return res
}
