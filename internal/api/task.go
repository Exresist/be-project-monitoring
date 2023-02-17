package api

import (
	"be-project-monitoring/internal/domain"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	CreateTaskReq struct {
		Name              string    `json:"name"`
		Description       string    `json:"description"`
		SuggestedEstimate *int      `json:"suggested_estimate"`
		CreatorUserID     uuid.UUID `json:"creator_user_id"`
		ParticipantUserID uuid.UUID `json:"paticipant_user_id"`
		Status            string    `json:"status"`
		ProjectID         int       `json:"project_id"`
	}
	taskResp struct {
		ID                int       `json:"id"`
		Name              string    `json:"name"`
		Description       string    `json:"description,omitempty"`
		SuggestedEstimate int       `json:"suggested_estimate,omitempty"`
		RealEstimate      int       `json:"real_estimate,omitempty"`
		ParticipantID     int       `json:"participant_id,omitempty"`
		CreatorID         int       `json:"creator_id,omitempty"`
		Status            string    `json:"status"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		ProjectID         int       `json:"project_id"`
	}
	GetTasksReq struct {
		ProjectID     int
		Name          *string
		ParticipantID *int
		Offset        int
		Limit         int
	}
	UpdateTaskReq struct {
		ID                int       `json:"id"`
		Name              *string   `json:"name"`
		Description       *string   `json:"description"`
		SuggestedEstimate *int      `json:"suggested_estimate"`
		RealEstimate      *int      `json:"real_estimate"`
		Status            *string   `json:"status"`
		ParticipantUserID uuid.UUID `json:"participant_user_id"`
		ChangeParticipant *bool     `json:"change_participant"`
	}
)

func (s *Server) getTasks(c *gin.Context) {
	taskReq := &GetTasksReq{}

	projectID, err := strconv.Atoi(c.Query("project_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ierr.ErrInvalidProjectID)
		return
	}
	taskReq.ProjectID = projectID

	if name := c.Query("name"); name != "" {
		taskReq.Name = &name
	}
	*taskReq.ParticipantID, _ = strconv.Atoi(c.Query("participant_id"))
	taskReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	taskReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	tasks, count, err := s.svc.GetTasks(c, taskReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		Tasks []taskResp `json:"tasks"`
		Count int        `json:"count"`
	}{
		Tasks: makeTasksResponses(tasks),
		Count: count,
	})

}
func (s *Server) createTask(c *gin.Context) {
	taskReq := &CreateTaskReq{
		CreatorUserID: c.MustGet(string(domain.UserIDCtx)).(uuid.UUID), //ОБЯЗАТЕЛЬНО ПРОВЕРИТЬ!
	}
	if err := json.NewDecoder(c.Request.Body).Decode(taskReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	task, err := s.svc.CreateTask(c, taskReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, makeTaskResponse(*task))
}
func (s *Server) updateTask(c *gin.Context) {
	taskReq := &UpdateTaskReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(taskReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	task, err := s.svc.UpdateTask(c, taskReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, makeTaskResponse(*task))
}
func (s *Server) deleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	if err := s.svc.DeleteTask(c, taskID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) getTaskInfo(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	taskInfo, err := s.svc.GetTaskInfo(c.Request.Context(), taskID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, struct {
		Task        *taskResp
		Creator     model.ShortUser
		Participant model.ShortUser
	}{
		Task:        makeTaskResponse(taskInfo.Task),
		Creator:     taskInfo.Creator,
		Participant: taskInfo.Participant,
	})
}

func makeTaskResponse(task model.Task) *taskResp {
	return &taskResp{
		ID:                task.ID,
		Name:              task.Name,
		Description:       task.Description.String,
		SuggestedEstimate: int(task.SuggestedEstimate.Int64),
		RealEstimate:      int(task.RealEstimate.Int64),
		ParticipantID:     int(task.ParticipantID.Int64),
		CreatorID:         int(task.CreatorID.Int64),
		Status:            string(task.Status),
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
		ProjectID:         task.ProjectID,
	}
}
func makeTasksResponses(tasks []model.Task) []taskResp {
	taskResponses := make([]taskResp, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses, *makeTaskResponse(task))
	}
	return taskResponses
}
func makeShortTasksResponses(tasks []model.ShortTask) []taskResp {
	taskResponses := make([]taskResp, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses,
			*makeTaskResponse(model.Task{
				ShortTask: task,
			}))
	}
	return taskResponses
}
