package model

import (
	"github.com/google/uuid"
)

const (
	Student        UserRole = "student"
	Admin          UserRole = "admin"
	ProjectManager UserRole = "project_manager"
	A
)

type (
	UserRole string
	User     struct {
		ShortUser
		HashedPassword string `json:"hashed_password"`
	}
	ShortUser struct {
		ID             uuid.UUID `json:"id"`
		Role           UserRole  `json:"role"`
		ColorCode      string    `json:"color_code"`
		Email          string    `json:"email"`
		Username       string    `json:"username"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Group          string    `json:"group"`
		GithubUsername string    `json:"github_username"`
	}
	Profile struct {
		ShortUser
		UserProjects []ShortProject
	}
)

var UserRoles = map[UserRole]struct{}{
	Student:        {},
	Admin:          {},
	ProjectManager: {},
}
