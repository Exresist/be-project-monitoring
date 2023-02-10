package model

import (
	"database/sql"
	"time"
)

const (
	TODO       TaskStatus = "TODO"
	InProgress TaskStatus = "In progress"
	InReview   TaskStatus = "In review"
	Testing    TaskStatus = "Testing"
	Done       TaskStatus = "Done"
)

type (
	TaskStatus string
	Task       struct {
		ShortTask
		SuggestedEstimate sql.NullInt64 `json:"suggested_estimate"`
		RealEstimate      sql.NullInt64 `json:"real_estimate"`
		CreatorID         sql.NullInt64 `json:"creator_id"`
		CreatedAt         time.Time     `json:"created_at"`
		UpdatedAt         time.Time     `json:"updated_at"`
		ProjectID         int           `json:"project_id"`
	}
	ShortTask struct {
		ID            int            `json:"id"`
		Name          string         `json:"name"`
		Description   sql.NullString `json:"description"`
		ParticipantID sql.NullInt64  `json:"participant_id"`
		Status        TaskStatus     `json:"status"`
	}
	TaskInfo struct {
		Task
		Creator     ShortUser
		Participant ShortUser
	}
)

var TaskStatuses = map[string]struct{}{
	"TODO":        {},
	"In progress": {},
	"In review":   {},
	"Testing":     {},
	"Done":        {},
}
