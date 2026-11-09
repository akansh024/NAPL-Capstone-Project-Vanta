package models

import "time"

type Finding struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Asset       string    `json:"asset"`
	Evidence    string    `json:"evidence"`
	Remediation string    `json:"remediation"`
	Source      string    `json:"source"`
	LastSeen    time.Time `json:"lastSeen"`
}
type Import struct {
	ID               int64     `json:"id"`
	Filename         string    `json:"filename"`
	Scanner          string    `json:"scanner"`
	Status           string    `json:"status"`
	AssetCount       int       `json:"assetCount"`
	ObservationCount int       `json:"observationCount"`
	ImportedAt       time.Time `json:"importedAt"`
}

type Observation struct {
	ID         int64     `json:"id"`
	Asset      string    `json:"asset"`
	Hostname   string    `json:"hostname"`
	Port       int       `json:"port"`
	Protocol   string    `json:"protocol"`
	Service    string    `json:"service"`
	State      string    `json:"state"`
	Evidence   string    `json:"evidence"`
	Filename   string    `json:"filename"`
	ObservedAt time.Time `json:"observedAt"`
}
