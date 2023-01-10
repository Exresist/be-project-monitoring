package model

import "time"

type Project struct {
	ID          int
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoURL    string `json:"photo_url"`
	ReportURL   string
	ReportName  string
	RepoURL     string    `json:"repo_url"`
	ActiveTo    time.Time `json:"active_to"`
}
