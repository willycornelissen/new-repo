package git_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"new-repo/internal/git"
)

func TestIsAvailable(t *testing.T) {
	available := git.IsAvailable()

	_, err := exec.LookPath("git")
	expected := err == nil

	if available != expected {
		t.Fatalf("IsAvailable() = %v, want %v", available, expected)
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()

	if err := git.Init(dir); err != nil {
		t.Fatal(err)
	}

	gitDir := filepath.Join(dir, ".git")
	if !dirExists(gitDir) {
		t.Fatal("expected .git directory after init")
	}
}

func dirExists(path string) bool {
	_, err := exec.Command("test", "-d", path).Output()
	return err == nil
}
