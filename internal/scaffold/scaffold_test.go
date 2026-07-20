package scaffold_test

import (
	"os"
	"path/filepath"
	"testing"

	"new-repo/internal/scaffold"
)

func TestCreateProjectDir(t *testing.T) {
	t.Run("creates directory", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "test-project")

		abs, err := scaffold.CreateProjectDir(name, false)
		if err != nil {
			t.Fatal(err)
		}

		if abs == "" {
			t.Fatal("expected non-empty path")
		}

		info, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			t.Fatal("expected directory")
		}
	})

	t.Run("errors on existing dir without force", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "existing")

		if err := os.Mkdir(name, 0755); err != nil {
			t.Fatal(err)
		}

		_, err := scaffold.CreateProjectDir(name, false)
		if err == nil {
			t.Fatal("expected error for existing directory")
		}
	})

	t.Run("succeeds on existing dir with force", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "existing")

		if err := os.Mkdir(name, 0755); err != nil {
			t.Fatal(err)
		}

		_, err := scaffold.CreateProjectDir(name, true)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestWriteGitignore(t *testing.T) {
	dir := t.TempDir()

	if err := scaffold.WriteGitignore(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty .gitignore")
	}
}

func TestWriteSkillsMD(t *testing.T) {
	dir := t.TempDir()
	content := "# Skills\n\n| Skill | Desc |\n|-------|------|\n| **foo** | bar |"

	if err := scaffold.WriteSkillsMD(dir, content); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "SKILLS.md"))
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Fatalf("got %q, want %q", data, content)
	}
}

func TestCreateOpenCodeDirs(t *testing.T) {
	dir := t.TempDir()

	if err := scaffold.CreateOpenCodeDirs(dir); err != nil {
		t.Fatal(err)
	}

	expected := []string{
		filepath.Join(dir, ".opencode", "skills"),
		filepath.Join(dir, ".opencode", "commands"),
		filepath.Join(dir, ".opencode", "docs"),
	}

	for _, path := range expected {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %s to be a directory", path)
		}
	}
}

func TestWriteAgentsMD(t *testing.T) {
	dir := t.TempDir()

	if err := scaffold.WriteAgentsMD(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty AGENTS.md")
	}
}

func TestWriteReadmeMD(t *testing.T) {
	dir := t.TempDir()

	if err := scaffold.WriteReadmeMD(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty README.md")
	}
}
