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
		ID                int
		Name              string
		Description       string
		SuggestedEstimate int
		RealEstimate      int
		ParticipantID     sql.NullInt64
		CreatorID         sql.NullInt64
		Status            TaskStatus
		CreatedAt         time.Time
		UpdatedAt         time.Time
		ProjectID         int
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
