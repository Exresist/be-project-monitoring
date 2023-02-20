package model

import (
	"database/sql"
	"time"
)

type (
	Project struct {
		ShortProject
		ReportURL  sql.NullString `json:"reportUrl"`
		ReportName sql.NullString `json:"reportName"`
		RepoURL    sql.NullString `json:"repo"`
	}
	ShortProject struct {
		ID          int            `json:"id"`
		Name        string         `json:"name"`
		Description sql.NullString `json:"description"`
		PhotoURL    sql.NullString `json:"avatar"`
		ActiveTo    time.Time      `json:"dueDate"`
	}
	ProjectInfo struct {
		Project
		Participants []Participant 
		Tasks        []ShortTask
	}
)
