package mathutil

import (
	"math"
)

// MaxInt finds the maximum interger from geven integers
// returns 0 if no integer passed
func MaxInt(ints ...int) int {
	if len(ints) == 0 {
		return 0
	}

	max := math.MinInt64
	for _, i := range ints {
		if i > max {
			max = i
		}
	}

	return max
}

// MinInt finds the minimum integer from given integers
// returns 0 if no integer passed
func MinInt(ints ...int) int {
	if len(ints) == 0 {
		return 0
	}

	min := math.MaxInt64
	for _, i := range ints {
		if i < min {
			min = i
		}
	}

	return min
}
