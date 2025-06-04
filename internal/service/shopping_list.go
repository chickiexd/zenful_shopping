package service

import (
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/store"
)

type ShoppingListService struct {
	storage *store.Storage
}

func (s *ShoppingListService) GetAll() ([]dto.ShoppingListResponse, error) {
	shopping_lists, err := s.storage.ShoppingLists.GetAll()
	if err != nil {
		return nil, err
	}
	var shopping_list_dto []dto.ShoppingListResponse
	for _, shopping_list := range shopping_lists {
		ingredients, err := s.storage.ShoppingLists.GetItemsByShoppingListID(shopping_list.ShoppingListID)
		if err != nil {
			return nil, err
		}
		var ingredients_dto []dto.ShoppingListItemResponse
		for _, ingredient := range ingredients {
			ingredients_dto = append(ingredients_dto, dto.ShoppingListItemResponse{
				ShoppingListItemID: ingredient.ShoppingListItemID,
				IngredientID:       ingredient.IngredientID,
				Quantity:           ingredient.Quantity,
				MeasurementUnitID:  ingredient.MeasurementUnitID,
			})
		}
		shopping_list_dto = append(shopping_list_dto, dto.ShoppingListResponse{
			ShoppingListID: shopping_list.ShoppingListID,
			Name:           shopping_list.Name,
			Items:          ingredients_dto,
		})
	}
	return shopping_list_dto, nil
}

func (s *ShoppingListService) RemoveItemFromShoppingList(shopping_list_item_id uint) error {
	err := s.storage.ShoppingLists.DeleteItemAssociationByID(shopping_list_item_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShoppingListService) RemoveAllItemsFromShoppingList(shopping_list_id uint) error {
	err := s.storage.ShoppingLists.DeleteAllItemsByShoppingListID(shopping_list_id)
	if err != nil {
		return err
	}
	return nil
}

