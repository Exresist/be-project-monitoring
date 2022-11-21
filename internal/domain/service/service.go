package service

import "be-project-monitoring/internal/domain"

type (
	userService struct {
		userStore domain.UserStore
	}
	projectService struct {
		projectStore domain.ProjectStore
	}
)

func NewUserService(store domain.UserStore) *userService {
	return &userService{userStore: store}
}

func NewProjectService(store domain.ProjectStore) *projectService {
	return &projectService{projectStore: store}
}
