package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/elton-bt/gotodolist/internal/todo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	config.MaxConns = 8
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) EnsureSchema(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tasks (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure postgres schema: %w", err)
	}

	return nil
}

func (s *Store) Create(ctx context.Context, task todo.Task) (todo.Task, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO tasks (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, completed, created_at, updated_at
	`, task.Title, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return todo.Task{}, mapStoreError(err)
	}

	return task, nil
}

func (s *Store) List(ctx context.Context) ([]todo.Task, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, mapStoreError(err)
	}
	defer rows.Close()

	tasks := make([]todo.Task, 0)
	for rows.Next() {
		var task todo.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, mapStoreError(err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, mapStoreError(err)
	}

	return tasks, nil
}

func (s *Store) Update(ctx context.Context, id int64, input todo.UpdateTaskInput) (todo.Task, error) {
	var task todo.Task
	err := s.pool.QueryRow(ctx, `
		UPDATE tasks
		SET
			title = COALESCE($2::text, title),
			description = COALESCE($3::text, description),
			completed = COALESCE($4::boolean, completed),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, completed, created_at, updated_at
	`, id, input.Title, input.Description, input.Completed).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return todo.Task{}, mapStoreError(err)
	}

	return task, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	commandTag, err := s.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return mapStoreError(err)
	}

	if commandTag.RowsAffected() == 0 {
		return todo.ErrNotFound
	}

	return nil
}

func (s *Store) Health(ctx context.Context) error {
	return mapStoreError(s.pool.Ping(ctx))
}

func mapStoreError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return todo.ErrNotFound
	}

	return todo.ErrUnavailable
}
