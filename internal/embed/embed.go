package embed

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed template
var templateFS embed.FS

const templateRoot = "template"

// dirRename maps embedded directory names to their on-disk names.
// This works around go:embed not matching directories starting with '.'.
var dirRename = map[string]string{
	"opencode": ".opencode",
}

func ExtractTemplate(dst string) error {
	return extractDir(templateFS, templateRoot, dst, []string{"README.md"})
}

func ListSkills() ([]string, error) {
	skillsDir := filepath.Join(templateRoot, "opencode", "skills")
	entries, err := fs.ReadDir(templateFS, skillsDir)
	if err != nil {
		return nil, fmt.Errorf("reading embedded skills: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func ExtractSkills(dstSkillsDir string, names []string) error {
	skillsRoot := filepath.Join(templateRoot, "opencode", "skills")
	for _, name := range names {
		src := filepath.Join(skillsRoot, name)
		info, err := fs.Stat(templateFS, src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skill %q not found in embedded template\n", name)
			continue
		}
		if !info.IsDir() {
			fmt.Fprintf(os.Stderr, "warning: skill %q is not a directory in embedded template\n", name)
			continue
		}
		dst := filepath.Join(dstSkillsDir, name)
		if err := extractDir(templateFS, src, dst, nil); err != nil {
			return fmt.Errorf("extracting skill %q: %w", name, err)
		}
	}
	return nil
}

func extractDir(srcFS fs.FS, srcPath, dstPath string, skipFiles []string) error {
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dstPath, err)
	}
	entries, err := fs.ReadDir(srcFS, srcPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", srcPath, err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if shouldSkip(name, skipFiles) {
			continue
		}

		dstName := name
		if renamed, ok := dirRename[name]; ok {
			dstName = renamed
		}

		srcEntry := filepath.Join(srcPath, name)
		dstEntry := filepath.Join(dstPath, dstName)

		if entry.IsDir() {
			if err := os.MkdirAll(dstEntry, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", dstEntry, err)
			}
			if err := extractDir(srcFS, srcEntry, dstEntry, skipFiles); err != nil {
				return err
			}
		} else {
			f, err := srcFS.Open(srcEntry)
			if err != nil {
				return fmt.Errorf("opening %s: %w", srcEntry, err)
			}
			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				return fmt.Errorf("reading %s: %w", srcEntry, err)
			}
			if err := os.WriteFile(dstEntry, data, 0644); err != nil {
				return fmt.Errorf("writing %s: %w", dstEntry, err)
			}
		}
	}
	return nil
}

func shouldSkip(name string, skipFiles []string) bool {
	for _, s := range skipFiles {
		if name == s {
			return true
		}
	}
	return false
}

func HasSkills() bool {
	skillsDir := filepath.Join(templateRoot, "opencode", "skills")
	entries, err := fs.ReadDir(templateFS, skillsDir)
	return err == nil && len(entries) > 0
}

func ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(templateRoot, path)
	return templateFS.ReadFile(fullPath)
}

func SkillExists(name string) bool {
	skillsDir := filepath.Join(templateRoot, "opencode", "skills", name)
	info, err := fs.Stat(templateFS, skillsDir)
	return err == nil && info.IsDir()
}

func ListAvailableSkillNames(skillNames []string) []string {
	available := make([]string, 0, len(skillNames))
	for _, name := range skillNames {
		if SkillExists(name) {
			available = append(available, name)
		}
	}
	return available
}
