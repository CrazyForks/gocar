package build

import (
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	gocarconfig "gocar/internal/config"
)

func TestBuilderCommandUsesConfig(t *testing.T) {
	cfg := NewConfig()
	cfg.Release = true
	cfg.Profile = "release"
	cfg.SetTarget("linux", "amd64")

	gcfg := gocarconfig.DefaultConfig()
	gcfg.Project.Version = "1.2.3"
	gcfg.Build.Entry = "cmd/api"
	gcfg.Build.Output = "dist"
	gcfg.Build.Tags = []string{"jsoniter", "sonic"}
	gcfg.Build.Ldflags = "-X main.commit=abc"

	builder := NewBuilder("/repo", "api", "standard", cfg, gcfg)

	if got := builder.GetRelativeOutputPath(); got != filepath.Join("dist", "release", "linux-amd64", "api") {
		t.Fatalf("GetRelativeOutputPath() = %q", got)
	}

	cmd := builder.buildCommand("/repo/dist/release/linux-amd64/api")
	args := cmd.Args

	if !slices.Contains(args, "-trimpath") {
		t.Fatalf("expected -trimpath in args: %#v", args)
	}
	if !slices.Contains(args, "-tags=jsoniter,sonic") {
		t.Fatalf("expected build tags in args: %#v", args)
	}
	if !slices.Contains(args, "./cmd/api") {
		t.Fatalf("expected normalized build entry in args: %#v", args)
	}

	var ldflags string
	for _, arg := range args {
		if len(arg) >= len("-ldflags=") && arg[:len("-ldflags=")] == "-ldflags=" {
			ldflags = arg
		}
	}
	for _, want := range []string{"-s -w", "-X main.version=1.2.3", "-X main.commit=abc"} {
		if !strings.Contains(ldflags, want) {
			t.Fatalf("ldflags %q missing %q", ldflags, want)
		}
	}
}

func TestBuilderWindowsOutputExtension(t *testing.T) {
	cfg := NewConfig()
	cfg.Profile = "ci"
	cfg.SetTarget("windows", "amd64")

	builder := NewBuilder("/repo", "api", "standard", cfg, nil)
	if got := builder.GetRelativeOutputPath(); got != filepath.Join("bin", "ci", "windows-amd64", "api.exe") {
		t.Fatalf("GetRelativeOutputPath() = %q", got)
	}
}

func TestBuildConfigCurrentPlatform(t *testing.T) {
	cfg := NewConfig()
	if cfg.TargetOS != runtime.GOOS || cfg.TargetArch != runtime.GOARCH {
		t.Fatalf("NewConfig target = %s/%s", cfg.TargetOS, cfg.TargetArch)
	}
	if !cfg.IsCurrentPlatform() {
		t.Fatal("new config should target current platform")
	}
}
