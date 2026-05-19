package todo

import (
	"context"
	"testing"
	"time"
)

func TestServiceCreateTask(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store)
	service.now = func() time.Time { return time.Date(2026, 5, 19, 10, 0, 0, 0, time.UTC) }

	task, err := service.Create(context.Background(), CreateTaskInput{
		Title:       "  Estudar Docker  ",
		Description: "  Configurar compose local  ",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if task.ID != 1 {
		t.Fatalf("expected id 1, got %d", task.ID)
	}
	if task.Title != "Estudar Docker" {
		t.Fatalf("expected trimmed title, got %q", task.Title)
	}
	if task.Description != "Configurar compose local" {
		t.Fatalf("expected trimmed description, got %q", task.Description)
	}
	if !task.CreatedAt.Equal(task.UpdatedAt) {
		t.Fatalf("expected created and updated times to match")
	}
}

func TestServiceListTasks(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store)

	_, _ = service.Create(context.Background(), CreateTaskInput{Title: "Primeira"})
	_, _ = service.Create(context.Background(), CreateTaskInput{Title: "Segunda"})

	tasks, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "Primeira" || tasks[1].Title != "Segunda" {
		t.Fatalf("unexpected ordering: %+v", tasks)
	}
}

func TestServiceUpdateTask(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store)

	created, _ := service.Create(context.Background(), CreateTaskInput{Title: "Original"})
	updatedTitle := "Atualizada"
	updatedDescription := "Mais detalhes"
	completed := true

	task, err := service.Update(context.Background(), created.ID, UpdateTaskInput{
		Title:       &updatedTitle,
		Description: &updatedDescription,
		Completed:   &completed,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if task.Title != updatedTitle {
		t.Fatalf("expected updated title, got %q", task.Title)
	}
	if task.Description != updatedDescription {
		t.Fatalf("expected updated description, got %q", task.Description)
	}
	if !task.Completed {
		t.Fatalf("expected task to be completed")
	}
}

func TestServiceDeleteTask(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store)

	created, _ := service.Create(context.Background(), CreateTaskInput{Title: "Remover"})

	if err := service.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	tasks, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected no tasks after delete, got %d", len(tasks))
	}
}
