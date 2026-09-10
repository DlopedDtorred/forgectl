# forgectl

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
![CI](https://img.shields.io/badge/CI-passing-brightgreen)
![Release](https://img.shields.io/github/v/release/forgectl/forgectl)

> **forge·ctl** — La herramienta universal de gestión de proyectos. Crea, inicializa y gestiona proyectos de desarrollo completos, localmente y en plataformas remotas (GitHub, GitLab, Gitea), con un solo comando.

`forgectl` es una CLI profesional inspirada en la `gh` CLI, pero **global** y no limitada a GitHub. Automatiza todas las tareas que un desarrollador hace al iniciar y mantener un proyecto: scaffolding, Git, repos remotos, CI/CD, licencias, changelogs, releases e issues.

---

## Características

| Característica | Descripción |
|---|---|
| **Wizard interactivo** | Asistentes guiados para `new`, `config`, `template`, `release`, `remote` y más |
| **Un comando para todo** | Crea carpeta local + repo Git + repo remoto + README + licencia + estructura + commit + push |
| **Multi-proveedor** | GitHub, GitLab y Gitea (incluidos self-hosted) vía sus APIs REST |
| **Organizaciones** | Soporte de organizaciones para repositorios remotos |
| **Plantillas de proyecto** | 14 plantillas: Go, Python, Node.js, TypeScript, Rust, Java, Kotlin, C#, PHP, Ruby, Swift, Dart, C/C++ y Minimal |
| **CI/CD generado** | GitHub Actions y GitLab CI listos según el lenguaje del proyecto |
| **Git hooks** | pre-commit, pre-push y commit-msg generados automáticamente |
| **Licencias** | 12 licencias: MIT, Apache, GPL, BSD, ISC, MPL, Unlicense, AGPL, LGPL, WTFPL, EPL y Artistic |
| **Changelog automático** | Generado desde el historial de Git y los tags |
| **Releases** | Tags locales + releases vía API en GitHub, GitLab y Gitea |
| **Issues** | Crear y listar issues directamente desde la terminal |
| **Archivos de proyecto** | Genera `.editorconfig`, `CODEOWNERS` y `SECURITY.md` |
| **Modo JSON** | Salida `--json` para scripting y automatización |
| **Diagnóstico** | `forgectl doctor` verifica Go, Git, Node, Python, Cargo y tu configuración |

---

## Instalación

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

## Quick Start

```bash
# 1. Configura tu proveedor remoto (interactivo)
forgectl remote add

# 2. Crea tu primer proyecto (wizard interactivo)
forgectl new

# 3. Trabaja con normalidad
cd mi-proyecto
forgectl status          # Estado del repo
forgectl sync            # pull + push con origin

# 4. Crea un release en el proveedor remoto
forgectl release new v1.0.0 --remote --push --changelog

# 5. Gestiona issues
forgectl issue new "Bug: error en login" --labels "bug,urgent"
forgectl issue list
```

---

## Lenguajes soportados

forgectl incluye 14 plantillas de proyecto con estructura estándar, `.gitignore` y configuración de CI/CD específicos para cada lenguaje:

| Plantilla | Lenguaje | Descripción |
|---|---|---|
| `go` | Go | Estructura estándar Go (`cmd/`, `internal/`, `pkg/`) |
| `python` | Python | Estructura estándar Python (`src/`, `tests/`) |
| `node` | Node.js | Estructura estándar Node.js (`src/`, `lib/`, `test/`) |
| `typescript` | TypeScript | Estructura estándar TypeScript (`src/`, `lib/`, `test/`) |
| `rust` | Rust | Estructura estándar Rust (`src/`, `crates/`, `benches/`) |
| `java` | Java | Estructura estándar Java (`src/main/java`, `src/test/java`) |
| `kotlin` | Kotlin | Estructura estándar Kotlin (`src/main/kotlin`, `src/test/kotlin`) |
| `csharp` | C# | Estructura estándar .NET (`src/`, `tests/`) |
| `php` | PHP | Estructura estándar PHP (`src/`, `public/`, `tests/`) |
| `ruby` | Ruby | Estructura estándar Ruby (`app/`, `lib/`, `spec/`) |
| `swift` | Swift | Estructura estándar Swift (`Sources/`, `Tests/`) |
| `dart` | Dart | Estructura estándar Dart (`lib/`, `bin/`, `test/`) |
| `c-cpp` | C/C++ | Estructura estándar C/C++ (`src/`, `include/`, `lib/`) |
| `minimal` | Minimal | Estructura mínima (`src/`, `docs/`) |

---

## Licencias soportadas

forgectl puede generar 12 tipos de licencias automáticamente:

| Licencia | Identificador | Descripción |
|---|---|---|
| MIT | `MIT` | Corta y permisiva |
| Apache 2.0 | `Apache` | Permisiva con concesión de patentes |
| GPL v3 | `GPL` | Copyleft, obras derivadas deben ser GPL |
| BSD 2-Clause | `BSD` | Permisiva, similar a MIT |
| ISC | `ISC` | Funcionalmente idéntica a MIT |
| MPL 2.0 | `MPL` | Copyleft débil, a nivel de archivo |
| Unlicense | `Unlicense` | Dedicación al dominio público |
| AGPL v3 | `AGPL` | GPL con cláusula de red |
| LGPL v3 | `LGPL` | GPL con excepción de enlazado |
| WTFPL | `WTFPL` | "Do What The Fuck You Want To" |
| EPL 2.0 | `EPL` | Licencia Eclipse |
| Artistic 2.0 | `Artistic` | Licencia Artística (Perl) |

---

## Guía de comandos

### Creación de proyectos

```bash
# Wizard interactivo completo
forgectl new

# Rápido con flags
forgectl new mi-app -t go -l MIT -d "Mi aplicación"
forgectl new mi-app --private --no-remote
forgectl new mi-app --ci github-actions --hooks

# Con archivos de proyecto
forgectl new mi-app --editorconfig --codeowners --security

# Salida JSON para scripting
forgectl new mi-app -t go -l MIT --json
```

`forgectl new` genera automáticamente:

```
mi-app/
├── .git/                  # Repositorio Git inicializado
├── .github/workflows/     # CI (si se solicita)
├── .gitignore             # Según el lenguaje
├── .editorconfig          # Configuración del editor (si se solicita)
├── CODEOWNERS             # Propietarios del código (si se solicita)
├── SECURITY.md            # Política de seguridad (si se solicita)
├── LICENSE                # Licencia elegida
├── README.md              # README profesional
├── cmd/                   # Estructura Go (ejemplo)
├── internal/
├── pkg/
├── src/
└── ...
```

### Gestión de issues

```bash
# Crear issue (interactivo)
forgectl issue new

# Crear issue directo
forgectl issue new "Bug: error en login"
forgectl issue new "Feature: modo oscuro" --body "Descripción detallada..."
forgectl issue new "Task" --labels "bug,help wanted"

# Listar issues del repositorio remoto
forgectl issue list

# Especificar proveedor remoto
forgectl issue new "Bug" --remote gitlab
forgectl issue list --remote gitea

# Salida JSON
forgectl issue list --json
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
# Wizard interactivo
forgectl release new

# Crear release con tag local
forgectl release new v1.0.0 --push
forgectl release new v1.0.0 --changelog

# Crear release en el proveedor remoto
forgectl release new v1.0.0 --remote --push

# Solo crear en el proveedor remoto (sin tag local)
forgectl release new v1.0.0 --remote-only

# Release draft o prerelease
forgectl release new v1.0.0-rc.1 --remote --draft
forgectl release new v1.0.0-rc.1 --remote --prerelease

# Listar releases locales
forgectl release list

# Listar releases del proveedor remoto
forgectl release list --remote

# Salida JSON
forgectl release list --json
forgectl release list --remote --json

# Generar changelog
forgectl changelog gen
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
forgectl config set remotes.github.organization mi-org
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

### Flags globales

Todos los comandos que interactúan con proveedores remotos soportan el flag `--remote` para especificar qué configuración remota utilizar. Muchos comandos también soportan `--json` para salida estructurada:

```bash
# Salida JSON para scripting
forgectl issue list --json | jq '.[].title'
forgectl release list --remote github --json | jq '.[].tag_name'
forgectl new mi-app -t go --json | jq '.path'
```

---

## Configuración

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
    organization: mi-org          # Organización (opcional)
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

### Campo `organization`

El campo `organization` en la configuración remota permite crear repositorios en organizaciones en lugar de en el usuario personal. Cuando se configura, los repositorios se crearán en la organización especificada:

```yaml
remotes:
  github:
    provider: github
    token: ghp_xxxxxxxx
    username: mi_user
    organization: mi-equipo         # Repos se crean en mi-equipo/
    base_url: https://github.com
```

> **Seguridad**: los tokens se guardan en texto plano en tu directorio `~/.forgectl`. Considera protegerlo con `chmod 700 ~/.forgectl`.

---

## Arquitectura

```
forgectl/
├── main.go                          # Entry point
├── cmd/                             # Comandos Cobra
│   ├── root.go                      # Root + registro de comandos
│   ├── new.go                       # Forgectl new (con wizard)
│   ├── config.go                    # Configuración (con wizard)
│   ├── remote.go                    # Gestión de remotes
│   ├── release.go                   # Releases y tags
│   ├── issue.go                     # Gestión de issues
│   ├── doctor.go                    # Diagnóstico
│   └── ...                          # sync, push, pull, status, tag,
│                                    #   changelog, license, template, project
├── internal/
│   ├── config/                      # Viper configuration
│   ├── git/                         # Operaciones Git (go-git + CLI)
│   ├── templates/                   # Plantillas de archivos
│   │   ├── templates.go             # Licencias y README
│   │   ├── scaffold.go              # Estructuras de proyecto y .gitignore
│   │   ├── ci.go                    # Configuración CI/CD
│   │   └── projectfiles.go          # .editorconfig, CODEOWNERS, SECURITY.md
│   └── project/                     # Manager de proyectos + Providers
│       ├── provider.go              # Interfaz Provider
│       ├── github.go                # API GitHub REST
│       ├── gitlab.go                # API GitLab v4
│       └── gitea.go                 # API Gitea v1
└── pkg/
    ├── errors/                      # Errores personalizados
    └── ui/                          # UI kit: prompts, spinner, colores, boxes
```

### Stack tecnológico

- **Go 1.22+** — tiempo de compilación, cero dependencias de runtime
- **[Cobra](https://github.com/spf13/cobra)** — CLI framework
- **[Viper](https://github.com/spf13/viper)** — configuración
- **[go-git](https://github.com/go-git/go-git)** — operaciones Git programáticas
- **[text/template](https://pkg.go.dev/text/template)** — plantillas de proyectos

---

## Extender forgectl

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
    CreateRelease(repoName, tagName, name, body string, draft, prerelease bool) (*Release, error)
    ListReleases(repoName string) ([]*Release, error)
    CreateIssue(repoName, title, body string, labels []string) (*Issue, error)
    ListIssues(repoName string) ([]*Issue, error)
}
```

Regístralo en `NewProvider` en `internal/project/provider.go`.

### Añadir una plantilla de proyecto

En `internal/templates/scaffold.go`, añade al mapa `ScaffoldTemplates`:

```go
"java": {
    Name:        "Java",
    Description: "Standard Java project layout",
    Directories: []string{"src/main/java", "src/test/java"},
},
```

Añade también su `.gitignore` en `GetGitignore`, plantillas de CI en `internal/templates/ci.go`, y `.editorconfig` en `GetEditorConfig` en `internal/templates/projectfiles.go`.

### Añadir una licencia

En `internal/templates/templates.go`, añade al mapa `licenses` de `GetLicense` y a `ListLicenses`.

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

---

## Roadmap

- [x] Wizard interactivo completo
- [x] Multi-proveedor (GitHub, GitLab, Gitea)
- [x] CI/CD generado (Actions, GitLab CI)
- [x] Git hooks automáticos
- [x] Changelog y releases
- [x] `doctor` de diagnóstico
- [x] Issues via API (crear y listar)
- [x] CODEOWNERS, .editorconfig, SECURITY.md
- [x] Soporte de organizaciones para remotes
- [x] Releases via API (GitHub, GitLab, Gitea)
- [x] Modo `--json` para scripting
- [ ] Rebase interactivo y gestión de PRs

---

## Contribuir

Las contribuciones son bienvenidas. Para contribuir:

1. Haz fork del repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Haz commit de tus cambios (`git commit -m 'Añadir nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

Por favor, asegúrate de que tu código pase los tests y el linter antes de enviar el PR:

```bash
make test
make lint
```

---

## Changelog

Los cambios se documentan en [CHANGELOG.md](CHANGELOG.md).

---

## Licencia

[MIT](LICENSE) © [forgectl contributors](https://github.com/forgectl/forgectl/graphs/contributors)
