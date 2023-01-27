package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"be-project-monitoring/internal/domain/model"

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
		Role           *string   `json:"role"`
		Username       *string   `json:"username"`
		FirstName      *string   `json:"first_name"`
		LastName       *string   `json:"last_name"`
		Group          *string   `json:"group"`
		GithubUsername *string   `json:"github_username"`
		Password       *string   `json:"password"`
	}
	deleteUserReq struct {
		ID uuid.UUID `json:"id"`
	}
	GetUserProfileResp struct {
		ID             uuid.UUID `json:"id"`
		ColorCode      string    `json:"color_code"`
		Email          string    `json:"email"`
		Role           string    `json:"role"`
		Username       string    `json:"username"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Group          string    `json:"group"`
		GithubUsername string    `json:"github_username"`
		UserProjects   []*UserProjectsResp
	}
	UserProjectsResp struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PhotoURL    string    `json:"photo_url"`
		ActiveTo    time.Time `json:"active_to"`
	}
)

func (s *Server) getUsers(ctx *gin.Context) {
	userReq := &GetUserReq{}

	userReq.Email = ctx.Query("email")
	userReq.Username = ctx.Query("username")
	userReq.Offset, _ = strconv.Atoi(ctx.Query("offset"))
	userReq.Limit, _ = strconv.Atoi(ctx.Query("limit"))

	users, count, err := s.svc.GetUsers(ctx.Request.Context(), userReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, getUserResp{Users: users, Count: count})
}

func (s *Server) updateUser(ctx *gin.Context) {
	userReq := &UpdateUserReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(userReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	user, err := s.svc.UpdateUser(ctx.Request.Context(), userReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (s *Server) deleteUser(ctx *gin.Context) {
	userReq := &deleteUserReq{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(userReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	if err := s.svc.DeleteUser(ctx.Request.Context(), userReq.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (s *Server) getUserProfile(ctx *gin.Context) {
	
	userProfileID, err := uuid.Parse(ctx.Query("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	userProfile, err := s.svc.GetUserProfile(ctx.Request.Context(), userProfileID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userProfile)
}
