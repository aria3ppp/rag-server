// TODO: remove this package
// Deprecated: remove this package
package slices

import "cmp"

func MaxIndex[T cmp.Ordered](collection []T) (maxValue T, maxIndex int) {
	if len(collection) == 0 {
		return maxValue, -1
	}

	maxValue = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item > maxValue {
			maxValue = item
			maxIndex = i
		}
	}

	return maxValue, maxIndex
}

// greaterFn(a, b) return true if a should be greater than b
func MaxIndexBy[T any](collection []T, greaterFn func(T, T) bool) (maxValue T, maxIndex int) {
	if len(collection) == 0 {
		return maxValue, -1
	}

	maxValue = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if greaterFn(item, maxValue) {
			maxValue = item
			maxIndex = i
		}
	}

	return maxValue, maxIndex
}
