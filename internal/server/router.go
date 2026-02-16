package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/davidkamere/htmx-go-scaffolding/internal/handlers"
	"github.com/davidkamere/htmx-go-scaffolding/internal/tasks"
	"github.com/davidkamere/htmx-go-scaffolding/internal/templates"
	"github.com/gorilla/mux"
)

func NewRouter(dbPath string) (http.Handler, io.Closer, error) {
	tmpl, err := templates.Parse()
	if err != nil {
		return nil, nil, err
	}

	store, err := tasks.NewFileStore(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("init file store: %w", err)
	}

	h := handlers.New(tmpl, store)

	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	h.Register(r)

	return r, store, nil
}
