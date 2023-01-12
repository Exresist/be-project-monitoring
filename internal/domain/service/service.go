package service

import (
	"be-project-monitoring/internal/domain"

	"github.com/google/go-github/v49/github"
)

type service struct {
	repo     domain.Repository
	githubCl *github.Client
}

func NewService(store domain.Repository, githubCl *github.Client) *service {
	return &service{
		repo:     store,
		githubCl: githubCl,
	}
}
