package api

import (
	"be-project-monitoring/internal/domain/model"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

	createProjectResp struct {
		*model.Project
	}

	GetProjReq struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"` //сколько записей опустить
		Limit  int    `json:"limit"`  //сколько записей подать
	}

	getProjResp struct {
		Projects []*model.Project
		Count    int
	}

	addParticipantReq struct {
		Role      int       `json:"role"`
		UserID    uuid.UUID `json:"user_id"`
		ProjectID int       `json:"project_id"`
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
	projReq := &GetProjReq{}
	projReq.Name = c.Query("name")
	projReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	projReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	projects, count, err := s.svc.GetProjects(c.Request.Context(), projReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, getProjResp{
		Projects: projects,
		Count:    count,
	})
}

func (s *Server) addParticipant(c *gin.Context) {
	req := &addParticipantReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	participants, err := s.svc.AddParticipant(c.Request.Context(), &model.Participant{
		Role:      model.ParticipantRole(req.Role),
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": participants})
}
