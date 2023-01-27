package api

import (
	"be-project-monitoring/internal/domain/model"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	CreateTaskReq struct {
		//ne ponyatno s participant id
		Name              string `json:"name"`
		Description       string `json:"description"`
		SuggestedEstimate int    `json:"suggested_estimate"`
		CreatorID         int    `json:"creator_id"`
		//ParticipantID int `json:"paticipant_id"`????
		//Status string `json:"status"`
	}
	GetTaskReq struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		CreatorID int       `json:"creator_id"`
		CreatedAt time.Time `json:"created_at"` //ne ponimayu kak za poslednuu nedelu
		Offset    int       `json:"offset"`
		Limit     int       `json:"limit"`
	}
	getTaskResp struct {
		Tasks []model.Task
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
		ParticipantID     *int      `json:"participant_id"`
		//CreatorID 		*int      `json:"creator_id"`?????????????????????/
	}
	DeleteTaskReq struct {
		ID int `json:"id"`
	}
)

func (s *Server) getTasks(ctx *gin.Context) {
	taskReq := &GetTaskReq{}

	taskReq.ID, _ = strconv.Atoi(ctx.Query("id"))
	taskReq.Name = ctx.Query("name")
	taskReq.Offset, _ = strconv.Atoi(ctx.Query("offset"))
	taskReq.Limit, _ = strconv.Atoi(ctx.Query("limit"))

	tasks, count, err := s.svc.GetTasks(ctx, taskReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, getTaskResp{
		Tasks: tasks,
		Count: count,
	})

}
func (s *Server) createTask(ctx *gin.Context) {
	taskReq := &CreateTaskReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(taskReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	task, err := s.svc.CreateTask(ctx, taskReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, task)
}
func (s *Server) updateTask(ctx *gin.Context) {
	taskReq := &UpdateTaskReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(taskReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	task, err := s.svc.UpdateTask(ctx, taskReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, task)
}
func (s *Server) deleteTask(ctx *gin.Context) {
	taskReq := &DeleteTaskReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(taskReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err := s.svc.DeleteTask(ctx, taskReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
