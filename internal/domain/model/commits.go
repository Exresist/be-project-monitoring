package model

import "time"

type CommitsInfo struct {
	ShortUser
	TotalCommits    int
	FirstCommitDate time.Time
	LastCommitDate  time.Time
}
