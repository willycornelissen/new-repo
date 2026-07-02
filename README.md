# new-repo

**new-repo** é uma CLI em Go que automatiza o scaffolding de projetos prontos para o [opencode](https://opencode.ai). Em um único comando, ela cria a estrutura de diretórios, configura o Git, instala skills e gera os arquivos de configuração necessários.

## Instalação

```bash
go build -o new-repo ./cmd/new-repo
```

Requer **Go 1.22+** e **Git** instalados no sistema.

## Uso

```bash
new-repo [flags] <project-name | .>
```

Exemplos:

```bash
# Criar um novo projeto em um subdiretório
new-repo meu-projeto

# Scaffold no diretório atual a partir de um template remoto
new-repo .
```

O comando `new-repo meu-projeto` cria a seguinte estrutura:

```
meu-projeto/
├── .gitignore
├── SKILLS.md
├── openspec/
│   ├── specs/              # Suas especificações (source of truth)
│   ├── changes/            # Mudanças propostas
│   │   └── archive/        # Mudanças arquivadas
│   └── config.yaml         # Configuração do OpenSpec
├── .opencode/
│   └── skills/
│       ├── coding-guidelines/
│       ├── context7/
│       ├── openspec-propose/       # /opsx:propose
│       ├── openspec-explore/       # /opsx:explore
│       ├── openspec-apply-change/  # /opsx:apply
│       ├── openspec-sync-specs/    # /opsx:sync
│       ├── openspec-archive-change/# /opsx:archive
│       └── ... (demais skills instaladas)
└── .git/
```

## Flags

| Flag | Descrição |
|------|-----------|
| `--force` | Sobrescreve diretório existente (apenas para `new-repo <name>`) |
| `--skills-file <path>` | Usa um arquivo `SKILLS.md` personalizado em vez do embutido (apenas para `new-repo <name>`) |
| `--help` | Exibe a ajuda |

### `new-repo .`

Quando usado com `.` como argumento, o comando não cria um novo diretório. Em vez disso, ele clona o template remoto [`ai-template`](https://github.com/willycornelissen/ai-template) diretamente no diretório atual, copia todo o conteúdo, remove o `.git` original e inicia um novo repositório Git.

Esse modo **ignora** as flags `--force` e `--skills-file`, e não passa pelas etapas de validação de nome, geração de `.gitignore` ou instalação de skills — tudo vem do template clonado.

```bash
new-repo .
# installed ai-template at /caminho/do/diretorio
```

## Funcionalidades

- **Validação de nome** — apenas caracteres alfanuméricos, hífen e underscore
- **`.gitignore` padrão** — entradas para Go, OS e IDEs
- **Skills embutidas** — skills opencode são copiadas para `.opencode/skills/`
- **Skills personalizadas** — use `--skills-file` para fornecer seu próprio conjunto
- **OpenSpec pronto** — estrutura `openspec/` criada com `specs/`, `changes/archive/` e `config.yaml`
- **Comandos OpenSpec** — skills `/opsx:propose`, `/opsx:explore`, `/opsx:apply`, `/opsx:sync` e `/opsx:archive` pré-instaladas
- **`git init` automático** — repositório iniciado logo após a criação

## Desenvolvimento

### Comandos

```bash
go build -o new-repo ./cmd/new-repo   # compilar
go test ./...                           # rodar testes
go vet ./...                            # análise estática
```

### Estrutura do projeto

```
new-repo/
├── cmd/new-repo/main.go     # ponto de entrada
├── internal/
│   ├── scaffold/            # criação de diretórios e arquivos
│   ├── skill/               # parsing e instalação de skills
│   ├── git/                 # inicialização do repositório
│   └── config/              # configurações e caminhos
├── SKILLS.md                # skills embutidas no binário
└── go.mod
```

O projeto usa **apenas a stdlib** — sem dependências externas.

## Licença

MIT
