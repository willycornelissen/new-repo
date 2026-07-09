package embed_test

import (
	"os"
	"path/filepath"
	"testing"

	"new-repo/internal/embed"
)

func TestListSkills(t *testing.T) {
	skills, err := embed.ListSkills()
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) == 0 {
		t.Fatal("expected at least one skill in embedded template")
	}
}

func TestSkillExists(t *testing.T) {
	if !embed.SkillExists("context7") {
		t.Fatal("expected context7 to exist")
	}
	if embed.SkillExists("nonexistent-skill") {
		t.Fatal("expected nonexistent-skill to not exist")
	}
}

func TestExtractSkills(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, ".opencode", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := embed.ExtractSkills(skillsDir, []string{"context7"}); err != nil {
		t.Fatal(err)
	}

	skillPath := filepath.Join(skillsDir, "context7", "SKILL.md")
	if _, err := os.Stat(skillPath); err != nil {
		t.Fatalf("expected context7/SKILL.md to exist: %v", err)
	}
}

func TestExtractSkills_Missing(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, ".opencode", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := embed.ExtractSkills(skillsDir, []string{"nonexistent"}); err != nil {
		t.Fatal(err)
	}
}

func TestExtractTemplate(t *testing.T) {
	dir := t.TempDir()

	if err := embed.ExtractTemplate(dir); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err != nil {
		t.Fatal("expected AGENTS.md to exist")
	}

	if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
		t.Fatal("expected README.md to be skipped")
	}

	skillsDir := filepath.Join(dir, ".opencode", "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one skill directory")
	}
}

func TestHasSkills(t *testing.T) {
	if !embed.HasSkills() {
		t.Fatal("expected HasSkills to be true")
	}
}

func TestExtractOpenCodeCommands(t *testing.T) {
	dir := t.TempDir()

	if err := embed.ExtractOpenCodeCommands(dir); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected command files to be extracted")
	}

	hasExplore := false
	for _, e := range entries {
		if e.Name() == "explore.md" {
			hasExplore = true
			break
		}
	}
	if !hasExplore {
		t.Fatal("expected explore.md in commands")
	}
}

func TestExtractOpenCodeDocs(t *testing.T) {
	dir := t.TempDir()

	if err := embed.ExtractOpenCodeDocs(dir); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected doc files to be extracted")
	}
}

func TestListAvailableSkillNames(t *testing.T) {
	available := embed.ListAvailableSkillNames([]string{
		"context7",
		"nonexistent",
		"mermaid-studio",
	})
	if len(available) != 2 {
		t.Fatalf("expected 2 available skills, got %d: %v", len(available), available)
	}
}

func TestReadFile(t *testing.T) {
	data, err := embed.ReadFile("SKILLS.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty SKILLS.md")
	}

	_, err = embed.ReadFile("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
