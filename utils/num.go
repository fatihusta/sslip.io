// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

// MinInt returns the minimum value between a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MakeIndices returns a slice of the given size, where each element has
// the same value of its own index position.
func MakeIndices(size int) []int {
	indices := make([]int, size)
	for i := range indices {
		indices[i] = i
	}
	return indices
}

// MakeIntMatrix returns a new 2-dimensional slice of int.
func MakeIntMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
	}
	return matrix
}

// ContainsInt returns whether the list contains the x-element, or not.
func ContainsInt(lst []int, x int) bool {
	for _, element := range lst {
		if element == x {
			return true
		}
	}
	return false
}

// IntSliceEqual returns whether the two slices are equal, or not.
func IntSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, va := range a {
		if va != b[i] {
			return false
		}
	}
	return true
}