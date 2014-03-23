package main

func evaluationFunc1(a []int) (sum float64)  {
	for i := 0; i < len(a); i++ {
		sum += float64(a[i]*(i+1))
	}
	return sum
}
