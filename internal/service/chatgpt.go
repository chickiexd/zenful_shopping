package service

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"zenful_shopping_backend/internal/dto"
	"zenful_shopping_backend/internal/env"
)

type OpenAIService struct {
	Client *openai.Client
}

func NewOpenAIService() *OpenAIService {
	api_key := env.GetString("OPENAI_API_KEY", "")
	client := openai.NewClient(api_key)
	return &OpenAIService{Client: client}
}

func (s *OpenAIService) ParseRecipe(recipeText string) (*dto.CreateRecipeRequest, error) {
	prompt := fmt.Sprintf(`
	Please parse this recipe:

	"%s"

	Return the recipe as an object with the keys: recipe (an object with the following keys: title, description, public which is always false, cook_time(int), servings),ingredients (an array of object with the keys: name, quantity, measurement_unit), instructions (an array of object with the keys: numbering, content), meal_type (an object with the key: name). 
	Return just the object without any strings or symbols around it. The measurement unit does not always have to be entered so dont try to force it when it does not make sense in the context, in that case just use an empty string for measurement_unit. Remember to always add "" around the keys for json parsing support. and the quantity should always be float value so if u get a range like 1-2 return the center so 1.5 in that case.
	`, recipeText)

	ctx := context.Background()
	response, err := s.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
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

	var recipe dto.CreateRecipeRequest
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}
