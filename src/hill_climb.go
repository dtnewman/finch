//package evaluation_func
package main

import "fmt"
import "math"
import "math/rand"
import "time"

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

func random_int(min, max int) int {

    return rand.Intn(max - min) + min
}

func create_random_start()([] int) {
	size := 5
	random_start := make([]int,size)
	for i:=0; i < len(random_start); i++ {
		random_start[i] = random_int(1,10)
	}
	return random_start
}


func hill_climb_go_routine(currentSolution [] int, ch chan<- float64) {
	time.Sleep(time.Duration(currentSolution[0] * 1e9))
	ch <- float64(currentSolution[0])
}

func random_restart_hill_climb(num_restarts int)(highest_score float64) {
	ch := make(chan float64)
	
	for i := 0; i < num_restarts; i++ {
		go hill_climb_go_routine(create_random_start(),ch)
	}

	var score float64;

	for i := 0; i < num_restarts; i++ {
        score = <-ch
        if score > highest_score {
        	highest_score = score
        }
    }
	return highest_score
}


func main() {
	rand.Seed(time.Now().Unix())
	var _ = math.Pi // delete
	p := []int{2, 3, 5, 7, 1}
	fmt.Println(p)
    fmt.Println(hill_climb(p))
    fmt.Println(create_random_start())
    fmt.Println(random_restart_hill_climb(2))
}

