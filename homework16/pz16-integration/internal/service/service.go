package service

import (
	"context"
	"pz16/internal/models"
	"pz16/internal/repo"
)

type Service struct{ Notes repo.NoteRepo }

func (s Service) Create(ctx context.Context, n *models.Note) error {
	return s.Notes.Create(ctx, n)
}
func (s Service) Get(ctx context.Context, id int64) (models.Note, error) {
	return s.Notes.Get(ctx, id)
}
