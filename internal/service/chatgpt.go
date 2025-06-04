package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/env"
	"github.com/chickiexd/zenful_shopping/internal/store"
	"github.com/chickiexd/zenful_shopping/utils"
)

type OpenAIService struct {
	Client  *openai.Client
	storage *store.Storage
}

func NewOpenAIService(storage *store.Storage) *OpenAIService {
	api_key := env.GetString("OPENAI_API_KEY", "")
	client := openai.NewClient(api_key)
	return &OpenAIService{Client: client, storage: storage}
}

func (s *OpenAIService) getShoppingListForFoodGroup(food_group_name string) (string, error) {
	shopping_lists, err := s.storage.ShoppingLists.GetAll()
	if err != nil {
		return "", err
	}
	shopping_lists_str := ""
	for _, shopping_list := range shopping_lists {
		shopping_lists_str += fmt.Sprintf("%s,", shopping_list.Name)
	}
	prompt := fmt.Sprintf(`
	You are given the following food group name:

	"%s"

	Return the shopping list name as an object with the keys: shopping_list_name. 
	The shopping_list_name is a string and can only be one of the following:
	"%s"
	Choose the one that makes sense for a recipe with the given food group.
	Return just the object without any strings or symbols around it.
	Remember to always add "" around the keys for json parsing support.
	`, food_group_name, shopping_lists_str)

	ctx := context.Background()
	response, err := s.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that returns structured data for a shoppinglist.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}
	var shopping_list dto.ShoppingListName
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &shopping_list)
	return shopping_list.Name, err
}

func (s *OpenAIService) createIngredient(ingredient_name string) (*dto.ParsedIngredientInformation, error) {
	food_groups, err := s.storage.FoodGroups.GetAll()
	if err != nil {
		return nil, err
	}
	food_groups_str := ""
	for _, food_group := range food_groups {
		food_groups_str += fmt.Sprintf("%s,", food_group.Name)
	}
	prompt := fmt.Sprintf(`
	You are given the following ingredient name:

	"%s"

	Return the ingredient as an object with the keys: ingredient_name, ingredient_description, measurement_units, food_groups. 
	The measurement_units is an array of strings and can be any amount of the following:
	gram,kilogram,liter,milliliter,teaspoon,tablespoon,cup,piece.
	Choose the ones that make sense for a recipe with the given ingredient.
	The food_groups are an array of strings and can be any of the following or u can also add more fitting new ones:
	"%s"
	Choose the ones that make sense for a recipe with the given ingredient.
	Return just the object without any strings or symbols around it.
	Remember to always add "" around the keys for json parsing support.
	`, ingredient_name, food_groups_str)

	ctx := context.Background()
	response, err := s.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		// Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that returns structured data for ingredients that are supposed to be in recipes.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var ingredient dto.ParsedIngredientInformation
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &ingredient)
	return &ingredient, err
}

func (s *OpenAIService) ParseRecipe(recipeText string) (*dto.ParsedRecipe, error) {
	prompt := fmt.Sprintf(`
	Please parse this recipe:

	"%s"

	First of all return everything in english and with correct Capitalisation. Even the ingredient names, measurement units and food groups.
	Return the recipe as an object with the keys: recipe (an object with the following keys: title, description, public which is always false, cook_time(int), servings),ingredients (an array of object with the keys: name, quantity, measurement_unit), instructions (an array of object with the keys: numbering, content), meal_type (an object with the key: name). 
	Return just the object without any strings or symbols around it. The measurement unit does not always have to be entered so dont try to force it when it does not make sense in the context, in that case just use an empty string for measurement_unit. The following measurement units can be used:
	gram,kilogram,liter,milliliter,teaspoon,tablespoon,cup,piece.
	The ingredient name should always be just the ingredient no any additional information.
	Remember to always add "" around the keys for json parsing support. and the quantity should always be float value so if u get a range like 1-2 return the center so 1.5 in that case.
	`, recipeText)

	ctx := context.Background()
	response, err := s.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		// Model: openai.GPT3Dot5Turbo,
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that returns structured data from recipes.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var recipe dto.ParsedRecipeInformation
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &recipe)
	if err != nil {
		return nil, err
	}

	var response_recipe dto.ParsedRecipe
	response_recipe.Recipe = recipe.Recipe

	for _, ingredient := range recipe.Ingredients {
		measurement_unit, err := s.storage.MeasurementUnits.GetByName(ingredient.MeasurementUnit)
		if err != nil {
			return nil, err
		}
		existing_ingredient, err := s.storage.Ingredients.GetByName(ingredient.Name)
		if err == gorm.ErrRecordNotFound {
			parsed_ingredient, parsed_foodgroups, err := s.ParseNewIngredient(ingredient.Name)
			if err != nil {
				return nil, err
			}
			parsed_ingredient.Quantity = ingredient.Quantity
			parsed_ingredient.MeasurementUnitID = measurement_unit.MeasurementUnitID
			response_recipe.NewIngredients = append(response_recipe.NewIngredients, *parsed_ingredient)
			response_recipe.NewFoodGrops = append(response_recipe.NewFoodGrops, parsed_foodgroups...)
		} else if err != nil {
			return nil, err
		} else {
			response_recipe.Ingredients = append(response_recipe.Ingredients, dto.RecipeIngredientResponse{
				IngredientID:      existing_ingredient.IngredientID,
				Quantity:          ingredient.Quantity,
				MeasurementUnitID: measurement_unit.MeasurementUnitID,
			})
		}
	}

	response_recipe.Instructions = recipe.Instructions
	response_recipe.NewFoodGrops = utils.RemoveDuplicate(response_recipe.NewFoodGrops)
	log.Println("Parsed recipe: ", response_recipe)
	return &response_recipe, nil
}

func (s *OpenAIService) ParseNewIngredient(ingredient_name string) (*dto.ParsedIngredient, []dto.ParsedFoodGroup, error) {
	parsed_ingredient := &dto.ParsedIngredient{}
	var new_parsed_food_groups []dto.ParsedFoodGroup
	parsed_ingredient_info, err := s.createIngredient(ingredient_name)
	if err != nil {
		return nil, nil, err
	}
	for _, food_group := range parsed_ingredient_info.FoodGroups {
		_, err := s.storage.FoodGroups.GetByName(food_group)
		if err == gorm.ErrRecordNotFound {
			shopping_list, err := s.getShoppingListForFoodGroup(food_group)
			new_food_group := &dto.ParsedFoodGroup{
				Name:          food_group,
				ShoppingLists: shopping_list,
			}
			if err != nil {
				return nil, nil, err
			}
			new_parsed_food_groups = append(new_parsed_food_groups, *new_food_group)
			parsed_ingredient.ParsedFoodGroups = append(parsed_ingredient.ParsedFoodGroups, *new_food_group)
		} else if err != nil {
			return nil, nil, err
		} else {
			parsed_food_group := &dto.ParsedFoodGroup{
				Name:          food_group,
				ShoppingLists: "",
			}
			parsed_ingredient.ParsedFoodGroups = append(parsed_ingredient.ParsedFoodGroups, *parsed_food_group)
		}
	}
	parsed_ingredient.Name = parsed_ingredient_info.IngredientName
	parsed_ingredient.MeasurementUnitNames = parsed_ingredient_info.MeasurementUnits
	return parsed_ingredient, new_parsed_food_groups, nil
}
