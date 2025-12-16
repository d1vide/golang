package service

import (
	"example.com/notes-api/internal/core"
	"example.com/notes-api/internal/repo"
)

type NoteService struct {
	repo *repo.NoteRepoMem
}

func NewNoteService(r *repo.NoteRepoMem) *NoteService {
	return &NoteService{repo: r}
}

func (s *NoteService) Create(title, content string) (*core.Note, error) {
	n := core.Note{
		Title:   title,
		Content: content,
	}

	return s.repo.Create(n)
}

func (s *NoteService) GetAll() ([]core.Note, error) {
	return s.repo.GetAll()
}

func (s *NoteService) GetByID(id int64) (*core.Note, error) {
	return s.repo.GetByID(id)
}

func (s *NoteService) Update(id int64, title, content string) (*core.Note, error) {
	return s.repo.Update(id, title, content)
}

func (s *NoteService) Delete(id int64) error {
	return s.repo.Delete(id)
}
