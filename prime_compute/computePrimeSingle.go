package main

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
