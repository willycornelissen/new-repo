package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

const gitignore = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Go workspace
go.work

# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
*.swp
*.swo
`

func CreateProjectDir(name string, force bool) (string, error) {
	info, err := os.Stat(name)
	if err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("%q already exists as a file", name)
		}
		if !force {
			return "", fmt.Errorf("directory %q already exists; use --force to overwrite", name)
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if err := os.MkdirAll(name, 0755); err != nil {
		return "", err
	}

	return filepath.Abs(name)
}

func WriteGitignore(dir string) error {
	path := filepath.Join(dir, ".gitignore")
	return os.WriteFile(path, []byte(gitignore), 0644)
}

func WriteSkillsMD(dir string, content string) error {
	path := filepath.Join(dir, "SKILLS.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func CreateOpenCodeDirs(dir string) error {
	skillsDir := filepath.Join(dir, ".opencode", "skills")
	return os.MkdirAll(skillsDir, 0755)
}
