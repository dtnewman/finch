package main


import "fmt"
import "time"
import "math/rand"
import "math"



// the recursive loop called by qsort_2d below
func qsort_inner(a []float64,  b []int) ([]float64, []int) {
	if len(a) < 2 { return a,b }

	left, right := 0, len(a) - 1

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
	qsort_inner(a[:left],b[:left])
	qsort_inner(a[left + 1:], b[left + 1:])
	
	return a,b
}



// takes in a 2d array and an index of the row to sort by
// and returns a 2d array with all rows sorted by the 
// index row. Assumes that the input array is square (nxn)
func qsort_2d(a_input [][]float64, idx int, ascend_or_desc string) [][]float64 {
	
	// copy a_input into a
	a := make([][]float64,len(a_input))

	for i := range a {
		a[i] = make([]float64, len(a_input[i]))
		copy(a[i],a_input[i])
	}

	// throw error message if ascend_or_desc is not set right
	// if ascend_or_desc = "descending", multiply every value in it
	// by -1 and then sort that ascending
	if ascend_or_desc == "ascending" {
	} else if ascend_or_desc == "descending" {
		for i:=0; i < len(a[idx]); i++ {
			a[idx][i] *= -1
		} 
	} else {
		fmt.Println("ERROR: ascend_or_desc in qsort_2d function must have value of 'ascending' or 'descending'")
	}

	b := make([]int, len(a[idx]))
	for i:=0; i < len(a[idx]); i++ {
		b[i] = i
	}

	// sort by the sorting row
	_,order := qsort_inner(a[idx],b)

	//sort all other rows
	for i:=0; i < len(a); i++ {
		if i != idx {
			temp := make([]float64,len(a[i]))
			copy(temp,a[i])
			for j:=0; j < len(order); j++ {
				a[i][j] = temp[order[j]]
			}
		}
	}

	// revert a[idx] row back to original values
	if ascend_or_desc == "descending" {
		for i:=0; i < len(a[idx]); i++ {
			a[idx][i] *= -1
		} 
	}
	return a
}

// check if two slices (of ints) are equal
func IntsEquals(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

// check if a candidate solution is in the tabu list
func inTabuList(candidate []int, tabuList [][]int) bool {
	for _, val := range tabuList {
		if IntsEquals(candidate, val) {
			return true
		}
	}
	return false
}


func tabu_search(initial_solution []int, evaluate func([]int) float64, 
				get_neighbors func([]int) [][]int) ([]int, float64){
	current_solution := make([]int,len(initial_solution))
	copy(current_solution,initial_solution)


	var max_fitness float64 = math.Inf(-1) // will hold max fitness in each generation
	var best_fitness_yet float64 = math.Inf(-1) // will hold highest fitness yet seen
	best_candidate_yet := make([]int,0) // will hold solution with highest fitness yet seen
	tabu_list := make([][]int,0)

	
	num_neighbors := 30
	tabu_list_max_size := 20
	num_iterations_no_improvement := 0
	max_iterations_no_improvement := 100


	for num_iterations_no_improvement < max_iterations_no_improvement {
		//fmt.Println("ITERATION:", 5001-stop_condition)
		neighbors := get_neighbors(current_solution)
		//fmt.Println("Num neighbors", len(neighbors))
	    
		// evaluate all the solutions and store values in fitness array
		fitnesses := make([][]float64,0)
		fitnesses = append(fitnesses,make([]float64,len(neighbors)))
		fitnesses = append(fitnesses,make([]float64,len(neighbors)))

		for i, _ := range fitnesses[0] {
			fitnesses[0][i] = evaluate(neighbors[i])
			fitnesses[1][i] = float64(i)
		}

		// sort the fitness array by value in the first row, which holds fitness scores
		fitnesses = qsort_2d(fitnesses,0, "descending")

		fitnesses[0] = fitnesses[0][:num_neighbors]
		fitnesses[1] = fitnesses[1][:num_neighbors]

		max_fitness = fitnesses[0][0]
		num_iterations_no_improvement += 1

		// update best_fitness_yet and best_candidate_yet
		if max_fitness > best_fitness_yet {
			best_fitness_yet = max_fitness
			best_candidate_yet = make([]int,len(neighbors[int(fitnesses[1][0])]))
			copy(best_candidate_yet,neighbors[int(fitnesses[1][0])])
			num_iterations_no_improvement = 0
		}

		// find the index of the next candidate in the fitness list. Keep iterating
		// until a candidate that is not in the tabu list is found. If none are found
		// then break
		next_candidate_idx := 0
		for inTabuList(neighbors[int(fitnesses[1][next_candidate_idx])],tabu_list) {
			next_candidate_idx += 1
			// if next_candidate_idx == num_neighbors means that all of the neighbors
			// are in the tabu list
			if next_candidate_idx == num_neighbors {
				fmt.Println("All candidates tabued, stuck on a local or global maximum")
				return best_candidate_yet, best_fitness_yet
			}
		}
		//fmt.Println("next_candidate_idx", next_candidate_idx)

		tabu_list = append(tabu_list, neighbors[int(fitnesses[1][next_candidate_idx])])
		if len(tabu_list)>tabu_list_max_size {
			tabu_list = tabu_list[1:]

		}


		//fmt.Println("tabu list size:", len(tabu_list))
		current_solution = neighbors[int(fitnesses[1][next_candidate_idx])]

		//fmt.Println("best fitness:", max_fitness)


	    
	}
	return best_candidate_yet, best_fitness_yet
}



// NOTE: functions below that are not found above can be found in sample_functions.go
func main() {
	rand.Seed(time.Now().Unix())

	// run the problem on our "simple" function, where we try take an array of values and try to set them to
	// values between 1 and 10, in order to maximize an objective function sum(x_i*i)
	//fmt.Println("\nRUN ON SIMPLE FUNCTION")
	//p := []int{2, 3, 5, 4, 1, 6}
	//fmt.Println("p",p)

	//best_solution, highest_score := hill_climb(p,simple_evaluation,simple_get_neighbors)
    //fmt.Println("hill climb results", best_solution, highest_score)
    //best_solution, highest_score = random_restart_hill_climb(1000,simple_evaluation,simple_create_random_start,simple_get_neighbors)
    //fmt.Println("random restart hill climb results", best_solution, highest_score)
    //best_solution, highest_score = stochastic_hill_climb(p,simple_evaluation,simple_get_neighbors)
    //fmt.Println("stochastic hill climb results", best_solution, highest_score, "\n")

    // Run on a travelling salesman problem with cities in the file tsp_data.csv (40 cities)
    fmt.Println("RUN ON TSP")
    tsp_setup_data()
    p2 := make([]int,len(g_data))
    for i := 0; i < len(g_data); i++ {
    	p2[i] = i
    }
    for i := 0; i < 10; i++ {
    	p2 = tsp_create_random_start()
	    fmt.Println("Initial distance:", -tsp_evaluation(p2))
	    best_solution, highest_score := hill_climb(p2,tsp_evaluation,tsp_get_neighbors)
	    fmt.Println("Optimized distance (regular hill climb):", -highest_score)
	    best_solution, highest_score = tabu_search(p2,tsp_evaluation,tsp_get_neighbors)
	    fmt.Println("tabu search results", best_solution, -highest_score, "\n")
    }

    /*best_solution, highest_score = random_restart_hill_climb(10,tsp_evaluation,tsp_create_random_start,tsp_get_neighbors)
    fmt.Println("Optimized distance (random restart hill climb):", -highest_score)
    best_solution, highest_score = stochastic_hill_climb(p2,tsp_evaluation,tsp_get_neighbors)
    fmt.Println("Optimized distance (stochastic hill climb):", -highest_score)*/


    /*
    tabulist := make([][]int, 3)
    tabulist[0] = make([]int, 3)
    tabulist[1] = make([]int, 3)
    tabulist[2] = make([]int, 3)
    tabulist[0][0] = 1
    tabulist[0][1] = 1
    tabulist[0][2] = 1
    tabulist[1][0] = 2
    tabulist[1][1] = 2
    tabulist[1][2] = 2
    tabulist[2][0] = 3
    tabulist[2][1] = 3
    tabulist[2][2] = 3

    test := make([]int,3)
    test[0] = 4
    test[1] = 4
    test[2] = 4
    fmt.Println(inTabuList(test,tabulist))
    */

}