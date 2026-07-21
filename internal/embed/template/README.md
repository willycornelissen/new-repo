# AI Template

Template para projetos de desenvolvimento de software assistido por IA usando [opencode](https://opencode.ai).

## Estrutura

```
├── .opencode/          # Configuração do opencode
│   ├── commands/       # Comandos personalizados (ex: explore)
│   ├── skills/         # Skills instaladas
│   ├── docs/           # Documentação de padrões de código
│   └── package.json    # Dependências do opencode
├── AGENTS.md           # Diretrizes comportamentais para o agente de IA
└── SKILLS.md           # Lista de skills disponíveis
```

## Skills Instaladas

| Skill | Descrição |
|-------|-----------|
| **code-review-skill** | Revisão de código estruturada para 20+ linguagens/frameworks |
| **context7** | Busca documentação atualizada de bibliotecas e frameworks |
| **docs-writer** | Escrita e revisão de documentação |
| **domain-analysis** | Mapeamento de domínios com DDD Strategic Design |
| **excalidraw-studio** | Geração de diagramas Excalidraw |
| **graphify** | Criação de grafos de conhecimento a partir de código/docs |
| **mermaid-studio** | Diagramas Mermaid (SVG/PNG/ASCII) |
| **modular-architecture** | Arquitetura modular com bounded contexts, facades e 10 princípios de design |
| **office-hours** | YC Office Hours: diagnóstico de produto para startups ou brainstorm para builders |
| **skill-architect** | Criação de novas skills |
| **spec-driven-eval** | Avaliação de implementação contra PRD/spec |
| **technical-design-doc-creator** | Criação de Documentos de Design Técnico |
| **tlc-spec-driven** | Planejamento em 4 fases: Spec → Design → Tasks → Execute |

## Comandos

| Comando | Descrição |
|---------|-----------|
| `/doc` | Gera documentação do projeto e código em `documentation/` |
| `/generate` | Executa o plano de implementação e gera código em `src/` a partir de `specification/plan.md` |
| `/new-feature` | Cria uma nova funcionalidade completa: TDD, PRD, roadmap, plano, código, review e documentação. Uso: `/new-feature <slug>: <nome-da-feature>` |
| `/plan` | Gera plano de implementação a partir de `specification/prd.md` e `specification/roadmap.md` |
| `/prd` | Cria um Product Requirements Document (PRD) para um projeto e salva em `specification/prd.md` |
| `/idea` | YC Office Hours — explora ideias, problemas ou conceitos de produto e salva a ideia bruta em `specification/idea.md` e o resumo em `specification/research.md` |
| `/review` | Revisa o código em `src/` comparando com `specification/prd.md` e `specification/plan.md` |
| `/roadmap` | Gera roadmap com features e tarefas a partir de `specification/tdd.md` e `specification/prd.md` |
| `/tdd` | Cria um Documento de Design Técnico (TDD) para um projeto e salva em `specification/tdd.md` |
