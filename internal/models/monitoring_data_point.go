package models

import "time"

// MonitoringDataPoint describes a single data point in a monitoring time series.
type MonitoringDataPoint struct {
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
