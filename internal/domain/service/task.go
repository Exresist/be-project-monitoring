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
	creatorID := &sql.NullInt64{}
	participantID := &sql.NullInt64{}

	if creator, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(taskReq.CreatorUserID).ByProjectID(taskReq.ProjectID)); err != nil {
		return nil, ierr.ErrTaskCreatorUserIDNotFound
	} else {
		creatorID.Scan(creator.ID)
	}

	if taskReq.ParticipantUserID != uuid.Nil {
		if participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
			ByUserID(taskReq.ParticipantUserID).ByProjectID(taskReq.ProjectID)); err != nil {
			return nil, ierr.ErrTaskParticipantUserIDNotFound
		} else {
			participantID.Scan(participant.ID)
		}
	}

	if strings.TrimSpace(taskReq.Name) == "" {
		return nil, ierr.ErrTaskNameIsInvalid
	}
	if strings.TrimSpace(taskReq.Status) == "" {
		taskReq.Status = string(model.TODO)
	}
	if _, ok := model.TaskStatuses[taskReq.Status]; !ok {
		return nil, ierr.ErrInvalidStatus
	}

	task := &model.Task{
		ShortTask: model.ShortTask{
			Name:          taskReq.Name,
			ParticipantID: *participantID,
			Status:        model.TaskStatus(taskReq.Status),
		},
		CreatorID: *creatorID,
		ProjectID: taskReq.ProjectID,
	}
	if strings.TrimSpace(taskReq.Description) != "" {
		task.Description.Scan(taskReq.Description)
	}
	if taskReq.SuggestedEstimate != nil {
		if *taskReq.SuggestedEstimate > 0 {
			task.SuggestedEstimate.Scan(*taskReq.SuggestedEstimate)
		} else {
			return nil, ierr.ErrTaskSuggestedEstimateIsInvalid
		}
	}

	return task, s.repo.InsertTask(ctx, task)
}
func (s *service) UpdateTask(ctx context.Context, taskReq *api.UpdateTaskReq) (*model.Task, error) {
	if taskReq.ID == 0 {
		return nil, ierr.ErrTaskIDIsInvalid
	}
	oldTask, err := s.repo.GetTask(ctx, repository.NewTaskFilter().
		ByID(taskReq.ID))
	if err != nil {
		return nil, err
	}

	participantID := &sql.NullInt64{}
	if taskReq.ChangeParticipant != nil && *taskReq.ChangeParticipant {
		if taskReq.ParticipantUserID != uuid.Nil {
			if participant, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
				ByUserID(taskReq.ParticipantUserID).ByProjectID(oldTask.ProjectID)); err != nil {
				return nil, ierr.ErrTaskParticipantUserIDNotFound
			} else {
				participantID.Scan(participant.ID)
			}
		}
	} else {
		participantID.Scan(&oldTask.ParticipantID)
	}

	newTask, err := mergeTaskFields(oldTask, taskReq, *participantID)
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
	// if _, err := s.repo.GetTask(ctx, repository.NewTaskFilter().ByID(id)); err != nil {
	// 	return nil, err
	// }
	return s.repo.GetTaskInfo(ctx, id)
}

func mergeTaskFields(oldTask *model.Task, taskReq *api.UpdateTaskReq, newParticipantID sql.NullInt64) (*model.Task, error) {
	newTask := &model.Task{
		ShortTask: model.ShortTask{
			ID:            taskReq.ID,
			ParticipantID: newParticipantID,
		},
		UpdatedAt: time.Now(),
		CreatorID: oldTask.CreatorID,
		ProjectID: oldTask.ProjectID,
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
	} else {
		newTask.Name = *taskReq.Name
	}
	if taskReq.Description == nil {
		newTask.Description = oldTask.Description
	} else {
		newTask.Description.Scan(*taskReq.Description)
	}

	if taskReq.SuggestedEstimate == nil {
		newTask.SuggestedEstimate = oldTask.SuggestedEstimate
	} else if *taskReq.SuggestedEstimate > 0 {
		newTask.SuggestedEstimate.Scan(*taskReq.SuggestedEstimate)
	} else {
		return nil, ierr.ErrTaskSuggestedEstimateIsInvalid
	}

	if taskReq.RealEstimate == nil {
		newTask.RealEstimate = oldTask.RealEstimate
	} else if *taskReq.RealEstimate > 0 {
		newTask.RealEstimate.Scan(*taskReq.RealEstimate)
	} else {
		return nil, ierr.ErrTaskRealEstimateIsInvalid
	}

	return newTask, nil
}
