package service

import (
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
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
