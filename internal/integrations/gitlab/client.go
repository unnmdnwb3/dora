package gitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
		return nil, errors.New(fmt.Sprintf("Could not find env:  %s", uriID))
	}

	authID := "GITLAB_BEARER"
	auth := os.Getenv(authID)
	if auth == "" {
		return nil, errors.New(fmt.Sprintf("Could not find env:  %s", authID))
	}

	return &Client{
		Uri: uri,
		Auth: auth,
	}, nil
}

type Project struct {
	CreatedAt string `json:"created_at"`
	DefaultBranch string `json:"default_branch"`
	NamespacedName string `json:"path_with_namespace"`
	Url string `json:"web_url"`
}

func (c *Client) GetProjects() (*[]Project, error) {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        log.Fatalln(err)
    }

	projects := []Project{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
        log.Fatalln(err)
    }

	return &projects, nil
}