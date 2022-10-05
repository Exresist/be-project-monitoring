package service

import "be-project-monitoring/internal/domain"

type service struct {
	store domain.UserStore
}

func NewService(store domain.UserStore) *service {
	return &service{store: store}
}
