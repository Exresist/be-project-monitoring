package service

import (
	"context"
	"errors"
	"time"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
)

func (s *service) GetProjects(ctx context.Context, projReq *api.GetProjectReq) ([]model.Project, int, error) {

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

	found, err := s.repo.GetProject(ctx, repository.NewProjectFilter().
		ByProjectNames(projectReq.Name))
	if err != nil && !errors.Is(err, ierr.ErrProjectNotFound) {
		return nil, err
	}

	if found != nil {
		return nil, ierr.ErrProjectNameAlreadyExists
	}
	//nado li kakie-nibud 10 min visrat?
	if projectReq.ActiveTo.Before(time.Now()) {
		return nil, ierr.ErrProjectDateIsNotValid
	}
	project := &model.Project{
		Name:        projectReq.Name,
		Description: projectReq.Description,
		ActiveTo:    projectReq.ActiveTo,
		PhotoURL:    projectReq.PhotoURL,
	}

	if err := s.repo.InsertProject(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *service) UpdateProject(ctx context.Context, projectReq *api.UpdateProjectReq) (*model.Project, error) {

	//nado li kakie-nibud minus 10 min visrat?
	// if projectReq.ActiveTo != nil && projectReq.ActiveTo.Before(time.Now()) {
	// 	return nil, ierr.ErrProjectDateIsNotValid
	// }

	oldProject, err := s.repo.GetProject(ctx, repository.NewProjectFilter().
		ByIDs(projectReq.ID))
	if err != nil {
		return nil, err
	}

	newProject := &model.Project{
		ID:          projectReq.ID,
		Name:        projectReq.Name,
		Description: projectReq.Description,
		PhotoURL:    projectReq.PhotoURL,
		ReportURL:   projectReq.ReportURL,
		ReportName:  projectReq.ReportName,
		RepoURL:     projectReq.RepoURL,
		ActiveTo:    projectReq.ActiveTo,
	}
	if err := mergeProjectFields(oldProject, newProject); err != nil {
		return nil, err
	}

	return newProject, s.repo.UpdateProject(ctx, newProject)
}

func (s *service) DeleteProject(ctx context.Context, project *model.Project) error {
	return nil
}

func mergeProjectFields(oldProject, newProject *model.Project) error {
	//c фронта всегда новый дескрипшн должен приходить
	if newProject.Name == "" {
		newProject.Name = oldProject.Name
	}
	//a popo photo?
	if newProject.PhotoURL == "" {
		newProject.PhotoURL = oldProject.PhotoURL
	}
	if newProject.ReportURL == "" {
		newProject.ReportURL = oldProject.ReportURL
	}
	if newProject.ReportName == "" {
		newProject.ReportName = oldProject.ReportName
	}
	if newProject.RepoURL == "" {
		newProject.RepoURL = oldProject.RepoURL
	}
	//escho datu nazad nizya
	// if newProject.ActiveTo == nil {
	// 	newProject.ActiveTo = oldProject.ActiveTo
	// }

	return nil
}
