package model

import "time"

type CommitsInfo struct {
	GithubUsername  string    `json:"githubUsername"`
	Username        string    `json:"username"`
	Total           int       `json:"total"`
	FirstCommitDate time.Time `json:"firstCommitDate"`
	LastCommitDate  time.Time `json:"lastCommitDate"`
}
