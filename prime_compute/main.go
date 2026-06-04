package main

import (
	"fmt"
	"os"
	pprof "runtime/pprof"
)

func main() {
	testValue := 1000_000

	f, _ := os.Create("pprof_computePrimeSingle")
	pprof.StartCPUProfile(f)
	computePrimeSingle(testValue)
	pprof.StopCPUProfile()
	fmt.Println("computed prime single")

	f1, _ := os.Create("pprof_computePrimesWithGoroutines")
	pprof.StartCPUProfile(f1)
	computePrimesWithGoroutines(testValue)
	pprof.StopCPUProfile()
	fmt.Println("computed primes with gorroutines")

	f2, _ := os.Create("pprof_computePrimesWithGoroutinesWithBoundedConcurrency")
	pprof.StartCPUProfile(f2)
	computePrimesWithGoroutinesBoundedConcurrency(testValue)
	pprof.StopCPUProfile()
	fmt.Println("computed primes with gorroutines and bouded concurrecy")

	f3, _ := os.Create("pprof_computePrimeGoroutinesAtomic")
	pprof.StartCPUProfile(f3)
	computePrimeGoroutinesAtomic(testValue)
	pprof.StopCPUProfile()
	fmt.Println("computed PrimeGoroutinesAtomic")

	f4, _ := os.Create("pprof_computePrimeGoroutinesAtomicBounded")
	pprof.StartCPUProfile(f4)
	computePrimeGoroutinesAtomicBounded(testValue)
	pprof.StopCPUProfile()
	fmt.Println("computed Prime GoroutinesAtomic Bounded")

	f5, _ := os.Create("pprof_computePrimeChuked")
	pprof.StartCPUProfile(f5)
	computePrimeChunked(testValue)
	pprof.StopCPUProfile()
	fmt.Println("computed Prime Chunked")
}
