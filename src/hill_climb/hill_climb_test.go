package main

import "testing"

/** Test the function simple_get_neighbors() to ensure that it gives
	the correct output for the input array {1,2} */
func Test_simple_get_neighbors(t *testing.T) {
	two_arrays_equal := true
	input := []int{1, 2}
	got := simple_get_neighbors(input)
	want := make([][]int,0)
	want = append(want,[]int{2,2})
	want = append(want,[]int{1,1})
	want = append(want,[]int{1,3})

	// check that arrays are the same size
	if len(want) != len(got) {
		two_arrays_equal = false
	} else {
		for i := 0; i < len(want); i++ {
			// check length of each subarray
			if len(want[i]) != len(got[i]) {
				two_arrays_equal = false
			} else {
				//now check the content of each subarray
				for j := 0; j < len(want[i]); j++ {
					if want[i][j] != got[i][j] {
						two_arrays_equal = false
					}
				}
			}
		}
	}
	if (!two_arrays_equal) {
		t.Errorf("simple_get_neighbors(%d) == %d, want %d", input, got, want)
	}
}

/** Test the function tsp_get_neighbors() to ensure that it gives
	the correct output for the input array {1,2,3} */
func Test_tsp_get_neighbors(t *testing.T) {
	two_arrays_equal := true
	input := []int{1, 2, 3}
	got := tsp_get_neighbors(input)
	want := make([][]int,0)
	want = append(want,[]int{2,1,3})
	want = append(want,[]int{3,2,1})
	want = append(want,[]int{1,3,2})

	// check that arrays are the same size
	if len(want) != len(got) {
		two_arrays_equal = false
	} else {
		for i := 0; i < len(want); i++ {
			// check length of each subarray
			if len(want[i]) != len(got[i]) {
				two_arrays_equal = false
			} else {
				//now check the content of each subarray
				for j := 0; j < len(want[i]); j++ {
					if want[i][j] != got[i][j] {
						two_arrays_equal = false
					}
				}
			}
		}
	}
	if (!two_arrays_equal) {
		t.Errorf("simple_get_neighbors(%d) == %d, want %d", input, got, want)
	}
}