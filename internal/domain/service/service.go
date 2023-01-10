package service

import "be-project-monitoring/internal/domain"

type service struct {
	repo domain.Repository
}

func NewService(store domain.Repository) *service {
	return &service{repo: store}
}
