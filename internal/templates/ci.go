package templates

func GetCITemplate(ciProvider, projectTemplate string) string {
	switch ciProvider {
	case "github-actions":
		return githubActionsTemplate(projectTemplate)
	case "gitlab-ci":
		return gitlabCITemplate(projectTemplate)
	}
	return ""
}

func githubActionsTemplate(projectTemplate string) string {
	switch projectTemplate {
	case "go":
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.21', '1.22']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Build
        run: go build -v ./...

      - name: Test
        run: go test -v -race ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
`
	case "python":
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: ['3.10', '3.11', '3.12']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: ${{ matrix.python-version }}

      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install -r requirements.txt
          pip install pytest pytest-cov

      - name: Lint
        run: |
          pip install ruff
          ruff check .

      - name: Test
        run: pytest --cov=. --cov-report=xml

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.xml
`
	case "node":
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [18, 20, 22]

    steps:
      - uses: actions/checkout@v4

      - name: Use Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node-version }}

      - name: Install dependencies
        run: npm ci

      - name: Lint
        run: npm run lint

      - name: Test
        run: npm test

      - name: Build
        run: npm run build
`
	case "rust":
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

env:
  CARGO_TERM_COLOR: always

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Rust
        uses: dtolnay/rust-toolchain@stable
        with:
          components: clippy, rustfmt

      - name: Format check
        run: cargo fmt --check

      - name: Clippy
        run: cargo clippy -- -D warnings

      - name: Build
        run: cargo build --verbose

      - name: Test
        run: cargo test --verbose
`
	default:
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run tests
        run: echo "Add your test commands here"
`
	}
}

func gitlabCITemplate(projectTemplate string) string {
	switch projectTemplate {
	case "go":
		return `stages:
  - lint
  - test
  - build

variables:
  GO_VERSION: "1.22"

lint:
  stage: lint
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run

test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go test -v -race ./...

build:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - go build -v ./...
  only:
    - main
    - master
`
	case "python":
		return `stages:
  - lint
  - test

lint:
  stage: lint
  image: python:3.12
  script:
    - pip install ruff
    - ruff check .

test:
  stage: test
  image: python:3.12
  script:
    - pip install -r requirements.txt
    - pip install pytest pytest-cov
    - pytest --cov=.
`
	case "node":
		return `stages:
  - lint
  - test
  - build

lint:
  stage: lint
  image: node:20
  script:
    - npm ci
    - npm run lint

test:
  stage: test
  image: node:20
  script:
    - npm ci
    - npm test

build:
  stage: build
  image: node:20
  script:
    - npm ci
    - npm run build
  only:
    - main
    - master
`
	case "rust":
		return `stages:
  - check
  - test

check:
  stage: check
  image: rust:latest
  script:
    - cargo fmt --check
    - cargo clippy -- -D warnings

test:
  stage: test
  image: rust:latest
  script:
    - cargo test --verbose
`
	default:
		return `stages:
  - build

build:
  stage: build
  image: ubuntu:latest
  script:
    - echo "Add your build commands here"
`
	}
}

func GetHook(name string) string {
	switch name {
	case "pre-commit":
		return `#!/bin/sh
# forgectl pre-commit hook
# Runs linting and formatting checks before each commit

set -e

echo "Running pre-commit checks..."

# Check if there are any files to commit
if git diff --cached --name-only | grep -q "."; then
    echo "Staged files found, running checks..."
else
    echo "No staged files found, skipping checks."
    exit 0
fi

# Run go vet if it's a Go project
if [ -f "go.mod" ]; then
    echo "Running go vet..."
    go vet ./...
fi

# Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
    echo "Running golangci-lint..."
    golangci-lint run
fi

# Run cargo fmt check if it's a Rust project
if [ -f "Cargo.toml" ]; then
    echo "Running cargo fmt check..."
    cargo fmt --check
fi

# Run eslint if it's a Node project
if [ -f "package.json" ]; then
    echo "Running npm lint..."
    npm run lint --if-present
fi

echo "Pre-commit checks passed!"
`
	case "commit-msg":
		return `#!/bin/sh
# forgectl commit-msg hook
# Validates commit message format

commit_msg_file=$1
commit_msg=$(cat "$commit_msg_file")

# Check if commit message is empty
if [ -z "$commit_msg" ]; then
    echo "Error: Commit message cannot be empty"
    exit 1
fi

# Check minimum length
if [ ${#commit_msg} -lt 3 ]; then
    echo "Error: Commit message must be at least 3 characters"
    exit 1
fi

echo "Commit message validation passed."
`
	case "pre-push":
		return `#!/bin/sh
# forgectl pre-push hook
# Runs tests before pushing

set -e

echo "Running pre-push checks..."

# Run tests if it's a Go project
if [ -f "go.mod" ]; then
    echo "Running tests..."
    go test ./...
fi

# Run tests if it's a Rust project
if [ -f "Cargo.toml" ]; then
    echo "Running cargo test..."
    cargo test
fi

echo "Pre-push checks passed!"
`
	}
	return ""
}
