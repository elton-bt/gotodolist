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
	handler := newHandler(config.Config{CORSAllowOrigin: "*", AppName: "gotodolist-api", Version: "1.2.3", InstanceName: "api-01", InstanceIP: "10.20.30.41"}, service)

	createBody := bytes.NewBufferString(`{"title":"Preparar aula","description":"Subir ambiente com compose"}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/api/tasks", createBody)
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createResponse.Code)
	}
	if createResponse.Header().Get("X-GoToDoList-Instance") != "api-01" {
		t.Fatalf("expected instance header on create response")
	}
	if createResponse.Header().Get("X-GoToDoList-IP") != "10.20.30.41" {
		t.Fatalf("expected IP header on create response")
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
	if deleteResponse.Header().Get("Access-Control-Expose-Headers") != "X-GoToDoList-Instance, X-GoToDoList-IP" {
		t.Fatalf("expected exposed instance headers on delete response")
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

func TestRootEndpointReturnsServiceMetadata(t *testing.T) {
	store := todo.NewMemoryStore()
	service := todo.NewService(store)
	handler := newHandler(config.Config{CORSAllowOrigin: "*", AppName: "gotodolist-api", Version: "1.2.3", InstanceName: "api-01", InstanceIP: "10.20.30.41"}, service)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		fatalf := t.Fatalf
		fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload struct {
		Service  string `json:"service"`
		Status   string `json:"status"`
		Version  string `json:"versao"`
		Instance string `json:"instancia"`
		IP       string `json:"ip"`
		Health   string `json:"health"`
		Tasks    string `json:"tasks"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal root response: %v", err)
	}

	if payload.Service != "gotodolist-api" {
		t.Fatalf("expected service gotodolist-api, got %q", payload.Service)
	}
	if payload.Status != "ok" {
		t.Fatalf("expected status ok, got %q", payload.Status)
	}
	if payload.Version != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %q", payload.Version)
	}
	if payload.Instance != "api-01" {
		t.Fatalf("expected instance api-01, got %q", payload.Instance)
	}
	if payload.IP != "10.20.30.41" {
		t.Fatalf("expected IP 10.20.30.41, got %q", payload.IP)
	}
	if payload.Health != "/health" {
		t.Fatalf("expected health endpoint /health, got %q", payload.Health)
	}
	if payload.Tasks != "/api/tasks" {
		t.Fatalf("expected tasks endpoint /api/tasks, got %q", payload.Tasks)
	}
}

func TestHealthEndpointReturnsDegraded(t *testing.T) {
	store := todo.NewMemoryStore()
	store.SetHealth(todo.ErrUnavailable)
	service := todo.NewService(store)
	handler := newHandler(config.Config{CORSAllowOrigin: "*", AppName: "gotodolist-api", Version: "1.2.3", InstanceName: "api-01", InstanceIP: "10.20.30.41"}, service)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, response.Code)
	}
	if response.Header().Get("X-GoToDoList-Instance") != "api-01" {
		t.Fatalf("expected instance header on health response")
	}
}

func TestRootEndpointReturnsDegradedStatus(t *testing.T) {
	store := todo.NewMemoryStore()
	store.SetHealth(todo.ErrUnavailable)
	service := todo.NewService(store)
	handler := newHandler(config.Config{CORSAllowOrigin: "*", AppName: "gotodolist-api", Version: "1.2.3", InstanceName: "api-01", InstanceIP: "10.20.30.41"}, service)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, response.Code)
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal degraded root response: %v", err)
	}

	if payload.Status != "degraded" {
		t.Fatalf("expected status degraded, got %q", payload.Status)
	}
}
