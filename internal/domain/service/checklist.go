package service

import (
	"context"

	"be-project-monitoring/internal/domain/model"
)

func (s *service) GetProjectChecklist(ctx context.Context, id int) ([]model.Checklist, error) {
	return s.repo.GetProjectChecklist(ctx, id)
}

func (s *service) AddProjectChecklist(ctx context.Context, id int, checklist []model.Checklist) ([]model.Checklist, error) {
	return s.repo.AddProjectChecklist(ctx, id, checklist)
}
func (s *service) UpdateProjectChecklist(ctx context.Context, id int, checklist *model.Checklist) ([]model.Checklist, error) {
	return s.repo.UpdateProjectChecklist(ctx, id, checklist)
}
func (s *service) DeleteProjectChecklist(ctx context.Context, id int, itemID int) ([]model.Checklist, error) {
	return s.repo.DeleteProjectChecklist(ctx, id, itemID)
}
