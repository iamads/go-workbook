package main

import "runtime"

func computePrimesWithGoroutinesBoundedConcurrency(n int) []int {

	pr := NewPrimeRegistry(n)
	sem := make(chan struct{}, runtime.NumCPU())

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
				sem <- struct{}{}
				go func(pr *PrimeRegistry, i int) {
					defer func() { <-sem }()

					markMultipleNotPrimeConcurrent(pr, i)
				}(pr, i)

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
