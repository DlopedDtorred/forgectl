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
	httpClient   *http.Client
	token        string
	username     string
	baseURL      string
	organization string
}

func newGiteaClient(cfg RemoteConfig) (Provider, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		return nil, fmt.Errorf("gitea requires a base_url in config. Run: forgectl config set remotes.<name>.base_url <url>")
	}
	return &giteaClient{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		token:        cfg.Token,
		username:     cfg.Username,
		baseURL:      baseURL,
		organization: cfg.Organization,
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

	path := "/user/repos"
	if g.organization != "" {
		path = "/orgs/" + g.organization + "/repos"
	}

	data, err := g.doRequest("POST", path, reqBody)
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
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}
	_, err := g.doRequest("DELETE", "/repos/"+owner+"/"+name, nil)
	return err
}

func (g *giteaClient) GetRepository(name string) (*Repository, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}
	data, err := g.doRequest("GET", "/repos/"+owner+"/"+name, nil)
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
	path := "/user/repos?limit=100&sort=updated"
	if g.organization != "" {
		path = "/orgs/" + g.organization + "/repos?limit=100&sort=updated"
	}
	data, err := g.doRequest("GET", path, nil)
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

type giteaReleaseRequest struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

type giteaReleaseResponse struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	HTMLURL    string `json:"html_url"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	CreatedAt  string `json:"created_at"`
}

func (g *giteaClient) CreateRelease(repoName, tagName, name, body string, draft, prerelease bool) (*Release, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	reqBody := giteaReleaseRequest{
		TagName:    tagName,
		Name:       name,
		Body:       body,
		Draft:      draft,
		Prerelease: prerelease,
	}

	data, err := g.doRequest("POST", "/repos/"+owner+"/"+repoName+"/releases", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create release failed: %w", err)
	}

	var rel giteaReleaseResponse
	if err := json.Unmarshal(data, &rel); err != nil {
		return nil, err
	}

	createdAt, _ := time.Parse(time.RFC3339, rel.CreatedAt)
	return &Release{
		Name:       rel.Name,
		TagName:    rel.TagName,
		Body:       rel.Body,
		Draft:      rel.Draft,
		Prerelease: rel.Prerelease,
		URL:        rel.HTMLURL,
		CreatedAt:  createdAt,
	}, nil
}

func (g *giteaClient) ListReleases(repoName string) ([]*Release, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	data, err := g.doRequest("GET", "/repos/"+owner+"/"+repoName+"/releases?limit=100", nil)
	if err != nil {
		return nil, err
	}

	var releases []giteaReleaseResponse
	if err := json.Unmarshal(data, &releases); err != nil {
		return nil, err
	}

	var result []*Release
	for _, r := range releases {
		createdAt, _ := time.Parse(time.RFC3339, r.CreatedAt)
		result = append(result, &Release{
			Name:       r.Name,
			TagName:    r.TagName,
			Body:       r.Body,
			Draft:      r.Draft,
			Prerelease: r.Prerelease,
			URL:        r.HTMLURL,
			CreatedAt:  createdAt,
		})
	}
	return result, nil
}

type giteaIssueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type giteaIssueResponse struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	HTMLURL string `json:"html_url"`
}

func (g *giteaClient) CreateIssue(repoName, title, body string, labels []string) (*Issue, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	reqBody := giteaIssueRequest{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	data, err := g.doRequest("POST", "/repos/"+owner+"/"+repoName+"/issues", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create issue failed: %w", err)
	}

	var issue giteaIssueResponse
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, err
	}

	return &Issue{
		Number: issue.Number,
		Title:  issue.Title,
		Body:   issue.Body,
		URL:    issue.HTMLURL,
		Labels: labels,
	}, nil
}

func (g *giteaClient) ListIssues(repoName string) ([]*Issue, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	data, err := g.doRequest("GET", "/repos/"+owner+"/"+repoName+"/issues?state=open&limit=100", nil)
	if err != nil {
		return nil, err
	}

	var issues []giteaIssueResponse
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, err
	}

	var result []*Issue
	for _, i := range issues {
		result = append(result, &Issue{
			Number: i.Number,
			Title:  i.Title,
			Body:   i.Body,
			URL:    i.HTMLURL,
		})
	}
	return result, nil
}
