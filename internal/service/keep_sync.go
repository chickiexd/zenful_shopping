package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/chickiexd/zenful_shopping/internal/env"
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
			// items come with Ingredient and MeasurementUnit preloaded
			unit := item.MeasurementUnit.Abbreviation
			if unit == "" {
				unit = item.MeasurementUnit.Name
			}
			ingredientData := IngredientSyncData{
				Name:     item.Ingredient.Name,
				Quantity: fmt.Sprintf("%g %s", item.Quantity, unit),
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

func (s *keepSyncService) SyncShoppingLists(ctx context.Context) error {
	data, err := s.prepareData()
	if err != nil {
		return fmt.Errorf("failed to prepare data for sync: %w", err)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to encode sync data: %w", err)
	}
	pythonPath := env.GetString("PYTHON_PATH", "python3")
	scriptPath := env.GetString("KEEP_SCRIPT_PATH", "./scripts/keep_updater.py")
	cmd := exec.CommandContext(ctx, pythonPath, scriptPath)
	cmd.Stdin = bytes.NewReader(jsonData)
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("keep sync script failed: %w: %s", err, stderr.String())
	}
	return nil
}
