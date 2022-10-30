package service

import "be-project-monitoring/internal/domain"

type userService struct {
	userStore domain.UserStore
}

func NewService(store domain.UserStore) *userService {
	return &userService{userStore: store}
}
