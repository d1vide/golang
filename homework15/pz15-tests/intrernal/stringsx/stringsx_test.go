package stringsx

import (
	"testing"
)

func TestClip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{"пустая строка, max 0", "", 0, ""},
		{"пустая строка, max 5", "", 5, ""},
		{"пустая строка, max -1", "", -1, ""},
		{"пустая строка, max -10", "", -10, ""},

		{"строка короче лимита", "hello", 10, "hello"},
		{"строка равна лимиту", "hello", 5, "hello"},
		{"один символ, лимит больше", "a", 5, "a"},

		{"обрезаем до 3 символов", "hello", 3, "hel"},
		{"обрезаем до 1 символа", "hello", 1, "h"},
		{"обрезаем до 0", "hello", 0, ""},

		{"max = -1", "hello", -1, ""},
		{"max = -5", "hello", -5, ""},
		{"max = очень большое число", "hello", 1000, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clip(tt.input, tt.max)
			if result != tt.expected {
				t.Errorf("Clip(%q, %d) = %q, ожидалось %q",
					tt.input, tt.max, result, tt.expected)
			}
		})
	}
}

func BenchmarkClip(b *testing.B) {
	testString := "This is a relatively long string for benchmarking purposes"
	maxLengths := []int{5, 10, 20, 30, 50}

	for _, max := range maxLengths {
		b.Run("max="+string(rune(max)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Clip(testString, max)
			}
		})
	}
}
