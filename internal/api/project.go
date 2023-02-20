package api

import (
	"encoding/json"
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
		ActiveTo    time.Time `json:"dueDate"`
		PhotoURL    string    `json:"avatar"`
	}
	CreateProjectResp struct {
		Project     *ProjectResp
		Participant partcipantResp
	}

	ProjectResp struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"avatar"`
		ReportURL   string    `json:"reportUrl"`
		ReportName  string    `json:"reportName"`
		RepoURL     string    `json:"repo"`
		ActiveTo    time.Time `json:"dueDate"`
	}
	shortProjectResp struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"avatar"`
		ActiveTo    time.Time `json:"dueDate"`
	}
	projectWithParticipantsResp struct {
		ProjectResp
		Participants []shortPartcipantResp `json:"participants"`
	}
	GetProjectsReq struct {
		Name   string
		Offset int
		Limit  int
	}
	getProjectResp struct {
		Projects []ProjectResp
		Count    int
	}

	UpdateProjectReq struct {
		ID          int       `json:"id"`
		Name        *string   `json:"name"`
		Description *string   `json:"description"`
		PhotoURL    *string   `json:"avatar"`
		ReportURL   *string   `json:"reportUrl"`
		ReportName  *string   `json:"reportName"`
		RepoURL     *string   `json:"repo"`
		ActiveTo    time.Time `json:"dueDate"`
	}

	projectInfoResp struct {
		ID           int                 `json:"id"`
		Name         string              `json:"name"`
		Description  string              `json:"description"`
		PhotoURL     string              `json:"avatar"`
		ReportURL    string              `json:"reportUrl"`
		ReportName   string              `json:"reportName"`
		RepoURL      string              `json:"repo"`
		ActiveTo     time.Time           `json:"dueDate"`
		Participants []model.Participant `json:"participants"`
		Tasks        []taskResp          `json:"tasks"`
	}
)

var (
	updatedProject   *UpdateProjectReq
	deletedProjectID *int
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

func (s *Server) parseBodyToUpdatedProject(c *gin.Context) {
	updatedProject = &UpdateProjectReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(updatedProject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.Set(string(domain.ProjectIDCtx), updatedProject.ID)
}
func (s *Server) updateProject(c *gin.Context) {
	project, err := s.svc.UpdateProject(c.Request.Context(), updatedProject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, makeProjectResponse(*project))
}
func (s *Server) parseBodyToDeletedProject(c *gin.Context) {
	if err := json.NewDecoder(c.Request.Body).Decode(deletedProjectID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.Set(string(domain.ProjectIDCtx), deletedProjectID)
}
func (s *Server) deleteProject(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("projectId"))
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

func (s *Server) getProjectInfo(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	projectInfo, err := s.svc.GetProjectInfo(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, projectInfoResp{
		ID:           projectInfo.Project.ID,
		Name:         projectInfo.Name,
		Description:  projectInfo.Description.String,
		PhotoURL:     projectInfo.PhotoURL.String,
		ReportURL:    projectInfo.ReportURL.String,
		ReportName:   projectInfo.ReportName.String,
		RepoURL:      projectInfo.RepoURL.String,
		ActiveTo:     projectInfo.ActiveTo,
		Participants: projectInfo.Participants,
		Tasks:        makeShortTasksResponses(projectInfo.Tasks),
	})
}

func makeProjectResponse(project model.Project) *ProjectResp {
	return &ProjectResp{
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
func makeProjectResponses(projects []model.Project) []ProjectResp {
	projectResponses := make([]ProjectResp, 0, len(projects))
	for _, project := range projects {
		projectResponses = append(projectResponses, *makeProjectResponse(project))
	}
	return projectResponses
}
func makeShortProjectResponse(project model.ShortProject) *ProjectResp {
	return &ProjectResp{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description.String,
		ActiveTo:    project.ActiveTo,
	}

}
func makeShortProjectResponses(projects []model.ShortProject) []ProjectResp {
	projectResponses := make([]ProjectResp, 0, len(projects))
	for _, project := range projects {
		projectResponses = append(projectResponses,
			*makeProjectResponse(model.Project{
				ShortProject: project,
			}))
	}
	return projectResponses
}
