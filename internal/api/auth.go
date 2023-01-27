package api

import (
	"encoding/json"
	"net/http"

	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"

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

func (s *Server) register(ctx *gin.Context) {
	userReq := &CreateUserReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(userReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	if !s.svc.FindGithubUser(ctx.Request.Context(), userReq.GithubUsername) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: ierr.ErrGithubUserNotFound.Error()})
	}

	user, token, err := s.svc.CreateUser(ctx.Request.Context(), userReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, userResp{
		User:  user,
		Token: token,
	})
}

func (s *Server) auth(ctx *gin.Context) {
	userReq := &authUserReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(userReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	token, err := s.svc.AuthUser(ctx.Request.Context(), userReq.Username, userReq.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userResp{Token: token})
}
