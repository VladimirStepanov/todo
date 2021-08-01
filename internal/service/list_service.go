package service

import "github.com/VladimirStepanov/todo-app/internal/models"

type ListService struct {
	repo models.ListRepository
}

func NewListService(repo models.ListRepository) models.ListService {
	return &ListService{
		repo: repo,
	}
}

func (ls *ListService) Create(title, description string, userID int64) (int64, error) {
	return 0, nil
}

func (ls *ListService) GrantRole(listID, fromUser, toUserID int64, role bool) error {
	return nil
}

func (ls *ListService) GetListByID(listID, userID int64) (*models.List, error) {
	return nil, nil
}

func (ls *ListService) GetUserLists(userID int64) ([]*models.List, error) {
	return nil, nil
}

func (ls *ListService) Delete(listID, userID int64) error {
	return nil
}

func (ls *ListService) Update(userID int64, list *models.List) error {
	return nil
}