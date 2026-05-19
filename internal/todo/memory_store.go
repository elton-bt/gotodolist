package todo

import (
	"context"
	"sort"
	"sync"
)

type MemoryStore struct {
	mu     sync.RWMutex
	nextID int64
	tasks  map[int64]Task
	health error
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID: 1,
		tasks:  make(map[int64]Task),
	}
}

func (s *MemoryStore) Create(_ context.Context, task Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task

	return task, nil
}

func (s *MemoryStore) List(_ context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

func (s *MemoryStore) Update(_ context.Context, id int64, input UpdateTaskInput) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	task.UpdatedAt = task.UpdatedAt.Add(1)
	s.tasks[id] = task

	return task, nil
}

func (s *MemoryStore) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}

	delete(s.tasks, id)
	return nil
}

func (s *MemoryStore) Health(_ context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.health
}

func (s *MemoryStore) SetHealth(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.health = err
}
