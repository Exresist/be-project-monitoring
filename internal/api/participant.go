package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"be-project-monitoring/internal/domain"
	"be-project-monitoring/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
)

type (
	AddedParticipant struct {
		Role      string    `json:"role"`
		UserID    uuid.UUID `json:"userId"`
		ProjectID int       `json:"projectId"`
	}

	ParticipantResp struct {
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
	addedParticipant  *AddedParticipant
	parsedParticipant *ParticipantResp

	sf = singleflight.Group{}
)

func (s *Server) parseBodyToAddedParticipant(c *gin.Context) {
	addedParticipant = &AddedParticipant{}

	if err := json.NewDecoder(c.Request.Body).Decode(addedParticipant); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	c.Set(string(domain.ProjectIDCtx), addedParticipant.ProjectID)
}
func (s *Server) addParticipant(c *gin.Context) {

	_, err, _ := sf.Do(
		fmt.Sprintf("%v-%v", addedParticipant.ProjectID, addedParticipant.UserID),
		func() (interface{}, error) {
			return s.svc.AddParticipant(c.Request.Context(), false, addedParticipant)
		})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	if _, err = s.svc.GetParticipants(c.Request.Context(), addedParticipant.ProjectID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}

	s.sendProjectInfoResponse(c, addedParticipant.ProjectID)
}
func (s *Server) parseBodyToParticipantResp(c *gin.Context) {
	parsedParticipant = &ParticipantResp{}
	if err := json.NewDecoder(c.Request.Body).Decode(parsedParticipant); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	participant, err := s.svc.GetParticipantByID(c.Request.Context(), parsedParticipant.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	parsedParticipant.ProjectID = participant.ProjectID
	parsedParticipant.User = participant.ShortUser
	c.Set(string(domain.ProjectIDCtx), participant.ProjectID)
}
func (s *Server) updateParticipant(c *gin.Context) {
	// userID, _ := uuid.Parse(c.Param("user_id"))
	// projectID, _ := strconv.Atoi(c.Param("project_id"))
	if _, err := s.svc.UpdateParticipantRole(c.Request.Context(), parsedParticipant); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	s.sendProjectInfoResponse(c, parsedParticipant.ProjectID)
	//c.JSON(http.StatusOK, participant)
}

//	func (s *Server) parseBodyToDeletedParticipant(c *gin.Context) {
//		deletedParticipant = &DeletedParticipant{}
//		if err := json.NewDecoder(c.Request.Body).Decode(deletedParticipant); err != nil {
//			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
//			return
//		}
//		participant, err := s.svc.GetParticipantByID(c.Request.Context(), deletedParticipant.ID)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
//			return
//		}
//		c.Set(string(domain.ProjectIDCtx), participant.ProjectID)
//	}
func (s *Server) deleteParticipant(c *gin.Context) {
	// userID, _ := uuid.Parse(c.Param("user_id"))
	// projectID, _ := strconv.Atoi(c.Param("project_id"))
	if err := s.svc.DeleteParticipant(c.Request.Context(), parsedParticipant.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errField: err.Error()})
		return
	}
	s.sendProjectInfoResponse(c, parsedParticipant.ProjectID)
	c.JSON(http.StatusOK, nil)
}
func makeParticipantResponse(participant model.Participant) *ParticipantResp {
	return &ParticipantResp{
		ID:   participant.ID,
		Role: string(participant.Role),
		User: participant.ShortUser,
	}
}
func castParticipants(participants []model.Participant) []ParticipantResp {
	participantResponses := make([]ParticipantResp, 0, len(participants))
	for _, participant := range participants {
		participantResponses = append(participantResponses, *makeParticipantResponse(participant))
	}
	return participantResponses
}
func makeShortParticipantResponse(participant model.Participant) *shortPartcipantResp {
	return &shortPartcipantResp{
		ID:   participant.ID,
		Role: string(participant.Role),
	}
}
func makeShortParticipantResponses(participants []model.Participant) []shortPartcipantResp {
	participantResponses := make([]shortPartcipantResp, 0, len(participants))
	for _, participant := range participants {
		participantResponses = append(participantResponses, *makeShortParticipantResponse(participant))
	}
	return participantResponses
}
