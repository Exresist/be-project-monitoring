package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	getProjReq struct {
		name string `json :"name"`
	}
	getProjResp struct {
		ID          int       `json :"id"`
		Name        string    `json :"name"`
		Description string    `json :"description"`
		PhotoURL    string    `json :"photo_url"`
		ReportURL   string    `json :"report_url"`
		ReportName  string    `json :"report_name"`
		RepoURL     string    `json :"repo_url"`
		ActiveTo    time.Time `json :"active_to"`
	}
)

func (s *Server) getProjects(c *gin.Context) {
	projReq := &getProjReq{}
	projReq.name = c.Query("name")

	projects, err := s.projSvc.GetProjects(c.Request.Context(), projReq.name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}
