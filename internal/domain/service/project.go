package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *projectService) GetProjects(ctx context.Context, projReq *api.GetProjReq) ([]*model.Project, int, error) {

	filter := repository.NewProjectFilter().ByProjectNames(projReq.Name)
	filter.Limit = uint64(projReq.Limit)
	filter.Offset = uint64(projReq.Offset)

	//ВОПРОС ПРО УЗКОЕ МЕСТО БД - ДЕЛАЕМ 2 ЗАПРОСА
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
	found, err := s.projectStore.GetProject(ctx, repository.NewProjectFilter().
		ByProjectNames(project.Name))
	if err != nil && !errors.Is(err, ierr.ErrProjectNotFound) {
		return nil, err
	}

	if found != nil {
		return nil, ierr.ErrProjectNameAlreadyExists
	}

	projectUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	project.ID = projectUUID
	//еще заполнить репорт юрл и нэйм

	if err = s.projectStore.Insert(ctx, project); err != nil {
		return nil, err
	}

	return project, err
}

func (s *projectService) UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *projectService) DeleteProject(ctx context.Context, project *model.Project) error {
	return nil
}
