package service

import "github.com/VladimirStepanov/todo-app/internal/models"

type ItemService struct {
	repo models.ItemRepository
}

func NewItemService(repo models.ItemRepository) models.ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (is *ItemService) Create(title, description string, listID int64) (int64, error) {
	return is.repo.Create(title, description, listID)
}

func (is *ItemService) GetItems(listID int64) ([]*models.Item, error) {
	return nil, nil
}

func (is *ItemService) GetItemByID(listID, itemID int64) (*models.Item, error) {
	return is.repo.GetItemByID(listID, itemID)
}

func (is *ItemService) Update(listID, itemID int64, item *models.UpdateItemReq) error {
	if item.Title == nil && item.Description == nil {
		return models.ErrUpdateEmptyArgs
	} else if len(*item.Title) < 5 {
		return models.ErrTitleTooShort
	}
	return is.repo.Update(listID, itemID, item)
}

func (is *ItemService) Done(listID, itemID int64) error {
	return nil
}

func (is *ItemService) Delete(listID, itemID int64) error {
	return is.repo.Delete(listID, itemID)
}
