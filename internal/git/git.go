package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Client struct {
	repo *git.Repository
	dir  string
}

func Init(dir string) (*Client, error) {
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		return nil, fmt.Errorf("git init failed: %w", err)
	}
	return &Client{repo: repo, dir: dir}, nil
}

func Open(dir string) (*Client, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("git open failed: %w", err)
	}
	return &Client{repo: repo, dir: dir}, nil
}

func OpenOrInit(dir string) (*Client, error) {
	client, err := Open(dir)
	if err == nil {
		return client, nil
	}
	return Init(dir)
}

func (c *Client) AddAll() error {
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = c.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitUser() (string, string) {
	name := getGitConfig("user.name")
	if name == "" {
		name = "forgectl"
	}
	email := getGitConfig("user.email")
	if email == "" {
		email = "forgectl@users.noreply.github.com"
	}
	return name, email
}

func getGitConfig(key string) string {
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (c *Client) Commit(message string) error {
	w, err := c.repo.Worktree()
	if err != nil {
		return err
	}

	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return err
	}

	name, email := gitUser()
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	return err
}

func (c *Client) AddAndCommit(message string) error {
	if err := c.AddAll(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}
	if err := c.Commit(message); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	return nil
}

func (c *Client) CreateBranch(name string) error {
	cmd := exec.Command("git", "checkout", "-b", name)
	cmd.Dir = c.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Checkout(name string) error {
	cmd := exec.Command("git", "checkout", name)
	cmd.Dir = c.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) RemoteAdd(name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Dir = c.dir
	return cmd.Run()
}

func (c *Client) RemoteSetURL(name, url string) error {
	cmd := exec.Command("git", "remote", "set-url", name, url)
	cmd.Dir = c.dir
	return cmd.Run()
}

func (c *Client) RemoteRemove(name string) error {
	cmd := exec.Command("git", "remote", "remove", name)
	cmd.Dir = c.dir
	return cmd.Run()
}

func (c *Client) RemoteAddOrSet(name, url string) error {
	if err := c.RemoteAdd(name, url); err != nil {
		return c.RemoteSetURL(name, url)
	}
	return nil
}

func (c *Client) Push(remote, branch string) error {
	cmd := exec.Command("git", "push", "-u", remote, branch)
	cmd.Dir = c.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) PushTags(remote string) error {
	cmd := exec.Command("git", "push", remote, "--tags")
	cmd.Dir = c.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Pull(remote, branch string) error {
	cmd := exec.Command("git", "pull", remote, branch)
	cmd.Dir = c.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Tag(tagName, message string) error {
	head, err := c.repo.Head()
	if err != nil {
		return err
	}
	commitHash := head.Hash()
	name, email := gitUser()
	if message == "" {
		_, err = c.repo.CreateTag(tagName, commitHash, nil)
	} else {
		_, err = c.repo.CreateTag(tagName, commitHash, &git.CreateTagOptions{
			Tagger: &object.Signature{
				Name:  name,
				Email: email,
				When:  time.Now(),
			},
			Message: message,
		})
	}
	return err
}

func (c *Client) Tags() ([]string, error) {
	tags, err := c.repo.Tags()
	if err != nil {
		return nil, err
	}
	var result []string
	tags.ForEach(func(ref *plumbing.Reference) error {
		result = append(result, ref.Name().Short())
		return nil
	})
	return result, nil
}

func (c *Client) IsClean() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = c.dir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(out))) == 0, nil
}

func (c *Client) Log(count int) ([]string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--oneline")
	cmd.Dir = c.dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

func (c *Client) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = c.dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c *Client) GetDir() string {
	return c.dir
}

func (c *Client) GetRemoteURL(remote string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", remote)
	cmd.Dir = c.dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c *Client) GenerateChangelog(fromTag, toTag string) (string, error) {
	var cmd *exec.Cmd
	if fromTag == "" {
		cmd = exec.Command("git", "log", "--pretty=format:- %s (%h)", toTag)
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:- %s (%h)", fromTag+".."+toTag)
	}
	cmd.Dir = c.dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func EnsureGitInstalled() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH: %w", err)
	}
	return nil
}

func FindRepoRoot(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not inside a git repository: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func EnsureRepoRoot(dir string) string {
	root, err := FindRepoRoot(dir)
	if err != nil {
		return dir
	}
	return filepath.Clean(root)
}
