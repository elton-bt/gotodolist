package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/elton-bt/gotodolist/internal/todo"
)

func TestMonolithCRUDFlow(t *testing.T) {
	service := todo.NewService(todo.NewMemoryStore())
	handler, err := newHandler(service)
	if err != nil {
		t.Fatalf("newHandler returned error: %v", err)
	}

	createRequest := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(url.Values{
		"title":       []string{"Preparar demo"},
		"description": []string{"Explicar compose e health check"},
	}.Encode()))
	createRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, createResponse.Code)
	}

	listResponse := httptest.NewRecorder()
	handler.ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/", nil))

	if listResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listResponse.Code)
	}
	if !strings.Contains(listResponse.Body.String(), "Preparar demo") {
		t.Fatalf("expected rendered page to contain created task")
	}

	updateRequest := httptest.NewRequest(http.MethodPost, "/tasks/1/edit", strings.NewReader(url.Values{
		"title":       []string{"Preparar demo final"},
		"description": []string{"Atualizar roteiro"},
	}.Encode()))
	updateRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	updateRequest.SetPathValue("id", "1")
	updateResponse := httptest.NewRecorder()
	handler.ServeHTTP(updateResponse, updateRequest)

	if updateResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, updateResponse.Code)
	}

	toggleRequest := httptest.NewRequest(http.MethodPost, "/tasks/1/toggle", strings.NewReader(url.Values{
		"completed": []string{"true"},
	}.Encode()))
	toggleRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	toggleRequest.SetPathValue("id", "1")
	toggleResponse := httptest.NewRecorder()
	handler.ServeHTTP(toggleResponse, toggleRequest)

	if toggleResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, toggleResponse.Code)
	}

	listAfterToggle := httptest.NewRecorder()
	handler.ServeHTTP(listAfterToggle, httptest.NewRequest(http.MethodGet, "/", nil))

	body := listAfterToggle.Body.String()
	if !strings.Contains(body, "Preparar demo final") {
		t.Fatalf("expected updated title to be rendered")
	}
	if !strings.Contains(body, "Concluida") {
		t.Fatalf("expected toggled state to be rendered")
	}

	deleteRequest := httptest.NewRequest(http.MethodPost, "/tasks/1/delete", nil)
	deleteRequest.SetPathValue("id", "1")
	deleteResponse := httptest.NewRecorder()
	handler.ServeHTTP(deleteResponse, deleteRequest)

	if deleteResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, deleteResponse.Code)
	}

	listAfterDelete := httptest.NewRecorder()
	handler.ServeHTTP(listAfterDelete, httptest.NewRequest(http.MethodGet, "/", nil))

	if !strings.Contains(listAfterDelete.Body.String(), "Nenhuma tarefa cadastrada") {
		t.Fatalf("expected empty state after delete")
	}
}

func TestMonolithHealthEndpointReturnsDegraded(t *testing.T) {
	store := todo.NewMemoryStore()
	store.SetHealth(todo.ErrUnavailable)

	service := todo.NewService(store)
	handler, err := newHandler(service)
	if err != nil {
		t.Fatalf("newHandler returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, response.Code)
	}
	if !strings.Contains(response.Body.String(), "degraded") {
		t.Fatalf("expected degraded payload in health response")
	}
}
