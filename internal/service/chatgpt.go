package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	// "strings"

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

type newIngredientFound struct {
	Name              string
	Quantity          float64
	MeasurementUnitID uint
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

func (s *OpenAIService) getIngredientsAndFoodGroups(ingredients []newIngredientFound) ([]dto.ParsedMultipleIngredientInformation, error) {
	food_groups, err := s.storage.FoodGroups.GetAll()
	if err != nil {
		return nil, err
	}
	food_groups_str := ""
	for _, food_group := range food_groups {
		food_groups_str += fmt.Sprintf("%s,", food_group.Name)
	}
	ingredients_str, err := json.MarshalIndent(ingredients, "", "  ")
	if err != nil {
		return nil, err
	}
	prompt := fmt.Sprintf(`
	You are given the following ingredients with the following type:
	type ingredientType struct {
		Name              string
		Quantity          float64
		MeasurementUnitID uint
	}
	:

	"%s"

	For each ingredient do the following and return them as an array of objects:
	Return the ingredient as an object with the keys: ingredient_name, ingredient_description, measurement_units, food_groups, quantity and measurement_unit_id. 
	The measurement_units is an array of strings and can be any amount of the following:
	gram,kilogram,liter,milliliter,teaspoon,tablespoon,cup,piece.
	Choose the ones that make sense for a recipe with the given ingredient.
	The food_groups are an array of strings and can be any of the following or u can also add more fitting new ones:
	"%s"
	Choose the ones that make sense for a recipe with the given ingredient.
	For quantity and measurement_unit_id use the values from the ingredientType struct.
	Return the array of parsed ingredient objects without any strings or symbols around and dont wrap it in a json codeblock, start with the [ and end with the ].
	Remember to always add "" around the keys for json parsing support.
	`, string(ingredients_str), food_groups_str)

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
	content := response.Choices[0].Message.Content

	// Strip code fences like ```json\n...\n```
	// content = strings.TrimPrefix(content, "```json\n")
	// content = strings.TrimSuffix(content, "```")
	// content = strings.TrimSpace(content) // in case there's whitespace

	log.Printf("Response from OpenAI: %s", response.Choices[0].Message.Content)
	var parsed_ingredients []dto.ParsedMultipleIngredientInformation
	err = json.Unmarshal([]byte(content), &parsed_ingredients)
	return parsed_ingredients, err
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

	var new_ingredients []newIngredientFound

	for _, ingredient := range recipe.Ingredients {
		measurement_unit, err := s.storage.MeasurementUnits.GetByName(ingredient.MeasurementUnit)
		if err != nil {
			log.Printf("Error getting measurement unit: %v for ingredient %s", err, ingredient.Name)
			return nil, err
		}
		existing_ingredient, err := s.storage.Ingredients.GetByName(ingredient.Name)
		if err == gorm.ErrRecordNotFound {
			new_ingredient := newIngredientFound{
				Name:              ingredient.Name,
				Quantity:          ingredient.Quantity,
				MeasurementUnitID: measurement_unit.MeasurementUnitID,
			}
			new_ingredients = append(new_ingredients, new_ingredient)
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
	log.Printf("New ingredients found: %v", new_ingredients)
	parsed_ingredients, parsed_foodgroups, err := s.ParseNewIngredients(new_ingredients)
	if err != nil {
		return nil, err
	}
	response_recipe.NewIngredients = append(response_recipe.NewIngredients, parsed_ingredients...)
	response_recipe.NewFoodGrops = append(response_recipe.NewFoodGrops, parsed_foodgroups...)
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

func (s *OpenAIService) ParseNewIngredients(new_ingredients []newIngredientFound) ([]dto.ParsedIngredient, []dto.ParsedFoodGroup, error) {
	var parsed_ingredients []dto.ParsedIngredient
	var new_parsed_food_groups []dto.ParsedFoodGroup
	parsed_ingredients_info, err := s.getIngredientsAndFoodGroups(new_ingredients)
	log.Printf("Parsed ingredients_info: %v", parsed_ingredients_info)
	if err != nil {
		return nil, nil, err
	}
	for _, parsed_ingredient_info := range parsed_ingredients_info {
		parsed_ingredient := dto.ParsedIngredient{
			Name:                 parsed_ingredient_info.IngredientName,
			Quantity:             parsed_ingredient_info.Quantity,
			MeasurementUnitNames: parsed_ingredient_info.MeasurementUnits,
			MeasurementUnitID:    parsed_ingredient_info.MeasurementUnitID,
			ParsedFoodGroups:     []dto.ParsedFoodGroup{},
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
		parsed_ingredients = append(parsed_ingredients, parsed_ingredient)
	}
	log.Printf("Parsed ingredients: %v", parsed_ingredients)
	log.Printf("New parsed food groups: %v", new_parsed_food_groups)
	return parsed_ingredients, new_parsed_food_groups, nil
}
