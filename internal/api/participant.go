package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	AddParticipantReq struct {
		Role      string    `json:"role"`
		UserID    uuid.UUID `json:"user_id"`
		ProjectID int       `json:"project_id"`
	}
)
var (
	deletedParticipantID int

)

func (s *Server) addParticipant(c *gin.Context) {
	req := &AddParticipantReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	_, err := s.svc.AddParticipant(c.Request.Context(), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	participants, err := s.svc.GetParticipants(c.Request.Context(), req.ProjectID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"participants": participants})
}

func (s *Server) deleteParticipant(c *gin.Context) {
	// userID, _ := uuid.Parse(c.Param("user_id"))
	// projectID, _ := strconv.Atoi(c.Param("project_id"))
	if err := s.svc.DeleteParticipant(c.Request.Context(), deletedParticipantID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
