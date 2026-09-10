package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type gitlabClient struct {
	httpClient *http.Client
	token      string
	username   string
	baseURL    string
}

func newGitLabClient(cfg RemoteConfig) (Provider, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://gitlab.com/api/v4"
	}
	return &gitlabClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      cfg.Token,
		username:   cfg.Username,
		baseURL:    baseURL,
	}, nil
}

func (g *gitlabClient) Name() string { return "gitlab" }

func (g *gitlabClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	u := g.baseURL + path
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
	req.Header.Set("PRIVATE-TOKEN", g.token)
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
		return nil, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

type gitlabRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	Initialize  bool   `json:"initialize_with_readme"`
}

type gitlabRepoResponse struct {
	Name     string `json:"name"`
	HTTPURL  string `json:"http_url_to_repo"`
	WebURL   string `json:"web_url"`
	Private  bool   `json:"private"`
	Visibility string `json:"visibility"`
}

func (g *gitlabClient) CreateRepository(name, description string, private bool) (*Repository, error) {
	visibility := "public"
	if private {
		visibility = "private"
	}

	reqBody := gitlabRepoRequest{
		Name:        name,
		Description: description,
		Visibility:  visibility,
		Initialize:  false,
	}

	data, err := g.doRequest("POST", "/projects", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create repo failed: %w", err)
	}

	var repo gitlabRepoResponse
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}

	return &Repository{
		Name:     repo.Name,
		URL:      repo.WebURL,
		CloneURL: repo.HTTPURL,
		Private:  repo.Visibility == "private",
	}, nil
}

func (g *gitlabClient) DeleteRepository(name string) error {
	encoded := url.PathEscape(g.username + "/" + name)
	_, err := g.doRequest("DELETE", "/projects/"+encoded, nil)
	return err
}

func (g *gitlabClient) GetRepository(name string) (*Repository, error) {
	encoded := url.PathEscape(g.username + "/" + name)
	data, err := g.doRequest("GET", "/projects/"+encoded, nil)
	if err != nil {
		return nil, err
	}
	var repo gitlabRepoResponse
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}
	return &Repository{
		Name:     repo.Name,
		URL:      repo.WebURL,
		CloneURL: repo.HTTPURL,
		Private:  repo.Visibility == "private",
	}, nil
}

func (g *gitlabClient) ListRepositories() ([]*Repository, error) {
	data, err := g.doRequest("GET", "/projects?owned=true&per_page=100&order_by=last_activity_at", nil)
	if err != nil {
		return nil, err
	}
	var repos []gitlabRepoResponse
	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}
	var result []*Repository
	for _, r := range repos {
		result = append(result, &Repository{
			Name:     r.Name,
			URL:      r.WebURL,
			CloneURL: r.HTTPURL,
			Private:  r.Visibility == "private",
		})
	}
	return result, nil
}

func (g *gitlabClient) GetCurrentUsername() (string, error) {
	if g.username != "" {
		return g.username, nil
	}
	data, err := g.doRequest("GET", "/user", nil)
	if err != nil {
		return "", err
	}
	var user struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(data, &user); err != nil {
		return "", err
	}
	return user.Username, nil
}
