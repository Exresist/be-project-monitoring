package repository

import (
	"context"
	"fmt"
	"strconv"

	"be-project-monitoring/internal/db"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func (r *Repository) GetProject(ctx context.Context, filter *ProjectFilter) (*model.Project, error) {
	projects, err := r.GetProjects(ctx, filter.WithPaginator(1, 0))
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	case len(projects) == 0:
		return nil, ierr.ErrProjectNotFound
	default:
		return &projects[0], nil
	}
}

func (r *Repository) GetProjects(ctx context.Context, filter *ProjectFilter) ([]model.Project, error) {
	filter.Limit = db.NormalizeLimit(filter.Limit)
	rows, err := r.sq.Select(
		"p.id", "p.name",
		"p.description", "p.photo_url",
		"p.report_url", "p.report_name",
		"p.repo_url", "p.active_to").
		From("projects p").
		Where(conditionsFromProjectFilter(filter)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}()

	projects := make([]model.Project, 0)
	for rows.Next() {
		project := model.Project{}
		if err = rows.Scan(
			&project.ID, &project.Name,
			&project.Description, &project.PhotoURL,
			&project.ReportURL, &project.ReportName,
			&project.RepoURL, &project.ActiveTo,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *Repository) GetProjectCountByFilter(ctx context.Context, filter *ProjectFilter) (int, error) {
	var count int
	if err := r.sq.Select("COUNT(1)").
		From("projects p").
		Where(conditionsFromProjectFilter(filter)).
		QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while scanning sql row: %w", err)
	}

	return count, nil
}

func (r *Repository) InsertProject(ctx context.Context, project *model.Project) error {
	row := r.sq.Insert("projects").
		Columns("name",
			"description", "photo_url",
			"active_to").
		Values(project.Name,
			project.Description, project.PhotoURL,
			project.ActiveTo).
		Suffix("RETURNING \"id\"").
		QueryRowContext(ctx)

	if err := row.Scan(&project.ID); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
	}
	return nil
}
func (r *Repository) UpdateProject(ctx context.Context, project *model.Project) error {
	_, err := r.sq.Update("projects").
		SetMap(map[string]interface{}{
			"name":        project.Name,
			"description": project.Description,
			"photo_url":   project.PhotoURL,
			"report_url":  project.ReportURL,
			"report_name": project.ReportName,
			"repo_url":    project.RepoURL,
			"active_to":   project.ActiveTo,
		}).Where(sq.Eq{"id": project.ID}).
		ExecContext(ctx)
	return err
}

func (r *Repository) DeleteProject(ctx context.Context, id int) error {
	_, err := r.sq.Delete("projects").
		Where(sq.Eq{"id": id}).ExecContext(ctx)
	return err
}

func (r *Repository) GetProjectInfo(ctx context.Context, id int) (*model.ProjectInfo, error) {
	query := `SELECT p.id, p.name, p.description, p.photo_url, p.report_url,
	 			p.report_name, p.repo_url, p.active_to,
				ARRAY_AGG (part.id) participants_ids,
				ARRAY_AGG (part.role) participants_roles,
				ARRAY_AGG (u.id) users_ids, ARRAY_AGG (u.role) users_roles,
				ARRAY_AGG (u.color_code) users_color_codes, ARRAY_AGG (u.email) users_emails,
				ARRAY_AGG (u.username) users_usernames, ARRAY_AGG (u.first_name) users_first_names,
				ARRAY_AGG (u.last_name) users_last_names, ARRAY_AGG (u."group") users_groups,
				ARRAY_AGG (u.github_username) users_github_usernames,
				ARRAY_AGG (t.id) tasks_ids, ARRAY_AGG (t.name) tasks_names,
				ARRAY_AGG (t.description) tasks_descriptions, ARRAY_AGG (t.participant_id) participants_ids,
				ARRAY_AGG (t.status) tasks_statuses
				FROM projects p
				  JOIN participants part ON part.project_id = p.id
				  JOIN users u ON part.user_id = u.id
				  LEFT JOIN tasks t ON t.project_id = p.id
				  WHERE p.id = $1
				  GROUP BY p.id, p.name, p.description, p.photo_url, p.report_url,
			 	  p.report_name, p.repo_url, p.active_to
				  `

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			r.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}()

	if rows.Next() {
		projectInfo := model.ProjectInfo{}
		participantIDs := make(pq.Int64Array, 0)
		participantRoles := make(pq.StringArray, 0)
		usersIDs := make(pq.StringArray, 0)
		usersRoles := make(pq.StringArray, 0)
		usersColorCodes := make(pq.StringArray, 0)
		usersEmails := make(pq.StringArray, 0)
		usersUsernames := make(pq.StringArray, 0)
		usersFirstNames := make(pq.StringArray, 0)
		usersLastNames := make(pq.StringArray, 0)
		usersGroups := make(pq.StringArray, 0)
		usersGithubUsernames := make(pq.StringArray, 0)
		tasksIDs := make(pq.ByteaArray, 0)
		tasksNames := make(pq.ByteaArray, 0)
		tasksDescriptions := make(pq.ByteaArray, 0)
		participantsIDs := make(pq.ByteaArray, 0)
		tasksStatuses := make(pq.ByteaArray, 0)
		params := []any{&projectInfo.Project.ID, &projectInfo.Project.Name,
			&projectInfo.Project.Description, &projectInfo.Project.PhotoURL,
			&projectInfo.Project.ReportURL, &projectInfo.Project.ReportName,
			&projectInfo.Project.RepoURL, &projectInfo.Project.ActiveTo,
			&participantIDs, &participantRoles,
			&usersIDs, &usersRoles, &usersColorCodes, &usersEmails,
			&usersUsernames, &usersFirstNames, &usersLastNames, &usersGroups,
			&usersGithubUsernames, &tasksIDs, &tasksNames,
			&tasksDescriptions, &participantsIDs,
			&tasksStatuses}

		if err = rows.Scan(params...); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}

		participants := make([]model.Participant, 0)
		for i := range participantIDs {

			userID, err := uuid.Parse(usersIDs[i])
			if err != nil {
				return nil, fmt.Errorf("error while parsing user id: %w", err)
			}

			participants = append(participants, model.Participant{
				ShortUser: model.ShortUser{
					ID:             userID,
					Role:           model.UserRole(usersRoles[i]),
					ColorCode:      usersColorCodes[i],
					Email:          usersEmails[i],
					Username:       usersUsernames[i],
					FirstName:      usersFirstNames[i],
					LastName:       usersLastNames[i],
					Group:          usersGroups[i],
					GithubUsername: usersGithubUsernames[i]},
				Role: model.ParticipantRole(participantRoles[i]),
				ID:   int(participantIDs[i]),
			})
		}
		projectInfo.Participants = participants

		tasks := make([]model.Task, 0)
		for i := range tasksIDs {
			if tasksIDs[i] != nil {

				taskID, err := strconv.Atoi(string(tasksIDs[i]))
				if err != nil {
					return nil, err
				}

				task := model.Task{
					ShortTask: model.ShortTask{
						ID:     taskID,
						Name:   string(tasksNames[i]),
						Status: model.TaskStatus(tasksStatuses[i]),
					},
				}
				task.Description.Scan(tasksDescriptions[i])
				if participantsIDs[i] != nil {
					participantID, err := strconv.Atoi(string(participantsIDs[i]))
					if err != nil {
						return nil, err
					}
					task.ParticipantID.Scan(participantID)
				}
				tasks = append(tasks, task)
			}
		}
		projectInfo.Tasks = tasks
		return &projectInfo, nil
	}
	return nil, ierr.ErrProjectNotFound
}
