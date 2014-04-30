package main

import (
	"fmt"
	"math"
)

func main(){
	
	for i := 1; i < 100; i++ {
		is_prime := true;
		for k := 2; k < int(math.Sqrt(float64(i))); k++ {
			if i % k == 0 {
				is_prime = false;
			}
		}

		if is_prime {
			fmt.Println(i)
		}
	}
}
