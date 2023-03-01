package api

import (
	"be-project-monitoring/internal/domain"
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
		ActiveTo    time.Time `json:"dueDate"`
		PhotoURL    string    `json:"avatar"`
	}
	CreateProjectResp struct {
		ProjectResp
		ParticipantResp `json:"participants"`
	}

	ProjectResp struct {
		ShortProjectResp
		ReportURL  string `json:"reportUrl"`
		ReportName string `json:"reportName"`
		RepoURL    string `json:"repo"`
	}
	ShortProjectResp struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"avatar"`
		ActiveTo    time.Time `json:"dueDate"`
	}
	projectWithParticipantsResp struct {
		ProjectResp
		Participants []ParticipantResp `json:"participants"`
	}
	projectWithShortParticipantsResp struct {
		ShortProjectResp
		Participants []shortPartcipantResp `json:"participants"`
	}
	GetProjectsReq struct {
		SearchText string
		// Offset int
		// Limit  int
	}
	getProjectsResp struct {
		Projects []projectWithParticipantsResp
		//Projects []projectWithParticipantsResp `json:"projects"`
		//Count    int                           `json:"count"`
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
		ProjectResp
		Participants []model.Participant `json:"participants"`
		Tasks        []ShortTaskResp     `json:"tasks"`
	}

	commitsInfoResp struct {
		GithubUsername string `json:"githubUsername"`
		Username       string `json:"username"`
		Total          int    `json:"total"`
		TotalHours     int    `json:"totalHours"`
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
	participant, err := s.svc.AddParticipant(c.Request.Context(), true, &AddedParticipant{
		Role:      string(model.RoleOwner),
		UserID:    c.MustGet(string(domain.UserIDCtx)).(uuid.UUID), //uuid.MustParse(c.MustGet(string(domain.UserIDCtx))), //как лучше?
		ProjectID: project.ID,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateProjectResp{
		ProjectResp: makeProjectResponse(*project),
		ParticipantResp: ParticipantResp{
			ID:        participant.ID,
			Role:      string(participant.Role),
			ProjectID: participant.ProjectID,
			User:      participant.ShortUser,
		}})
}

func (s *Server) getProjects(c *gin.Context) {
	projReq := &GetProjectsReq{}
	projReq.SearchText = c.Query("searchParam")
	// projReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	// projReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	projects, _, err := s.svc.GetProjects(c.Request.Context(), projReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	projectsResp := make([]projectWithParticipantsResp, 0)
	for _, v := range projects {
		participants, err := s.svc.GetParticipants(c.Request.Context(), v.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
			return
		}
		projectsResp = append(projectsResp, projectWithParticipantsResp{
			ProjectResp:  makeProjectResponse(v),
			Participants: makeParticipantResponses(participants),
		})
	}
	// c.JSON(http.StatusOK, getProjectsResp{
	// 	Projects: projectsResp,
	// 	Count:    count,
	// })
	c.JSON(http.StatusOK, projectsResp)
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
	participants, err := s.svc.GetParticipants(c.Request.Context(), project.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectWithParticipantsResp{
		ProjectResp:  makeProjectResponse(*project),
		Participants: makeParticipantResponses(participants),
	})
}
func (s *Server) parseBodyToDeletedProject(c *gin.Context) {
	if err := json.NewDecoder(c.Request.Body).Decode(deletedProjectID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.Set(string(domain.ProjectIDCtx), deletedProjectID)
}
func (s *Server) deleteProject(c *gin.Context) {
	// projectID, err := strconv.Atoi(c.Param("projectId"))
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
	// 	return
	// }

	if err := s.svc.DeleteProject(c.Request.Context(), *deletedProjectID); err != nil {
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
	s.sendProjectInfoResponse(c, projectID)
}
func (s *Server) getProjectCommits(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	commitsInfo, err := s.svc.GetProjectCommits(c.Request.Context(), projectID)
	if err != nil {
		return
	}
	resp := make([]commitsInfoResp, 0, len(commitsInfo))
	for _, info := range commitsInfo {
		resp = append(resp,
			commitsInfoResp{
				GithubUsername: info.GithubUsername,
				Username:       info.Username,
				Total:          info.Total,
				TotalHours:     int(info.LastCommitDate.Sub(info.FirstCommitDate).Hours()),
			},
		)
	}

	c.JSON(http.StatusOK, resp)
}
func (s *Server) sendProjectInfoResponse(c *gin.Context, projectID int) {
	projectInfo, err := s.svc.GetProjectInfo(c.Request.Context(), projectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	shortTasksResponse := make([]ShortTaskResp, 0)
	for _, v := range projectInfo.Tasks {
		shortTasksResponse = append(shortTasksResponse, makeShortTaskResponse(v.ShortTask))
	}
	c.JSON(http.StatusOK, projectInfoResp{
		ProjectResp: ProjectResp{
			ShortProjectResp: ShortProjectResp{
				ID:          projectInfo.Project.ID,
				Name:        projectInfo.Name,
				Description: projectInfo.Description.String,
				PhotoURL:    projectInfo.PhotoURL.String,
				ActiveTo:    projectInfo.ActiveTo,
			},
			ReportURL:  projectInfo.ReportURL.String,
			ReportName: projectInfo.ReportName.String,
			RepoURL:    projectInfo.RepoURL.String,
		},
		Participants: projectInfo.Participants,
		Tasks:        shortTasksResponse,
	})
}
func makeProjectResponse(project model.Project) ProjectResp {
	return ProjectResp{
		ShortProjectResp: ShortProjectResp{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description.String,
			PhotoURL:    project.PhotoURL.String,
			ActiveTo:    project.ActiveTo,
		},
		ReportURL:  project.ReportURL.String,
		ReportName: project.ReportName.String,
		RepoURL:    project.RepoURL.String,
	}
}
func makeProjectResponses(projects []model.Project) []ProjectResp {
	projectResponses := make([]ProjectResp, 0, len(projects))
	for _, project := range projects {
		projectResponses = append(projectResponses, makeProjectResponse(project))
	}
	return projectResponses
}
func makeShortProjectResponse(project model.ShortProject) ShortProjectResp {
	return ShortProjectResp{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description.String,
		PhotoURL:    project.PhotoURL.String,
		ActiveTo:    project.ActiveTo,
	}

}
func makeShortProjectResponses(projects []model.ShortProject) []ShortProjectResp {
	projectResponses := make([]ShortProjectResp, 0, len(projects))
	for _, project := range projects {
		projectResponses = append(projectResponses,
			makeShortProjectResponse(project))
	}
	return projectResponses
}
