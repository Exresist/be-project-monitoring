package api

import (
	"net/http"
	"strings"

	"be-project-monitoring/internal/domain/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware(toAllow ...model.UserRole) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		parts := strings.Split(strings.TrimSpace(token), " ")
		if len(parts) < 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid token")
			return
		}

		err := s.svc.VerifyToken(ctx, parts[1], toAllow...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
			return
		}
	}
}

// func (s *Server) updateUserMiddleware(toAllow ...model.UserRole) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()

// 		err := s.svc.VerifyToken(ctx, parts[1], toAllow...)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
// 			return
// 		}

// 	}
// }