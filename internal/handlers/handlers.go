package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/davidkamere/htmx-go-scaffolding/internal/tasks"
	"github.com/gorilla/mux"
)

type Handler struct {
	templates *template.Template
	store     tasks.Repository
}

type IndexPageData struct {
	Tasks []tasks.Task
}

func New(tmpl *template.Template, store tasks.Repository) *Handler {
	return &Handler{templates: tmpl, store: store}
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/", h.home).Methods(http.MethodGet)
	r.HandleFunc("/tasks", h.listTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks", h.createTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id}", h.deleteTask).Methods(http.MethodDelete)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	allTasks, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	data := IndexPageData{Tasks: allTasks}
	h.renderTemplate(w, "index.gohtmx", data)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	allTasks, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "task_list.gohtmx", allTasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "task title is required", http.StatusUnprocessableEntity)
		return
	}

	if _, err := h.store.Create(r.Context(), title); err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	allTasks, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "task_list.gohtmx", allTasks)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	ok, err := h.store.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	allTasks, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "task_list.gohtmx", allTasks)
}

func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, fmt.Sprintf("template render error: %v", err), http.StatusInternalServerError)
	}
}
