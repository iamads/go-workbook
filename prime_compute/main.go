package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	pprof "runtime/pprof"
	// "runtime"
	"sync"
	// "golang.org/x/sync/semaphore"
)

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)

	fmt.Println("Prime till n=%n", computePrimeChunkedAndAtomic(1000000))
	pprof.StopCPUProfile()
}

func markAllMultipleAsNotPrime(isPrimeArr []bool, toCheck int) {

	cur := toCheck * toCheck
	for cur < len(isPrimeArr) {
		if isPrimeArr[cur] {
			isPrimeArr[cur] = false
		}
		cur += toCheck
	}
}

func computePrimeSingle(n int) []int {
	isPrimeArr := make([]bool, n+1)

	for i := 0; i < len(isPrimeArr); i++ {
		if i == 0 || i == 1 {
			isPrimeArr[i] = false
		} else {
			isPrimeArr[i] = true
		}
	}
	i := 2 // first prime

	for i <= n {
		if isPrimeArr[i] && i*i < n {
			markAllMultipleAsNotPrime(isPrimeArr, i)
		}
		i++
	}

	res := []int{}
	for i, v := range isPrimeArr {
		if v {
			res = append(res, i)
		}
	}

	return res
}

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
	// // fmt.Println("p")
	cur := i + i
	// fmt.Println("q")
	// fmt.Println("s")
	for cur < len(pr.isPrimeArr) {
		pr.mu.Lock()
		pr.isPrimeArr[cur] = false
		pr.mu.Unlock()
		cur += i
	}
	// fmt.Println("t")
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
		// 	fmt.Println("a")
		pr.mu.RLock()
		// 	fmt.Println("b")
		isPrime := pr.isPrimeArr[i]
		// 	fmt.Println("c")
		pr.mu.RUnlock()
		// 	fmt.Println("d")
		if isPrime {
			// 		fmt.Println("e")
			if isReallyPrime(i) {
				// 			fmt.Println("f")
				pr.wg.Add(1)
				go markMultipleNotPrimeConcurrent(pr, i)
				// 			fmt.Println("g")
			}
			// 		fmt.Println("h")
		}
		i++
	}
	// fmt.Println("i")

	pr.wg.Wait()

	res := []int{}
	for i, v := range pr.isPrimeArr {
		if v {
			res = append(res, i)
		}
	}

	return res
}

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
		// 	fmt.Println("a")
		pr.mu.RLock()
		// 	fmt.Println("b")
		isPrime := pr.isPrimeArr[i]
		// 	fmt.Println("c")
		pr.mu.RUnlock()
		// 	fmt.Println("d")
		if isPrime {
			// 		fmt.Println("e")
			if isReallyPrime(i) {
				// 			fmt.Println("f")
				pr.wg.Add(1)
				sem <- struct{}{}
				go func(pr *PrimeRegistry, i int) {
					defer func() { <-sem }()

					markMultipleNotPrimeConcurrent(pr, i)
					// 			fmt.Println("g")
				}(pr, i)

			}
			// 		fmt.Println("h")
		}
		i++
	}
	// fmt.Println("i")

	pr.wg.Wait()

	res := []int{}
	for i, v := range pr.isPrimeArr {
		if v {
			res = append(res, i)
		}
	}

	return res
}

// I thought of writing ta sharded memory implementation
// but I still there I will be facing locking issues and
// will lead to bad performance
func compute_prime_sharded_memory() {}
