package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

const openspecConfig = `schema: spec-driven

# Project context (optional)
# This is shown to AI when creating artifacts.
# Add your tech stack, conventions, style guides, domain knowledge, etc.

# Per-artifact rules (optional)
# Add custom rules for specific artifacts.
`

var openSpecSkills = map[string]string{
	"openspec-explore": join(
		`---`,
		`name: openspec-explore`,
		`description: Enter explore mode - a thinking partner for exploring ideas and investigating problems.`,
		`license: MIT`,
		`compatibility: Requires openspec CLI.`,
		`metadata:`,
		`  author: openspec`,
		`  version: "1.0"`,
		`  generatedBy: "new-repo"`,
		`---`,
		``,
		`Enter explore mode. Think deeply. Visualize freely. Follow the conversation wherever it goes.`,
		``,
		`IMPORTANT: Explore mode is for thinking, not implementing. Read files and investigate, but do NOT write code or implement features.`,
		``,
		`**The Stance**`,
		`- Curious, not prescriptive - ask questions that emerge naturally`,
		`- Visual - use diagrams liberally when they clarify thinking`,
		`- Adaptive - follow interesting threads, pivot when new info emerges`,
		`- Grounded - explore the actual codebase, don't just theorize`,
		``,
		`**What You Might Do**`,
		`- Explore the problem space`,
		`- Investigate the codebase`,
		`- Compare options`,
		`- Surface risks and unknowns`,
		``,
		`**OpenSpec Awareness**`,
		`- Check "openspec list --json" for context`,
		`- Reference existing artifacts in conversation`,
		`- Offer to capture insights when decisions are made`,
		``,
		`**Guardrails**`,
		`- Never write code or implement features`,
		`- Don't fake understanding - dig deeper`,
		`- Don't force structure - let patterns emerge`,
		`- Do visualize - a good diagram is worth many paragraphs`,
	),
}

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

func join(lines ...string) string {
	var result string
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

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

func CreateOpenSpecDirs(dir string) error {
	dirs := []string{
		filepath.Join(dir, "openspec", "specs"),
		filepath.Join(dir, "openspec", "changes", "archive"),
		filepath.Join(dir, ".opencode", "skills", "openspec-explore"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func WriteOpenSpecConfig(dir string) error {
	path := filepath.Join(dir, "openspec", "config.yaml")
	return os.WriteFile(path, []byte(openspecConfig), 0644)
}

func WriteOpenSpecSkills(dir string) error {
	for name, content := range openSpecSkills {
		path := filepath.Join(dir, ".opencode", "skills", name, "SKILL.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}
