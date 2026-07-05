package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const TemplateRepoURL = "https://github.com/willycornelissen/ai-template"

func Init(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git init failed: %w\n%s", err, out)
	}
	return nil
}

func CloneTemplate(dir string) error {
	tmpDir, err := os.MkdirTemp("", "new-repo-template-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth", "1", TemplateRepoURL, tmpDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w\n%s", err, out)
	}

	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		return fmt.Errorf("removing template .git: %w", err)
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		src := filepath.Join(tmpDir, entry.Name())
		dst := filepath.Join(dir, entry.Name())
		if err := copyRecursive(src, dst); err != nil {
			return fmt.Errorf("copying %s: %w", entry.Name(), err)
		}
	}

	return Init(dir)
}

func CloneTemplateSkills() (skillsDir string, cleanup func(), err error) {
	tmpDir, err := os.MkdirTemp("", "new-repo-skills-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup = func() { os.RemoveAll(tmpDir) }

	cmd := exec.Command("git", "clone", "--depth", "1", TemplateRepoURL, tmpDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning template: %w\n%s", err, out)
	}

	skillsDir = filepath.Join(tmpDir, ".opencode", "skills")
	if _, err := os.Stat(skillsDir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("skills dir not found in template: %w", err)
	}

	return skillsDir, cleanup, nil
}

func copyRecursive(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := copyRecursive(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func IsAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
