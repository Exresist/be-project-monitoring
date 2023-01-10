package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	"context"
)

func (s *service) GetProjects(ctx context.Context, projReq *api.GetProjReq) ([]*model.Project, int, error) {

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

func (s *service) CreateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *service) UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *service) DeleteProject(ctx context.Context, project *model.Project) error {
	return nil
}
