package todo

import (
	"context"
	"strings"
	"time"
)

type Store interface {
	Create(ctx context.Context, task Task) (Task, error)
	List(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, id int64, input UpdateTaskInput) (Task, error)
	Delete(ctx context.Context, id int64) error
	Health(ctx context.Context) error
}

type Service struct {
	store Store
	now   func() time.Time
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
		now:   time.Now,
	}
}

func (s *Service) Create(ctx context.Context, input CreateTaskInput) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Task{}, ErrInvalidTitle
	}

	now := s.now().UTC()
	return s.store.Create(ctx, Task{
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (s *Service) List(ctx context.Context) ([]Task, error) {
	return s.store.List(ctx)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateTaskInput) (Task, error) {
	if input.Title != nil {
		trimmed := strings.TrimSpace(*input.Title)
		if trimmed == "" {
			return Task{}, ErrInvalidTitle
		}
		input.Title = &trimmed
	}

	if input.Description != nil {
		trimmed := strings.TrimSpace(*input.Description)
		input.Description = &trimmed
	}

	return s.store.Update(ctx, id, input)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) Health(ctx context.Context) error {
	return s.store.Health(ctx)
}
