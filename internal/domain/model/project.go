package model

import "time"

type Project struct {
	ID          int
	Name        string
	Description string
	PhotoURL    string
	ReportURL   string
	ReportName  string
	RepoURL     string
	ActiveTo    time.Time
}
