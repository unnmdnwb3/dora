package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/unnmdnwb3/dora/internal/models"
)

// Client represents a Gitlab API client
type Client struct {
	URI  string
	Auth string
}

// NewClient creates a new Gitlab API client
func NewClient(URI string, auth string) *Client {
	return &Client{
		URI:  URI,
		Auth: auth,
	}
}

// GetOrganisations gets all organisations readable with the bearer token provided
func (c *Client) GetOrganisations() (*[]models.Organisation, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/groups", c.URI)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	organisations := []models.Organisation{}
	err = json.Unmarshal(body, &organisations)
	if err != nil {
		return nil, err
	}

	return &organisations, nil
}

// GetRepositories gets all repositories readable with the bearer token provided
func (c *Client) GetRepositories() (*[]models.Repository, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/projects", c.URI)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("owned", "true")
	q.Add("simple", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	repositories := []models.Repository{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		return nil, err
	}

	return &repositories, nil
}

// GetPullRequests gets all pull requests of a repository
func (c *Client) GetPullRequests(projectID int, targetBranch string) (*[]models.PullRequest, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/projects/%s/merge_requests", c.URI, strconv.Itoa(projectID))
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	// TODO add "since" parameter
	q := req.URL.Query()
	q.Add("state", "merged")
	q.Add("target_branch", targetBranch)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pullRequests := []models.PullRequest{}
	err = json.Unmarshal(body, &pullRequests)
	if err != nil {
		return nil, err
	}

	return &pullRequests, nil
}

// GetCommits gets all commits of a repository
func (c *Client) GetCommits(projectID int, referenceBranch string) (*[]models.Commit, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/projects/%s/repository/commits", c.URI, strconv.Itoa(projectID))
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	// TODO add "since" parameter
	q := req.URL.Query()
	q.Add("order", "default") // asc
	q.Add("ref_name", referenceBranch)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	commits := []models.Commit{}
	err = json.Unmarshal(body, &commits)
	if err != nil {
		return nil, err
	}

	return &commits, nil
}

// GetPipelineRuns gets all workflow runs of a project
func (c *Client) GetPipelineRuns(projectID int, referenceBranch string) (*[]models.PipelineRun, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/projects/%s/pipelines", c.URI, strconv.Itoa(projectID))
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("ref", referenceBranch)
	q.Add("sort", "asc")
	q.Add("source", "push")
	q.Add("status", "success")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pipelineRuns []models.PipelineRun
	err = json.Unmarshal(body, &pipelineRuns)
	if err != nil {
		return nil, err
	}

	return &pipelineRuns, nil
}
