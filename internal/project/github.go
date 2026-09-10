package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type githubClient struct {
	httpClient *http.Client
	token      string
	username   string
	baseURL    string
}

func newGitHubClient(cfg RemoteConfig) (Provider, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &githubClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      cfg.Token,
		username:   cfg.Username,
		baseURL:    baseURL,
	}, nil
}

func (g *githubClient) Name() string { return "github" }

func (g *githubClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	url := g.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

type githubRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

type githubRepoResponse struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
	HTMLURL  string `json:"html_url"`
	Private  bool   `json:"private"`
}

func (g *githubClient) CreateRepository(name, description string, private bool) (*Repository, error) {
	reqBody := githubRepoRequest{
		Name:        name,
		Description: description,
		Private:     private,
		AutoInit:    false,
	}

	data, err := g.doRequest("POST", "/user/repos", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create repo failed: %w", err)
	}

	var repo githubRepoResponse
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}

	return &Repository{
		Name:     repo.Name,
		URL:      repo.HTMLURL,
		CloneURL: repo.CloneURL,
		Private:  repo.Private,
	}, nil
}

func (g *githubClient) DeleteRepository(name string) error {
	_, err := g.doRequest("DELETE", "/repos/"+g.username+"/"+name, nil)
	return err
}

func (g *githubClient) GetRepository(name string) (*Repository, error) {
	data, err := g.doRequest("GET", "/repos/"+g.username+"/"+name, nil)
	if err != nil {
		return nil, err
	}
	var repo githubRepoResponse
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}
	return &Repository{
		Name:     repo.Name,
		URL:      repo.HTMLURL,
		CloneURL: repo.CloneURL,
		Private:  repo.Private,
	}, nil
}

func (g *githubClient) ListRepositories() ([]*Repository, error) {
	data, err := g.doRequest("GET", "/user/repos?per_page=100&sort=updated", nil)
	if err != nil {
		return nil, err
	}
	var repos []githubRepoResponse
	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}
	var result []*Repository
	for _, r := range repos {
		result = append(result, &Repository{
			Name:     r.Name,
			URL:      r.HTMLURL,
			CloneURL: r.CloneURL,
			Private:  r.Private,
		})
	}
	return result, nil
}

func (g *githubClient) GetCurrentUsername() (string, error) {
	if g.username != "" {
		return g.username, nil
	}
	data, err := g.doRequest("GET", "/user", nil)
	if err != nil {
		return "", err
	}
	var user struct {
		Login string `json:"login"`
	}
	if err := json.Unmarshal(data, &user); err != nil {
		return "", err
	}
	return user.Login, nil
}
