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

const agentsMD = `# AGENTS.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
` + "```" + `
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
` + "```" + `

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
`

const readmeMD = `# AI Template

Template para projetos de desenvolvimento de software assistido por IA usando [opencode](https://opencode.ai).

## Estrutura

` + "```" + `
├── .opencode/          # Configuração do opencode
│   ├── commands/       # Comandos personalizados (ex: opsx-explore)
│   ├── skills/         # Skills instaladas
│   └── package.json    # Dependências do opencode
├── openspec/           # OpenSpec — framework de desenvolvimento spec-driven
│   └── config.yaml     # Configuração do OpenSpec
├── SKILLS.md           # Lista de skills disponíveis
└── README.md
` + "```" + `

## Skills Instaladas

| Skill | Descrição |
|-------|-----------|
| **coding-guidelines** | Diretrizes para reduzir erros comuns de codificação em LLMs |
| **context7** | Busca documentação atualizada de bibliotecas e frameworks |
| **docs-writer** | Escrita e revisão de documentação |
| **domain-analysis** | Mapeamento de domínios com DDD Strategic Design |
| **excalidraw-studio** | Geração de diagramas Excalidraw |
| **graphify** | Criação de grafos de conhecimento a partir de código/docs |
| **mermaid-studio** | Diagramas Mermaid (SVG/PNG/ASCII) |
| **skill-architect** | Criação de novas skills |
| **spec-driven-eval** | Avaliação de implementação contra PRD/spec |
| **technical-design-doc-creator** | Criação de Documentos de Design Técnico |
| **code-review-skill** | Revisão de código estruturada para 20+ linguagens/frameworks |
| **tlc-spec-driven** | Planejamento em 4 fases: Spec → Design → Tasks → Execute |

## OpenSpec

O diretório openspec/ contém a configuração do framework **OpenSpec**, que permite:

- **Propostas** — Definição de escopo e objetivos
- **Designs** — Decisões técnicas e arquiteturais
- **Tasks** — Divisão de trabalho em tarefas atômicas
- **Specs** — Especificações detalhadas de funcionalidades

## Comandos

| Comando | Descrição |
|---------|-----------|
| ` + "`" + `/opsx-explore` + "`" + ` | Modo exploratório para pensar, investigar e esclarecer requisitos sem implementar |

## Uso

1. Ative uma skill relevante para sua tarefa (ex: ` + "`" + `tlc-spec-driven` + "`" + ` para planejamento)
2. Use o OpenSpec para capturar requisitos e decisões
3. Implemente com verificação atômica por tarefa
`

func WriteAgentsMD(dir string) error {
	path := filepath.Join(dir, "AGENTS.md")
	return os.WriteFile(path, []byte(agentsMD), 0644)
}

func WriteReadmeMD(dir string) error {
	path := filepath.Join(dir, "README.md")
	return os.WriteFile(path, []byte(readmeMD), 0644)
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
	dirs := []string{
		filepath.Join(dir, ".opencode", "skills"),
		filepath.Join(dir, ".opencode", "commands"),
		filepath.Join(dir, ".opencode", "docs"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
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
