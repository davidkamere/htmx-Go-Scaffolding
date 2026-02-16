package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, title string) (Task, error)
	Delete(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context) ([]Task, error)
	Close() error
}

type FileStore struct {
	mu     sync.RWMutex
	path   string
	nextID int64
	tasks  map[int64]Task
}

type persistedState struct {
	NextID int64  `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}

func NewFileStore(path string) (*FileStore, error) {
	if err := ensureDir(path); err != nil {
		return nil, err
	}

	s := &FileStore{
		path:   path,
		nextID: 1,
		tasks:  make(map[int64]Task),
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *FileStore) Create(_ context.Context, title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("task title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:        s.nextID,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	s.tasks[task.ID] = task
	s.nextID++

	if err := s.persistLocked(); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *FileStore) Delete(_ context.Context, id int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return false, nil
	}
	delete(s.tasks, id)

	if err := s.persistLocked(); err != nil {
		return false, err
	}

	return true, nil
}

func (s *FileStore) List(_ context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].ID > result[j].ID
		}
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func (s *FileStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.persistLocked()
}

func (s *FileStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read store file: %w", err)
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("unmarshal store file: %w", err)
	}

	if state.NextID < 1 {
		state.NextID = 1
	}

	s.nextID = state.NextID
	for _, t := range state.Tasks {
		s.tasks[t.ID] = t
		if t.ID >= s.nextID {
			s.nextID = t.ID + 1
		}
	}

	return nil
}

func (s *FileStore) persistLocked() error {
	state := persistedState{
		NextID: s.nextID,
		Tasks:  make([]Task, 0, len(s.tasks)),
	}
	for _, t := range s.tasks {
		state.Tasks = append(state.Tasks, t)
	}

	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal store state: %w", err)
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return fmt.Errorf("write temp store file: %w", err)
	}
	if err := os.Rename(tmpPath, s.path); err != nil {
		return fmt.Errorf("replace store file: %w", err)
	}

	return nil
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create store directory: %w", err)
	}
	return nil
}
