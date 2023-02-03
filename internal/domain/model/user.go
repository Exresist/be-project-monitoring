package model

import (
	"github.com/google/uuid"
)

const (
	Student        UserRole = "student"
	Admin          UserRole = "admin"
	ProjectManager UserRole = "project_manager"
)

type (
	UserRole string
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
	ShortUserInfo struct {
		ID             uuid.UUID `json:"user_id"`
		Role           string    `json:"role"`
		ColorCode      string    `json:"color_code"`
		Username       string    `json:"username"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Group          string    `json:"group"`
		GithubUsername string    `json:"github_username"`
	}
)

var UserRoles = map[string]struct{}{
	"student":         {},
	"admin":           {},
	"project_manager": {},
}
