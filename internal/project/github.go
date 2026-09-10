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
	httpClient   *http.Client
	token        string
	username     string
	baseURL      string
	organization string
}

func newGitHubClient(cfg RemoteConfig) (Provider, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &githubClient{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		token:        cfg.Token,
		username:     cfg.Username,
		baseURL:      baseURL,
		organization: cfg.Organization,
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

	path := "/user/repos"
	if g.organization != "" {
		path = "/orgs/" + g.organization + "/repos"
	}

	data, err := g.doRequest("POST", path, reqBody)
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
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}
	_, err := g.doRequest("DELETE", "/repos/"+owner+"/"+name, nil)
	return err
}

func (g *githubClient) GetRepository(name string) (*Repository, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}
	data, err := g.doRequest("GET", "/repos/"+owner+"/"+name, nil)
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
	path := "/user/repos?per_page=100&sort=updated"
	if g.organization != "" {
		path = "/orgs/" + g.organization + "/repos?per_page=100&sort=updated"
	}
	data, err := g.doRequest("GET", path, nil)
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

type githubReleaseRequest struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

type githubReleaseResponse struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	HTMLURL    string `json:"html_url"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	CreatedAt  string `json:"created_at"`
}

func (g *githubClient) CreateRelease(repoName, tagName, name, body string, draft, prerelease bool) (*Release, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	reqBody := githubReleaseRequest{
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

	var rel githubReleaseResponse
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

func (g *githubClient) ListReleases(repoName string) ([]*Release, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	data, err := g.doRequest("GET", "/repos/"+owner+"/"+repoName+"/releases?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var releases []githubReleaseResponse
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

type githubIssueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type githubIssueResponse struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	HTMLURL string `json:"html_url"`
}

func (g *githubClient) CreateIssue(repoName, title, body string, labels []string) (*Issue, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	reqBody := githubIssueRequest{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	data, err := g.doRequest("POST", "/repos/"+owner+"/"+repoName+"/issues", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create issue failed: %w", err)
	}

	var issue githubIssueResponse
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

func (g *githubClient) ListIssues(repoName string) ([]*Issue, error) {
	owner := g.username
	if g.organization != "" {
		owner = g.organization
	}

	data, err := g.doRequest("GET", "/repos/"+owner+"/"+repoName+"/issues?state=open&per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var issues []githubIssueResponse
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
