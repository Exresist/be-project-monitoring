package api

import (
	"be-project-monitoring/internal/domain/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	GetProjReq struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"` //сколько записей опустить
		Limit  int    `json:"limit"`  //сколько записей подать
	}
	getProjResp struct {
		Projects []*model.Project
		Count    int
	}
)

func (s *Server) getProjects(c *gin.Context) {
	projReq := &GetProjReq{}
	projReq.Name = c.Query("name")
	projReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	projReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	projects, count, err := s.projSvc.GetProjects(c.Request.Context(), projReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, getProjResp{
		Projects: projects,
		Count:    count,
	})
}
