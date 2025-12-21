package dto

import "time"

type NoteResponse struct {
	ID        int64      `json:"id" example:"1"`
	Title     string     `json:"title" example:"TestTitle"`
	Content   string     `json:"content" example:"TestContent"`
	CreatedAt time.Time  `json:"createdAt" example:"2025-12-15T12:00:00Z"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" example:""`
}
