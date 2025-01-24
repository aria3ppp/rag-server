// TODO: remove this package
// Deprecated: remove this package
package slices_test

import (
	"testing"

	"github.com/aria3ppp/rag-server/pkg/slices"
)

func TestMaxIndex(t *testing.T) {
	tests := []struct {
		name       string
		collection any
		wantValue  any
		wantIndex  int
	}{
		{
			name:       "empty int slice",
			collection: []int{},
			wantValue:  0,
			wantIndex:  -1,
		},
		{
			name:       "single int element",
			collection: []int{42},
			wantValue:  42,
			wantIndex:  0,
		},
		{
			name:       "multiple ints",
			collection: []int{1, 5, 3, 2, 4},
			wantValue:  5,
			wantIndex:  1,
		},
		{
			name:       "multiple ints with duplicate max",
			collection: []int{1, 5, 3, 5, 4},
			wantValue:  5,
			wantIndex:  1, // Should return first occurrence
		},
		{
			name:       "negative ints",
			collection: []int{-5, -3, -1, -4, -2},
			wantValue:  -1,
			wantIndex:  2,
		},
		{
			name:       "float64 values",
			collection: []float64{1.1, 2.2, 3.3, 2.8, 1.5},
			wantValue:  3.3,
			wantIndex:  2,
		},
		{
			name:       "strings",
			collection: []string{"apple", "zebra", "banana", "yak"},
			wantValue:  "zebra",
			wantIndex:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch collection := tt.collection.(type) {
			case []int:
				gotValue, gotIndex := slices.MaxIndex(collection)
				if gotValue != tt.wantValue.(int) || gotIndex != tt.wantIndex {
					t.Errorf("slices.MaxIndex() = (%v, %v), want (%v, %v)",
						gotValue, gotIndex, tt.wantValue, tt.wantIndex)
				}
			case []float64:
				gotValue, gotIndex := slices.MaxIndex(collection)
				if gotValue != tt.wantValue.(float64) || gotIndex != tt.wantIndex {
					t.Errorf("slices.MaxIndex() = (%v, %v), want (%v, %v)",
						gotValue, gotIndex, tt.wantValue, tt.wantIndex)
				}
			case []string:
				gotValue, gotIndex := slices.MaxIndex(collection)
				if gotValue != tt.wantValue.(string) || gotIndex != tt.wantIndex {
					t.Errorf("slices.MaxIndex() = (%v, %v), want (%v, %v)",
						gotValue, gotIndex, tt.wantValue, tt.wantIndex)
				}
			}
		})
	}
}

func TestMaxIndexBy(t *testing.T) {
	tests := []struct {
		name       string
		collection any
		greaterFn  any
		wantValue  any
		wantIndex  int
	}{
		{
			name:       "empty int slice",
			collection: []int{},
			greaterFn:  func(a, b int) bool { return a > b },
			wantValue:  0,
			wantIndex:  -1,
		},
		{
			name:       "single int element",
			collection: []int{42},
			greaterFn:  func(a, b int) bool { return a > b },
			wantValue:  42,
			wantIndex:  0,
		},
		{
			name:       "multiple ints",
			collection: []int{1, 5, 3, 2, 4},
			greaterFn:  func(a, b int) bool { return a > b },
			wantValue:  5,
			wantIndex:  1,
		},
		{
			name:       "multiple ints with duplicate max",
			collection: []int{1, 5, 3, 5, 4},
			greaterFn:  func(a, b int) bool { return a > b },
			wantValue:  5,
			wantIndex:  1, // Should return first occurrence
		},
		{
			name:       "negative ints",
			collection: []int{-5, -3, -1, -4, -2},
			greaterFn:  func(a, b int) bool { return a > b },
			wantValue:  -1,
			wantIndex:  2,
		},
		{
			name:       "float64 values",
			collection: []float64{1.1, 2.2, 3.3, 2.8, 1.5},
			greaterFn:  func(a, b float64) bool { return a > b },
			wantValue:  3.3,
			wantIndex:  2,
		},
		{
			name:       "strings",
			collection: []string{"apple", "zebra", "banana", "yak"},
			greaterFn:  func(a, b string) bool { return a > b },
			wantValue:  "zebra",
			wantIndex:  1,
		},
		{
			name:       "custom comparison for strings (length)",
			collection: []string{"apple", "zebra", "banana", "yak"},
			greaterFn:  func(a, b string) bool { return len(a) > len(b) },
			wantValue:  "banana",
			wantIndex:  2,
		},
		{
			name:       "custom comparison for ints (absolute value)",
			collection: []int{-5, 3, -1, 4, -2},
			greaterFn: func(a, b int) bool {
				abs := func(x int) int {
					if x < 0 {
						return -x
					}
					return x
				}
				return abs(a) > abs(b)
			},
			wantValue: -5,
			wantIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch collection := tt.collection.(type) {
			case []int:
				greaterFn := tt.greaterFn.(func(int, int) bool)
				gotValue, gotIndex := slices.MaxIndexBy(collection, greaterFn)
				if gotValue != tt.wantValue.(int) || gotIndex != tt.wantIndex {
					t.Errorf("slices.MaxIndexBy() = (%v, %v), want (%v, %v)",
						gotValue, gotIndex, tt.wantValue, tt.wantIndex)
				}
			case []float64:
				greaterFn := tt.greaterFn.(func(float64, float64) bool)
				gotValue, gotIndex := slices.MaxIndexBy(collection, greaterFn)
				if gotValue != tt.wantValue.(float64) || gotIndex != tt.wantIndex {
					t.Errorf("slices.MaxIndexBy() = (%v, %v), want (%v, %v)",
						gotValue, gotIndex, tt.wantValue, tt.wantIndex)
				}
			case []string:
				greaterFn := tt.greaterFn.(func(string, string) bool)
				gotValue, gotIndex := slices.MaxIndexBy(collection, greaterFn)
				if gotValue != tt.wantValue.(string) || gotIndex != tt.wantIndex {
					t.Errorf("slices.MaxIndexBy() = (%v, %v), want (%v, %v)",
						gotValue, gotIndex, tt.wantValue, tt.wantIndex)
				}
			}
		})
	}
}
