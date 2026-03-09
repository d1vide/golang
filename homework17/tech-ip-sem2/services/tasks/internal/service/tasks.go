package service

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type TaskSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type CreateInput struct {
	Title       string
	Description string
	DueDate     string
}

type UpdateInput struct {
	Title       *string
	Description *string
	DueDate     *string
	Done        *bool
}

var ErrNotFound = errors.New("task not found")

type TasksService struct {
	mu      sync.RWMutex
	store   map[string]*Task
	counter int
}

func New() *TasksService {
	return &TasksService{store: make(map[string]*Task)}
}

func (s *TasksService) nextID() string {
	s.counter++
	return fmt.Sprintf("t_%03d", s.counter)
}

func (s *TasksService) Create(in CreateInput) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID()
	due := in.DueDate
	if due == "" {
		due = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}
	t := &Task{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		DueDate:     due,
		Done:        false,
	}
	s.store[id] = t
	return *t
}

func (s *TasksService) List() []TaskSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]TaskSummary, 0, len(s.store))
	for _, t := range s.store {
		out = append(out, TaskSummary{ID: t.ID, Title: t.Title, Done: t.Done})
	}
	return out
}

func (s *TasksService) Get(id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.store[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return *t, nil
}

func (s *TasksService) Update(id string, in UpdateInput) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.store[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	if in.Title != nil {
		t.Title = *in.Title
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if in.DueDate != nil {
		t.DueDate = *in.DueDate
	}
	if in.Done != nil {
		t.Done = *in.Done
	}
	return *t, nil
}

func (s *TasksService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[id]; !ok {
		return ErrNotFound
	}
	delete(s.store, id)
	return nil
}
