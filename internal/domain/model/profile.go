package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	Profile struct {
		ID             uuid.UUID `json:"id"`
		ColorCode      string    `json:"color_code"`
		Email          string    `json:"email"`
		Role           string    `json:"role"`
		Username       string    `json:"username"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Group          string    `json:"group"`
		GithubUsername string    `json:"github_username"`
		UserProjects   []UserProjects
	}
	UserProjects struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"photo_url"`
		ActiveTo    time.Time `json:"active_to"`
	}
)
