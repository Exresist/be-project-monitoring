package service

import (
	"be-project-monitoring/internal/domain/model"
	"context"
)

func (s *service) AddParticipant(ctx context.Context, participant *model.Participant) ([]model.Participant, error) {
	return s.repo.AddParticipant(ctx, participant)
}

func (s *service) GetParticipants(ctx context.Context, projectID int) ([]model.Participant, error) {
	return s.repo.GetParticipants(ctx, projectID)
}
