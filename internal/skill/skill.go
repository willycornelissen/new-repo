package skill

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed SKILLS.md
var embeddedSkillsMD string

func ParseSkills(data string) ([]string, error) {
	var skills []string
	lines := strings.Split(data, "\n")
	inTable := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "|-") {
			inTable = true
			continue
		}

		if !inTable || !strings.HasPrefix(trimmed, "|") {
			continue
		}

		parts := strings.Split(trimmed, "|")
		if len(parts) < 3 {
			continue
		}

		name := strings.TrimSpace(parts[1])
		name = strings.Trim(name, "*")

		if name == "" || name == "Skill" {
			continue
		}

		skills = append(skills, name)
	}

	return skills, nil
}

func InstallSkills(dstDir string, srcDirs []string, names []string) error {
	for _, name := range names {
		found := false
		for _, srcDir := range srcDirs {
			src := filepath.Join(srcDir, name)

			info, err := os.Stat(src)
			if err != nil {
				continue
			}

			if !info.IsDir() {
				continue
			}

			dst := filepath.Join(dstDir, name)
			if err := copyDir(src, dst); err != nil {
				return fmt.Errorf("copy skill %q: %w", name, err)
			}
			found = true
			break
		}
		if !found {
			fmt.Fprintf(os.Stderr, "warning: skill %q not found in any source directory\n", name)
		}
	}
	return nil
}

func ReadEmbedded() (string, error) {
	return embeddedSkillsMD, nil
}

func EmbeddedSkills() ([]string, error) {
	return ParseSkills(embeddedSkillsMD)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
