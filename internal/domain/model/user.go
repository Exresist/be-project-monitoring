package model

import (
	"github.com/google/uuid"
)

const (
	Student        UserRole = "STUDENT"
	Admin          UserRole = "ADMIN"
	ProjectManager UserRole = "PROJECT_MANAGER"
)

type (
	UserRole string
	User     struct {
		ShortUser
		HashedPassword string `json:"hashedPassword"`
	}
	ShortUser struct {
		ID             uuid.UUID `json:"id"`
		Role           UserRole  `json:"role"`
		ColorCode      string    `json:"avatarColor"`
		Email          string    `json:"email"`
		Username       string    `json:"username"`
		FirstName      string    `json:"firstName"`
		LastName       string    `json:"lastName"`
		Group          string    `json:"group"`
		GithubUsername string    `json:"ghUsername"`
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
