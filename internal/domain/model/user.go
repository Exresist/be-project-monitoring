package model

import "github.com/google/uuid"

const (
	Student UserRole = iota
	Admin
	ProjectManager
)

type (
	UserRole int
	User     struct {
		ID             uuid.UUID
		Role           UserRole
		ColorCode      string
		Email          string
		Username       string
		FirstName      string
		LastName       string
		Group          string
		GithubUsername string
		HashedPassword string
	}
)

var RolesToString = map[UserRole]string{
	Student:        "student",
	Admin:          "admin",
	ProjectManager: "project_manager",
}
