package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	"context"
)

func (s *projectService) GetProjects(ctx context.Context, projReq *api.GetProjReq) ([]*model.Project, int, error) {
	
	filter := repository.NewProjectFilter().ByProjectNames(projReq.Name)
	filter.Limit = uint64(projReq.Limit)
	filter.Offset = uint64(projReq.Offset)
	count, err := s.projectStore.GetCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	projects, err := s.projectStore.GetProjects(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return projects, count, nil
}

func (s *projectService) CreateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *projectService) UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *projectService) DeleteProject(ctx context.Context, project *model.Project) error {
	return nil
}
