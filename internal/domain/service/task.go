package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *service) GetTasks(ctx context.Context, taskReq *api.GetTasksReq) ([]model.Task, int, error) {
	filter := repository.NewTaskFilter().
		WithPaginator(uint64(taskReq.Limit), uint64(taskReq.Offset)).
		ByProjectID(taskReq.ProjectID).ByParticipantID(*taskReq.ParticipantID).ByTaskName(*taskReq.Name)
	count, err := s.repo.GetTaskCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	tasks, err := s.repo.GetTasks(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return tasks, count, nil
}

func (s *service) CreateTask(ctx context.Context, taskReq *api.CreateTaskReq) (*model.Task, error) {
	var (
		creatorID,
		participantID sql.NullInt64
	)

	if creator, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(taskReq.CreatorUserID).ByProjectID(taskReq.ProjectID)); err != nil {
		return nil, ierr.ErrTaskCreatorUserIDNotFound
	} else {
		creatorID.Scan(&creator.ID)
	}

	if taskReq.ParticipantUserID != uuid.Nil {
		if participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
			ByUserID(taskReq.ParticipantUserID).ByProjectID(taskReq.ProjectID)); err != nil {
			return nil, ierr.ErrTaskParticipantUserIDNotFound
		} else {
			participantID.Scan(&participant.ID)
		}
	}

	if strings.TrimSpace(taskReq.Name) == "" {
		return nil, ierr.ErrTaskNameIsInvalid
	}
	if taskReq.SuggestedEstimate < 0 {
		return nil, ierr.ErrTaskSuggestedEstimateIsInvalid
	}
	if strings.TrimSpace(taskReq.Status) == "" {
		taskReq.Status = string(model.TODO)
	}
	if _, ok := model.TaskStatuses[taskReq.Status]; !ok {
		return nil, ierr.ErrInvalidStatus
	}

	task := &model.Task{
		Name:              taskReq.Name,
		Description:       taskReq.Description,
		SuggestedEstimate: taskReq.SuggestedEstimate,
		ParticipantID:     participantID,
		CreatorID:         creatorID,
		Status:            model.TaskStatus(taskReq.Status),
		ProjectID:         taskReq.ProjectID,
	}
	return task, s.repo.InsertTask(ctx, task)
}
func (s *service) UpdateTask(ctx context.Context, taskReq *api.UpdateTaskReq) (*model.Task, error) {
	oldTask, err := s.repo.GetTask(ctx, repository.NewTaskFilter().
		ByID(taskReq.ID))
	if err != nil {
		return nil, err
	}

	var participantID sql.NullInt64
	if taskReq.ChangeParticipant != nil && *taskReq.ChangeParticipant {
		if taskReq.ParticipantUserID != uuid.Nil {
			if participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
				ByUserID(taskReq.ParticipantUserID).ByProjectID(oldTask.ProjectID)); err != nil {
				return nil, ierr.ErrTaskParticipantUserIDNotFound
			} else {
				participantID.Scan(&participant.ID)
			}
		}
	} else {
		participantID.Scan(&oldTask.ParticipantID)
	}

	newTask, err := mergeTaskFields(oldTask, taskReq, participantID)
	if err != nil {
		return nil, err
	}
	return newTask, s.repo.UpdateTask(ctx, newTask)
}

func (s *service) DeleteTask(ctx context.Context, id int) error {
	if _, err := s.repo.GetTask(ctx, repository.NewTaskFilter().ByID(id)); err != nil {
		return err
	}
	return s.repo.DeleteTask(ctx, id)
}

func (s *service) GetTaskInfo(ctx context.Context, id int) (*model.TaskInfo, error) {
	if _, err := s.repo.GetTask(ctx, repository.NewTaskFilter().ByID(id)); err != nil {
		return nil, err
	}
	return s.repo.GetTaskInfo(ctx, id)
}

func mergeTaskFields(oldTask *model.Task, taskReq *api.UpdateTaskReq, newParticipantID sql.NullInt64) (*model.Task, error) {
	newTask := &model.Task{
		ID:                taskReq.ID,
		Name:              *taskReq.Name,
		Description:       *taskReq.Description,
		SuggestedEstimate: *taskReq.SuggestedEstimate,
		RealEstimate:      *taskReq.RealEstimate,
		Status:            model.TaskStatus(*taskReq.Status),
		UpdatedAt:         time.Now(),
		CreatorID:         oldTask.CreatorID,
		ParticipantID:     newParticipantID,
		ProjectID:         oldTask.ProjectID,
	}
	if _, ok := model.TaskStatuses[*taskReq.Status]; ok {
		newTask.Status = model.TaskStatus(*taskReq.Status)
	} else {
		newTask.Status = oldTask.Status
	}
	if !newParticipantID.Valid {
		newTask.ParticipantID = oldTask.ParticipantID
	}

	if taskReq.Name == nil {
		newTask.Name = oldTask.Name
	}
	if taskReq.Description == nil {
		newTask.Description = oldTask.Description
	}
	if taskReq.SuggestedEstimate == nil || *taskReq.SuggestedEstimate < 0 {
		newTask.SuggestedEstimate = oldTask.SuggestedEstimate
	}
	if taskReq.RealEstimate == nil || *taskReq.RealEstimate < 0 {
		newTask.RealEstimate = oldTask.RealEstimate
	}

	return newTask, nil
}
