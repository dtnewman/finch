package main

import "fmt"
import "time"
import "math/rand"
import "math"

func make_changes(current_solution []int, make_change func([]int) []int, num_changes int)([] int) {
	new_solution := make([]int, len(current_solution))
	copy(new_solution, current_solution)
	for i:=0;i<num_changes;i++ {
		new_solution = make_change(new_solution)
	}
	return new_solution
}


/*  A simulated annealing algorithm that makes random changes in search of a global maximum
*	@param[in] initial_solution The starting point
*	@param[in] evaluate The function to evaluate a solution. This function
*	takes in an slice of integers and evaluates it, returning a float64 score
*	@param[in] make_change The function that takes a given solution (which is
*	an slice of integers) and returns a solution modified by 1
	@param[in] init_temp The starting temperature for the algorithm
	@param[in] thermostat The thermostat level for determine how quickly to decrease temperature
	@param[in] itol The max number of iterations. Once this is reached, the algorithm stops
	@param[in] reannealing The number of rounds before trying reannealing 
*	@return Returns a slice containing the best solution as well as a float64 with
*	that solution's evaluated score
 */
func simulated_annealing(initial_solution []int, evaluate func([]int) float64, make_change func([]int) []int,
						 init_temp float64, thermostat float64, itol int, reannealing int)([]int, float64) {

	current_solution := make([]int, len(initial_solution))
	proposed_solution := make([]int, len(initial_solution))
	return_solution := make([]int, len(initial_solution))
	copy(current_solution, initial_solution)
	copy(proposed_solution, initial_solution)
	copy(return_solution, initial_solution)

	var prev_score float64
	var new_score float64
	var delta_score float64
	highest_score := evaluate(current_solution)
	prev_score = highest_score

	temperature := init_temp 
	L := 0
	num_acceptances := 0


	for i := 0; i < itol; i++ {
		L = int(math.Sqrt(temperature))
		proposed_solution = make_changes(current_solution, make_change, L)
		new_score = evaluate(proposed_solution)
		delta_score = new_score-prev_score

		// If the proposed solution is better than the current solution, then
		// accept it. Otherwise, accept with probability E^(delta_score/temperature)
		if (new_score > prev_score) {
			copy(current_solution, proposed_solution)
			prev_score = new_score 
			if (new_score > highest_score) {
				highest_score = new_score
				copy(return_solution, proposed_solution)
			}
			num_acceptances++
		} else if (rand.Float64() < math.Exp(delta_score/temperature)) {
			copy(current_solution, proposed_solution)
			prev_score = new_score 
			num_acceptances++
		} 

		// adjust temperature after the number of acceptances determined by reannealing parameter
		if (num_acceptances % reannealing == 0) {
			temperature = thermostat * temperature
			// If the temperature goes below 1, reset it to the initial temperature
			if (temperature < 1){
				temperature = init_temp
			}
		}


	} 

	// Return the best solution and the highest score
	return return_solution, highest_score
}


// NOTE: functions below that are not found above can be found in sample_functions.go
func main() {
	rand.Seed(time.Now().Unix())

	// run the problem on our "simple" function, where we try take an array of values and try to set them to
	// values between 1 and 10, in order to maximize an objective function sum(x_i*i)
	fmt.Println("\nRUN ON SIMPLE FUNCTION")
	p := []int{2, 3, 5, 4, 1, 6}
	fmt.Println("p", p)

	best_solution, highest_score := simulated_annealing(p, simple_evaluation,simple_make_change,10,0.9,1000,10)
	fmt.Println("Simulated annealing results", best_solution, highest_score)
	
	// Run on a travelling salesman problem with cities in the file tsp_data.csv (40 cities)
	/*fmt.Println("RUN ON TSP")
	tsp_setup_data()
	p2 := make([]int, len(g_data))
	for i := 0; i < len(g_data); i++ {
		p2[i] = i
	}
	fmt.Println("Initial distance:", -tsp_evaluation(p2))
	best_solution, highest_score = simulated_annealing(p2, tsp_evaluation)
	fmt.Println("Simulated annealing results):", -highest_score)*/
}
