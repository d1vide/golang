package stringsx

// Clip возвращает s, обрезанную до max (UTF-8 безопасность можно опустить для ПЗ).
func Clip(s string, max int) string {
	if max < 0 {
		max = 0
	}
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// go test ./...
// go test -cover ./...
// go test -coverprofile=coverage.out ./...
// go tool cover -html=coverage.out ./...
