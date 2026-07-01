package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreatorCreatesStandardApplicationByDefault(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	creator := NewCreator("myapp")
	if err := creator.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	for _, path := range []string{
		filepath.Join("myapp", "go.mod"),
		filepath.Join("myapp", "cmd", "myapp", "main.go"),
		filepath.Join("myapp", "internal"),
		filepath.Join("myapp", "README.md"),
		filepath.Join("myapp", ".gitignore"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	for _, path := range []string{
		filepath.Join("myapp", "main.go"),
		filepath.Join("myapp", "pkg"),
		filepath.Join("myapp", "test"),
		filepath.Join("myapp", "bin"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s not to exist", path)
		}
	}
}
