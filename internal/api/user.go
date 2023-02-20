package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"be-project-monitoring/internal/domain"
	"be-project-monitoring/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	GetUserReq struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		Username    string    `json:"username"`
		IsOnProject bool      `json:"isOnProject"` //описать значения для парс.бул
		ProjectID   int       `json:"projectId"`
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
		FirstName      *string   `json:"firstName"`
		LastName       *string   `json:"lastName"`
		Group          *string   `json:"group"`
		GithubUsername *string   `json:"ghUsername"`
		Password       *string   `json:"password"`
	}
	GetUserResp struct {
		ID             uuid.UUID     `json:"id"`
		Role           string        `json:"role"`
		Email          string        `json:"email"`
		Username       string        `json:"username"`
		FirstName      string        `json:"firstName"`
		LastName       string        `json:"lastName"`
		ColorCode      string        `json:"avatarColor"`
		Group          string        `json:"group"`
		GithubUsername string        `json:"ghUsername"`
		Projects       []projectResp `json:"projects"`
	}
)

var (
	updatedUser *UpdateUserReq
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
	userReq.IsOnProject, _ = strconv.ParseBool(c.Query("isOnProject")) //мб сразу true выставлять здесь в паршиалЮзерс?
	userReq.ProjectID, _ = strconv.Atoi(c.Query("projectId"))
	userReq.Offset, _ = strconv.Atoi(c.Query("offset"))
	userReq.Limit, _ = strconv.Atoi(c.Query("limit"))
	users, count, err := s.svc.GetPartialUsers(c.Request.Context(), userReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, getShortUserResp{Users: users, Count: count})
}
func (s *Server) parseBodyToUpdatedUser(c *gin.Context) {
	updatedUser = &UpdateUserReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(updatedUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.Set(string(domain.UserIDCtx), updatedUser.ID)
}
func (s *Server) updateUser(c *gin.Context) {

	user, err := s.svc.UpdateUser(c.Request.Context(), updatedUser)
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
	c.JSON(http.StatusOK, struct {
		User         model.ShortUser `json:"user"`
		UserProjects []projectResp
	}{
		User:         userProfile.ShortUser,
		UserProjects: makeShortProjectResponses(userProfile.UserProjects),
	})
}
func (s *Server) getUserProfileFromToken(c *gin.Context) {
	userID := c.MustGet(string(domain.UserIDCtx)).(uuid.UUID)
	userProfile, err := s.svc.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetUserResp{
		ID:             userProfile.ShortUser.ID,
		Email:          userProfile.Email,
		Username:       userProfile.ShortUser.Username,
		FirstName:      userProfile.ShortUser.FirstName,
		LastName:       userProfile.ShortUser.LastName,
		Group:          userProfile.ShortUser.Group,
		GithubUsername: userProfile.ShortUser.GithubUsername,
		ColorCode:      userProfile.ShortUser.ColorCode,
		Role:           string(userProfile.ShortUser.Role),
		Projects:       makeShortProjectResponses(userProfile.UserProjects),
	})
}

func (s *Server) getUserProjects(c *gin.Context) {
	userID := c.MustGet(string(domain.UserIDCtx)).(uuid.UUID)

	userProfile, err := s.svc.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	projectResponses := make([]projectWithParticipantsResp, 0)
	for _, v := range userProfile.UserProjects {
		participants, err := s.svc.GetParticipants(c.Request.Context(), v.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
			return
		}
		projectResponses = append(projectResponses, projectWithParticipantsResp{
			participants: makeParticipantResponses(participants),
			projectResp:  *makeShortProjectResponse(v),
		})
	}
	c.JSON(http.StatusOK, projectResponses)
}
