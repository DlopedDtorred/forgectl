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

func TestGetGitignore(t *testing.T) {
	for _, tmpl := range []string{"go", "python", "node", "rust", "minimal"} {
		content := GetGitignore(tmpl)
		if content == "" {
			t.Errorf("GetGitignore(%s) returned empty", tmpl)
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

func TestGetHook(t *testing.T) {
	for _, hook := range []string{"pre-commit", "pre-push", "commit-msg"} {
		content := GetHook(hook)
		if content == "" {
			t.Errorf("GetHook(%s) returned empty", hook)
		}
	}
}