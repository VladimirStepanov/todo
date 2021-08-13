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
	return ls.repo.Create(title, description, userID)
}

func (ls *ListService) IsListAdmin(ListID, userID int64) error {
	return ls.repo.IsListAdmin(ListID, userID)
}

func (ls *ListService) EditRole(listID, userID int64, role bool) error {
	return ls.repo.EditRole(listID, userID, role)
}

func (ls *ListService) GetListByID(listID, userID int64) (*models.List, error) {
	return ls.repo.GetListByID(listID, userID)
}

func (ls *ListService) GetUserLists(userID int64) ([]*models.List, error) {
	return nil, nil
}

func (ls *ListService) Delete(listID int64) error {
	return ls.repo.Delete(listID)
}

func (ls *ListService) Update(listID int64, list *models.UpdateListReq) error {
	if list.Title == nil && list.Description == nil {
		return models.ErrUpdateEmptyArgs
	} else if len(*list.Title) < 5 {
		return models.ErrTitleTooShort
	}
	return ls.repo.Update(listID, list)
}
