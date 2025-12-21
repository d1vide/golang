package repo

import (
	"errors"
	"sync"
	"time"

	"example.com/notes-api/internal/core"
)

var ErrNotFound = errors.New("note not found")

type NoteRepoMem struct {
	mu    sync.Mutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) Create(n core.Note) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.next++
	now := time.Now().UTC()

	n.ID = r.next
	n.CreatedAt = now

	r.notes[n.ID] = &n
	return &n, nil
}

func (r *NoteRepoMem) GetAll() ([]core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]core.Note, 0, len(r.notes))
	for _, n := range r.notes {
		res = append(res, *n)
	}
	return res, nil
}

func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (r *NoteRepoMem) Update(id int64, title, content string) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	if !ok {
		return nil, ErrNotFound
	}

	n.Title = title
	n.Content = content
	now := time.Now().UTC()
	n.UpdatedAt = &now

	return n, nil
}

func (r *NoteRepoMem) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return ErrNotFound
	}
	delete(r.notes, id)
	return nil
}
