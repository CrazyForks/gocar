package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigDoesNotOverrideBuiltIns(t *testing.T) {
	cfg := DefaultConfig()

	for _, name := range []string{"fmt", "vet", "test", "check"} {
		if _, ok := cfg.Commands[name]; ok {
			t.Fatalf("default custom commands should not override built-in %s", name)
		}
	}
}

func TestLoadMergesProjectConfig(t *testing.T) {
	root := t.TempDir()
	content := `
[project]
name = "api"
version = "1.2.3"

[build]
entry = "cmd/api"
output = "dist"
tags = ["jsoniter"]

[profile.release]
trimpath = false
cgo_enabled = true

[profile.ci]
race = true

[commands]
lint = "golangci-lint run"
`
	if err := os.WriteFile(filepath.Join(root, ConfigFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(root)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Project.Name != "api" || cfg.Project.Version != "1.2.3" {
		t.Fatalf("project config not merged: %+v", cfg.Project)
	}
	if got := cfg.GetBuildEntry(); got != "cmd/api" {
		t.Fatalf("GetBuildEntry() = %q", got)
	}
	if got := cfg.GetBuildOutputRoot(); got != "dist" {
		t.Fatalf("GetBuildOutputRoot() = %q", got)
	}
	if len(cfg.Build.Tags) != 1 || cfg.Build.Tags[0] != "jsoniter" {
		t.Fatalf("tags not merged: %#v", cfg.Build.Tags)
	}
	release, ok := cfg.GetProfile("release")
	if !ok {
		t.Fatal("release profile not found")
	}
	if release.Trimpath == nil || *release.Trimpath {
		t.Fatal("release trimpath should be overridden to false")
	}
	if release.CgoEnabled == nil || !*release.CgoEnabled {
		t.Fatal("release cgo_enabled should be overridden to true")
	}
	ci, ok := cfg.GetProfile("ci")
	if !ok || !ci.Race {
		t.Fatalf("custom ci profile not merged: %#v", cfg.Profile.Profiles)
	}
	if cfg.Commands["lint"] == "" {
		t.Fatalf("commands not merged: %#v", cfg.Commands)
	}
}

func TestStandardBuildEntryUsesDefaultAppName(t *testing.T) {
	cfg := DefaultConfig()

	if got := cfg.GetBuildEntryForApp("myapp"); got != "cmd/myapp" {
		t.Fatalf("GetBuildEntryForApp() = %q", got)
	}
	if got := cfg.GetRunEntryForApp("myapp"); got != "cmd/myapp" {
		t.Fatalf("GetRunEntryForApp() = %q", got)
	}
}

func TestResolveBuildOutputDirRejectsDangerousPaths(t *testing.T) {
	root := t.TempDir()

	cfg := DefaultConfig()
	cfg.Build.Output = "."
	if _, err := cfg.ResolveBuildOutputDir(root); err == nil {
		t.Fatal("expected error for project root output")
	}

	cfg.Build.Output = root
	if _, err := cfg.ResolveBuildOutputDir(root); err == nil {
		t.Fatal("expected error for absolute project root output")
	}
}
