package prometheus

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/unnmdnwb3/dora/internal/models"
)

// Client represents a Gitlab API client.
type Client struct {
	URI   string
	Auth  string
	Query string
	Start time.Time
	End   time.Time
	Step  string
}

// NewClient creates a new Gitlab API client.
func NewClient(URI string, auth string, query string, start time.Time, end time.Time, step string) *Client {
	return &Client{
		URI:   URI,
		Auth:  auth,
		Query: query,
		Start: start,
		End:   end,
		Step:  step,
	}
}

// GetMonitoringDataPoints gets all monitoring data points.
func (c *Client) GetMonitoringDataPoints() (*[]models.MonitoringDataPoint, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/api/v1/query_range", c.URI)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := fmt.Sprintf("Bearer %s", c.Auth)
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	q.Add("query", c.Query)
	q.Add("start", strconv.FormatInt(c.Start.Unix(), 10))
	q.Add("end", strconv.FormatInt(c.End.Unix(), 10))
	q.Add("step", c.Step)
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

	var queryRange QueryRangeResponse
	err = json.Unmarshal(body, &queryRange)
	if err != nil {
		log.Fatalln(err)
	}

	return c.CreateMonitoringDataPoints(queryRange)
}

// QueryRangeResponse represents a Prometheus query response.
type QueryRangeResponse struct {
	Data struct {
		Result []struct {
			Metric struct {
				Name string `json:"__name__"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// CreateMonitoringDataPoints creates MonitoringDataPoints from a QueryRangeResponse.
func (c *Client) CreateMonitoringDataPoints(queryRangeResponse QueryRangeResponse) (*[]models.MonitoringDataPoint, error) {
	monitoringDataPoints := []models.MonitoringDataPoint{}
	for _, result := range queryRangeResponse.Data.Result {
		for _, dataPoint := range result.Values {
			unix := int64(dataPoint[0].(float64))
			createdAt := time.Unix(unix, 0).UTC()
			value, err := strconv.ParseFloat(dataPoint[1].(string), 64)
			if err != nil {
				return nil, err
			}
			monitoringDataPoints = append(monitoringDataPoints, models.MonitoringDataPoint{
				CreatedAt: createdAt,
				Value:     value,
			})
		}
	}

	return &monitoringDataPoints, nil
}
