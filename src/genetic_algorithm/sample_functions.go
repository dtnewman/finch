package main

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

/** return a random integer between (and including both) min and max */
func random_int(min, max int) int {
	return rand.Intn(max+1-min) + min
}


/** Take in a solution for the "simple" example and make one change to it */
func simple_make_change(current_solution []int) []int {
	minval := 1
	maxval := 10

	length := len(current_solution)
	return_value := make([]int, len(current_solution))
	copy(return_value, current_solution)

	value_to_change := random_int(0,length-1)
	return_value[value_to_change] = random_int(minval,maxval)

	return return_value
}

/** evaluate a solution by adding together all elements in the array
*	multiplied by those elements' positions */
func simple_evaluation(a []int) (sum float64) {
	for i := 0; i < len(a); i++ {
		sum += float64(a[i] * (i + 1))
	}
	return sum
}

/** return an array with 6 random integers between 1 and 10 */
func simple_create_random_start() []int {
	size := 10
	random_start := make([]int, size)
	for i := 0; i < len(random_start); i++ {
		random_start[i] = random_int(1, 10)
	}
	return random_start
}

/** a mutation function that takes in an individual, and mutates it according
/*  to the given mutation_rate */
func simple_mutate(input_individual []int, mutation_rate float64)([]int) {
	output_individual := make([]int, len(input_individual))
	var rand_num float64
	for i := 0; i < len(input_individual); i++ {
		rand_num = rand.Float64()
		if (rand_num < mutation_rate) {
			output_individual[i] = random_int(1,10)
		} else {
			output_individual[i] = input_individual[i]
		}
	}
	return output_individual
}

/** Uses uniform crossover to create offspring from two parents */
func simple_crossover(parent1 []int, parent2 []int)([]int) {
	var rand_num int
	child := make([]int, len(parent1))
	for i := 0; i < len(child) ; i++ {
		rand_num = rand.Intn(2)
		if (rand_num == 0) {
			child[i] = parent1[i]
		} else {
			child[i] = parent2[i]
		}
	}
	return child
}

/** Take in a solution for the TSP example and make one change to it */
func tsp_make_change(current_solution []int) []int {
	length := len(current_solution)
	return_value := make([]int, len(current_solution))
	copy(return_value, current_solution)
	random_int_1 := random_int(0,length-1)
	random_int_2 := random_int(0,length-1)
	temp := return_value[random_int_1]
	return_value[random_int_1] = return_value[random_int_2]
	return_value[random_int_2] = temp

	return return_value
}

/** a mutation function that takes in an individual, and mutates it according
/*  to the given mutation_rate */
func tsp_mutate(input_individual []int, mutation_rate float64)([]int) {
	output_individual := make([]int, len(input_individual))
	copy(output_individual,input_individual)
	var rand_num float64
	for i := 0; i < len(input_individual); i++ {
		rand_num = rand.Float64()
		if (rand_num < mutation_rate) {
			output_individual = tsp_make_change(output_individual)
		} 
	}
	return output_individual
}

func intInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


/** Uses crossover to create offspring from two parents */
func tsp_crossover(parent1 []int, parent2 []int)([]int) {
	
	// create an empty slice to hold the child
	child := make([]int,len(parent1))

	for i := 0; i < len(child); i++ {
		child[i] = -1
	}
	
	// pick a starting point and ending point from parent 1
	start := random_int(0,len(parent1)-1)
	end := random_int(0,len(parent1)-1)

	// switch them if end < start
	if (end < start) {
		temp := end
		end = start
		start = temp
	}

	// set the elements from start to end in child to those elements in parent1
	for i := start; i <= end; i++ {
		child[i] = parent1[i]
	}

	// now add objects not in the child currently to the child in the second
	// parent's order
	position := 0
	for i := 0; i < len(child); i++{
		if (i < start || i > end) {
			for (intInSlice(parent2[position],child)) {
				position += 1
			}
			child[i] = parent2[position]
			position += 1
		}
	}
	return child
}

var g_data [][]int

func tsp_setup_data() {
	file, err := os.Open("tsp_data.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		g_data = append(g_data, make([]int, 2))
		//g_data[len(g_data)-1][0] =
		value1, err := strconv.Atoi(strings.TrimSpace(record[0]))
		value2, err := strconv.Atoi(strings.TrimSpace(record[1]))
		g_data[len(g_data)-1][0] = value1
		g_data[len(g_data)-1][1] = value2
	}
}

func norm(x1 int, y1 int, x2 int, y2 int) float64 {
	return math.Sqrt(float64(x1-x2)*float64(x1-x2) + float64(y1-y2)*float64(y1-y2))
}

func tsp_evaluation(path []int) (sum float64) {
	for i := 0; i < len(path); i++ {
		sum -= norm(g_data[path[i]][0], g_data[path[i]][1], g_data[path[(i+1)%len(path)]][0], g_data[path[(i+1)%len(path)]][1])
	}
	return sum
}

func tsp_create_random_start() []int {
	size := len(g_data)
	return rand.Perm(size)
}

// Takes in an array of which order to visit cities and plots the route. Assumes
// that g_data global variable has already been initialized with tsp_setup_data
// function.
func plotTSP(city_order []int, file_name string) {
	pts := make(plotter.XYs, len(city_order))
	for i := 0; i < len(city_order); i++ {
		pts[i].X = float64(g_data[city_order[i]][0])
		pts[i].Y = float64(g_data[city_order[i]][1])
	}

	p, _ := plot.New()
	p.Title.Text = "TSP best solution"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	plotutil.AddLinePoints(p, "", pts)
	p.X.Min = 0
	p.X.Max = 1800
	p.Y.Min = 0
	p.Y.Max = 1300
	p.Save(4, 4, file_name)
}