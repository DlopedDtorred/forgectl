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
	case "java":
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
        java-version: ['17', '21']

    steps:
      - uses: actions/checkout@v4

      - name: Set up JDK ${{ matrix.java-version }}
        uses: actions/setup-java@v4
        with:
          java-version: ${{ matrix.java-version }}
          distribution: temurin

      - name: Setup Gradle
        uses: gradle/actions/setup-gradle@v3

      - name: Build with Gradle
        run: gradle build

      - name: Test with Gradle
        run: gradle test
`
	case "csharp":
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
        dotnet-version: ['6.0.x', '8.0.x']

    steps:
      - uses: actions/checkout@v4

      - name: Setup .NET ${{ matrix.dotnet-version }}
        uses: actions/setup-dotnet@v4
        with:
          dotnet-version: ${{ matrix.dotnet-version }}

      - name: Restore dependencies
        run: dotnet restore

      - name: Build
        run: dotnet build --no-restore --configuration Release

      - name: Test
        run: dotnet test --no-build --configuration Release --verbosity normal
`
	case "typescript":
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

      - name: TypeScript compile
        run: tsc --noEmit

      - name: Test
        run: npm test
`
	case "php":
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
        php-version: ['8.1', '8.2', '8.3']

    steps:
      - uses: actions/checkout@v4

      - name: Set up PHP ${{ matrix.php-version }}
        uses: shivammathur/setup-php@v2
        with:
          php-version: ${{ matrix.php-version }}
          extensions: mbstring, xml, ctype, json, bcmath

      - name: Install dependencies
        run: composer install --no-progress --prefer-dist

      - name: Run tests
        run: vendor/bin/phpunit
`
	case "ruby":
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
        ruby-version: ['3.1', '3.2', '3.3']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Ruby ${{ matrix.ruby-version }}
        uses: ruby/setup-ruby@v1
        with:
          ruby-version: ${{ matrix.ruby-version }}
          bundler-cache: true

      - name: Install dependencies
        run: bundle install

      - name: Run tests
        run: bundle exec rspec
`
	case "swift":
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        swift-version: ['5.9', '5.10']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Swift ${{ matrix.swift-version }}
        uses: swift-actions/setup-swift@v2
        with:
          swift-version: ${{ matrix.swift-version }}

      - name: Build
        run: swift build

      - name: Test
        run: swift test
`
	case "kotlin":
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
        java-version: ['17', '21']

    steps:
      - uses: actions/checkout@v4

      - name: Set up JDK ${{ matrix.java-version }}
        uses: actions/setup-java@v4
        with:
          java-version: ${{ matrix.java-version }}
          distribution: temurin

      - name: Setup Gradle
        uses: gradle/actions/setup-gradle@v3

      - name: Build with Gradle
        run: gradle build

      - name: Test with Gradle
        run: gradle test
`
	case "dart":
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
        dart-version: ['3.3', '3.4']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Dart ${{ matrix.dart-version }}
        uses: dart-lang/setup-dart@v1
        with:
          sdk: ${{ matrix.dart-version }}

      - name: Install dependencies
        run: dart pub get

      - name: Analyze
        run: dart analyze

      - name: Test
        run: dart test

      - name: Compile
        run: dart compile exe bin/main.dart
`
	case "c-cpp":
		return `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        compiler: [gcc-12, gcc-13]

    steps:
      - uses: actions/checkout@v4

      - name: Install ${{ matrix.compiler }}
        run: |
          sudo apt-get update
          sudo apt-get install -y ${{ matrix.compiler }} cmake

      - name: Configure
        run: cmake -B build -DCMAKE_C_COMPILER=${{ matrix.compiler }}

      - name: Build
        run: cmake --build build

      - name: Test
        run: ctest --test-dir build --output-on-failure
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
	case "java":
		return `stages:
  - build
  - test

variables:
  GRADLE_OPTS: "-Dorg.gradle.daemon=false"

build:
  stage: build
  image: gradle:8.5-jdk17
  script:
    - gradle build

test:
  stage: test
  image: gradle:8.5-jdk17
  script:
    - gradle test
`
	case "csharp":
		return `stages:
  - build
  - test

build:
  stage: build
  image: mcr.microsoft.com/dotnet/sdk:8.0
  script:
    - dotnet restore
    - dotnet build --configuration Release

test:
  stage: test
  image: mcr.microsoft.com/dotnet/sdk:8.0
  script:
    - dotnet test --configuration Release --verbosity normal
`
	case "typescript":
		return `stages:
  - lint
  - test
  - build

lint:
  stage: lint
  image: node:20
  script:
    - npm ci
    - npx tsc --noEmit

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
    - npx tsc
  only:
    - main
    - master
`
	case "php":
		return `stages:
  - lint
  - test

lint:
  stage: lint
  image: php:8.3
  script:
    - composer install --no-progress --prefer-dist
    - vendor/bin/phpunit

test:
  stage: test
  image: php:8.3
  script:
    - composer install --no-progress --prefer-dist
    - vendor/bin/phpunit
`
	case "ruby":
		return `stages:
  - lint
  - test

lint:
  stage: lint
  image: ruby:3.3
  script:
    - bundle install
    - bundle exec rspec

test:
  stage: test
  image: ruby:3.3
  script:
    - bundle install
    - bundle exec rspec
`
	case "swift":
		return `stages:
  - build
  - test

build:
  stage: build
  image: swift:5.10
  script:
    - swift build

test:
  stage: test
  image: swift:5.10
  script:
    - swift test
`
	case "kotlin":
		return `stages:
  - build
  - test

variables:
  GRADLE_OPTS: "-Dorg.gradle.daemon=false"

build:
  stage: build
  image: gradle:8.5-jdk17
  script:
    - gradle build

test:
  stage: test
  image: gradle:8.5-jdk17
  script:
    - gradle test
`
	case "dart":
		return `stages:
  - analyze
  - test

analyze:
  stage: analyze
  image: dart:stable
  script:
    - dart pub get
    - dart analyze

test:
  stage: test
  image: dart:stable
  script:
    - dart pub get
    - dart test
`
	case "c-cpp":
		return `stages:
  - build
  - test

build:
  stage: build
  image: gcc:13
  script:
    - apt-get update && apt-get install -y cmake
    - cmake -B build
    - cmake --build build

test:
  stage: test
  image: gcc:13
  script:
    - apt-get update && apt-get install -y cmake
    - cmake -B build
    - cmake --build build
    - ctest --test-dir build --output-on-failure
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

# Run tsc if it's a TypeScript project
if [ -f "tsconfig.json" ]; then
    echo "Running TypeScript type check..."
    npx tsc --noEmit
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
