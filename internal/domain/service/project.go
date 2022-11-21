package service

import (
	"be-project-monitoring/internal/domain/model"
	"context"
)

func (s *projectService) GetProjects(ctx context.Context, name string) ([]*model.Project, error) {

	return nil, nil
	// filter := domain.NewUserFilter().ByUsernames(userReq.Username).ByEmails(userReq.Email)
	// filter.Limit = uint64(userReq.Limit)
	// filter.Offset = uint64(userReq.Offset)
	// count, err := s.projectStore.GetCountByFilter(ctx, filter)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// users, err := s.userStore.GetUsers(ctx, filter)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// return users, count, nil
}

func (s *projectService) CreateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *projectService) UpdateProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	return nil, nil
}

func (s *projectService) DeleteProject(ctx context.Context, project *model.Project) error {
	return nil
}
