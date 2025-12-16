package dto

type CreateNoteRequest struct {
	Title   string `json:"title" example:"TestTitle"`
	Content string `json:"content" example:"TestContent"`
}
