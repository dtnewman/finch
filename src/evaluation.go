//package evaluation_func
package main

import "fmt"
import "math"


func EvaluationFunc1(a []int) (sum float64)  {
	for i := 0; i < len(a); i++ {
		sum += float64(a[i]*(i+1))
	}
	return sum

}

/** get all neighbors of a given solution */
func get_neighbors(currentSolution []int)( [][]int) {
	minval := 1
	maxval := 10

	neighbors := make([][]int,0)
	temp := make([]int,len(currentSolution))

	for i ,value := range currentSolution {
		if value > minval {
			copy(temp,currentSolution)
			temp[i] = value-1
			neighbors = append(neighbors,make([]int,len(currentSolution)))
			copy(neighbors[len(neighbors)-1],temp)
		}
		if value < maxval {
			copy(temp,currentSolution)
			temp[i] = value+1
			neighbors = append(neighbors,make([]int,len(currentSolution)))
			copy(neighbors[len(neighbors)-1],temp)
		}
	}
	return neighbors
}


func hill_climb(currentSolution []int)(highest_score float64){
    var score float64
    highest_score = EvaluationFunc1(currentSolution)
	var highest_score_position int
	
	for {
		highest_score_position = -1
		neighbors := get_neighbors(currentSolution)
	    for i ,value := range neighbors {
	    	score = EvaluationFunc1(value)
	    	if score > highest_score {
	    		highest_score = score
	    		highest_score_position = i
	    	}
	    }
	    if highest_score_position < 0 {
	    	break
	    } else {
	    	copy(currentSolution,neighbors[highest_score_position])
	    }
	    
	}
    fmt.Println(currentSolution, highest_score)
	return highest_score

}



func main() {
	var _ = math.Pi // delete
	p := []int{2, 3, 5, 7, 1}
	fmt.Println(p)
    fmt.Println(hill_climb(p))
}

