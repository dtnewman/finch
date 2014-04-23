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
	get_neighbors func([]int) [][]int) ([]int, float64) {
	current_solution := make([]int, len(initial_solution))
	copy(current_solution, initial_solution)

	var score float64
	highest_score := evaluate(current_solution)
	var highest_score_position int

	for {
		highest_score_position = -1
		neighbors := get_neighbors(current_solution)
		for i, value := range neighbors {
			score = evaluate(value)
			if score > highest_score {
				highest_score = score
				highest_score_position = i
			}
		}
		if highest_score_position < 0 {
			break
		} else {
			copy(current_solution, neighbors[highest_score_position])
		}
	}
	return current_solution, highest_score
}

/*  This function calls the hill climb function given a starting solution. This function is made to be
*	called as a goroutine. It returns the best solution slice and that slice's evaluated score through communication
*	channels
 */
func hill_climb_go_routine(current_solution []int, evaluate func([]int) float64, get_neighbors func([]int) [][]int, ch1 chan<- []int, ch2 chan<- float64) {
	best_solution, highest_score := hill_climb(current_solution, evaluate, get_neighbors)
	ch1 <- best_solution
	ch2 <- highest_score
}

/*  This function calls the hill climbing function through goroutines (using hill_climb_go_routine()). For each
*	call, it starts from a new randomly generated starting point. It will run hill climbing with @a num_restarts
*	starting points
 */
func random_restart_hill_climb(num_restarts int, evaluate func([]int) float64, create_random func() []int, get_neighbors func([]int) [][]int) ([]int, float64) {
	var highest_score float64
	highest_score = math.Inf(-1)
	ch1 := make(chan []int)
	ch2 := make(chan float64)

	for i := 0; i < num_restarts; i++ {
		go hill_climb_go_routine(create_random(), evaluate, get_neighbors, ch1, ch2)
	}

	var score float64
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

/*	This function is for the stochastic hill climbing algorithm. In this algorithm, we check all neighbors of a
*	given solution and evaluate how much they improve the score. Once we do that, we choose a random neighbor, but
*	we weight it such that the neighbors that the chances of a neighbor being picked are proportional to home much
*	that neighbor improves on the solution. The algorithm stops once it gets to a local maximum.
 */
func stochastic_hill_climb(initial_solution []int, evaluate func([]int) float64,
	get_neighbors func([]int) [][]int) ([]int, float64) {

	current_solution := make([]int, len(initial_solution))
	copy(current_solution, initial_solution)

	current_score := evaluate(current_solution)
	var total_score_improvement_sum float64
	var score float64
	var randnum float64
	var tempSum float64
	var pickPosition int

	// keep looping until it hits a local maximum
	for {
		total_score_improvement_sum = 0.0
		neighbors := get_neighbors(current_solution)
		scores_slice := make([]float64, len(neighbors))
		// iterate through each neighbor and keep track of how much it improves the solution in "scores_slice".
		// If it doesn't improve the solution, then the corresponding value in scores_slice remains 0.
		for i, value := range neighbors {
			score = evaluate(value)
			if score > current_score {
				scores_slice[i] = score - current_score
				total_score_improvement_sum += score - current_score //keep track of the sum of improvements
			}
		}
		// if there were no solutions that were improvements, then total_score_improvement_sum will be 0, which
		// means that we are at a local minimum, so this is our stopping condition
		if total_score_improvement_sum == 0.0 {
			break
		}
		// divide each element of scores_slice by total_score_improvement_sum. Now, this slice holds the
		// probabilites of each neighbor being picked
		for i, value := range scores_slice {
			scores_slice[i] = value / total_score_improvement_sum
		}

		// pick a random number between 0 and 1
		randnum = rand.Float64()
		tempSum = 0.0
		// iterate through scores_slice until the sum > randnum. At that point, pickPosition will be holding
		// the position of the neighbor that we're going to go to
		for i := 0; tempSum < randnum; i++ {
			pickPosition = i
			tempSum += scores_slice[i]
		}

		// assign the appropriate neighbor to the current_solution and then run the same process again until
		// the stopping condition is met
		copy(current_solution, neighbors[pickPosition])
		current_score = evaluate(current_solution)

	}
	return current_solution, evaluate(current_solution)
}

// NOTE: functions below that are not found above can be found in sample_functions.go
func main() {
	rand.Seed(time.Now().Unix())

	// run the problem on our "simple" function, where we try take an array of values and try to set them to
	// values between 1 and 10, in order to maximize an objective function sum(x_i*i)
	fmt.Println("\nRUN ON SIMPLE FUNCTION")
	p := []int{2, 3, 5, 4, 1, 6}
	fmt.Println("p", p)

	best_solution, highest_score := hill_climb(p, simple_evaluation, simple_get_neighbors)
	fmt.Println("hill climb results", best_solution, highest_score)
	best_solution, highest_score = random_restart_hill_climb(1000, simple_evaluation, simple_create_random_start, simple_get_neighbors)
	fmt.Println("random restart hill climb results", best_solution, highest_score)
	best_solution, highest_score = stochastic_hill_climb(p, simple_evaluation, simple_get_neighbors)
	fmt.Println("stochastic hill climb results", best_solution, highest_score, "\n")

	// Run on a travelling salesman problem with cities in the file tsp_data.csv (40 cities)
	fmt.Println("RUN ON TSP")
	tsp_setup_data()
	p2 := make([]int, len(g_data))
	for i := 0; i < len(g_data); i++ {
		p2[i] = i
	}
	fmt.Println("Initial distance:", -tsp_evaluation(p2))
	best_solution, highest_score = hill_climb(p2, tsp_evaluation, tsp_get_neighbors)
	fmt.Println("Optimized distance (regular hill climb):", -highest_score)
	best_solution, highest_score = random_restart_hill_climb(10, tsp_evaluation, tsp_create_random_start, tsp_get_neighbors)
	fmt.Println("Optimized distance (random restart hill climb):", -highest_score)
	best_solution, highest_score = stochastic_hill_climb(p2, tsp_evaluation, tsp_get_neighbors)
	fmt.Println("Optimized distance (stochastic hill climb):", -highest_score)
}
