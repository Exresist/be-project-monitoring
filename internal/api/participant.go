package api

import (
	"be-project-monitoring/internal/domain/model"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	AddParticipantReq struct {
		Role      string    `json:"role"`
		UserID    uuid.UUID `json:"userId"`
		ProjectID int       `json:"projectId"`
	}
	partcipantResp struct {
		ID        int             `json:"id"`
		Role      string          `json:"role"`
		ProjectID int             `json:"projectId,omitempty"`
		User      model.ShortUser `json:"user,omitempty"`
	}
	shortPartcipantResp struct {
		ID   int    `json:"id"`
		Role string `json:"role"`
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

func makeShortParticipantResponse(participant model.Participant) *shortPartcipantResp {
	return &shortPartcipantResp{
		ID:        participant.ID,
		Role:      string(participant.Role),
	}
}
func makeShortParticipantResponses(participants []model.Participant) []shortPartcipantResp {
	participantResponses := make([]shortPartcipantResp, 0, len(participants))
	for _, participant := range participants {
		participantResponses = append(participantResponses, *makeShortParticipantResponse(participant))
	}
	return participantResponses
}
