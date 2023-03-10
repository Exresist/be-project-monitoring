package model

import "time"

type CommitsInfo struct {
	ShortUser
	TotalCommits       int
	TotalTasksDone     int
	TotalTasksEstimate int
	FirstCommitDate    time.Time
	LastCommitDate     time.Time
}
