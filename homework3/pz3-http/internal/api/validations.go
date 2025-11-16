package api

const (
	MAX_TITLE_LEN = 140
	MIN_TITLE_LEN = 1
)

func validateTitleLen(title string) bool {
	if len(title) < MIN_TITLE_LEN || len(title) > MAX_TITLE_LEN {
		return false
	}
	return true
}
