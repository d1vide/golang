package mathx

import (
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -5, 10, 5},
		{"zero", 0, 0, 0},
		{"positive and zero", 7, 0, 7},
		{"negative and zero", -3, 0, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Sum(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	successTests := []struct {
		name           string
		a, b           int
		expectedResult int
		expectError    bool
	}{
		{"positive division", 10, 2, 5, false},
		{"negative division", -10, 2, -5, false},
		{"division with remainder", 7, 3, 2, false},
		{"division by negative", 10, -2, -5, false},
		{"both negative", -10, -2, 5, false},
		{"zero divided", 0, 5, 0, false},
	}

	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if tt.expectError && err == nil {
				t.Errorf("Divide(%d, %d) expected error, got nil", tt.a, tt.b)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
			}

			if result != tt.expectedResult {
				t.Errorf("Divide(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expectedResult)
			}
		})
	}

	t.Run("divide by zero", func(t *testing.T) {
		result, err := Divide(10, 0)

		if err == nil {
			t.Errorf("Divide(10, 0) expected error, got nil")
		}

		if result != 0 {
			t.Errorf("Divide(10, 0) should return 0 on error, got %d", result)
		}

		expectedError := "divide by zero"
		if err.Error() != expectedError {
			t.Errorf("Divide(10, 0) error = %v, expected %v", err.Error(), expectedError)
		}
	})

	t.Run("zero divided by zero", func(t *testing.T) {
		result, err := Divide(0, 0)

		if err == nil {
			t.Errorf("Divide(0, 0) expected error, got nil")
		}

		if result != 0 {
			t.Errorf("Divide(0, 0) should return 0 on error, got %d", result)
		}
	})
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum(123, 456)
	}
}

func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Divide(100, 25)
	}
}
