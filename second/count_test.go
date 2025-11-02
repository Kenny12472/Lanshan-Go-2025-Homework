package main

import (
	"reflect"
	"testing"
)

func TestCountArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected map[int]int
	}{
		{
			name:     "基本测试",
			input:    []int{1, 2, 2, 3, 3, 3},
			expected: map[int]int{1: 1, 2: 2, 3: 3},
		},
		{
			name:     "重复元素",
			input:    []int{5, 5, 5, 5},
			expected: map[int]int{5: 4},
		},
		{
			name:     "空数组",
			input:    []int{},
			expected: map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countarray(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("countarray(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
