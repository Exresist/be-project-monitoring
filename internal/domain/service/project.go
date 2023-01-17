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

	if projectReq.ActiveTo.IsZero() || projectReq.ActiveTo.Before(time.Now()) {
		return nil, ierr.ErrProjectActiveToIsInvalid
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

	oldProject, err := s.repo.GetProject(ctx, repository.NewProjectFilter().
		ByIDs(projectReq.ID))
	if err != nil {
		return nil, err
	}

	newProject, err := mergeProjectFields(oldProject, projectReq)
	if err != nil {
		return nil, err
	}

	return newProject, s.repo.UpdateProject(ctx, newProject)
}

func (s *service) DeleteProject(ctx context.Context, projectReq *api.DeleteProjectReq) error {
	if _, err := s.repo.GetProject(ctx, repository.NewProjectFilter().ByIDs(projectReq.ID)); err != nil {
		return err
	}

	return s.repo.DeleteProject(ctx, projectReq.ID)
}

func mergeProjectFields(oldProject *model.Project, projectReq *api.UpdateProjectReq) (*model.Project, error) {

	newProject := &model.Project{
		ID:          projectReq.ID,
		Name:        *projectReq.Name,
		Description: *projectReq.Description,
		PhotoURL:    *projectReq.PhotoURL,
		ReportURL:   *projectReq.ReportURL,
		ReportName:  *projectReq.ReportName,
		RepoURL:     *projectReq.RepoURL,
		ActiveTo:    projectReq.ActiveTo,
	}

	if projectReq.Name == nil {
		newProject.Name = oldProject.Name
	}
	if projectReq.Description == nil {
		newProject.Description = oldProject.Description
	}
	if projectReq.PhotoURL == nil {
		newProject.PhotoURL = oldProject.PhotoURL
	}
	if projectReq.ReportURL == nil {
		newProject.ReportURL = oldProject.ReportURL
	}
	if projectReq.ReportName == nil {
		newProject.ReportName = oldProject.ReportName
	}
	if projectReq.RepoURL == nil {
		newProject.RepoURL = oldProject.RepoURL
	}
	if projectReq.ActiveTo.IsZero() || projectReq.ActiveTo.Before(time.Now()) {
		newProject.ActiveTo = oldProject.ActiveTo
	}

	return newProject, nil
}
