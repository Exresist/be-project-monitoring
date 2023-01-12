package service

import (
	"context"
	"errors"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
)

func (s *service) GetProjects(ctx context.Context, projReq *api.GetProjReq) ([]model.Project, int, error) {

	filter := repository.NewProjectFilter().ByProjectNames(projReq.Name)
	filter.Limit = uint64(projReq.Limit)
	filter.Offset = uint64(projReq.Offset)
	count, err := s.repo.GetProjectCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	projects, err := s.repo.GetProjects(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return projects, count, nil
}

func (s *service) CreateProject(ctx context.Context, projectReq *api.CreateProjectReq) (*model.Project, error) {

	project := &model.Project{
		Name:        projectReq.Name,
		Description: projectReq.Description,
		ActiveTo:    projectReq.ActiveTo,
		PhotoURL:    projectReq.PhotoURL,
	}

	found, err := s.repo.GetProject(ctx, repository.NewProjectFilter().
		ByProjectNames(project.Name))
	if err != nil && !errors.Is(err, ierr.ErrProjectNotFound) {
		return nil, err
	}

	if found != nil {
		return nil, ierr.ErrProjectNameAlreadyExists
	}

	return s.repo.InsertProject(ctx, project)
}

func (s *service) UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *service) DeleteProject(ctx context.Context, project *model.Project) error {
	return nil
}
