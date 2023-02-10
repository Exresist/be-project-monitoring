package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
)

func (s *service) GetProjects(ctx context.Context, projectReq *api.GetProjectsReq) ([]model.Project, int, error) {

	filter := repository.NewProjectFilter().
		WithPaginator(uint64(projectReq.Limit), uint64(projectReq.Offset)).
		ByProjectName(strings.TrimSpace(projectReq.Name))

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
	if strings.TrimSpace(projectReq.Name) == "" {
		return nil, ierr.ErrProjectNameIsInvalid
	}
	if projectReq.ActiveTo.IsZero() || projectReq.ActiveTo.Before(time.Now()) {
		return nil, ierr.ErrProjectActiveToIsInvalid
	}

	found, err := s.repo.GetProject(ctx, repository.NewProjectFilter().
		ByProjectName(projectReq.Name))
	if err != nil && !errors.Is(err, ierr.ErrProjectNotFound) {
		return nil, err
	}
	if found != nil {
		return nil, ierr.ErrProjectNameAlreadyExists
	}

	project := &model.Project{
		ShortProject: model.ShortProject{
			Name:     projectReq.Name,
			ActiveTo: projectReq.ActiveTo,
		}}
	project.Description.Scan(projectReq.Description)
	project.PhotoURL.Scan(projectReq.PhotoURL)

	return project, s.repo.InsertProject(ctx, project)
}

func (s *service) UpdateProject(ctx context.Context, projectReq *api.UpdateProjectReq) (*model.Project, error) {
	oldProject, err := s.repo.GetProject(ctx, repository.NewProjectFilter().
		ByID(projectReq.ID))
	if err != nil {
		return nil, err
	}

	newProject, err := mergeProjectFields(oldProject, projectReq)
	if err != nil {
		return nil, err
	}

	return newProject, s.repo.UpdateProject(ctx, newProject)
}

func (s *service) DeleteProject(ctx context.Context, id int) error {
	if _, err := s.repo.GetProject(ctx, repository.NewProjectFilter().ByID(id)); err != nil {
		return err
	}
	return s.repo.DeleteProject(ctx, id)
}

func (s *service) GetProjectInfo(ctx context.Context, id int) (*model.ProjectInfo, error) {
	if _, err := s.repo.GetProject(ctx, repository.NewProjectFilter().ByID(id)); err != nil {
		return nil, err
	}
	return s.repo.GetProjectInfo(ctx, id)
}

func mergeProjectFields(oldProject *model.Project, projectReq *api.UpdateProjectReq) (*model.Project, error) {
	newProject := &model.Project{
		ShortProject: model.ShortProject{
			ID: 	 projectReq.ID,
			ActiveTo: projectReq.ActiveTo,
		}}

	if projectReq.ActiveTo.IsZero() || projectReq.ActiveTo.Before(time.Now()) {
		newProject.ActiveTo = oldProject.ActiveTo
	}

	if projectReq.Name == nil {
		newProject.Name = oldProject.Name
	} else {
		newProject.Name = *projectReq.Name
	}
	if projectReq.Description == nil {
		newProject.Description = oldProject.Description
	} else {
		newProject.Description.Scan(*projectReq.Description)
	}
	if projectReq.PhotoURL == nil {
		newProject.PhotoURL = oldProject.PhotoURL
	} else {
		newProject.PhotoURL.Scan(*projectReq.PhotoURL)
	}
	if projectReq.ReportURL == nil {
		newProject.ReportURL = oldProject.ReportURL
	} else {
		newProject.ReportURL.Scan(*projectReq.ReportURL)
	}
	if projectReq.ReportName == nil {
		newProject.ReportName = oldProject.ReportName
	} else {
		newProject.ReportName.Scan(*projectReq.ReportName)
	}
	if projectReq.RepoURL == nil {
		newProject.RepoURL = oldProject.RepoURL
	} else {
		newProject.RepoURL.Scan(*projectReq.RepoURL)
	}
	return newProject, nil
}
