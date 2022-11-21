package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unnmdnwb3/dora/internal/models"
)

// Client represents a Gitlab API client
type Client struct {
	URI  string
	Auth string
}

// NewClient creates a new Gitlab API client
func NewClient() (*Client, error) {
	uriID := "GITLAB_URI"
	uri := os.Getenv(uriID)
	if uri == "" {
		return nil, fmt.Errorf("could not find env:  %s", uriID)
	}

	authID := "GITLAB_BEARER"
	auth := os.Getenv(authID)
	if auth == "" {
		return nil, fmt.Errorf("could not find env:  %s", authID)
	}

	return &Client{
		URI:  uri,
		Auth: auth,
	}, nil
}

// GetRepositories gets all repositories readable with the bearer token provided
func (c *Client) GetRepositories() (*[]models.Repository, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/projects", c.URI)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("owned", "true")
	q.Add("simple", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	repositories := []models.Repository{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		log.Fatalln(err)
	}

	return &repositories, nil
}

// GetOrganisations gets all organisations readable with the bearer token provided
func (c *Client) GetOrganisations() (*[]models.Organisation, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/groups", c.URI)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	organisations := []models.Organisation{}
	err = json.Unmarshal(body, &organisations)
	if err != nil {
		log.Fatalln(err)
	}

	return &organisations, nil
}
