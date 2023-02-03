package api

import (
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	CreateTaskReq struct {
		Name              string    `json:"name"`
		Description       string    `json:"description"`
		SuggestedEstimate int       `json:"suggested_estimate"`
		CreatorUserID     uuid.UUID `json:"creator_user_id"`
		ParticipantUserID uuid.UUID `json:"paticipant_user_id"`
		Status            string    `json:"status"`
		ProjectID         int       `json:"project_id"`
	}
	GetTasksReq struct {
		ProjectID     int
		Name          *string
		ParticipantID *int
		Offset        int
		Limit         int
	}
	getTaskResp struct {
		Tasks []TaskResp
		Count int
	}
	UpdateTaskReq struct {
		ID                int       `json:"id"`
		Name              *string   `json:"name"`
		Description       *string   `json:"description"`
		SuggestedEstimate *int      `json:"suggested_estimate"`
		RealEstimate      *int      `json:"real_estimate"`
		Status            *string   `json:"status"`
		UpdatedAt         time.Time `json:"updated_at"`
		ParticipantUserID uuid.UUID `json:"participant_user_id"`
	}
	TaskResp struct {
		ID                int
		Name              string
		Description       string
		SuggestedEstimate int
		RealEstimate      int
		ParticipantID     int
		CreatorID         int
		Status            string
		CreatedAt         time.Time
		UpdatedAt         time.Time
		ProjectID         int
	}
	DeleteTaskReq struct {
		ID int `json:"id"`
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

	c.JSON(http.StatusOK, getTaskResp{
		Tasks: makeTasksResponse(tasks),
		Count: count,
	})

}
func (s *Server) createTask(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	parts := strings.Split(strings.TrimSpace(token), " ")

	id, err := s.svc.GetUserIDFromToken(c.Request.Context(), parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	taskReq := &CreateTaskReq{
		CreatorUserID: id,
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

	c.JSON(http.StatusOK, makeTaskResponse(task))
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

	c.JSON(http.StatusOK, makeTaskResponse(task))
}
func (s *Server) deleteTask(c *gin.Context) {
	taskReq := &DeleteTaskReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(taskReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err := s.svc.DeleteTask(c, taskReq); err != nil {
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

	c.JSON(http.StatusOK, taskInfo)
}

func makeTaskResponse(task *model.Task) *TaskResp {
	return &TaskResp{
		ID:                task.ID,
		Name:              task.Name,
		Description:       task.Description,
		SuggestedEstimate: task.SuggestedEstimate,
		RealEstimate:      task.RealEstimate,
		ParticipantID:     int(task.ParticipantID.Int64),
		CreatorID:         task.CreatorID,
		Status:            string(task.Status),
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
		ProjectID:         task.ProjectID,
	}
}
func makeTasksResponse(tasks []model.Task) []TaskResp {
	var taskResponses []TaskResp
	for _, task := range tasks {
		taskResponses = append(taskResponses, *makeTaskResponse(&task))
	}
	return taskResponses
}
