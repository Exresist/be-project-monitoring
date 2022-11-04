package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be-project-monitoring/internal/domain/model"
)

func (s *Server) authMiddleware(toAllow ...model.UserRole) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token := c.Request.Header.Get("Authorization")

		err := s.userSvc.VerifyToken(ctx, token, toAllow...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{errField: err.Error()})
			return
		}
	}
}
