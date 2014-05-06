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
	size := 6
	random_start := make([]int, size)
	for i := 0; i < len(random_start); i++ {
		random_start[i] = random_int(1, 10)
	}
	return random_start
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