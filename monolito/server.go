package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/elton-bt/gotodolist/internal/config"
	"github.com/elton-bt/gotodolist/internal/httpx"
	"github.com/elton-bt/gotodolist/internal/todo"
)

//go:embed templates/*.html static/*
var monolithAssets embed.FS

type pageData struct {
	Tasks            []todo.Task
	Message          string
	DraftTitle       string
	DraftDescription string
	Version          string
}

type server struct {
	templates *template.Template
	service   *todo.Service
	version   string
}

func newHTTPServer(cfg config.Config, service *todo.Service) (*http.Server, error) {
	handler, err := newHandler(cfg, service)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              cfg.Address(),
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}, nil
}

func newHandler(cfg config.Config, service *todo.Service) (http.Handler, error) {
	templates, err := template.New("index.html").Funcs(template.FuncMap{
		"formatTime": func(value time.Time) string {
			if value.IsZero() {
				return "agora"
			}

			return value.Local().Format("02/01/2006 15:04")
		},
	}).ParseFS(monolithAssets, "templates/*.html")
	if err != nil {
		return nil, err
	}

	staticFiles, err := fs.Sub(monolithAssets, "static")
	if err != nil {
		return nil, err
	}

	server := &server{
		templates: templates,
		service:   service,
		version:   cfg.Version,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(server.handleIndex))
	mux.Handle("POST /tasks", http.HandlerFunc(server.handleCreate))
	mux.Handle("POST /tasks/{id}/edit", http.HandlerFunc(server.handleUpdate))
	mux.Handle("POST /tasks/{id}/toggle", http.HandlerFunc(server.handleToggle))
	mux.Handle("POST /tasks/{id}/delete", http.HandlerFunc(server.handleDelete))
	mux.Handle("GET /health", http.HandlerFunc(server.handleHealth))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFiles)))

	return mux, nil
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.service.List(r.Context())
	if err != nil {
		status, message := httpx.ErrorStatusAndMessage(err)
		s.renderIndex(w, status, pageData{Message: message})
		return
	}

	s.renderIndex(w, http.StatusOK, pageData{Tasks: tasks})
}

func (s *server) handleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderPage(w, r, http.StatusBadRequest, "Nao foi possivel ler o formulario.", r.FormValue("title"), r.FormValue("description"))
		return
	}

	_, err := s.service.Create(r.Context(), todo.CreateTaskInput{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	})
	if err != nil {
		status, message := httpx.ErrorStatusAndMessage(err)
		s.renderPage(w, r, status, message, r.FormValue("title"), r.FormValue("description"))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		s.renderPage(w, r, http.StatusBadRequest, "Identificador de tarefa invalido.", "", "")
		return
	}

	if err := r.ParseForm(); err != nil {
		s.renderPage(w, r, http.StatusBadRequest, "Nao foi possivel ler o formulario.", "", "")
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	_, err := s.service.Update(r.Context(), id, todo.UpdateTaskInput{
		Title:       &title,
		Description: &description,
	})
	if err != nil {
		status, message := httpx.ErrorStatusAndMessage(err)
		s.renderPage(w, r, status, message, title, description)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleToggle(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		s.renderPage(w, r, http.StatusBadRequest, "Identificador de tarefa invalido.", "", "")
		return
	}

	if err := r.ParseForm(); err != nil {
		s.renderPage(w, r, http.StatusBadRequest, "Nao foi possivel ler o formulario.", "", "")
		return
	}

	completed, err := strconv.ParseBool(r.FormValue("completed"))
	if err != nil {
		s.renderPage(w, r, http.StatusBadRequest, "Estado de conclusao invalido.", "", "")
		return
	}

	_, err = s.service.Update(r.Context(), id, todo.UpdateTaskInput{Completed: &completed})
	if err != nil {
		status, message := httpx.ErrorStatusAndMessage(err)
		s.renderPage(w, r, status, message, "", "")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		s.renderPage(w, r, http.StatusBadRequest, "Identificador de tarefa invalido.", "", "")
		return
	}

	if err := s.service.Delete(r.Context(), id); err != nil {
		status, message := httpx.ErrorStatusAndMessage(err)
		s.renderPage(w, r, status, message, "", "")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
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

func (s *server) renderPage(w http.ResponseWriter, r *http.Request, status int, message, title, description string) {
	tasks, err := s.service.List(r.Context())
	if err != nil && message == "" {
		status, message = httpx.ErrorStatusAndMessage(err)
	}

	s.renderIndex(w, status, pageData{
		Tasks:            tasks,
		Message:          message,
		DraftTitle:       title,
		DraftDescription: description,
	})
}

func (s *server) renderIndex(w http.ResponseWriter, status int, data pageData) {
	data.Version = s.version
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func parseID(raw string) (int64, bool) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
