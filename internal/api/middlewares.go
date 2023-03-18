package api

import (
	"net/http"
	"strconv"
	"strings"

	"be-project-monitoring/internal/domain"
	"be-project-monitoring/internal/domain/model"
	ierr "be-project-monitoring/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) authMiddleware(toAllow ...model.UserRole) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := getTokenFromHeader(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}

		ctx := c.Request.Context()
		if err := s.svc.VerifyToken(ctx, token, toAllow...); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}

		id, err := s.svc.GetUserIDFromToken(ctx, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}

		c.Set(string(domain.UserIDCtx), id)
	}
}

func (s *Server) updateMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		token, err := getTokenFromHeader(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}

		id, err := s.svc.GetUserIDFromToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}

		if id != updatedUser.ID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: ierr.ErrAccessDeniedAnotherUser.Error()})
			return
		}
	}
}
func (s *Server) verifyParticipantMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			ctx          = c.Request.Context()
			userID       = c.MustGet(string(domain.UserIDCtx)).(uuid.UUID)
			projectID, _ = strconv.Atoi(c.Param("projectId"))
		)
		if _, err := s.svc.VerifyParticipant(ctx, userID, projectID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}
	}
}
func (s *Server) verifyParticipantRoleMiddleware(toAllow ...model.ParticipantRole) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			ctx       = c.Request.Context()
			userID    = c.MustGet(string(domain.UserIDCtx)).(uuid.UUID)
			projectID = c.MustGet(string(domain.ProjectIDCtx)).(int)
		)

		if err := s.svc.VerifyParticipantRole(ctx, userID, projectID, toAllow...); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}
	}
}
func getTokenFromHeader(tokenHeader string) (string, error) {
	if tokenHeader == "" {
		return "", ierr.ErrTokenHeaderIsEmpty
	}
	parts := strings.Split(strings.TrimSpace(tokenHeader), " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ierr.ErrInvalidToken
	}
	return parts[1], nil
}
