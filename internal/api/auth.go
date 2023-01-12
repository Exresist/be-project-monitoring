package api

import (
	"be-project-monitoring/internal/domain/model"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	CreateUserReq struct {
		Email          string `json:"email"`
		Username       string `json:"username"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Group          string `json:"group"`
		GithubUsername string `json:"github_username"`
		Password       string `json:"password"`
		Role           string `json:"role"`
	}
	authUserReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	userResp struct {
		User  *model.User `json:"user,omitempty"`
		Token string      `json:"token"`
	}
)

var errField = "error"

func (s *Server) register(c *gin.Context) {
	userReq := &CreateUserReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	user, token, err := s.svc.CreateUser(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userResp{
		User:  user,
		Token: token,
	})
}

func (s *Server) auth(c *gin.Context) {
	userReq := &authUserReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	token, err := s.svc.AuthUser(c.Request.Context(), userReq.Username, userReq.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResp{Token: token})
}
