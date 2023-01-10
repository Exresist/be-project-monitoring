package api

import (
	"be-project-monitoring/internal/domain/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	createProjectReq struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ActiveTo    time.Time `json:"active_to"`
		PhotoURL    string    `json:"photo_url"`
		RepoURL     string    `json:"repo_url"`
	}
)

func (s *Server) createProject(c *gin.Context) {
	project := &model.Project{}
	if err := json.NewDecoder(c.Request.Body).Decode(project); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

}
