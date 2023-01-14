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
	CreateProjectReq struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ActiveTo    time.Time `json:"active_to"`
		PhotoURL    string    `json:"photo_url"`
	}

	createProjectResp struct {
		*model.Project
	}

	GetProjectReq struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"` //сколько записей опустить
		Limit  int    `json:"limit"`  //сколько записей подать
	}

	getProjectResp struct {
		Projects []model.Project
		Count    int
	}

	UpdateProjectReq struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"photo_url"`
		ReportURL   string    `json:"report_url"`
		ReportName  string    `json:"report_name"`
		RepoURL     string    `json:"repo_url"`
		ActiveTo    time.Time `json:"active_to"`
	}
)

func (s *Server) createProject(c *gin.Context) {
	newProject := &CreateProjectReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(newProject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	project, err := s.svc.CreateProject(c.Request.Context(), newProject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createProjectResp{
		project,
	})

}

func (s *Server) getProjects(c *gin.Context) {
	projReq := &GetProjectReq{}
	projReq.Name = c.Query("name")
	projReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	projReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	projects, count, err := s.svc.GetProjects(c.Request.Context(), projReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, getProjectResp{
		Projects: projects,
		Count:    count,
	})
}

func (s *Server) updateProject(c *gin.Context) {
	newProj := &UpdateProjectReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(newProj); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	project, err := s.svc.UpdateProject(c.Request.Context(), newProj)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}
