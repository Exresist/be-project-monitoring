package model

import "time"

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
		Status            TaskStatus
		CreatedAt         time.Time
		UpdatedAt         time.Time
		ParticipantID     int
		CreatorID         int
	}
)

var TaskStatuses = map[string]struct{}{
	"TODO":        {},
	"In progress": {},
	"In review":   {},
	"Testing":     {},
	"Done":        {},
}
