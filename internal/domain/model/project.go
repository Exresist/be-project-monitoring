package model

import (
	"database/sql"
	"time"
)

type (
	Project struct {
		ShortProject
		ReportURL  sql.NullString `json:"report_url"`
		ReportName sql.NullString `json:"report_name"`
		RepoURL    sql.NullString `json:"repo_url"`
	}
	ShortProject struct {
		ID          int            `json:"id"`
		Name        string         `json:"name"`
		Description sql.NullString `json:"description"`
		PhotoURL    sql.NullString `json:"photo_url"`
		ActiveTo    time.Time      `json:"active_to"`
	}
	ProjectInfo struct {
		Project
		Participants []Participant 
		Tasks        []ShortTask
	}
)
