package service

import (
	"time"

	"example.com/notes-api/internal/core"
	"example.com/notes-api/internal/repo"
)

type NoteService struct {
	repo *repo.NoteDBRepo
}

func NewNoteService(r *repo.NoteDBRepo) *NoteService {
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

func (s *NoteService) GetNotesPaginated(limit int, offset int64, cursorID int64, cursorTime time.Time, mode string) ([]core.Note, error) {
	if mode == "keyset" {
		return s.repo.GetPaginated(limit, cursorID, cursorTime)
	} else {
		return s.repo.GetPaginatedOffset(limit, offset)
	}
}

func (s *NoteService) GetNotesBatch(ids []int64) (map[int64]*core.Note, error) {
	return s.repo.GetByIDs(ids)
}

func (s *NoteService) SearchNotes(query string, limit int) ([]core.Note, error) {
	return s.repo.SearchByTitle(query, limit)
}
