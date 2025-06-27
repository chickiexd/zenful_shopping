package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"github.com/chickiexd/zenful_shopping/internal/store"
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
	Name     string `json:"name"`
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
			ingredient, err := s.storage.Ingredients.GetByID(item.IngredientID)
			if err != nil {
				return nil, err
			}
			measurementUnit, err := s.storage.MeasurementUnits.GetByID(item.MeasurementUnitID)
			if err != nil {
				return nil, err
			}
			// TODO: use abbreviation for measurement unit
			ingredientData := IngredientSyncData{
				Name:     ingredient.Name,
				Quantity: fmt.Sprintf("%g %s", item.Quantity, measurementUnit.Name),
			}
			ingredientsData = append(ingredientsData, ingredientData)
		}
		listData := ShoppingListSyncData{
			Title:       list.Name,
			Ingredients: ingredientsData,
			Color:       list.Color,
		}
		syncData = append(syncData, listData)
	}
	return syncData, nil
}

func (s *keepSyncService) SyncShoppingLists() error {
	data, err := s.prepareData()
	if err != nil {
		return fmt.Errorf("failed to prepare data for sync: %w", err)
	}
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
