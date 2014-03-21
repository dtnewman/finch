//package evaluation_func
package main

import "fmt"
//import "math"
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


func hill_climb(current_solution []int, evaluate func([]int) float64)([]int, float64){
    var score float64
    highest_score := evaluate(current_solution)
	var highest_score_position int
	
	for {
		highest_score_position = -1
		neighbors := get_neighbors(current_solution)
	    for i ,value := range neighbors {
	    	score = evaluate(value)
	    	if score > highest_score {
	    		highest_score = score
	    		highest_score_position = i
	    	}
	    }
	    if highest_score_position < 0 {
	    	break
	    } else {
	    	copy(current_solution,neighbors[highest_score_position])
	    }
	    
	}
    //fmt.Println(current_solution, highest_score)
	return current_solution, highest_score
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


func hill_climb_go_routine(current_solution [] int, evaluate func([]int) float64, ch1 chan<- []int, ch2 chan <- float64) {
	best_solution, highest_score := hill_climb(current_solution, evaluate)
	ch1 <- best_solution
	ch2 <- highest_score
}

func random_restart_hill_climb(num_restarts int, evaluate func([]int) float64, create_random func() []int)([]int, float64) {
	
	var highest_score float64
	ch1 := make(chan []int)
	ch2 := make(chan float64)

	for i := 0; i < num_restarts; i++ {
		go hill_climb_go_routine(create_random(),evaluate,ch1, ch2)
	}

	var score float64;
	var best_solution []int

	for i := 0; i < num_restarts; i++ {
        best_solution = <-ch1
        score = <-ch2
        if score > highest_score {
        	highest_score = score
        }
    }

	return best_solution,highest_score
}


func main() {
	rand.Seed(time.Now().Unix())
	//var _ = math.Pi // delete
	p := []int{2, 3, 5, 7, 1}
	fmt.Println("p",p)
	best_solution, highest_score := hill_climb(p,EvaluationFunc1)
    fmt.Println("hill climb results", best_solution, highest_score)
    best_solution, highest_score = random_restart_hill_climb(2,EvaluationFunc1,create_random_start)
    fmt.Println("random restart hill climb results", best_solution, highest_score)

}

