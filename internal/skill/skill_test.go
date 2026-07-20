package skill_test

import (
	"os"
	"path/filepath"
	"strings"
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

	if err := skill.InstallSkills(dst, []string{src}, []string{"my-skill"}); err != nil {
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

	// Should not error, warn via stderr
	if err := skill.InstallSkills(dst, []string{src}, []string{"nonexistent"}); err != nil {
		t.Fatal(err)
	}
}

func TestInstallSkills_FallbackSrc(t *testing.T) {
	dst := t.TempDir()
	src1 := t.TempDir()
	src2 := t.TempDir()

	skillDir := filepath.Join(src2, "fallback-skill")
	if err := os.Mkdir(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Fallback"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := skill.InstallSkills(dst, []string{src1, src2}, []string{"fallback-skill"}); err != nil {
		t.Fatal(err)
	}

	installedFile := filepath.Join(dst, "fallback-skill", "SKILL.md")
	data, err := os.ReadFile(installedFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# Fallback" {
		t.Fatalf("got %q, want %q", data, "# Fallback")
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

func TestFindGStackSkills_And_GenerateSkillsMD(t *testing.T) {
	tempDir := t.TempDir()

	// Write root SKILL.md (should be skipped)
	if err := os.WriteFile(filepath.Join(tempDir, "SKILL.md"), []byte("---\nname: gstack\n---"), 0644); err != nil {
		t.Fatal(err)
	}

	// Write skill 1
	skill1Dir := filepath.Join(tempDir, "skill1")
	if err := os.Mkdir(skill1Dir, 0755); err != nil {
		t.Fatal(err)
	}
	skill1Content := `---
name: skill1
description: Description 1
---
Some body`
	if err := os.WriteFile(filepath.Join(skill1Dir, "SKILL.md"), []byte(skill1Content), 0644); err != nil {
		t.Fatal(err)
	}

	// Write skill 2
	skill2Dir := filepath.Join(tempDir, "skill2")
	if err := os.Mkdir(skill2Dir, 0755); err != nil {
		t.Fatal(err)
	}
	skill2Content := `---
name: "skill2"
description: "Description 2"
---`
	if err := os.WriteFile(filepath.Join(skill2Dir, "SKILL.md"), []byte(skill2Content), 0644); err != nil {
		t.Fatal(err)
	}

	skills, err := skill.FindGStackSkills(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	var s1, s2 skill.GStackSkill
	for _, s := range skills {
		if s.Name == "skill1" {
			s1 = s
		} else if s.Name == "skill2" {
			s2 = s
		}
	}

	if s1.Name != "skill1" || s1.Description != "Description 1" {
		t.Errorf("skill1 mismatch: %+v", s1)
	}
	if s2.Name != "skill2" || s2.Description != "Description 2" {
		t.Errorf("skill2 mismatch: %+v", s2)
	}

	md := skill.GenerateSkillsMD(skills)
	if !strings.Contains(md, "| **skill1** | Description 1 |") {
		t.Errorf("generated markdown missing skill1: %s", md)
	}
	if !strings.Contains(md, "| **skill2** | Description 2 |") {
		t.Errorf("generated markdown missing skill2: %s", md)
	}
}
