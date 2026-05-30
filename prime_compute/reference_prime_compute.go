package main

// computePrimeReference is a straightforward sieve of Eratosthenes.
// It is intended to be the correctness reference all different version
func computePrimeReference(n int) []int {
	if n < 2 {
		return []int{}
	}

	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	for p := 2; p*p <= n; p++ {
		if !isPrime[p] {
			continue
		}
		for multiple := p * p; multiple <= n; multiple += p {
			isPrime[multiple] = false
		}
	}

	primes := make([]int, 0, n/2)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}
