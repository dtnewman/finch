package main

import "fmt"
import "time"
import "math/rand"
import "math"

// the recursive loop called by qsort_2d below
func qsort_inner(a []float64, b []int) ([]float64, []int) {
	if len(a) < 2 {
		return a, b
	}

	left, right := 0, len(a)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(a)

	// Move the pivot to the right
	a[pivotIndex], a[right] = a[right], a[pivotIndex]
	b[pivotIndex], b[right] = b[right], b[pivotIndex]

	// Pile elements smaller than the pivot on the left
	for i := range a {
		if a[i] < a[right] {
			a[i], a[left] = a[left], a[i]
			b[i], b[left] = b[left], b[i]
			left++
		}
	}
	// Place the pivot after the last smaller element
	a[left], a[right] = a[right], a[left]
	b[left], b[right] = b[right], b[left]

	// Go down the rabbit hole
	qsort_inner(a[:left], b[:left])
	qsort_inner(a[left+1:], b[left+1:])

	return a, b
}

// takes in a 2d array and an index of the row to sort by
// and returns a 2d array with all rows sorted by the
// index row. Assumes that the input array is square (nxn)
func qsort_2d(a_input [][]float64, idx int, ascend_or_desc string) [][]float64 {

	// copy a_input into a
	a := make([][]float64, len(a_input))

	for i := range a {
		a[i] = make([]float64, len(a_input[i]))
		copy(a[i], a_input[i])
	}

	// throw error message if ascend_or_desc is not set right
	// if ascend_or_desc = "descending", multiply every value in it
	// by -1 and then sort that ascending
	if ascend_or_desc == "ascending" {
	} else if ascend_or_desc == "descending" {
		for i := 0; i < len(a[idx]); i++ {
			a[idx][i] *= -1
		}
	} else {
		fmt.Println("ERROR: ascend_or_desc in qsort_2d function must have value of 'ascending' or 'descending'")
	}

	b := make([]int, len(a[idx]))
	for i := 0; i < len(a[idx]); i++ {
		b[i] = i
	}

	// sort by the sorting row
	_, order := qsort_inner(a[idx], b)

	//sort all other rows
	for i := 0; i < len(a); i++ {
		if i != idx {
			temp := make([]float64, len(a[i]))
			copy(temp, a[i])
			for j := 0; j < len(order); j++ {
				a[i][j] = temp[order[j]]
			}
		}
	}

	// revert a[idx] row back to original values
	if ascend_or_desc == "descending" {
		for i := 0; i < len(a[idx]); i++ {
			a[idx][i] *= -1
		}
	}
	return a
}


func genetic_algorithm(population_size int, evaluate func([]int) float64,
	create_random func() []int, get_neighbors func([]int) [][]int) ([]int, float64) {
	// decalre variables
	var next_generation_candidates [][]int
	var neighbors [][]int
	var max_fitness float64 = math.Inf(-1)

	// create a population of population_size random starting points
	population := make([][]int, 0)
	for i := 0; i < population_size; i++ {
		population = append(population, create_random())
	}
	for {
		// get all the neighbors for all beams and put them into next_generation_candidates
		next_generation_candidates = make([][]int, 0)
		for i := 0; i < num_beams; i++ {
			neighbors = get_neighbors(population[i])
			for _, value := range neighbors {
				next_generation_candidates = append(next_generation_candidates, value)
			}
		}

		// evaluate all the solutions and store values in fitness array
		fitnesses := make([][]float64, 0)
		fitnesses = append(fitnesses, make([]float64, len(next_generation_candidates)))
		fitnesses = append(fitnesses, make([]float64, len(next_generation_candidates)))

		for i, _ := range fitnesses[0] {
			fitnesses[0][i] = evaluate(next_generation_candidates[i])
			fitnesses[1][i] = float64(i)
		}

		// sort the fitness array by value in the first row, which holds fitness scores
		fitnesses = qsort_2d(fitnesses, 0, "descending")
		// Stop the loop if the highest fitness is not greater than max_fitness
		if fitnesses[0][0] <= max_fitness {
			break
		}

		// Now fill up the next generation
		for i := 0; i < num_beams; i++ {
			copy(population[i], next_generation_candidates[int(fitnesses[1][i])])
		}

		// Set max_fitness to (-1) times the first element in the first row of fitnesses
		// since that indicates the leading value
		max_fitness = fitnesses[0][0]
	}

	return population[0], max_fitness
}


// NOTE: functions below that are not found above can be found in sample_functions.go
func main() {
	rand.Seed(time.Now().Unix())
	// run the problem on our "simple" function, where we try take an array of values and try to set them to
	// values between 1 and 10, in order to maximize an objective function sum(x_i*i)
	/*fmt.Println("\nRUN ON SIMPLE FUNCTION")
	best_solution, highest_score := beam_search(5,simple_evaluation,simple_create_random_start,simple_get_neighbors)
	fmt.Println("beam search results", best_solution, highest_score)
	best_solution, highest_score = stochastic_beam_search(2,simple_evaluation,simple_create_random_start,simple_get_neighbors)
	fmt.Println("stochastic beam search results", best_solution, highest_score)
	*/
	fmt.Println("\nRUN ON TSP")
	tsp_setup_data()
	best_solution, highest_score := genetic_algorithm(50, tsp_evaluation, tsp_create_random_start, tsp_get_neighbors)
	fmt.Println("beam search results", best_solution, -highest_score)
	plotTSP(best_solution, "tsp_beam_search.png")
	best_solution, highest_score = stochastic_beam_search(50, 200, tsp_evaluation, tsp_create_random_start, tsp_get_neighbors)
	fmt.Println("stochastic beam search results", best_solution, -highest_score)
	plotTSP(best_solution, "tsp_stochastic_beam_search.png")

}
