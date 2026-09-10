package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type giteaClient struct {
	httpClient *http.Client
	token      string
	username   string
	baseURL    string
}

func newGiteaClient(cfg RemoteConfig) (Provider, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		return nil, fmt.Errorf("gitea requires a base_url in config. Run: forgectl config set remotes.<name>.base_url <url>")
	}
	return &giteaClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      cfg.Token,
		username:   cfg.Username,
		baseURL:    baseURL,
	}, nil
}

func (g *giteaClient) Name() string { return "gitea" }

func (g *giteaClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	u := g.baseURL + "/api/v1" + path
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+g.token)
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
		return nil, fmt.Errorf("Gitea API error (%d): %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

type giteaRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

type giteaRepoResponse struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
	HTMLURL  string `json:"html_url"`
	Private  bool   `json:"private"`
}

func (g *giteaClient) CreateRepository(name, description string, private bool) (*Repository, error) {
	reqBody := giteaRepoRequest{
		Name:        name,
		Description: description,
		Private:     private,
		AutoInit:    false,
	}

	data, err := g.doRequest("POST", "/user/repos", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create repo failed: %w", err)
	}

	var repo giteaRepoResponse
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

func (g *giteaClient) DeleteRepository(name string) error {
	_, err := g.doRequest("DELETE", "/repos/"+g.username+"/"+name, nil)
	return err
}

func (g *giteaClient) GetRepository(name string) (*Repository, error) {
	data, err := g.doRequest("GET", "/repos/"+g.username+"/"+name, nil)
	if err != nil {
		return nil, err
	}
	var repo giteaRepoResponse
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

func (g *giteaClient) ListRepositories() ([]*Repository, error) {
	data, err := g.doRequest("GET", "/user/repos?limit=100&sort=updated", nil)
	if err != nil {
		return nil, err
	}
	var repos []giteaRepoResponse
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

func (g *giteaClient) GetCurrentUsername() (string, error) {
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
