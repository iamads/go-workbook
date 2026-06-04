package main

import (
	"fmt"
	"log"
	"os"
	pprof "runtime/pprof"
)

func main() {
	testValue := 1000_000

	profiles := []struct {
		fileName string
		label    string
		compute  func()
	}{
		{
			fileName: "pprof_computePrimeSingle",
			label:    "computed prime single",
			compute:  func() { computePrimeSingle(testValue) },
		},
		{
			fileName: "pprof_computePrimesWithGoroutines",
			label:    "computed primes with goroutines",
			compute:  func() { computePrimesWithGoroutines(testValue) },
		},
		{
			fileName: "pprof_computePrimesWithGoroutinesWithBoundedConcurrency",
			label:    "computed primes with goroutines and bounded concurrency",
			compute:  func() { computePrimesWithGoroutinesBoundedConcurrency(testValue) },
		},
		{
			fileName: "pprof_computePrimeGoroutinesAtomic",
			label:    "computed prime goroutines atomic",
			compute:  func() { computePrimeGoroutinesAtomic(testValue) },
		},
		{
			fileName: "pprof_computePrimeGoroutinesAtomicBounded",
			label:    "computed prime goroutines atomic bounded",
			compute:  func() { computePrimeGoroutinesAtomicBounded(testValue) },
		},
		{
			fileName: "pprof_computePrimeChunked",
			label:    "computed prime chunked",
			compute:  func() { computePrimeChunked(testValue) },
		},
		{
			fileName: "pprof_computePrimeChunkedAtomic",
			label:    "computed prime chunked atomic",
			compute:  func() { computePrimeChunkedAtomic(testValue) },
		},
	}

	for _, profile := range profiles {
		if err := runCPUProfile(profile.fileName, profile.compute); err != nil {
			log.Fatalf("failed to write %s: %v", profile.fileName, err)
		}
		fmt.Println(profile.label)
	}
}

func runCPUProfile(fileName string, compute func()) (err error) {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create profile file: %w", err)
	}

	profileStarted := false
	fileClosed := false
	defer func() {
		if profileStarted {
			pprof.StopCPUProfile()
		}
		if !fileClosed {
			if closeErr := f.Close(); err == nil && closeErr != nil {
				err = fmt.Errorf("close profile file: %w", closeErr)
			}
		}
	}()

	if err := pprof.StartCPUProfile(f); err != nil {
		return fmt.Errorf("start cpu profile: %w", err)
	}
	profileStarted = true

	compute()

	pprof.StopCPUProfile()
	profileStarted = false

	if err := f.Close(); err != nil {
		return fmt.Errorf("close profile file: %w", err)
	}
	fileClosed = true

	return nil
}
