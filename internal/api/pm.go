package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"be-project-monitoring/internal/domain/model"
)

type (
	// createProjectReq struct {
	// 	Name        string    `json:"name"`
	// 	Description string    `json:"description"`
	// 	ActiveTo    time.Time `json:"active_to"`
	// 	PhotoURL    string    `json:"photo_url"`
	// 	RepoURL     string    `json:"repo_url"`
	// }

	projectResp struct {
		Project  *model.Project `json:"project,omitempty"` //vopros po omitu
		//nado li escho chto to?
	}
)

func (s *Server) createProject(c *gin.Context) {
	newProject := &model.Project{}
	if err := json.NewDecoder(c.Request.Body).Decode(newProject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	// newUser := &model.User{
	// 	Role:           userReq.Role,
	// 	Email:          userReq.Email,
	// 	Username:       userReq.Username,
	// 	FirstName:      userReq.FirstName,
	// 	LastName:       userReq.LastName,
	// 	Group:          userReq.Group,
	// 	GithubUsername: userReq.GithubUsername,
	// 	HashedPassword: hashPass(userReq.Password),
	// }

	project, err := s.projSvc.CreateProject(c.Request.Context(), newProject)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, projectResp{
		Project: project,
	})

}
