# new-repo

**new-repo** é uma CLI em Go que automatiza o scaffolding de projetos prontos para o [opencode](https://opencode.ai). Em um único comando, ela cria a estrutura de diretórios, configura o Git, instala skills e gera os arquivos de configuração necessários.

## Instalação

```bash
go build -o new-repo ./cmd/new-repo
```

Requer **Go 1.22+** e **Git** instalados no sistema.

## Uso

```bash
new-repo [flags] <project-name>
```

Exemplo:

```bash
new-repo meu-projeto
```

Isso cria a seguinte estrutura:

```
meu-projeto/
├── .gitignore
├── SKILLS.md
├── .opencode/
│   └── skills/
│       ├── coding-guidelines/
│       ├── context7/
│       └── ... (demais skills instaladas)
└── .git/
```

## Flags

| Flag | Descrição |
|------|-----------|
| `--force` | Sobrescreve diretório existente |
| `--skills-file <path>` | Usa um arquivo `SKILLS.md` personalizado em vez do embutido |
| `--help` | Exibe a ajuda |

## Funcionalidades

- **Validação de nome** — apenas caracteres alfanuméricos, hífen e underscore
- **`.gitignore` padrão** — entradas para Go, OS e IDEs
- **Skills` embutidas** — 13 skills opencode são copiadas para `.opencode/skills/`
- **Skills personalizadas** — use `--skills-file` para fornecer seu próprio conjunto
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
