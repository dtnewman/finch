package main



/** get all neighbors of a given solution. Neighbors in this function
	are defined as all values where just one integer in the array is 
	incremented or decremented by 1, and only where integer values remain
	between 1 and 10 */
func simple_get_neighbors(currentSolution []int)([][]int) {
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