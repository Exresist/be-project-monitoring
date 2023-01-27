package service

import (
	"be-project-monitoring/internal/api"
	"be-project-monitoring/internal/domain/model"
	"be-project-monitoring/internal/domain/repository"
	ierr "be-project-monitoring/internal/errors"
	"context"
	"strings"
	"time"
)

func (s *service) GetTasks(ctx context.Context, taskReq *api.GetTaskReq) ([]model.Task, int, error) {
	filter := repository.NewTaskFilter().
		WithPaginator(uint64(taskReq.Limit), uint64(taskReq.Offset)).
		ByIDs(taskReq.ID).ByTaskNames(taskReq.Name).ByCreatedAt() //FIX THIS!!!!

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
	if taskReq.CreatorID == 0 {
		return nil, ierr.ErrTaskCreatorIDIsInvalid
	}
	if strings.Trim(taskReq.Name, " \n\r\t") == "" {
		return nil, ierr.ErrTaskNameIsInvalid
	}
	if taskReq.SuggestedEstimate < 0 { //nado li chekat? mozhet li pustoye
		taskReq.SuggestedEstimate = 0
	}
	task := &model.Task{
		Name:              taskReq.Name,
		Description:       taskReq.Description,
		SuggestedEstimate: taskReq.SuggestedEstimate,
		CreatorID:         taskReq.CreatorID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		//participant, status i td
	}

	if err := s.repo.InsertTask(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *service) UpdateTask(ctx context.Context, taskReq *api.UpdateTaskReq) (*model.Task, error) {
	oldProject, err := s.repo.GetTask(ctx, repository.NewTaskFilter().
		ByIDs(taskReq.ID))
	if err != nil {
		return nil, err
	}

	newTask, err := mergeTaskFields(oldProject, taskReq)
	if err != nil {
		return nil, err
	}
	return newTask, s.repo.UpdateTask(ctx, newTask)
}

func (s *service) DeleteTask(ctx context.Context, taskReq *api.DeleteTaskReq) error {
	if _, err := s.repo.GetTask(ctx, repository.NewTaskFilter().ByIDs(taskReq.ID)); err != nil {
		return err
	}
	return s.repo.DeleteTask(ctx, taskReq.ID)
}

func mergeTaskFields(oldTask *model.Task, taskReq *api.UpdateTaskReq) (*model.Task, error) {

	newTask := &model.Task{
		ID:                taskReq.ID,
		Name:              *taskReq.Name,
		Description:       *taskReq.Description,
		SuggestedEstimate: *taskReq.SuggestedEstimate,
		RealEstimate:      *taskReq.RealEstimate,
		Status:            model.TaskStatus(*taskReq.Status),
		UpdatedAt:         time.Now(),
		ParticipantID:     *taskReq.ParticipantID,
	}

	if _, ok := model.TaskStatuses[*taskReq.Status]; ok {
		newTask.Status = model.TaskStatus(*taskReq.Status)
	} else {
		newTask.Status = oldTask.Status
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
	if taskReq.ParticipantID == nil {
		newTask.ParticipantID = oldTask.ParticipantID
	}

	return newTask, nil
}
