# forgectl

![Version](https://img.shields.io/badge/version-0.2.0-blue)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> **forge·ctl** — La herramienta universal de gestión de proyectos. Crea, inicializa y gestiona proyectos de desarrollo completos, localmente y en plataformas remotas (GitHub, GitLab, Gitea), con un solo comando.

`forgectl` es una CLI profesional inspirada en la `gh` CLI, pero **global** y no limitada a GitHub. Automatiza todas las tareas que un desarrollador hace al iniciar y mantener un proyecto: scaffolding, Git, repos remotos, CI/CD, licencias, changelogs y releases.

---

## ✨ Características

| Característica | Descripción |
|---|---|
| **Wizard interactivo** | Asistentes guiados para `new`, `config`, `template`, `release`, `remote` y más |
| **Un comando para todo** | Crea carpeta local + repo Git + repo remoto + README + licencia + estructura + commit + push |
| **Multi-proveedor** | GitHub, GitLab y Gitea (incluidos self-hosted) vía sus APIs REST |
| **Plantillas de proyecto** | Go, Python, Node.js, Rust y Minimal — con estructura estándar por lenguaje |
| **CI/CD generado** | GitHub Actions y GitLab CI listos según el lenguaje del proyecto |
| **Git hooks** | pre-commit, pre-push y commit-msg generados automáticamente |
| **Licencias** | MIT, Apache 2.0, GPL v3, BSD, ISC, MPL 2.0 y Unlicense |
| **Changelog automático** | Generado desde el historial de Git y los tags |
| **Releases y tags** | Con opción de auto-generar changelog y hacer push |
| **Diagnóstico** | `forgectl doctor` verifica Go, Git, Node, Python, Cargo y tu configuración |

---

## 🚀 Instalación

### Con `install.sh` (recomendado)

```bash
# Desde el repositorio
./install.sh                    # Instala en /usr/local/bin
./install.sh --prefix ~/.local  # Instala en ~/.local/bin
./install.sh --bin-dir ~/bin    # Instala en un directorio personalizado

# Desinstalar
./install.sh --remove
```

### Compilar desde el código fuente

```bash
git clone https://github.com/forgectl/forgectl.git
cd forgectl
make build            # o: go build -o forgectl .
sudo mv forgectl /usr/local/bin/
```

### Instalar con Go

```bash
go install github.com/forgectl/forgectl@latest
```

### Autocompletado (shell)

```bash
# Bash
forgectl completion bash > ~/.bash_completion.d/forgectl

# Zsh
forgectl completion zsh > "${fpath[1]}/_forgectl"

# Fish
forgectl completion fish > ~/.config/fish/completions/forgectl.fish
```

---

## 🏁 Quick Start

```bash
# 1. Configura tu proveedor remoto (interactivo)
forgectl remote add

# 2. Crea tu primer proyecto (wizard interactivo)
forgectl new

# 3. Trabaja con normalidad
cd mi-proyecto
forgectl status          # Estado del repo
forgectl sync            # pull + push con origin
forgectl release new v1.0.0 --push --changelog
```

---

## 📖 Guía de comandos

### Creación de proyectos

```bash
# Wizard interactivo completo
forgectl new

# Rápido con flags
forgectl new mi-app -t go -l MIT -d "Mi aplicación"
forgectl new mi-app --private --no-remote
forgectl new mi-app --ci github-actions --hooks
```

`forgectl new` genera automáticamente:

```
mi-app/
├── .git/                  # Repositorio Git inicializado
├── .github/workflows/     # CI (si se solicita)
├── .gitignore             # Según el lenguaje
├── LICENSE                # Licencia elegida
├── README.md              # README profesional
├── cmd/                   # Estructura Go
├── internal/
├── pkg/
├── src/
└── ...
```

### Git

```bash
forgectl init                          # Init repo local
forgectl sync                          # pull + push
forgectl push --tags                   # Push con tags
forgectl pull origin main              # Pull de rama específica
forgectl status                        # Branch, remote, clean, tags, commits
forgectl tag v1.0.0 -m "Versión 1.0"   # Crear tag
```

### Releases y changelog

```bash
forgectl release new                   # Wizard
forgectl release new v1.0.0 --push     # Directo + push
forgectl release new v1.0.0 --changelog  # + auto-genera CHANGELOG.md
forgectl release list                  # Ver releases
forgectl changelog gen                 # Generar changelog desde git
```

### Plantillas

```bash
forgectl template list                 # Ver plantillas disponibles
forgectl template info go              # Detalle de una plantilla
forgectl template apply rust           # Aplicar en directorio actual
```

### Licencias y README

```bash
forgectl license add MIT --author "Nombre"
forgectl license list
forgectl readme gen                     # Wizard con nombre, desc, licencia
```

### Configuración

```bash
forgectl config set                     # Wizard interactivo
forgectl config show                    # Ver toda la config
forgectl config get default_provider    # Obtener un valor
forgectl config reset                   # Resetear a defaults

# Configuración manual directa:
forgectl config set remotes.github.token ghp_xxx
forgectl config set remotes.github.username mi_user
forgectl config set default_license Apache
```

### Gestión de remotes y proyectos

```bash
forgectl remote list                    # Listar remotes configurados
forgectl remote add gitlab              # Añadir remoto

forgectl project list                   # Proyectos trackeados
forgectl project info mi-app            # Info detallada
forgectl project remove mi-app          # Quitar del tracking

forgectl search terraform               # Buscar repos en tu provider
forgectl open                           # Abrir repo en el navegador
forgectl open --dir                     # Abrir gestor de archivos
```

### Diagnóstico

```bash
forgectl doctor                         # Verifica todo tu entorno
forgectl doctor --fix                   # Intenta arreglar issues
```

---

## ⚙️ Configuración

La configuración se guarda en `~/.forgectl/config.yaml`:

```yaml
default_provider: github
default_license: MIT
default_template: go
remotes:
  github:
    provider: github
    token: ghp_xxxxxxxx
    username: mi_user
    base_url: https://github.com
    default: true
  gitlab:
    provider: gitlab
    token: glpat-xxxxxxxx
    username: mi_user
    base_url: https://gitlab.com
  gitea:
    provider: gitea
    token: xxxxxxxx
    username: admin
    base_url: https://gitea.example.com
projects:
  mi-app:
    name: mi-app
    description: Mi aplicación
    license: MIT
    template: go
    provider: github
    local_path: /home/user/mi-app
```

> **Seguridad**: los tokens se guardan en texto plano en tu directorio `~/.forgectl`. Considera protegerlo con `chmod 700 ~/.forgectl`.

---

## 🏗️ Arquitectura

```
forgectl/
├── main.go                          # Entry point
├── cmd/                             # Comandos Cobra
│   ├── root.go                      # Root + registro de comandos
│   ├── new.go                       # Forgectl new (con wizard)
│   ├── config.go                    # Configuración (con wizard)
│   ├── remote.go                    # Gestión de remotes
│   ├── remote.go                    # Doctor, open, search
│   └── ...                          # sync, push, pull, status, tags,
│                                    #   release, changelog, license, template, project
├── internal/
│   ├── config/                      # Viper configuration
│   ├── git/                         # Operaciones Git (go-git + CLI)
│   ├── templates/                   # README, licencias, scaffolding, CI/CD, hooks
│   └── project/                     # Manager de proyectos + Providers
│       ├── provider.go              # Interfaz Provider
│       ├── github.go                # API GitHub REST
│       ├── gitlab.go                # API GitLab v4
│       └── gitea.go                 # API Gitea v1
└── pkg/
    └── ui/                          # UI kit: prompts, spinner, colores, boxes
```

### Stack tecnológico

- **Go 1.22+** — tiempo de compilación, cero dependencias de runtime
- **[Cobra](https://github.com/spf13/cobra)** — CLI framework
- **[Viper](https://github.com/spf13/viper)** — configuración
- **[go-git](https://github.com/go-git/go-git)** — operaciones Git programáticas
- **[text/template](https://pkg.go.dev/text/template)** — plantillas de proyectos

---

## 🧩 Extender forgectl

### Añadir un nuevo proveedor

Crea `internal/project/<nombre>.go` e implementa la interfaz:

```go
func new<x>Client(cfg RemoteConfig) (Provider, error) { ... }
```

```go
type Provider interface {
    Name() string
    CreateRepository(name, desc string, private bool) (*Repository, error)
    DeleteRepository(name string) error
    GetRepository(name string) (*Repository, error)
    ListRepositories() ([]*Repository, error)
    GetCurrentUsername() (string, error)
}
```

Regístralo en `NewProvider` en `internal/project/provider.go`.

### Añadir una plantilla de proyecto

En `internal/templates/scaffold.go`, añade al mapa `ScaffoldTemplates`:

```go
"java": {
    Name:        "Java",
    Description: "Standard Java project layout",
    Directories: []string{"src/main/java", "src/test/java", "...",},
},
```

Añade también su `.gitignore` y plantillas de CI en `internal/templates/ci.go`.

### Añadir un comando

Crea `cmd/mi_comando.go` y regístralo en `cmd/root.go`:

```go
var miComandoCmd = &cobra.Command{
    Use:   "mi-comando",
    Short: "Hace algo increíble",
    RunE: func(cmd *cobra.Command, args []string) error {
        // ¡Tu lógica aquí!
        return nil
    },
}

func init() {
    rootCmd.AddCommand(miComandoCmd)
}
```

### Añadir una licencia

En `internal/templates/templates.go`, añade al mapa `licenses` de `GetLicense` y a `ListLicenses`.

---

## 🛠️ Roadmap

- [x] Wizard interactivo completo
- [x] Multi-proveedor (GitHub, GitLab, Gitea)
- [x] CI/CD generado (Actions, GitLab CI)
- [x] Git hooks automáticos
- [x] Changelog y releases
- [x] `doctor` de diagnóstico
- [ ] Issues iniciales automáticos
- [ ] Codeowners, .editorconfig, SECURITY.md
- [ ] Soporte de organizaciones para remotes
- [ ] Rebase interactivo y gestión de PRs
- [ ] GitHub/GitLab/Gitea releases via API
- [ ] Modo `--json` para scripting

---

## 📄 Licencia

[MIT](LICENSE) © [forgectl contributors](https://github.com/forgectl/forgectl/graphs/contributors)

---

*Made with ❤️ and Go.*