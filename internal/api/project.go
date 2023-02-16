package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"be-project-monitoring/internal/domain"
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
	CreateProjectResp struct {
		Project     *projectResp
		Participant partcipantResp
	}
	partcipantResp struct {
		ID        int             `json:"id"`
		Role      string          `json:"participant_role"`
		ProjectID int             `json:"project_id"`
		User      model.ShortUser `json:"user"`
	}
	projectResp struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description,omitempty"`
		PhotoURL    string    `json:"photo_url,omitempty"`
		ReportURL   string    `json:"report_url,omitempty"`
		ReportName  string    `json:"report_name,omitempty"`
		RepoURL     string    `json:"repo_url,omitempty"`
		ActiveTo    time.Time `json:"active_to,omitempty"`
	}
	GetProjectsReq struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"`
		Limit  int    `json:"limit"`
	}

	getProjectResp struct {
		Projects []projectResp
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
	AddParticipantReq struct {
		Role      string    `json:"role"`
		UserID    uuid.UUID `json:"user_id"`
		ProjectID int       `json:"project_id"`
	}
)

func (s *Server) createProject(c *gin.Context) {
	projectReq := &CreateProjectReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(projectReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	project, err := s.svc.CreateProject(c.Request.Context(), projectReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	participant, err := s.svc.AddParticipant(c.Request.Context(), &AddParticipantReq{
		Role:      string(model.RoleOwner),
		UserID:    c.MustGet(string(domain.UserIDCtx)).(uuid.UUID), //uuid.MustParse(c.MustGet(string(domain.UserIDCtx))), //как лучше?
		ProjectID: project.ID,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	fmt.Println(participant.ShortUser.ID)
	c.JSON(http.StatusCreated, CreateProjectResp{
		Project: makeProjectResponse(*project),
		Participant: partcipantResp{
			ID:        participant.ID,
			Role:      string(participant.Role),
			ProjectID: participant.ProjectID,
			User:      participant.ShortUser,
		}})

}

func (s *Server) getProjects(c *gin.Context) {
	projReq := &GetProjectsReq{}

	projReq.Name = c.Query("name")
	projReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	projReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	projects, count, err := s.svc.GetProjects(c.Request.Context(), projReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, getProjectResp{
		Projects: makeProjectResponses(projects),
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
	c.JSON(http.StatusOK, makeProjectResponse(*project))
}

func (s *Server) deleteProject(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	if err := s.svc.DeleteProject(c.Request.Context(), userID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addParticipant(c *gin.Context) {
	req := &AddParticipantReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	_, err := s.svc.AddParticipant(c.Request.Context(), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	participants, err := s.svc.GetParticipants(c.Request.Context(), req.ProjectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"participants": participants})
}

func (s *Server) deleteParticipant(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("user_id"))
	projectID, _ := strconv.Atoi(c.Param("id"))
	if err := s.svc.DeleteParticipant(c.Request.Context(), userID, projectID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
func (s *Server) getProjectInfo(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	projectInfo, err := s.svc.GetProjectInfo(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, struct {
		Project *projectResp
		Users   []model.ShortUser
		Tasks   []taskResp
	}{
		Project: makeProjectResponse(projectInfo.Project),
		Users:   projectInfo.Users,
		Tasks:   makeShortTasksResponses(projectInfo.Tasks),
	})
}
func makeProjectResponse(project model.Project) *projectResp {
	return &projectResp{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description.String,
		PhotoURL:    project.PhotoURL.String,
		ReportURL:   project.ReportURL.String,
		ReportName:  project.ReportName.String,
		RepoURL:     project.RepoURL.String,
		ActiveTo:    project.ActiveTo,
	}
}
func makeProjectResponses(projects []model.Project) []projectResp {
	var projectResponses []projectResp
	for _, project := range projects {
		projectResponses = append(projectResponses, *makeProjectResponse(project))
	}
	return projectResponses
}
func makeShortProjectResponses(projects []model.ShortProject) []projectResp {
	var projectResponses []projectResp
	for _, project := range projects {
		projectResponses = append(projectResponses,
			*makeProjectResponse(model.Project{
				ShortProject: project,
			}))
	}
	return projectResponses
}
