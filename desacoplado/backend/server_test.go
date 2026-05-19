package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elton-bt/gotodolist/internal/config"
	"github.com/elton-bt/gotodolist/internal/todo"
)

func TestCRUDEndpoints(t *testing.T) {
	store := todo.NewMemoryStore()
	service := todo.NewService(store)
	handler := newHandler(config.Config{CORSAllowOrigin: "*"}, service)

	createBody := bytes.NewBufferString(`{"title":"Preparar aula","description":"Subir ambiente com compose"}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/api/tasks", createBody)
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createResponse.Code)
	}

	var created todo.Task
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	listResponse := httptest.NewRecorder()
	handler.ServeHTTP(listResponse, listRequest)

	if listResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listResponse.Code)
	}

	var listed struct {
		Tasks []todo.Task `json:"tasks"`
	}
	if err := json.Unmarshal(listResponse.Body.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list response: %v", err)
	}
	if len(listed.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(listed.Tasks))
	}

	updateBody := bytes.NewBufferString(`{"title":"Preparar aula 2","description":"Atualizar roteiro","completed":true}`)
	updateRequest := httptest.NewRequest(http.MethodPut, "/api/tasks/1", updateBody)
	updateRequest.SetPathValue("id", "1")
	updateRequest.Header.Set("Content-Type", "application/json")
	updateResponse := httptest.NewRecorder()
	handler.ServeHTTP(updateResponse, updateRequest)

	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateResponse.Code)
	}

	var updated todo.Task
	if err := json.Unmarshal(updateResponse.Body.Bytes(), &updated); err != nil {
		t.Fatalf("unmarshal update response: %v", err)
	}
	if !updated.Completed {
		t.Fatalf("expected updated task to be completed")
	}
	if updated.Title != "Preparar aula 2" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/tasks/1", nil)
	deleteRequest.SetPathValue("id", "1")
	deleteResponse := httptest.NewRecorder()
	handler.ServeHTTP(deleteResponse, deleteRequest)

	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, deleteResponse.Code)
	}

	finalListRequest := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	finalListResponse := httptest.NewRecorder()
	handler.ServeHTTP(finalListResponse, finalListRequest)

	if err := json.Unmarshal(finalListResponse.Body.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal final list response: %v", err)
	}
	if len(listed.Tasks) != 0 {
		t.Fatalf("expected 0 tasks after delete, got %d", len(listed.Tasks))
	}

	if created.ID == 0 {
		t.Fatalf("expected created task id to be populated")
	}
}

func TestHealthEndpointReturnsDegraded(t *testing.T) {
	store := todo.NewMemoryStore()
	store.SetHealth(todo.ErrUnavailable)
	service := todo.NewService(store)
	handler := newHandler(config.Config{CORSAllowOrigin: "*"}, service)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, response.Code)
	}
}
