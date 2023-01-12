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
)

var UserRoles = map[string]struct{}{
	"student":         {},
	"admin":           {},
	"project_manager": {},
}
