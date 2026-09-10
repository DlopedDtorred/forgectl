package templates

import (
	"strings"
	"testing"
)

func TestGenerateReadme(t *testing.T) {
	data := ReadmeData{
		Name:        "test-project",
		Description: "A test project",
		License:     "MIT",
		Year:        2026,
		Author:      "Test Author",
	}

	content, err := GenerateReadme(data)
	if err != nil {
		t.Fatalf("GenerateReadme returned error: %v", err)
	}

	required := []string{
		"# test-project",
		"A test project",
		"MIT License",
		"Test Author",
		"2026",
		"forgectl",
	}

	for _, r := range required {
		if !strings.Contains(content, r) {
			t.Errorf("README content missing %q", r)
		}
	}
}

func TestGenerateReadmeEmptyLicense(t *testing.T) {
	data := ReadmeData{
		Name:    "test-project",
		License: "none",
	}

	content, err := GenerateReadme(data)
	if err != nil {
		t.Fatalf("GenerateReadme returned error: %v", err)
	}

	if !strings.Contains(content, "No license specified") {
		t.Errorf("Expected 'No license specified' for empty license, got: %s", content)
	}
}

func TestGenerateReadmeNoDescription(t *testing.T) {
	data := ReadmeData{
		Name: "my-app",
	}

	content, err := GenerateReadme(data)
	if err != nil {
		t.Fatalf("GenerateReadme returned error: %v", err)
	}

	if !strings.Contains(content, "A project managed with forgectl") {
		t.Error("Expected default description when none provided")
	}
}

func TestGenerateReadmeNoAuthor(t *testing.T) {
	data := ReadmeData{
		Name: "my-app",
	}

	content, err := GenerateReadme(data)
	if err != nil {
		t.Fatalf("GenerateReadme returned error: %v", err)
	}

	if !strings.Contains(content, "Generated with") {
		t.Error("Expected generated author attribution when no author provided")
	}
}

func TestGetLicense(t *testing.T) {
	licenses := ListLicenses()
	for _, name := range licenses {
		content, err := GetLicense(name, 2026)
		if err != nil {
			t.Errorf("GetLicense(%s) returned error: %v", name, err)
		}
		if content == "" {
			t.Errorf("GetLicense(%s) returned empty content", name)
		}
	}

	if _, err := GetLicense("NONEXISTENT", 2026); err == nil {
		t.Error("Expected error for unknown license")
	}
}

func TestLicenseAGPL(t *testing.T) {
	content, err := GetLicense("AGPL", 2026)
	if err != nil {
		t.Fatalf("GetLicense(AGPL) returned error: %v", err)
	}
	for _, want := range []string{"GNU AFFERO GENERAL PUBLIC LICENSE", "2026", "Version 3"} {
		if !strings.Contains(content, want) {
			t.Errorf("AGPL license missing %q", want)
		}
	}
}

func TestLicenseLGPL(t *testing.T) {
	content, err := GetLicense("LGPL", 2025)
	if err != nil {
		t.Fatalf("GetLicense(LGPL) returned error: %v", err)
	}
	for _, want := range []string{"GNU LESSER GENERAL PUBLIC LICENSE", "2025", "Version 3"} {
		if !strings.Contains(content, want) {
			t.Errorf("LGPL license missing %q", want)
		}
	}
}

func TestLicenseWTFPL(t *testing.T) {
	content, err := GetLicense("WTFPL", 2026)
	if err != nil {
		t.Fatalf("GetLicense(WTFPL) returned error: %v", err)
	}
	for _, want := range []string{"DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE", "Sam Hocevar"} {
		if !strings.Contains(content, want) {
			t.Errorf("WTFPL license missing %q", want)
		}
	}
}

func TestLicenseEPL(t *testing.T) {
	content, err := GetLicense("EPL", 2026)
	if err != nil {
		t.Fatalf("GetLicense(EPL) returned error: %v", err)
	}
	for _, want := range []string{"Eclipse Public License", "2026", "EPL-2.0"} {
		if !strings.Contains(content, want) {
			t.Errorf("EPL license missing %q", want)
		}
	}
}

func TestLicenseArtistic(t *testing.T) {
	content, err := GetLicense("Artistic", 2026)
	if err != nil {
		t.Fatalf("GetLicense(Artistic) returned error: %v", err)
	}
	for _, want := range []string{"The Artistic License 2.0", "2026", "General Public License"} {
		if !strings.Contains(content, want) {
			t.Errorf("Artistic license missing %q", want)
		}
	}
}

func TestLicenseMIT(t *testing.T) {
	content, err := GetLicense("MIT", 2024)
	if err != nil {
		t.Fatalf("GetLicense(MIT) returned error: %v", err)
	}
	for _, want := range []string{"MIT License", "2024", "Permission is hereby granted"} {
		if !strings.Contains(content, want) {
			t.Errorf("MIT license missing %q", want)
		}
	}
}

func TestLicenseApache(t *testing.T) {
	content, err := GetLicense("Apache", 2023)
	if err != nil {
		t.Fatalf("GetLicense(Apache) returned error: %v", err)
	}
	for _, want := range []string{"Apache License", "2023", "Version 2.0"} {
		if !strings.Contains(content, want) {
			t.Errorf("Apache license missing %q", want)
		}
	}
}

func TestLicenseGPL(t *testing.T) {
	content, err := GetLicense("GPL", 2022)
	if err != nil {
		t.Fatalf("GetLicense(GPL) returned error: %v", err)
	}
	for _, want := range []string{"GNU GENERAL PUBLIC LICENSE", "2022", "Version 3"} {
		if !strings.Contains(content, want) {
			t.Errorf("GPL license missing %q", want)
		}
	}
}

func TestLicenseBSD(t *testing.T) {
	content, err := GetLicense("BSD", 2021)
	if err != nil {
		t.Fatalf("GetLicense(BSD) returned error: %v", err)
	}
	for _, want := range []string{"BSD 2-Clause License", "2021"} {
		if !strings.Contains(content, want) {
			t.Errorf("BSD license missing %q", want)
		}
	}
}

func TestLicenseISC(t *testing.T) {
	content, err := GetLicense("ISC", 2020)
	if err != nil {
		t.Fatalf("GetLicense(ISC) returned error: %v", err)
	}
	for _, want := range []string{"ISC License", "2020"} {
		if !strings.Contains(content, want) {
			t.Errorf("ISC license missing %q", want)
		}
	}
}

func TestLicenseMPL(t *testing.T) {
	content, err := GetLicense("MPL", 2019)
	if err != nil {
		t.Fatalf("GetLicense(MPL) returned error: %v", err)
	}
	for _, want := range []string{"Mozilla Public License Version 2.0", "2019"} {
		if !strings.Contains(content, want) {
			t.Errorf("MPL license missing %q", want)
		}
	}
}

func TestLicenseUnlicense(t *testing.T) {
	content, err := GetLicense("Unlicense", 2026)
	if err != nil {
		t.Fatalf("GetLicense(Unlicense) returned error: %v", err)
	}
	for _, want := range []string{"public domain", "unencumbered software"} {
		if !strings.Contains(content, want) {
			t.Errorf("Unlicense missing %q", want)
		}
	}
}

func TestListLicensesCount(t *testing.T) {
	licenses := ListLicenses()
	if len(licenses) != 12 {
		t.Errorf("ListLicenses() returned %d licenses, want 12", len(licenses))
	}
}

func TestScaffoldTemplates(t *testing.T) {
	templates := []string{"go", "python", "node", "rust", "minimal"}
	for _, name := range templates {
		tmpl, ok := GetScaffoldTemplate(name)
		if !ok {
			t.Errorf("GetScaffoldTemplate(%s) returned not ok", name)
			continue
		}
		if tmpl.Name == "" {
			t.Errorf("Template %s has empty name", name)
		}
		if len(tmpl.Directories) == 0 {
			t.Errorf("Template %s has no directories", name)
		}
	}
}

func TestScaffoldTemplatesNewLanguages(t *testing.T) {
	newTemplates := []string{"java", "csharp", "typescript", "php", "ruby", "swift", "kotlin", "dart", "c-cpp"}
	for _, name := range newTemplates {
		tmpl, ok := GetScaffoldTemplate(name)
		if !ok {
			t.Errorf("GetScaffoldTemplate(%s) returned not ok", name)
			continue
		}
		if tmpl.Name == "" {
			t.Errorf("Template %s has empty name", name)
		}
		if len(tmpl.Directories) == 0 {
			t.Errorf("Template %s has no directories", name)
		}
		if tmpl.Description == "" {
			t.Errorf("Template %s has empty description", name)
		}
	}
}

func TestScaffoldTemplateNotFound(t *testing.T) {
	_, ok := GetScaffoldTemplate("nonexistent")
	if ok {
		t.Error("Expected not ok for nonexistent template")
	}
}

func TestListScaffoldTemplates(t *testing.T) {
	names := ListScaffoldTemplates()
	if len(names) < 5 {
		t.Errorf("ListScaffoldTemplates() returned %d templates, want at least 5", len(names))
	}
}

func TestGetGitignore(t *testing.T) {
	for _, tmpl := range []string{"go", "python", "node", "rust", "minimal"} {
		content := GetGitignore(tmpl)
		if content == "" {
			t.Errorf("GetGitignore(%s) returned empty", tmpl)
		}
	}
}

func TestGetGitignoreNewLanguages(t *testing.T) {
	for _, tmpl := range []string{"java", "csharp", "typescript", "php", "ruby", "swift", "kotlin", "dart", "c-cpp"} {
		content := GetGitignore(tmpl)
		if content == "" {
			t.Errorf("GetGitignore(%s) returned empty", tmpl)
		}
	}
}

func TestGetGitignoreUnknownReturnsMinimal(t *testing.T) {
	content := GetGitignore("unknown")
	minimal := GetGitignore("minimal")
	if content != minimal {
		t.Error("GetGitignore(unknown) should return minimal gitignore")
	}
}

func TestGetGitignoreJava(t *testing.T) {
	content := GetGitignore("java")
	for _, want := range []string{"*.class", "*.jar", "build/", ".gradle/"} {
		if !strings.Contains(content, want) {
			t.Errorf("Java gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreCSharp(t *testing.T) {
	content := GetGitignore("csharp")
	for _, want := range []string{"bin/", "obj/", "*.dll", "*.exe"} {
		if !strings.Contains(content, want) {
			t.Errorf("C# gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreTypeScript(t *testing.T) {
	content := GetGitignore("typescript")
	for _, want := range []string{"node_modules/", "dist/", "build/"} {
		if !strings.Contains(content, want) {
			t.Errorf("TypeScript gitignore missing %q", want)
		}
	}
}

func TestGetGitignorePHP(t *testing.T) {
	content := GetGitignore("php")
	for _, want := range []string{"vendor/", "composer.lock", "*.phar"} {
		if !strings.Contains(content, want) {
			t.Errorf("PHP gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreRuby(t *testing.T) {
	content := GetGitignore("ruby")
	for _, want := range []string{"*.gem", ".bundle/", "tmp/"} {
		if !strings.Contains(content, want) {
			t.Errorf("Ruby gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreSwift(t *testing.T) {
	content := GetGitignore("swift")
	for _, want := range []string{".build/", ".swiftpm/", "*.xcuserdata"} {
		if !strings.Contains(content, want) {
			t.Errorf("Swift gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreKotlin(t *testing.T) {
	content := GetGitignore("kotlin")
	for _, want := range []string{"build/", ".gradle/", "*.class"} {
		if !strings.Contains(content, want) {
			t.Errorf("Kotlin gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreDart(t *testing.T) {
	content := GetGitignore("dart")
	for _, want := range []string{".dart_tool/", "pubspec.lock", "*.g.dart"} {
		if !strings.Contains(content, want) {
			t.Errorf("Dart gitignore missing %q", want)
		}
	}
}

func TestGetGitignoreCCpp(t *testing.T) {
	content := GetGitignore("c-cpp")
	for _, want := range []string{"*.o", "*.so", "*.dylib", "CMakeFiles/"} {
		if !strings.Contains(content, want) {
			t.Errorf("C/C++ gitignore missing %q", want)
		}
	}
}

func TestGetCITemplate(t *testing.T) {
	for _, ci := range []string{"github-actions", "gitlab-ci"} {
		for _, tmpl := range []string{"go", "python", "node", "rust", "unknown"} {
			content := GetCITemplate(ci, tmpl)
			if content == "" {
				t.Errorf("GetCITemplate(%s, %s) returned empty", ci, tmpl)
			}
		}
	}

	if content := GetCITemplate("unknown", "go"); content != "" {
		t.Errorf("GetCITemplate(unknown) should be empty, got %s", content)
	}
}

func TestGetCIGitHubActionsNewLanguages(t *testing.T) {
	newLangs := []string{"java", "csharp", "typescript", "php", "ruby", "swift", "kotlin", "dart", "c-cpp"}
	for _, lang := range newLangs {
		content := GetCITemplate("github-actions", lang)
		if content == "" {
			t.Errorf("GetCITemplate(github-actions, %s) returned empty", lang)
		}
		if !strings.Contains(content, "name: CI") {
			t.Errorf("GitHub Actions %s CI missing 'name: CI'", lang)
		}
		if !strings.Contains(content, "actions/checkout") {
			t.Errorf("GitHub Actions %s CI missing actions/checkout", lang)
		}
	}
}

func TestGetCIGitLabCINewLanguages(t *testing.T) {
	newLangs := []string{"java", "csharp", "typescript", "php", "ruby", "swift", "kotlin", "dart", "c-cpp"}
	for _, lang := range newLangs {
		content := GetCITemplate("gitlab-ci", lang)
		if content == "" {
			t.Errorf("GetCITemplate(gitlab-ci, %s) returned empty", lang)
		}
		if !strings.Contains(content, "stages:") {
			t.Errorf("GitLab CI %s missing 'stages:'", lang)
		}
	}
}

func TestGetCIGitHubActionsJava(t *testing.T) {
	content := GetCITemplate("github-actions", "java")
	for _, want := range []string{"setup-java", "temurin", "gradle"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions Java CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsCSharp(t *testing.T) {
	content := GetCITemplate("github-actions", "csharp")
	for _, want := range []string{"setup-dotnet", "dotnet restore", "dotnet build"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions C# CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsTypeScript(t *testing.T) {
	content := GetCITemplate("github-actions", "typescript")
	for _, want := range []string{"setup-node", "tsc --noEmit"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions TypeScript CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsPHP(t *testing.T) {
	content := GetCITemplate("github-actions", "php")
	for _, want := range []string{"setup-php", "phpunit", "composer"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions PHP CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsRuby(t *testing.T) {
	content := GetCITemplate("github-actions", "ruby")
	for _, want := range []string{"setup-ruby", "rspec", "bundle"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions Ruby CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsSwift(t *testing.T) {
	content := GetCITemplate("github-actions", "swift")
	for _, want := range []string{"setup-swift", "swift build", "swift test"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions Swift CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsKotlin(t *testing.T) {
	content := GetCITemplate("github-actions", "kotlin")
	for _, want := range []string{"setup-java", "temurin", "gradle"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions Kotlin CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsDart(t *testing.T) {
	content := GetCITemplate("github-actions", "dart")
	for _, want := range []string{"setup-dart", "dart pub get", "dart test"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions Dart CI missing %q", want)
		}
	}
}

func TestGetCIGitHubActionsCCpp(t *testing.T) {
	content := GetCITemplate("github-actions", "c-cpp")
	for _, want := range []string{"cmake", "gcc"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitHub Actions C/C++ CI missing %q", want)
		}
	}
}

func TestGetCIGitLabCIJava(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "java")
	for _, want := range []string{"gradle", "jdk17"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI Java missing %q", want)
		}
	}
}

func TestGetCIGitLabCICSharp(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "csharp")
	for _, want := range []string{"dotnet", "sdk"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI C# missing %q", want)
		}
	}
}

func TestGetCIGitLabCITypeScript(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "typescript")
	for _, want := range []string{"node", "tsc"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI TypeScript missing %q", want)
		}
	}
}

func TestGetCIGitLabCIPHP(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "php")
	for _, want := range []string{"php", "phpunit", "composer"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI PHP missing %q", want)
		}
	}
}

func TestGetCIGitLabCIRuby(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "ruby")
	for _, want := range []string{"ruby", "rspec", "bundle"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI Ruby missing %q", want)
		}
	}
}

func TestGetCIGitLabCISwift(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "swift")
	for _, want := range []string{"swift", "swift build", "swift test"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI Swift missing %q", want)
		}
	}
}

func TestGetCIGitLabCIKotlin(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "kotlin")
	for _, want := range []string{"gradle", "jdk17"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI Kotlin missing %q", want)
		}
	}
}

func TestGetCIGitLabCIDart(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "dart")
	for _, want := range []string{"dart", "dart pub get", "dart test"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI Dart missing %q", want)
		}
	}
}

func TestGetCIGitLabCICCpp(t *testing.T) {
	content := GetCITemplate("gitlab-ci", "c-cpp")
	for _, want := range []string{"gcc", "cmake"} {
		if !strings.Contains(content, want) {
			t.Errorf("GitLab CI C/C++ missing %q", want)
		}
	}
}

func TestGetHook(t *testing.T) {
	for _, hook := range []string{"pre-commit", "pre-push", "commit-msg"} {
		content := GetHook(hook)
		if content == "" {
			t.Errorf("GetHook(%s) returned empty", hook)
		}
	}
}

func TestGetHookUnknown(t *testing.T) {
	content := GetHook("unknown-hook")
	if content != "" {
		t.Errorf("GetHook(unknown) should be empty, got content")
	}
}

func TestGetHookPreCommit(t *testing.T) {
	content := GetHook("pre-commit")
	for _, want := range []string{"#!/bin/sh", "pre-commit checks", "go vet", "golangci-lint", "cargo fmt", "npm run lint"} {
		if !strings.Contains(content, want) {
			t.Errorf("pre-commit hook missing %q", want)
		}
	}
}

func TestGetHookCommitMsg(t *testing.T) {
	content := GetHook("commit-msg")
	for _, want := range []string{"#!/bin/sh", "commit_msg_file", "minimum length"} {
		if !strings.Contains(content, want) {
			t.Errorf("commit-msg hook missing %q", want)
		}
	}
}

func TestGetHookPrePush(t *testing.T) {
	content := GetHook("pre-push")
	for _, want := range []string{"#!/bin/sh", "pre-push checks", "go test", "cargo test"} {
		if !strings.Contains(content, want) {
			t.Errorf("pre-push hook missing %q", want)
		}
	}
}

func TestGetEditorConfig(t *testing.T) {
	templates := []string{"go", "python", "node", "rust", "java", "csharp", "typescript", "php", "ruby", "swift", "kotlin", "dart", "c-cpp", "minimal"}
	for _, tmpl := range templates {
		content := GetEditorConfig(tmpl)
		if content == "" {
			t.Errorf("GetEditorConfig(%s) returned empty", tmpl)
			continue
		}
		if !strings.Contains(content, "root = true") {
			t.Errorf("GetEditorConfig(%s) missing 'root = true'", tmpl)
		}
		if !strings.Contains(content, "charset = utf-8") {
			t.Errorf("GetEditorConfig(%s) missing 'charset = utf-8'", tmpl)
		}
		if !strings.Contains(content, "end_of_line = lf") {
			t.Errorf("GetEditorConfig(%s) missing 'end_of_line = lf'", tmpl)
		}
		if !strings.Contains(content, "insert_final_newline = true") {
			t.Errorf("GetEditorConfig(%s) missing 'insert_final_newline = true'", tmpl)
		}
	}
}

func TestGetEditorConfigGo(t *testing.T) {
	content := GetEditorConfig("go")
	if !strings.Contains(content, "indent_style = tab") {
		t.Error("Go editorconfig missing tab indent style")
	}
}

func TestGetEditorConfigPython(t *testing.T) {
	content := GetEditorConfig("python")
	if !strings.Contains(content, "*.py") {
		t.Error("Python editorconfig missing *.py section")
	}
}

func TestGetEditorConfigNode(t *testing.T) {
	content := GetEditorConfig("node")
	if !strings.Contains(content, "*.js") {
		t.Error("Node editorconfig missing *.js section")
	}
}

func TestGetEditorConfigTypeScript(t *testing.T) {
	content := GetEditorConfig("typescript")
	if !strings.Contains(content, "*.ts") {
		t.Error("TypeScript editorconfig missing *.ts section")
	}
}

func TestGetEditorConfigJava(t *testing.T) {
	content := GetEditorConfig("java")
	if !strings.Contains(content, "*.java") {
		t.Error("Java editorconfig missing *.java section")
	}
}

func TestGetEditorConfigCSharp(t *testing.T) {
	content := GetEditorConfig("csharp")
	if !strings.Contains(content, "*.cs") {
		t.Error("C# editorconfig missing *.cs section")
	}
}

func TestGetEditorConfigPHP(t *testing.T) {
	content := GetEditorConfig("php")
	if !strings.Contains(content, "*.php") {
		t.Error("PHP editorconfig missing *.php section")
	}
}

func TestGetEditorConfigRuby(t *testing.T) {
	content := GetEditorConfig("ruby")
	if !strings.Contains(content, "*.rb") {
		t.Error("Ruby editorconfig missing *.rb section")
	}
}

func TestGetEditorConfigSwift(t *testing.T) {
	content := GetEditorConfig("swift")
	if !strings.Contains(content, "*.swift") {
		t.Error("Swift editorconfig missing *.swift section")
	}
}

func TestGetEditorConfigKotlin(t *testing.T) {
	content := GetEditorConfig("kotlin")
	if !strings.Contains(content, "*.kt") {
		t.Error("Kotlin editorconfig missing *.kt section")
	}
}

func TestGetEditorConfigDart(t *testing.T) {
	content := GetEditorConfig("dart")
	if !strings.Contains(content, "*.dart") {
		t.Error("Dart editorconfig missing *.dart section")
	}
}

func TestGetEditorConfigCCpp(t *testing.T) {
	content := GetEditorConfig("c-cpp")
	if !strings.Contains(content, "*.c") || !strings.Contains(content, "*.cpp") {
		t.Error("C/C++ editorconfig missing *.c or *.cpp section")
	}
}

func TestGetEditorConfigDefault(t *testing.T) {
	content := GetEditorConfig("unknown")
	if !strings.Contains(content, "root = true") {
		t.Error("Default editorconfig missing 'root = true'")
	}
}

func TestGetCodeowners(t *testing.T) {
	content := GetCodeowners("")
	if !strings.Contains(content, "* @owner") {
		t.Error("Default CODEOWNERS missing '* @owner'")
	}
	if !strings.Contains(content, "Code Owners") {
		t.Error("CODEOWNERS missing header")
	}
}

func TestGetCodeownersCustom(t *testing.T) {
	content := GetCodeowners("@team/core")
	if !strings.Contains(content, "@team/core") {
		t.Error("Custom CODEOWNERS missing custom content")
	}
}

func TestGetSecurityMD(t *testing.T) {
	content := GetSecurityMD("my-project", "")
	if !strings.Contains(content, "Security Policy") {
		t.Error("SECURITY.md missing header")
	}
	if !strings.Contains(content, "my-project") {
		t.Error("SECURITY.md missing project name")
	}
	if !strings.Contains(content, "security@example.com") {
		t.Error("SECURITY.md missing default contact")
	}
}

func TestGetSecurityMDWithContact(t *testing.T) {
	content := GetSecurityMD("my-project", "admin@example.com")
	if !strings.Contains(content, "admin@example.com") {
		t.Error("SECURITY.md missing custom contact")
	}
	if !strings.Contains(content, "my-project") {
		t.Error("SECURITY.md missing project name")
	}
}

func TestGetSecurityMDContainsResponseProcess(t *testing.T) {
	content := GetSecurityMD("test", "")
	for _, want := range []string{"Acknowledgment", "Assessment", "Resolution"} {
		if !strings.Contains(content, want) {
			t.Errorf("SECURITY.md missing %q in response process", want)
		}
	}
}
