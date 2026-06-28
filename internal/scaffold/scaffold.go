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
	"openspec-propose": join(
		`---`,
		`name: openspec-propose`,
		`description: Propose a new change with all artifacts generated in one step.`,
		`license: MIT`,
		`compatibility: Requires openspec CLI.`,
		`metadata:`,
		`  author: openspec`,
		`  version: "1.0"`,
		`  generatedBy: "new-repo"`,
		`---`,
		``,
		`Propose a new change - create the change and generate all artifacts in one step.`,
		``,
		`Creates a change with artifacts: proposal.md (what & why), design.md (how), tasks.md (implementation steps).`,
		``,
		`When ready to implement, run /opsx:apply`,
		``,
		`**Steps**`,
		``,
		`1. If no clear input provided, ask what they want to build`,
		`2. Create the change directory: openspec new change "<name>"`,
		`3. Get the artifact build order: openspec status --change "<name>" --json`,
		`4. Create artifacts in dependency order until apply-ready`,
		`5. Show final status`,
		``,
		`**Guardrails**`,
		`- Create ALL artifacts needed for implementation`,
		`- Always read dependency artifacts before creating a new one`,
		`- If context is unclear, ask the user`,
		`- Verify each artifact file exists after writing`,
	),
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
	"openspec-apply-change": join(
		`---`,
		`name: openspec-apply-change`,
		`description: Implement tasks from an OpenSpec change.`,
		`license: MIT`,
		`compatibility: Requires openspec CLI.`,
		`metadata:`,
		`  author: openspec`,
		`  version: "1.0"`,
		`  generatedBy: "new-repo"`,
		`---`,
		``,
		`Implement tasks from an OpenSpec change.`,
		``,
		`**Steps**`,
		``,
		`1. Select the change (infer from context, auto-select if one active change, or list options)`,
		`2. Check status: openspec status --change "<name>" --json`,
		`3. Get apply instructions: openspec instructions apply --change "<name>" --json`,
		`4. Read all context files`,
		`5. Show current progress`,
		`6. Implement tasks in a loop until done or blocked`,
		`7. On completion or pause, show status`,
		``,
		`**Guardrails**`,
		`- Keep going through tasks until done or blocked`,
		`- Always read context files before starting`,
		`- If task is ambiguous, pause and ask`,
		`- Keep code changes minimal and scoped to each task`,
		`- Update task checkboxes immediately after completing each task`,
		`- Pause on errors, blockers, or unclear requirements`,
	),
	"openspec-sync-specs": join(
		`---`,
		`name: openspec-sync-specs`,
		`description: Sync specifications from proposals and design artifacts into the specs directory.`,
		`license: MIT`,
		`compatibility: Requires openspec CLI.`,
		`metadata:`,
		`  author: openspec`,
		`  version: "1.0"`,
		`  generatedBy: "new-repo"`,
		`---`,
		``,
		`Sync specifications from proposals and design artifacts into the specs directory.`,
		``,
		`**Steps**`,
		``,
		`1. List active changes: openspec list --json`,
		`2. For each change with artifacts to sync, run: openspec instructions sync --change "<name>" --json`,
		`3. Read the proposal, design, and tasks artifacts for context`,
		`4. Extract capability specifications and update specs/<capability>/spec.md`,
		``,
		`**Guardrails**`,
		`- Only sync finalized decisions, not in-progress thinking`,
		`- Maintain separation of concerns between specs, design, and implementation`,
		`- Update the spec index when adding new capabilities`,
	),
	"openspec-archive-change": join(
		`---`,
		`name: openspec-archive-change`,
		`description: Archive a completed change.`,
		`license: MIT`,
		`compatibility: Requires openspec CLI.`,
		`metadata:`,
		`  author: openspec`,
		`  version: "1.0"`,
		`  generatedBy: "new-repo"`,
		`---`,
		``,
		`Archive a completed change.`,
		``,
		`**Steps**`,
		``,
		`1. Select the change to archive`,
		`2. Confirm all tasks are complete`,
		`3. Sync any remaining specs: openspec instructions sync --change "<name>" --json`,
		`4. Archive the change: openspec archive change "<name>"`,
		``,
		`**Guardrails**`,
		`- Only archive changes where all tasks are complete`,
		`- Sync specs before archiving to capture all decisions`,
		`- Verify the archive was created successfully`,
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
		filepath.Join(dir, ".opencode", "skills", "openspec-propose"),
		filepath.Join(dir, ".opencode", "skills", "openspec-explore"),
		filepath.Join(dir, ".opencode", "skills", "openspec-apply-change"),
		filepath.Join(dir, ".opencode", "skills", "openspec-sync-specs"),
		filepath.Join(dir, ".opencode", "skills", "openspec-archive-change"),
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
