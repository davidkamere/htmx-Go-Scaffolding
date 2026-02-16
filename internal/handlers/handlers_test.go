package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davidkamere/htmx-go-scaffolding/internal/server"
)

func TestHomeRouteReturnsHTML(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "HTMX Go Scaffolding") {
		t.Fatalf("expected page title content in response")
	}
}

func TestTaskLifecycle(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	form := url.Values{}
	form.Set("title", "Ship scaffold v2")

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(form.Encode()))
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusOK {
		body, _ := io.ReadAll(createRec.Body)
		t.Fatalf("expected create status 200, got %d: %s", createRec.Code, string(body))
	}

	if !strings.Contains(createRec.Body.String(), "Ship scaffold v2") {
		t.Fatalf("created task not present in create response")
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected delete status 200, got %d", deleteRec.Code)
	}

	if strings.Contains(deleteRec.Body.String(), "Ship scaffold v2") {
		t.Fatalf("deleted task still present in response")
	}
}

func TestTaskPersistsAcrossRouterRestart(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app.db")

	router1, cleanup1 := newRouterWithDB(t, dbPath)
	form := url.Values{}
	form.Set("title", "Persisted task")
	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(form.Encode()))
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createRec := httptest.NewRecorder()
	router1.ServeHTTP(createRec, createReq)
	cleanup1()

	router2, cleanup2 := newRouterWithDB(t, dbPath)
	defer cleanup2()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	router2.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Persisted task") {
		t.Fatalf("expected persisted task after restart")
	}
}

func newTestRouter(t *testing.T) (http.Handler, func()) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "app.db")
	return newRouterWithDB(t, dbPath)
}

func newRouterWithDB(t *testing.T, dbPath string) (http.Handler, func()) {
	t.Helper()
	router, closer, err := server.NewRouter(dbPath)
	if err != nil {
		t.Fatalf("build router: %v", err)
	}
	return router, func() {
		_ = closer.Close()
	}
}
