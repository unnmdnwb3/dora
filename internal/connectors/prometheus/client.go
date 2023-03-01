package prometheus

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/unnmdnwb3/dora/internal/models"
)

// Client represents a Gitlab API client.
type Client struct {
	URI   string
	Auth  string
	Query string
}

// NewClient creates a new Gitlab API client.
func NewClient(URI string, auth string, query string) *Client {
	return &Client{
		URI:   URI,
		Auth:  auth,
		Query: query,
	}
}

// QueryResponse represents a Prometheus query response.
type QueryResponse struct {
	Data struct {
		Result []struct {
			Metric struct {
				Name string `json:"__name__"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// GetAlerts gets all alerts.
func (c *Client) GetAlerts() (*[]models.Alert, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/api/v1/query", c.URI)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("query", c.Query)
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

	var queryResponse QueryResponse
	err = json.Unmarshal(body, &queryResponse)
	if err != nil {
		log.Fatalln(err)
	}

	alerts, err := c.CreateAlerts(queryResponse)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

// CreateAlerts creates Alerts from a QueryResponse.
func (c *Client) CreateAlerts(queryResponse QueryResponse) (*[]models.Alert, error) {
	alerts := []models.Alert{}
	for _, result := range queryResponse.Data.Result {
		for _, dataPoint := range result.Values {
			alerts = append(alerts, models.Alert{CreatedAt: time.Unix(int64(dataPoint[0].(float64)), 0)})
		}
	}

	return &alerts, nil
}
