package model

import (
	"time"
)

type (
	Project struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"photo_url"`
		ReportURL   string    `json:"report_url"`
		ReportName  string    `json:"report_name"`
		RepoURL     string    `json:"repo_url"`
		ActiveTo    time.Time `json:"active_to"`
	}

	ProjectInfo struct {
		Project
		Users []ShortUserInfo
		Tasks []ProjectTask
	}

	ProjectTask struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		ParticipantID int    `json:"participant_id"`
		Status        string `json:"status"`
	}
)
