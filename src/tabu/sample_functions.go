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

/** return a random integer between min and max */
func random_int(min, max int) int {
	return rand.Intn(max-min) + min
}

/** get all neighbors of a given solution. Neighbors in this function
are defined as all values where just one integer in the array is
incremented or decremented by 1, and only where integer values remain
between 1 and 10 */
func simple_get_neighbors(currentSolution []int) [][]int {
	minval := 1
	maxval := 10

	neighbors := make([][]int, 0)
	temp := make([]int, len(currentSolution))

	for i, value := range currentSolution {
		if value > minval {
			copy(temp, currentSolution)
			temp[i] = value - 1
			neighbors = append(neighbors, make([]int, len(currentSolution)))
			copy(neighbors[len(neighbors)-1], temp)
		}
		if value < maxval {
			copy(temp, currentSolution)
			temp[i] = value + 1
			neighbors = append(neighbors, make([]int, len(currentSolution)))
			copy(neighbors[len(neighbors)-1], temp)
		}
	}
	return neighbors
}

func simple_evaluation(a []int) (sum float64) {
	for i := 0; i < len(a); i++ {
		sum += float64(a[i]) //*(i+1))
	}
	return sum
}

/** return an array with 5 random integers between 1 and 10 */
func simple_create_random_start() []int {
	size := 3
	random_start := make([]int, size)
	for i := 0; i < len(random_start); i++ {
		random_start[i] = random_int(1, 10)
	}
	return random_start
}

func tsp_get_neighbors(currentSolution []int) [][]int {
	neighbors := make([][]int, 0)
	temp_neighbor := make([]int, len(currentSolution))
	var temp_int int

	for i := 0; i < len(currentSolution); i++ {
		for j := i + 1; j < len(currentSolution); j++ {
			copy(temp_neighbor, currentSolution)
			temp_int = temp_neighbor[i]
			temp_neighbor[i] = temp_neighbor[j]
			temp_neighbor[j] = temp_int
			neighbors = append(neighbors, make([]int, len(currentSolution)))
			copy(neighbors[len(neighbors)-1], temp_neighbor)
		}
	}
	return neighbors
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

/*  Hill climbing algorithm for comparison
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
		//fmt.Println(highest_score)

	}
	return current_solution, highest_score
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
