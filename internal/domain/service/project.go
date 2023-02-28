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
		// WithPaginator(uint64(projectReq.Limit), uint64(projectReq.Offset)).
		ByProjectNameLike(projectReq.SearchText)

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
		return nil, ierr.ErrInvalidProjectName
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
	if strings.TrimSpace(projectReq.Description) != "" {
		project.Description.Scan(projectReq.Description)
	}
	if strings.TrimSpace(projectReq.PhotoURL) != "" {
		project.PhotoURL.Scan(projectReq.PhotoURL)
	}

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
	// if _, err := s.repo.GetProject(ctx, repository.NewProjectFilter().ByID(id)); err != nil {
	// 	return err
	// }
	return s.repo.DeleteProject(ctx, id)
}
func (s *service) GetProjectCommits(ctx context.Context, id int) ([]model.CommitsInfo, error) {
	project, err := s.repo.GetProject(ctx, repository.NewProjectFilter().ByID(id))
	if err != nil {
		return nil, err
	}

	users, err := s.repo.GetPartialUsers(ctx, repository.NewUserFilter().ByAtProject(id))
	if err != nil {
		return nil, err
	}

	usersCommitsInfo := make(map[string]model.CommitsInfo, len(users))
	for _, user := range users {
		usersCommitsInfo[user.GithubUsername] = model.CommitsInfo{
			GithubUsername: user.GithubUsername,
			Username:       user.Username,
		}
	}

	if !project.RepoURL.Valid {
		return nil, ierr.ErrRepositoryURLIsEmpty
	}
	repoURL := strings.Split(project.RepoURL.String, "/") //https://github.com/Exresist/be-project-monitoring
	if len(repoURL) != 5 {
		return nil, ierr.ErrRepositoryURLWrongFormat
	}
	// owner := repoURL[3]
	// repoName := repoURL[4]
	commits, _, err := s.githubCl.Repositories.ListCommits(ctx, repoURL[3], repoURL[4], nil)
	if err != nil {
		return nil, err
	}

	for _, commit := range commits {
		ghUsername := commit.Author.GetLogin()
		if info, ok := usersCommitsInfo[ghUsername]; ok {
			curDate := commit.Commit.Author.GetDate()
			//checking only first because otherwise we're updating both
			if !info.FirstCommitDate.IsZero() {
				if info.FirstCommitDate.After(curDate) {
					info.FirstCommitDate = curDate
				}
				if info.LastCommitDate.Before(curDate) {
					info.LastCommitDate = curDate
				}
			} else {
				info.FirstCommitDate = curDate
				info.LastCommitDate = curDate
			}

			info.FirstCommitDate = commit.Commit.Author.GetDate()
			info.LastCommitDate = commit.Commit.Author.GetDate()
			info.Total++
			usersCommitsInfo[ghUsername] = info
		}
	}

	res := make([]model.CommitsInfo, 0, len(usersCommitsInfo))

	for _, commitInfo := range usersCommitsInfo {
		res = append(res, commitInfo)
	}

	return res, nil
}
func (s *service) GetProjectInfo(ctx context.Context, id int) (*model.ProjectInfo, error) {
	project, err := s.repo.GetProject(ctx, repository.NewProjectFilter().ByID(id))
	if err != nil {
		return nil, err
	}
	participants, err := s.repo.GetParticipants(ctx, repository.NewParticipantFilter().ByProjectID(id))
	if err != nil {
		return nil, err
	}
	tasks, err := s.repo.GetTasks(ctx, repository.NewTaskFilter().ByProjectID(id))
	if err != nil {
		return nil, err
	}
	projectInfo := &model.ProjectInfo{
		Project:      *project,
		Participants: participants,
		Tasks:        tasks,
	}
	return projectInfo, nil
	//return s.repo.GetProjectInfo(ctx, id)
}

func mergeProjectFields(oldProject *model.Project, projectReq *api.UpdateProjectReq) (*model.Project, error) {
	newProject := &model.Project{
		ShortProject: model.ShortProject{
			ID:       projectReq.ID,
			ActiveTo: projectReq.ActiveTo,
		}}

	if projectReq.ActiveTo.IsZero() {
		newProject.ActiveTo = oldProject.ActiveTo
	} else if projectReq.ActiveTo.Before(time.Now()) {
		return nil, ierr.ErrInvalidActiveTo
	}
	if projectReq.ActiveTo.IsZero() || projectReq.ActiveTo.Before(time.Now()) {
		newProject.ActiveTo = oldProject.ActiveTo
	}

	if projectReq.Name == nil {
		newProject.Name = oldProject.Name
	} else if strings.TrimSpace(*projectReq.Name) == "" {
		return nil, ierr.ErrInvalidProjectName
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
