package model

type CommitsInfo struct {
	ShortUser
	TotalCommits       int
	TotalTasksDone     int
	TotalTasksEstimate int
	NumberOfAdditions  int
	NumberOfDeletions  int
}
