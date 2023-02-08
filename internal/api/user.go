package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"be-project-monitoring/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	GetUserReq struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		Username    string    `json:"username"`
		IsOnProject bool      `json:"is_on_project"` //описать значения для парс.бул
		ProjectID   int       `json:"project_id"`
		Offset      int       `json:"offset"`
		Limit       int       `json:"limit"`
	}

	getUserResp struct {
		Users []model.User `json:"users"`
		Count int          `json:"count"`
	}
	getShortUserResp struct {
		Users []model.ShortUser `json:"users"`
		Count int               `json:"count"`
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
)

func (s *Server) getFullUsers(c *gin.Context) {
	userReq := &GetUserReq{}
	userReq.Email = c.Query("email")
	userReq.Username = c.Query("username")
	userReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	userReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	users, count, err := s.svc.GetFullUsers(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, getUserResp{Users: users, Count: count})
}
func (s *Server) getPartialUsers(c *gin.Context) {
	userReq := &GetUserReq{}
	userReq.Email = c.Query("email")
	userReq.Username = c.Query("username")
	userReq.IsOnProject, _ = strconv.ParseBool(c.Query("is_on_project")) //мб сразу true выставлять здесь в паршиалЮзерс?
	userReq.ProjectID, _ = strconv.Atoi(c.Query("project_id"))
	userReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	userReq.Limit, _ = strconv.Atoi(c.Query("limit"))

	users, count, err := s.svc.GetPartialUsers(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, getShortUserResp{Users: users, Count: count})
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

func (s *Server) deleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	if err := s.svc.DeleteUser(c.Request.Context(), userID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) getUserProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	userProfile, err := s.svc.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}
