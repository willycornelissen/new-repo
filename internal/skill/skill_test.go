package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"new-repo/internal/skill"
)

func TestParseSkills(t *testing.T) {
	t.Run("parses valid table", func(t *testing.T) {
		md := `# Skills

| Skill | Desc |
|-------|------|
| **foo** | description |
| **bar-baz** | another |
`
		skills, err := skill.ParseSkills(md)
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"foo", "bar-baz"}
		if len(skills) != len(want) {
			t.Fatalf("got %v, want %v", skills, want)
		}
		for i, s := range skills {
			if s != want[i] {
				t.Fatalf("skill[%d] = %q, want %q", i, s, want[i])
			}
		}
	})

	t.Run("handles empty content", func(t *testing.T) {
		skills, err := skill.ParseSkills("")
		if err != nil {
			t.Fatal(err)
		}
		if len(skills) != 0 {
			t.Fatalf("expected 0 skills, got %d", len(skills))
		}
	})

	t.Run("ignores header row", func(t *testing.T) {
		md := `| Skill | Desc |
|-------|------|
| **test** | value |
`
		skills, err := skill.ParseSkills(md)
		if err != nil {
			t.Fatal(err)
		}
		if len(skills) != 1 || skills[0] != "test" {
			t.Fatalf("got %v, want [test]", skills)
		}
	})

	t.Run("handles no skills", func(t *testing.T) {
		md := `# Only a title
no table here
`
		skills, err := skill.ParseSkills(md)
		if err != nil {
			t.Fatal(err)
		}
		if len(skills) != 0 {
			t.Fatalf("expected 0 skills, got %d", len(skills))
		}
	})
}

func TestInstallSkills(t *testing.T) {
	dst := t.TempDir()
	src := t.TempDir()

	// Create source skill directory with a file
	skillDir := filepath.Join(src, "my-skill")
	if err := os.Mkdir(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# My Skill"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := skill.InstallSkills(dst, src, []string{"my-skill"}); err != nil {
		t.Fatal(err)
	}

	// Verify it was copied
	installedFile := filepath.Join(dst, "my-skill", "SKILL.md")
	data, err := os.ReadFile(installedFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# My Skill" {
		t.Fatalf("got %q, want %q", data, "# My Skill")
	}
}

func TestInstallSkills_MissingSkill(t *testing.T) {
	dst := t.TempDir()
	src := t.TempDir()

	// Should not error, just warn
	if err := skill.InstallSkills(dst, src, []string{"nonexistent"}); err != nil {
		t.Fatal(err)
	}
}

func TestEmbeddedSkills_Content(t *testing.T) {
	content, err := skill.ReadEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty embedded SKILLS.md")
	}

	skills, err := skill.EmbeddedSkills()
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) == 0 {
		t.Fatal("expected at least one skill in embedded SKILLS.md")
	}
}
