package service

import (
	"log"
	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/store"

	"gorm.io/gorm"
)

type FoodGroupService struct {
	storage *store.Storage
}

func (s *FoodGroupService) Create(food_group *dto.CreateFoodGroupRequest) (*store.FoodGroup, error) {
	var created_food_group *store.FoodGroup
	err := s.storage.DB.Transaction(func(tx *gorm.DB) error {
		foodGroupsRepo := s.storage.FoodGroups.(*store.FoodGroupRepository).WithTransaction(tx)
		shoppingListsRepo := s.storage.ShoppingLists.(*store.ShoppingListRepository).WithTransaction(tx)

		created_food_group = &store.FoodGroup{
			Name: food_group.Name,
		}
		if err := foodGroupsRepo.Create(created_food_group); err != nil {
			return err
		}

		shopping_list, err := shoppingListsRepo.GetByID(food_group.ShoppingListID)
		if err != nil {
			return err
		}
		err = shoppingListsRepo.CreateFoodGroupAssociation(shopping_list, created_food_group)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Error creating food group: ", err)
		return nil, err
	}
	return created_food_group, nil
}

func (s *FoodGroupService) GetAll() ([]store.FoodGroup, error) {
	food_groups, err := s.storage.FoodGroups.GetAll()
	return food_groups, err
}

// func (s *FoodGroupService) GetByID(id uint) (*store.Measurement, error) {
// 	ingredient, err := s.storage.FoodGroups.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ingredient, nil
// }
