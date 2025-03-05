package service

import (
	"log"
	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
)

type ingredientService struct {
	storage *store.Storage
}

func (s *ingredientService) Create(ingredient *dto.CreateIngredientRequest) error {
	log.Println("service create ingredient_id")
	return s.storage.Ingredients.Create(nil)
}

func (s *ingredientService) GetAll() ([]store.Ingredient, error) {
	// get food group
	// get measurements
	// ingredients, err := s.storage.Ingredients.GetAll()
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println(ingredients)
	//
	// var ingredientsResponse []dto.IngredientResponse
	//
	// for _, ingredient := range ingredients {
	// 	var ingredientResponse dto.IngredientResponse
	// 	measurements, err := s.storage.Ingredients.GetMeasurementUnitsByID(ingredient.IngredientID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	foodGroups, err := s.storage.Ingredients.GetFoodGroupsByID(ingredient.IngredientID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	//
	// 	ingredientResponse.IngredientID = ingredient.IngredientID
	// 	ingredientResponse.Name = ingredient.Name
	// 	ingredientResponse.Measurements = measurements
	// 	ingredientResponse.FoodGroups = foodGroups
	//
	// 	ingredientsResponse = append(ingredientsResponse, ingredientResponse)
	// }
	//
	// log.Println(ingredientsResponse)
	// return ingredientsResponse, nil

	ingredients, err := s.storage.Ingredients.GetAllTest()

	return ingredients, err


	// measurements, err := s.storage.Ingredients.GetMeasurementUnitsByID(1)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println(measurements)

	// ingredients = append(ingredients, dto.IngredientResponse{
	// 	IngredientID: 1,
	// 	Name:         "Tomato",
	// 	Measurements: []dto.IngredientResponseMeasurement{
	// 		{MeasurementID: 1, Name: "Cup"},
	// 		{MeasurementID: 2, Name: "Ounce"},
	// 	},
	// 	FoodGroup: dto.IngredientResponseFoodGroup{
	// 		FoodGroupID: 1,
	// 		Name:        "Vegetables",
	// 	},
	// })
	// log.Println(ingredients)

	// return ingredients, nil
	return nil, nil
}

func (s *ingredientService) GetByID(id uint) (*store.Ingredient, error) {
	ingredient, err := s.storage.Ingredients.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}
