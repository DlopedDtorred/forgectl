package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndCommit(t *testing.T) {
	dir := t.TempDir()

	client, err := Init(dir)
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	// Create a file
	filePath := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("Cannot write file: %v", err)
	}

	// Commit
	if err := client.AddAndCommit("Initial commit"); err != nil {
		t.Fatalf("AddAndCommit returned error: %v", err)
	}

	clean, err := client.IsClean()
	if err != nil {
		t.Fatalf("IsClean returned error: %v", err)
	}
	if !clean {
		t.Error("Expected clean working tree after commit")
	}

	// Open existing repo
	client2, err := Open(dir)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	if client2.GetDir() != dir {
		t.Errorf("Expected dir %s, got %s", dir, client2.GetDir())
	}
}

func TestOpenOrInit(t *testing.T) {
	dir := t.TempDir()

	client, err := OpenOrInit(dir)
	if err != nil {
		t.Fatalf("OpenOrInit returned error: %v", err)
	}
	if client == nil {
		t.Fatal("OpenOrInit returned nil client")
	}
}

func TestTag(t *testing.T) {
	dir := t.TempDir()

	client, err := Init(dir)
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := client.AddAndCommit("commit"); err != nil {
		t.Fatal(err)
	}

	if err := client.Tag("v1.0.0", "Release v1.0.0"); err != nil {
		t.Fatalf("Tag returned error: %v", err)
	}

	tags, err := client.Tags()
	if err != nil {
		t.Fatalf("Tags returned error: %v", err)
	}
	if len(tags) != 1 || tags[0] != "v1.0.0" {
		t.Errorf("Expected tag v1.0.0, got %v", tags)
	}
}

func TestFindRepoRoot(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "sub", "dir")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	client, err := Init(dir)
	if err != nil {
		t.Fatal(err)
	}
	_ = client

	root, err := FindRepoRoot(subDir)
	if err != nil {
		t.Fatalf("FindRepoRoot returned error: %v", err)
	}
	expected, _ := filepath.Abs(dir)
	if root != expected {
		t.Errorf("Expected %s, got %s", expected, root)
	}
}