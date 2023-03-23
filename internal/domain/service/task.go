package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"be-project-monitoring/internal/repository"

	"github.com/google/uuid"
)

func (s *service) GetTasks(ctx context.Context, taskReq *api.GetTasksReq) ([]model.Task, int, error) {
	filter := repository.NewTaskFilter().
		WithPaginator(uint64(taskReq.Limit), uint64(taskReq.Offset)).
		ByProjectID(taskReq.ProjectID)
	if taskReq.ParticipantID != nil {
		filter.ByParticipantID(*taskReq.ParticipantID)
	}
	if taskReq.Name != nil {
		filter.ByTaskName(*taskReq.Name)
	}
	if taskReq.Approved != nil {
		filter.ByApproved(*taskReq.Approved)
	}

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

func (s *service) CreateTask(ctx context.Context, creatorUserID uuid.UUID, taskReq *api.CreateTaskReq) (*model.Task, error) {
	creatorID := &sql.NullInt64{}
	participantID := &sql.NullInt64{}

	creator, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
		ByUserID(creatorUserID).ByProjectID(taskReq.ProjectID))
	if err != nil {
		return nil, err
	}
	creatorID.Scan(creator.ID)
	fmt.Println(creatorID)
	// if _, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
	// 	ByID(creator.ID).ByProjectID(taskReq.ProjectID)); err != nil {
	// 	return nil, ierr.ErrTaskCreatorUserIDNotFound
	// } else {
	// 	creatorID.Scan(taskReq.CreatorID)
	// }

	if taskReq.ParticipantID != nil {
		if _, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
			ByID(*taskReq.ParticipantID).ByProjectID(taskReq.ProjectID)); err != nil {
			return nil, ierr.ErrTaskParticipantIDNotFound
		} else {
			//fmt.Println(participant)
			participantID.Scan(*taskReq.ParticipantID)
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
			CreatorID:     *creatorID,
			Status:        model.TaskStatus(taskReq.Status),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		ProjectID: taskReq.ProjectID,
	}

	if strings.TrimSpace(taskReq.Description) != "" {
		task.Description.Scan(taskReq.Description)
	}

	if taskReq.SuggestedEstimate != 0 {
		task.Estimate.Scan(taskReq.SuggestedEstimate)
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
	if taskReq.ParticipantID != nil {
		if *taskReq.ParticipantID != 0 {
			if _, err := s.repo.GetParticipant(ctx, repository.NewParticipantFilter().
				ByID(*taskReq.ParticipantID).ByProjectID(taskReq.ProjectID)); err != nil {
				return nil, ierr.ErrTaskParticipantIDNotFound
			} else {
				participantID.Scan(*taskReq.ParticipantID)
			}
		}
	} else if oldTask.ParticipantID.Valid {
		participantID.Scan(oldTask.ParticipantID.Int64)
	}

	newTask, err := mergeTaskFields(oldTask, taskReq, *participantID)
	if err != nil {
		return nil, err
	}

	return newTask, s.repo.UpdateTask(ctx, newTask)
}

func (s *service) DeleteTask(ctx context.Context, id int) error {
	return s.repo.DeleteTask(ctx, id)
}

func (s *service) GetTaskInfo(ctx context.Context, id int) (*model.TaskInfo, error) {
	return s.repo.GetTaskInfo(ctx, id)
}

func mergeTaskFields(oldTask *model.Task, taskReq *api.UpdateTaskReq, newParticipantID sql.NullInt64) (*model.Task, error) {
	newTask := &model.Task{
		ShortTask: model.ShortTask{
			ID:            taskReq.ID,
			ParticipantID: newParticipantID,
			CreatorID:     oldTask.CreatorID,
			CreatedAt:     oldTask.CreatedAt,
			UpdatedAt:     time.Now(),
		},
		ProjectID: oldTask.ProjectID,
	}

	if taskReq.Status != nil {
		if _, ok := model.TaskStatuses[*taskReq.Status]; ok {
			newTask.Status = model.TaskStatus(*taskReq.Status)
		}
	} else {
		newTask.Status = oldTask.Status
	}

	if taskReq.Name == nil {
		newTask.Name = oldTask.Name
	} else if strings.TrimSpace(*taskReq.Name) != "" {
		newTask.Name = *taskReq.Name
	} else {
		return nil, ierr.ErrTaskNameIsInvalid
	}
	if taskReq.Description == nil {
		newTask.Description = oldTask.Description
	} else {
		newTask.Description.Scan(*taskReq.Description)
	}

	if taskReq.SuggestedEstimate == nil {
		newTask.Estimate = oldTask.Estimate
	} else {
		newTask.Estimate.Scan(*taskReq.SuggestedEstimate)
	}

	if taskReq.Approved == nil {
		newTask.Approved = oldTask.Approved
	} else {
		newTask.Approved.Scan(*taskReq.Approved)
	}

	return newTask, nil
}
