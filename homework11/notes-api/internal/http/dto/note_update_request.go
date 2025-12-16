package dto

type UpdateNoteRequest struct {
	Title   string `json:"title" example:"TestTitleUpdated"`
	Content string `json:"content" example:"TestContentUpdated"`
}
