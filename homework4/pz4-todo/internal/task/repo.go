package task

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]*Task
	file  string
}

func NewRepo(filename string) *Repo {
	repo := &Repo{
		items: make(map[int64]*Task),
		file:  filename,
	}

	if data, err := os.ReadFile(filename); err == nil {
		var tasks []*Task
		if json.Unmarshal(data, &tasks) == nil {
			for _, task := range tasks {
				repo.items[task.ID] = task
				if task.ID > repo.seq {
					repo.seq = task.ID
				}
			}
		}
	}

	return repo
}

func (r *Repo) save() {
	tasks := make([]*Task, 0, len(r.items))
	for _, task := range r.items {
		tasks = append(tasks, task)
	}

	data, _ := json.MarshalIndent(tasks, "", "  ")
	os.WriteFile(r.file, data, 0644)
}

func (r *Repo) List() []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		out = append(out, t)
	}
	return out
}

func (r *Repo) ListWithOptions(options ListOptions) ([]*Task, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filteredTasks := make([]*Task, 0, len(r.items))
	for _, task := range r.items {
		if options.Done != nil && task.Done != *options.Done {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}

	total := len(filteredTasks)

	sort.Slice(filteredTasks, func(i, j int) bool {
		return filteredTasks[i].ID < filteredTasks[j].ID
	})

	offset := (options.Page - 1) * options.Limit
	if offset >= total {
		return []*Task{}, total
	}

	end := offset + options.Limit
	if end > total {
		end = total
	}

	result := make([]*Task, end-offset)
	copy(result, filteredTasks[offset:end])

	return result, total
}

func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	now := time.Now()
	t := &Task{ID: r.seq, Title: title, CreatedAt: now, UpdatedAt: now, Done: false}
	r.items[t.ID] = t
	r.save()
	return t
}

func (r *Repo) Update(id int64, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	r.save()
	return t, nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	r.save()
	return nil
}
