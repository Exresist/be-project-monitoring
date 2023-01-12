package api

import (
	"be-project-monitoring/internal/domain/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	GetUserReq struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
	}

	getUserResp struct {
		Users []model.User `json:"users"`
		Count int          `json:"count"`
	}

	UpdateUserReq struct {
		ID             uuid.UUID `json:"id"`
		Role           string    `json:"role"`
		Username       string    `json:"username"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Group          string    `json:"group"`
		GithubUsername string    `json:"github_username"`
		Password       string    `json:"password"`
	}
)

func (s *Server) getUsers(c *gin.Context) {
	userReq := &GetUserReq{}

	userReq.Email = c.Query("email")
	userReq.Username = c.Query("username")
	userReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	userReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	users, count, err := s.svc.GetUsers(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, getUserResp{Users: users, Count: count})
}

func (s *Server) updateUser(c *gin.Context) {
	userReq := &UpdateUserReq{}

	if err := json.NewDecoder(c.Request.Body).Decode(userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	user, err := s.svc.UpdateUser(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// func (s *Server) updateUserRole(c *gin.Context) {
// 	userReq := &UpdateUserReq{}
// 	if err := json.NewDecoder(c.Request.Body).Decode(userReq); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
// 		return
// 	}

// 	user, err := s.svc.UpdateUser(c.Request.Context(), userReq)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, user)
// }
