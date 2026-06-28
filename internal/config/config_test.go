package config_test

import (
	"path/filepath"
	"testing"

	"new-repo/internal/config"
)

func TestNew(t *testing.T) {
	cfg := config.New("my-project", false)

	if cfg.ProjectName != "my-project" {
		t.Fatalf("ProjectName = %q, want %q", cfg.ProjectName, "my-project")
	}
	if cfg.ProjectDir != "my-project" {
		t.Fatalf("ProjectDir = %q, want %q", cfg.ProjectDir, "my-project")
	}
	if cfg.Force {
		t.Fatal("expected Force = false")
	}
	if cfg.GitBin != "git" {
		t.Fatalf("GitBin = %q, want %q", cfg.GitBin, "git")
	}
}

func TestNew_WithForce(t *testing.T) {
	cfg := config.New("p", true)
	if !cfg.Force {
		t.Fatal("expected Force = true")
	}
}

func TestSkillsDir(t *testing.T) {
	cfg := config.New("my-app", false)

	want := filepath.Join("my-app", ".opencode", "skills")
	got := cfg.SkillsDir()

	if got != want {
		t.Fatalf("SkillsDir() = %q, want %q", got, want)
	}
}
