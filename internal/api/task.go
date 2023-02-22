package api

import (
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	CreateTaskReq struct {
		Name              string `json:"title"`
		Description       string `json:"description"`
		SuggestedEstimate string `json:"estimatedTime"`
		CreatorID         int    `json:"creatorId"`
		ParticipantID     *int   `json:"asignee"`
		Status            string `json:"status"`
		ProjectID         int    `json:"projectId"`
	}
	ShortTaskResp struct {
		ID                int         `json:"id"`
		Name              string      `json:"title"`
		Status            string      `json:"status"`
		Description       string      `json:"description"`
		Estimate          string      `json:"estimatedTime"`
		CreatedAt         time.Time   `json:"createdAt"`
		UpdatedAt         time.Time   `json:"updatedAt"`
		ParticipantID     int         `json:"asignee,omitempty"`
		ParticipantIDNull interface{} `json:"asignee,omitempty"`
	}
	TaskResp struct {
		ShortTaskResp
		CreatorID int `json:"creatorId"`
		//ProjectID int       `json:"projectId"`
	}
	taskInfoResp struct {
		TaskResp
		// Creator     model.ShortUser `json:"creator"`
		// Participant model.ShortUser `json:"asignee"`
	}
	GetTasksReq struct {
		ProjectID     int
		Name          *string
		ParticipantID *int
		Offset        int
		Limit         int
	}
	UpdateTaskReq struct {
		ID                int     `json:"id"`
		Name              *string `json:"title"`
		Description       *string `json:"description"`
		SuggestedEstimate *string `json:"estimatedTime"`
		Status            *string `json:"status"`
		ParticipantID     *int    `json:"asignee"`
		ProjectID         int     `json:"projectId"`
		//ChangeParticipant *bool   `json:"change_participant"`
	}
)

func (s *Server) getTasks(c *gin.Context) {
	taskReq := &GetTasksReq{}

	projectID, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ierr.ErrInvalidProjectID)
		return
	}
	taskReq.ProjectID = projectID

	if name := c.Query("name"); name != "" {
		taskReq.Name = &name
	}
	*taskReq.ParticipantID, _ = strconv.Atoi(c.Query("asignee"))
	taskReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	taskReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	tasks, count, err := s.svc.GetTasks(c, taskReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		Tasks []TaskResp `json:"tasks"`
		Count int        `json:"count"`
	}{
		Tasks: makeTasksResponses(tasks),
		Count: count,
	})

}
func (s *Server) createTask(c *gin.Context) {
	// taskReq := &CreateTaskReq{
	// 	CreatorID: c.MustGet(string(domain.UserIDCtx)).(uuid.UUID), //ОБЯЗАТЕЛЬНО ПРОВЕРИТЬ!
	// }
	taskReq := &CreateTaskReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(taskReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	taskReq.ProjectID = projectID

	task, err := s.svc.CreateTask(c, taskReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, makeTaskResponse(*task))
}
func (s *Server) updateTask(c *gin.Context) {
	taskReq := &UpdateTaskReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(taskReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	taskReq.ProjectID = projectID

	task, err := s.svc.UpdateTask(c, taskReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, makeTaskResponse(*task))
}
func (s *Server) deleteTask(c *gin.Context) {
	deletedTask := &struct {
		ID int `json:"id"`
	}{}
	//taskID, err := strconv.Atoi(c.Param("id"))
	if err := json.NewDecoder(c.Request.Body).Decode(deletedTask); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	if err := s.svc.DeleteTask(c, deletedTask.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) getTaskInfo(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	taskInfo, err := s.svc.GetTaskInfo(c.Request.Context(), taskID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskInfoResp{
		TaskResp: makeTaskResponse(taskInfo.Task),
		// Creator:     taskInfo.Creator,
		// Participant: taskInfo.Participant,
	})
}

func makeTaskResponse(task model.Task) TaskResp {
	taskResp := TaskResp{
		ShortTaskResp: ShortTaskResp{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description.String,
			Estimate:    task.Estimate.String,
			Status:      string(task.Status),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
		CreatorID: int(task.CreatorID.Int64),
		//ProjectID:     task.ProjectID,
	}
	if task.ParticipantID.Valid {
		taskResp.ParticipantID = int(task.ParticipantID.Int64)
	} else {
		taskResp.ParticipantIDNull = nil
	}
	return taskResp
}
func makeShortTaskResponse(task model.ShortTask) ShortTaskResp {
	shortTaskResp := ShortTaskResp{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description.String,
		Estimate:    task.Estimate.String,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if task.ParticipantID.Valid {
		shortTaskResp.ParticipantID = int(task.ParticipantID.Int64)
	} else {
		shortTaskResp.ParticipantIDNull = nil
	}
	return shortTaskResp
}
func makeTasksResponses(tasks []model.Task) []TaskResp {
	taskResponses := make([]TaskResp, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses, makeTaskResponse(task))
		// Creator:     taskInfo.Creator,
		// Participant: taskInfo.Participant,)
	}
	return taskResponses
}
func makeShortTasksResponses(tasks []model.ShortTask) []ShortTaskResp {
	taskResponses := make([]ShortTaskResp, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses, makeShortTaskResponse(task))
	}
	return taskResponses
}
