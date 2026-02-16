package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func Parse() (*template.Template, error) {
	glob, err := templateGlob()
	if err != nil {
		return nil, err
	}
	return template.ParseGlob(glob)
}

func templateGlob() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("read cwd: %w", err)
	}

	dir := wd
	for {
		if fileExists(filepath.Join(dir, "go.mod")) {
			return filepath.Join(dir, "web", "templates", "*.gohtmx"), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root containing go.mod")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
