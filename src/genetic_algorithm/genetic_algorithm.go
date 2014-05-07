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





/** Genetic algorithm implementation */
func genetic_algorithm(population_size int, evaluate func([]int) float64,
	create_random func() []int , mutate func([]int, float64) []int,
	crossover func([]int,[]int) []int, 	mutation_rate float64, num_iterations int,
	num_elite int) ([]int, float64) {
	// declare variables
	var max_fitness float64 = math.Inf(-1)
	
	var sum_fitnesses float64
	var cum_probability float64
	var rand_float float64


	// create a population of population_size random starting points
	population := make([][]int, 0)
	
	// create a 2d array that holds fitnesses in the first row and numbers identifying which
	// individual in population that each fitness refers to
	fitnesses := make([][]float64, 0)
	fitnesses = append(fitnesses, make([]float64, population_size))
	fitnesses = append(fitnesses, make([]float64, population_size))

	// create a starting population
	for i := 0; i < population_size; i++ {
		population = append(population, create_random())
	}

	// create a variable to hold the best individual
	best_individual := make([]int, len(population[0]))

	// create variables to hold chromosomes for parents 
	parent1 := make([]int, len(population[0]))
	parent2 := make([]int, len(population[0]))

	// cycle through the problem for num_iterations
	for iter := 0; iter < num_iterations; iter++ {
		// start out with a cleared next generation
		next_generation_population := make([][]int, 0)

		for i := 0; i < population_size; i++ {
			fitnesses[0][i] = evaluate(population[i])
			fitnesses[1][i] = float64(i)
		}
		
		// sort the fitness array by value in the first row, which holds fitness scores
		fitnesses = qsort_2d(fitnesses, 0, "descending")
		
		// update best solution
		if (fitnesses[0][0] > max_fitness) {
			max_fitness = fitnesses[0][0]
			copy(best_individual,population[int(fitnesses[1][0])])
		}
		
		// move the elite solutions to the next generation
		for i := 0; i < num_elite; i++ {
			next_generation_population = append(next_generation_population,population[int(fitnesses[1][i])])
		}


		// we're gonna pick based on order, with items first getting higher probability. First, let's get
		// selection probabilities for each individual in the population
		sum_fitnesses = float64(len(fitnesses[0])*(len(fitnesses[0])+1)) / 2.0

		selection_probability := make([]float64, len(fitnesses[0]))
		cum_probability = 0.0

		for i, _ := range fitnesses[0] {
			selection_probability[i] = cum_probability + float64(len(fitnesses[0])-i)/sum_fitnesses
			cum_probability += float64(len(fitnesses[0])-i) / sum_fitnesses
		}

		// now fill in the rest of the population with breeding
		for i := num_elite; i < population_size; i++{
			// pick the parents
			rand_float = rand.Float64()
			j := 0
			for j = 0; selection_probability[j] < rand_float; j++ {
			}
			parent1 = population[int(fitnesses[1][j])]
			rand_float = rand.Float64()
			j = 0
			for j = 0; selection_probability[j] < rand_float; j++ {
			}
			parent2 = population[int(fitnesses[1][j])]
		

			// generate a child with crossover breeding
			child := crossover(parent1,parent2)
			
			// mutate the child
			child = mutate(child, mutation_rate)
			next_generation_population = append(next_generation_population,child)

		}

		// copy the next generation into the population variable
		for i := 0; i < population_size; i++ {
			copy(population[i], next_generation_population[i])
		}
	}

	return best_individual, max_fitness
}


// NOTE: functions below that are not found above can be found in sample_functions.go
func main() {
	rand.Seed(time.Now().Unix())
	// run the problem on our "simple" function, where we try take an array of values and try to set them to
	// values between 1 and 10, in order to maximize an objective function sum(x_i*i)
	fmt.Println("\nRUN ON SIMPLE FUNCTION")
	best_solution, highest_score := genetic_algorithm(10,simple_evaluation,simple_create_random_start,simple_mutate,simple_crossover,0.1,1000,1)
	fmt.Println("genetic algorithm results", best_solution, highest_score)

	
	fmt.Println("\nRUN ON TSP")

	tsp_setup_data()
	best_solution, highest_score = genetic_algorithm(5, tsp_evaluation, tsp_create_random_start, tsp_mutate,tsp_crossover,0.02,10000,2)
	fmt.Println("genetic algorithm results", best_solution, -highest_score)
	plotTSP(best_solution, "tsp_ga.png")
}
