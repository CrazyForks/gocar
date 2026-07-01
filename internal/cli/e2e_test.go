package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestE2ENewBuildRunDoctorCheck(t *testing.T) {
	root := findRepoRoot(t)
	tmp := t.TempDir()
	bin := filepath.Join(tmp, "gocar")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	runCmd(t, root, "go", "build", "-o", bin, "./cmd/gocar")

	out, err := runCmdOutput(tmp, bin, "new", "myapp")
	if err != nil {
		t.Fatalf("gocar new failed: %v\n%s", err, out)
	}

	appRoot := filepath.Join(tmp, "myapp")
	for _, path := range []string{
		filepath.Join(appRoot, "go.mod"),
		filepath.Join(appRoot, "cmd", "myapp", "main.go"),
		filepath.Join(appRoot, "internal"),
		filepath.Join(appRoot, "README.md"),
		filepath.Join(appRoot, ".gitignore"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(appRoot, ".gocar.toml")); !os.IsNotExist(err) {
		t.Fatal("new projects should not include .gocar.toml by default")
	}

	for _, tc := range [][]string{
		{"doctor"},
		{"fmt"},
		{"vet"},
		{"check", "--no-test"},
		{"test"},
		{"build"},
		{"run"},
	} {
		out, err := runCmdOutput(appRoot, bin, tc...)
		if err != nil {
			t.Fatalf("gocar %s failed: %v\n%s", strings.Join(tc, " "), err, out)
		}
	}

	out, err = runCmdOutput(appRoot, bin, "init")
	if err != nil {
		t.Fatalf("gocar init failed: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(appRoot, ".gocar.toml")); err != nil {
		t.Fatalf("expected .gocar.toml after init: %v", err)
	}

	out, err = runCmdOutput(tmp, bin, "new", "oldapp", "--mode", "project")
	if err == nil {
		t.Fatalf("expected --mode to fail, output:\n%s", out)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func runCmd(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	out, err := runCmdOutput(dir, name, args...)
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, out)
	}
}

func runCmdOutput(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return string(out), err
}
