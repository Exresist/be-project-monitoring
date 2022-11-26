package api

import (
	"be-project-monitoring/internal/domain/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	GetUserReq struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Offset   int    `json:"offset"` //сколько записей опустить
		Limit    int    `json:"limit"`  //сколько записей подать
	}

	getUserResp struct {
		Users []*model.User `json:"users"`
		Count int           `json:"count"`
	}
)

func (s *Server) getUsers(c *gin.Context) {
	userReq := &GetUserReq{}

	userReq.Email = c.Query("email")
	userReq.Username = c.Query("username")
	userReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	userReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	users, count, err := s.userSvc.GetUsers(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, getUserResp{Users: users, Count: count})

}
