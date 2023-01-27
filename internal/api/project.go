package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"be-project-monitoring/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	CreateProjectReq struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ActiveTo    time.Time `json:"active_to"`
		PhotoURL    string    `json:"photo_url"`
	}

	GetProjectReq struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"`
		Limit  int    `json:"limit"`
	}

	getProjectResp struct {
		Projects []model.Project
		Count    int
	}

	UpdateProjectReq struct {
		ID          int       `json:"id"`
		Name        *string   `json:"name"`
		Description *string   `json:"description"`
		PhotoURL    *string   `json:"photo_url"`
		ReportURL   *string   `json:"report_url"`
		ReportName  *string   `json:"report_name"`
		RepoURL     *string   `json:"repo_url"`
		ActiveTo    time.Time `json:"active_to"`
	}
	DeleteProjectReq struct {
		ID int `json:"id"`
	}
	addParticipantReq struct {
		Role      int       `json:"role"`
		UserID    uuid.UUID `json:"user_id"`
		ProjectID int       `json:"project_id"`
	}
)

func (s *Server) createProject(ctx *gin.Context) {
	newProject := &CreateProjectReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(newProject); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	project, err := s.svc.CreateProject(ctx.Request.Context(), newProject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, project)

}

func (s *Server) getProjects(ctx *gin.Context) {
	projReq := &GetProjectReq{}
	projReq.Name = ctx.Query("name")
	projReq.Offset, _ = strconv.Atoi(ctx.Query("offset"))
	projReq.Limit, _ = strconv.Atoi(ctx.Query("limit"))

	projects, count, err := s.svc.GetProjects(ctx.Request.Context(), projReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, getProjectResp{
		Projects: projects,
		Count:    count,
	})
}

func (s *Server) updateProject(ctx *gin.Context) {
	newProj := &UpdateProjectReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(newProj); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	project, err := s.svc.UpdateProject(ctx.Request.Context(), newProj)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

func (s *Server) deleteProject(ctx *gin.Context) {
	projectReq := &DeleteProjectReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(projectReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	err := s.svc.DeleteProject(ctx.Request.Context(), projectReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (s *Server) addParticipant(ctx *gin.Context) {
	req := &addParticipantReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	participants, err := s.svc.AddParticipant(ctx.Request.Context(), &model.Participant{
		Role:      model.ParticipantRole(req.Role),
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"participants": participants})
}
