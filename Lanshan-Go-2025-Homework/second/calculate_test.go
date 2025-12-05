package main

import (
	"testing"
)

func TestCalculationFactory(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		a, b     int
		expected int
		nilFunc  bool
	}{
		{"加法", "add", 3, 5, 8, false},
		{"减法", "subtract", 10, 4, 6, false},
		{"乘法", "multiply", 2, 7, 14, false},
		{"除法", "divide", 20, 5, 4, false},
		{"除以0", "divide", 5, 0, 0, false},
		{"未知操作", "unknown", 1, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := CalculationFactory(tt.op)
			if fn == nil {
				if !tt.nilFunc {
					t.Errorf("CalculationFactory(%q) 返回 nil，但不应该", tt.op)
				}
				return
			}
			if tt.nilFunc {
				t.Errorf("CalculationFactory(%q) 应该返回 nil，但返回了函数", tt.op)
				return
			}
			result := fn(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("%s(%d, %d) = %d; expected %d", tt.op, tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
