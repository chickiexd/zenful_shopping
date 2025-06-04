package service

import (
	"fmt"
	"bytes"
	"encoding/json"
	// "log"
	"os"
	"os/exec"
	// "zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/store"
	// "zenful_shopping_backend/utils"
	// "gorm.io/gorm"
)

type keepSyncService struct {
	storage *store.Storage
}

type ShoppingListSyncData struct {
	Title       string               `json:"title"`
	Ingredients []IngredientSyncData `json:"ingredients"`
	Color       string               `json:"color"`
}

type IngredientSyncData struct {
	Name     string  `json:"name"`
	Quantity string `json:"quantity"`
}

func (s *keepSyncService) prepareData() ([]ShoppingListSyncData, error) {
	shoppingLists, err := s.storage.ShoppingLists.GetAll()
	if err != nil {
		return nil, err
	}
	var syncData []ShoppingListSyncData
	for _, list := range shoppingLists {
		sl_items, err := s.storage.ShoppingLists.GetItemsByShoppingListID(list.ShoppingListID)
		if err != nil {
			return nil, err
		}
		var ingredientsData []IngredientSyncData
		for _, item := range sl_items {
			// Fetch the ingredient name from the Ingredient table
			ingredient, err := s.storage.Ingredients.GetByID(item.IngredientID)
			if err != nil {
				return nil, err
			}
			// fetch measurement unit
			measurementUnit, err := s.storage.MeasurementUnits.GetByID(item.MeasurementUnitID)
			if err != nil {
				return nil, err
			}
			ingredientData := IngredientSyncData{
				Name:     ingredient.Name,
				Quantity: fmt.Sprintf("%g %s", item.Quantity, measurementUnit.Name),
			}
			ingredientsData = append(ingredientsData, ingredientData)
		}
		listData := ShoppingListSyncData{
			Title:       list.Name,
			Ingredients: ingredientsData,
			Color:       "GREEN",
		}
		syncData = append(syncData, listData)
	}
	return syncData, nil

}

func (s *keepSyncService) SyncShoppingLists() error {

	data, err := s.prepareData()

	pythonPath := "./scripts/.venv/bin/python3"
	scriptPath := "./scripts/keep_updater.py"

	jsonData, _ := json.Marshal(data)

	cmd := exec.Command(pythonPath, scriptPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewReader(jsonData)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
