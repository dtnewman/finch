package main

import "fmt"
import "time"
import "math/rand"
import "math"

/*  A basic hill climbing algorithm that looks at neighbors and then
*	goes to the neighbor at each step that most increases the objective
*	function 
*	@param[in] current_solution The starting point
*	@param[in] evaluate The function to evaluate a solution. This function
*	takes in an slice of integers and evaluates it, returning a float64 score
*	@param[in] get_neighbors The function that takes a given solution (which is 
*	an slice of integers) and returns a 2d slice containing all of that solution's
*	neighbors
*	@return Returns a slice containing the best solution as well as a float64 with
*	that solution's evaluated score
*/
func hill_climb(initial_solution []int, evaluate func([]int) float64, 
				get_neighbors func([]int) [][]int) ([]int, float64){
	current_solution := make([]int,len(initial_solution))
	copy(current_solution,initial_solution)

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
	return current_solution, highest_score
}

/*  This function actually calls the hill climb function given a starting solution. This function is made to be
*	called as a goroutine. It returns the best solution slice and that slice's evaluated score through communication
*	channels 
*/
func hill_climb_go_routine(current_solution [] int, evaluate func([]int) float64, get_neighbors func([]int) [][]int, ch1 chan<- []int, ch2 chan <- float64) {
	best_solution, highest_score := hill_climb(current_solution, evaluate, get_neighbors)
	ch1 <- best_solution
	ch2 <- highest_score
}

/*  This function calls the hill climbing function through goroutines (using hill_climb_go_routine()). For each
*	call, it starts from a new randomly generated starting point. It will run hill climbing with @a num_restarts
*	starting points
*/
func random_restart_hill_climb(num_restarts int, evaluate func([]int) float64, create_random func() []int, get_neighbors func([]int) [][]int)([]int, float64) {
	var highest_score float64
	highest_score = math.Inf(-1)
	ch1 := make(chan []int)
	ch2 := make(chan float64)

	for i := 0; i < num_restarts; i++ {
		go hill_climb_go_routine(create_random(),evaluate, get_neighbors, ch1, ch2)
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

	return best_solution, highest_score
}






func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("\nRUN ON SIMPLE FUNCTION")
	p := []int{2, 3, 5, 4, 1, 6}
	fmt.Println("p",p)

	best_solution, highest_score := hill_climb(p,simple_evaluation,simple_get_neighbors)
    fmt.Println("hill climb results", best_solution, highest_score)
    best_solution, highest_score = random_restart_hill_climb(1000,simple_evaluation,simple_create_random_start,simple_get_neighbors)
    fmt.Println("random restart hill climb results", best_solution, highest_score, "\n")

    fmt.Println("RUN ON TSP")
    tsp_setup_data()
    p2 := make([]int,len(g_data))
    for i := 0; i < len(g_data); i++ {
    	p2[i] = i
    }
    fmt.Println("Initial distance:", -tsp_evaluation(p2))
    best_solution, highest_score = hill_climb(p2,tsp_evaluation,tsp_get_neighbors)
    fmt.Println("Optimized distance (regular hill climb):", -highest_score)
    best_solution, highest_score = random_restart_hill_climb(1000,tsp_evaluation,tsp_create_random_start,tsp_get_neighbors)
    fmt.Println("Optimized distance (random restart hill climb):", -highest_score)
}