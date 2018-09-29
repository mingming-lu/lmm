package mathutil

import (
	"math"
)

// MinInt finds the minimum integer from given integers
// returns 0 if no integer passed
func MinInt(ints ...int) int {
	if len(ints) == 0 {
		return 0
	}

	min := math.MaxInt32
	for _, i := range ints {
		if i < min {
			min = i
		}
	}

	return min
}
