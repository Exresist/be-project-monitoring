package api

import (
	"be-project-monitoring/internal/domain/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	// profileReq struct{
	// 	ID uuid.UUID
	// }
	profileResp struct {
		User     *model.User      `json:"user,omitempty"`
		Projects []*model.Project `json:"projects,omitempty"` //omiti
	}
)

func (s *Server) getProfile(c *gin.Context) {
	ID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{errField: err.Error()})
		return
	}

	profile := s.SERVICE.getprofile(c.Request.Context(), ID)

	c.JSON(http.StatusOK, profileResp{
		Profile: profile,
	})
}
