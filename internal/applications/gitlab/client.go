package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Client struct{
	Uri string
	Auth string
}

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
		Uri: uri,
		Auth: auth,
	}, nil
}

type Repository struct {
	Id string `json:"id"`
	CreatedAt string `json:"created_at"`
	DefaultBranch string `json:"default_branch"`
	NamespacedName string `json:"path_with_namespace"`
	Url string `json:"web_url"`
}

func (c *Client) GetRepositories() (*[]Repository, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/projects", c.Uri)
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

	repositories := []Repository{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
        log.Fatalln(err)
    }

	return &repositories, nil
}
