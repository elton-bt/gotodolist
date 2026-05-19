package main

import (
	"net/http"
	"strconv"

	"github.com/elton-bt/gotodolist/internal/config"
	"github.com/elton-bt/gotodolist/internal/httpx"
	"github.com/elton-bt/gotodolist/internal/todo"
)

type apiServer struct {
	service         *todo.Service
	allowOriginHost string
}

func newHTTPServer(cfg config.Config, service *todo.Service) *http.Server {
	return &http.Server{
		Addr:              cfg.Address(),
		Handler:           newHandler(cfg, service),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}
}

func newHandler(cfg config.Config, service *todo.Service) http.Handler {
	server := &apiServer{
		service:         service,
		allowOriginHost: cfg.CORSAllowOrigin,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /api/tasks", http.HandlerFunc(server.handleListTasks))
	mux.Handle("POST /api/tasks", http.HandlerFunc(server.handleCreateTask))
	mux.Handle("PUT /api/tasks/{id}", http.HandlerFunc(server.handleUpdateTask))
	mux.Handle("DELETE /api/tasks/{id}", http.HandlerFunc(server.handleDeleteTask))
	mux.Handle("GET /health", http.HandlerFunc(server.handleHealth))

	return server.withCORS(mux)
}

func (s *apiServer) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.service.List(r.Context())
	if err != nil {
		s.writeError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string][]todo.Task{
		"tasks": tasks,
	})
}

func (s *apiServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request todo.CreateTaskInput
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido."})
		return
	}

	task, err := s.service.Create(r.Context(), request)
	if err != nil {
		s.writeError(w, err)
		return
	}

	w.Header().Set("Location", "/api/tasks/"+strconv.FormatInt(task.ID, 10))
	httpx.WriteJSON(w, http.StatusCreated, task)
}

func (s *apiServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Identificador invalido."})
		return
	}

	defer r.Body.Close()

	var request todo.UpdateTaskInput
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido."})
		return
	}

	if request.Title == nil && request.Description == nil && request.Completed == nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Informe ao menos um campo para atualizar."})
		return
	}

	task, err := s.service.Update(r.Context(), id, request)
	if err != nil {
		s.writeError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, task)
}

func (s *apiServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Identificador invalido."})
		return
	}

	if err := s.service.Delete(r.Context(), id); err != nil {
		s.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *apiServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.service.Health(r.Context()); err != nil {
		httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":  "degraded",
			"message": "Servico temporariamente indisponivel. Verifique a conexao com o banco.",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *apiServer) writeError(w http.ResponseWriter, err error) {
	status, message := httpx.ErrorStatusAndMessage(err)
	httpx.WriteJSON(w, status, map[string]string{"error": message})
}

func (s *apiServer) withCORS(next http.Handler) http.Handler {
	allowOrigin := s.allowOriginHost
	if allowOrigin == "" {
		allowOrigin = "*"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseID(raw string) (int64, bool) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
