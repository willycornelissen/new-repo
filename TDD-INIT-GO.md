# TDD — Gerador de Projeto OpenCode CLI

| Campo            | Valor                                            |
| ---------------- | ------------------------------------------------ |
| Tech Lead        | @dev                                             |
| Equipe           | 1 desenvolvedor                                  |
| Status           | Draft                                            |
| Criado           | 2026-06-27                                       |
| Última alteração | 2026-06-27                                       |

---

## Contexto

Atualmente, iniciar um novo projeto com opencode exige uma série de passos manuais: criar o diretório, inicializar o repositório Git, copiar as skills do opencode, configurar o arquivo SKILLS.md e garantir que tudo esteja funcional. Esse processo repetitivo consome tempo e é propenso a erros de configuração.

Este projeto faz parte do domínio de **automação de setup de desenvolvimento** (scaffolding). O público-alvo são desenvolvedores que utilizam o opencode como ferramenta principal e desejam bootstrap rápido de novos projetos.

---

## Definição do Problema e Motivação

### Problemas a Resolver

- **Setup manual repetitivo**: Criar diretório, `.gitignore`, inicializar repositório, baixar skills → ~5 minutos por projeto.
- **Inconsistência entre projetos**: Cada desenvolvedor pode esquecer de copiar skills ou configurar o SKILLS.md de forma diferente.
- **Falta de padronização**: Sem um scaffold, não há garantia de que a estrutura de diretórios e arquivos de configuração siga o mesmo padrão.

### Por Que Agora?

- O uso do opencode está crescendo na equipe/comunidade.
- Projetos novos surgem com frequência — a automação reduz fricção e erros.
- Ferramenta complementar ao ecossistema opencode.

### Impacto de Não Resolver

- Perda de ~5 min por setup manual × N projetos = dezenas de minutos acumulados.
- Risco de projetos sem as skills corretas, levando a comandos não reconhecidos.
- Dificuldade de onboarding de novos membros.

---

## Escopo

### ✅ Incluso (MVP)

- CLI em Go que aceita um parâmetro (nome do projeto).
- Cria o diretório com o nome do projeto.
- Inicializa repositório Git (`git init`).
- Gera `.gitignore` padrão para projetos Go.
- Cria um arquivo `SKILLS.md` no diretório gerado.
- Lê as skills de um arquivo `SKILLS.md` de origem (embutido ou via flag) e as copia para o diretório de skills do opencode (`~/.config/opencode/skills/` ou diretório local `.opencode/skills/`).
- Gera estrutura de diretórios `.opencode/` configurada.
- Gera estrutura `openspec/` com `specs/`, `changes/archive/` e `config.yaml` (pronto para OpenSpec).
- Instala skills OpenSpec (`openspec-propose`, `openspec-explore`, `openspec-apply-change`, `openspec-sync-specs`, `openspec-archive-change`) em `.opencode/skills/`.
- Flag `--help` documentada.

### ❌ Excluído (V1)

- Suporte a múltiplos templates de projeto (apenas opencode).
- Integração com GitHub (criar repositório remoto automaticamente).
- Flags interativas (`--interactive`).
- Templates de CI/CD.
- Suporte a plugins ou extensões.

### 🔮 Versões Futuras

- Flag `--template` para diferentes tipos de projeto.
- Publicação via `go install`.
- Suporte a GitHub CLI (`gh repo create`).

---

## Solução Técnica

### Arquitetura Geral

```
[Usuário] → new-repo <nome> → [CLI em Go]
                                 │
                                 ├─ Cria diretório <nome>/
                                 ├─ Cria .gitignore
                                 ├─ Cria SKILLS.md
                                 ├─ Cria .opencode/skills/ (com skills locais)
                                 ├─ Cria openspec/ (specs, changes, config.yaml)
                                 ├─ Instala skills OpenSpec
                                 └─ Executa git init
```

### Fluxo de Dados

1. Usuário executa: `new-repo meu-projeto`
2. CLI valida o nome (não vazio, sem caracteres especiais).
3. Cria diretório `meu-projeto/` com permissão 0755.
4. Copia SKILLS.md de origem (recurso embutido ou `--skills-file`) para `meu-projeto/SKILLS.md`.
5. Lê o SKILLS.md, extrai os nomes das skills, copia cada skill do diretório de skills local para `meu-projeto/.opencode/skills/<skill>/`.
6. Cria `.gitignore` com entradas padrão Go.
7. Cria `openspec/` com `specs/`, `changes/archive/`, `config.yaml` e skills OpenSpec em `.opencode/skills/openspec-*/`.
8. Executa `git init` dentro do diretório gerado.
9. Exibe mensagem de sucesso com caminho absoluto.

### Estrutura de Arquivos Gerada

```
meu-projeto/
├── .gitignore
├── SKILLS.md
├── openspec/
│   ├── specs/
│   ├── changes/
│   │   └── archive/
│   └── config.yaml
├── .opencode/
│   └── skills/
│       ├── openspec-propose/
│       │   └── SKILL.md
│       ├── openspec-explore/
│       │   └── SKILL.md
│       ├── openspec-apply-change/
│       │   └── SKILL.md
│       ├── openspec-sync-specs/
│       │   └── SKILL.md
│       ├── openspec-archive-change/
│       │   └── SKILL.md
│       ├── skill-a/
│       │   └── SKILL.md
│       ├── skill-b/
│       │   └── SKILL.md
│       └── ...
└── .git/
```

### Comandos da CLI

| Comando             | Descrição                                        |
| ------------------- | ------------------------------------------------ |
| `new-repo <nome>`     | Gera projeto com nome especificado               |
| `new-repo --help`     | Exibe ajuda                                      |
| `new-repo --skills-file <path>` | Especifica arquivo SKILLS.md personalizado |

### Decisões Técnicas

- **Linguagem**: Go 1.22+ — binário único, sem dependência de runtime, ideal para CLI.
- **Pacotes**: Apenas stdlib (`os`, `os/exec`, `flag`, `path/filepath`, `fmt`, `io`). Sem dependências externas para manter o binário mínimo.
- **Ponto de entrada**: `cmd/new-repo/main.go`.
- **Organização**: Responsabilidades separadas em pacotes (`scaffold`, `git`, `config`, `skill`).

### Diagrama de Pacotes

```
new-repo/
├── main.go              # Ponto de entrada, parsing de flags
├── internal/
│   ├── scaffold/
│   │   └── scaffold.go  # Criação de diretórios e arquivos
│   ├── skill/
│   │   └── skill.go     # Leitura/instalação de skills do SKILLS.md
│   ├── git/
│   │   └── git.go       # Inicialização do repositório
│   └── config/
│       └── config.go    # Configurações (caminhos, templates)
└── go.mod
```

---

## Riscos

| Risco                        | Impacto | Probabilidade | Mitigação                                           |
| ---------------------------- | ------- | ------------- | --------------------------------------------------- |
| Nome de projeto inválido     | Baixo   | Média         | Validação com regex (apenas alfanumérico + hífen)   |
| Diretório já existe          | Médio   | Média         | Perguntar antes de sobrescrever (`--force`)         |
| Git não instalado no sistema | Alto    | Baixa         | Verificar `git --version` antes, erro claro         |
| Skills de origem não encontradas | Médio | Baixa | Fallback para skills embutidas, warning no output   |
| Permissão negada ao criar dir | Alto   | Baixa         | Mensagem de erro clara com diagnóstico              |

---

## Estratégia de Testes

| Tipo de Teste        | Escopo                       | Cobertura Alvo | Abordagem                      |
| -------------------- | ---------------------------- | -------------- | ------------------------------ |
| Testes Unitários     | Validação, scaffold, skill   | > 80%          | `go test` com tabelas          |
| Testes de Integração | Criação de diretório + git   | Fluxo crítico  | Diretório temporário (`t.TempDir`) |
| Testes de CLI        | Parsing de flags             | Cobertura total | `os.Args` simulado            |

### Cenários Críticos

- Nome vazio → erro.
- Nome com caracteres especiais → erro.
- Diretório já existente → confirmação.
- `SKILLS.md` vazio → projeto criado sem skills (warning).
- `git init` falha → erro específico.
- Flag `--help` → exibe uso.

---

## Plano de Implementação

| Fase               | Tarefa                        | Descrição                               | Estimativa |
| ------------------ | ----------------------------- | --------------------------------------- | ---------- |
| **Fase 1 — Setup** | Inicializar módulo Go         | `go mod init`, estrutura de diretórios  | 0.5d       |
|                    | Criar SKILLS.md de origem     | Skills a serem copiadas nos projetos    | 0.5d       |
| **Fase 2 — Core**  | Scaffold do diretório         | Criar diretório e `.gitignore`          | 1d         |
|                    | Parser de SKILLS.md           | Extrair skills e copiar para destino    | 1d         |
|                    | Inicialização Git             | Executar `git init` (pacote `os/exec`)  | 0.5d       |
| **Fase 3 — CLI**   | Ponto de entrada              | `main.go` com `flag` package            | 0.5d       |
|                    | Validação de parâmetros       | Regex, mensagens de erro                | 0.5d       |
|                    | Flag `--help` e `--force`     | Documentação inline                     | 0.5d       |
| **Fase 4 — Testes**| Testes unitários              | Scaffold, validação, skill              | 1d         |
|                    | Testes de integração          | Criação temporária + git init           | 1d         |
| **Fase 5 — Documentação** | README, exemplos       | Uso, flags, exemplos                    | 0.5d       |

**Total estimado**: ~7 dias úteis

**Dependências**:
- Go 1.22+ instalado na máquina de desenvolvimento.
- Git instalado no ambiente de testes.
- SKILLS.md de origem definido antes da Fase 2.

---

## Decisões Registradas

| Decisão | Opção Escolhida | Justificativa |
|---------|----------------|---------------|
| Localização das skills | `.opencode/skills/` dentro do projeto | Skills locais ao projeto, portabilidade sem depender da máquina hospedeira |
| Origem do SKILLS.md | Embutido no binário via `//go:embed` | Binário autocontido, sem dependência de arquivo externo, distribuição simplificada |
| Flag `--no-git` | Não implementar | Não há caso de uso identificado que justifique a flag |
| Scaffold OpenSpec | Integrado no `new-repo` | Projeto já sai com `openspec/` estruturado e skills `/opsx:*` instaladas, sem necessidade de rodar `openspec init` separadamente |
| Skills OpenSpec | Geradas nativamente em `.opencode/skills/openspec-*/` | Skills auto-contidas no binário via funções Go (sem depender de arquivos externos ou `openspec init`) |
