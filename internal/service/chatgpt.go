package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/env"
	"github.com/chickiexd/zenful_shopping/internal/logger"
	"github.com/chickiexd/zenful_shopping/internal/store"
)

// ErrRecipeRefused is returned when the model declines to parse the given text.
var ErrRecipeRefused = errors.New("the recipe could not be parsed from the given text")

// completer returns a JSON document matching schema; abstracted for tests.
type completer interface {
	completeJSON(ctx context.Context, system, user string, schema *jsonschema.Definition) (string, error)
}

type OpenAIService struct {
	llm     completer
	storage *store.Storage
}

func NewOpenAIService(storage *store.Storage) *OpenAIService {
	config := openai.DefaultConfig(env.GetString("OPENAI_API_KEY", ""))
	// lets the API be reached through a proxy or a compatible service
	if baseURL := env.GetString("OPENAI_BASE_URL", ""); baseURL != "" {
		config.BaseURL = baseURL
	}
	return &OpenAIService{
		llm: &openAICompleter{
			client: openai.NewClientWithConfig(config),
			model:  env.GetString("OPENAI_MODEL", openai.GPT4oMini),
		},
		storage: storage,
	}
}

type openAICompleter struct {
	client *openai.Client
	model  string
}

func (c *openAICompleter) completeJSON(ctx context.Context, system, user string, schema *jsonschema.Definition) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		// Temperature is omitempty, so a literal 0 would fall back to the API default
		Temperature: math.SmallestNonzeroFloat32,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "recipe",
				Schema: schema,
				Strict: true,
			},
		},
	}
	response, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil && retryable(ctx, err) {
		logger.Log.Warnw("retrying openai request", "error", err)
		response, err = c.client.CreateChatCompletion(ctx, req)
	}
	if err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	message := response.Choices[0].Message
	if message.Refusal != "" {
		logger.Log.Warnw("openai refused to parse recipe", "refusal", message.Refusal)
		return "", ErrRecipeRefused
	}
	return message.Content, nil
}

// retryable reports whether a failed request is worth one more attempt.
func retryable(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatusCode == http.StatusTooManyRequests || apiErr.HTTPStatusCode >= 500
	}
	return true
}

// shape of the model's answer; mirrors recipeSchema
type llmRecipe struct {
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	CookTime      int             `json:"cook_time"`
	Servings      int             `json:"servings"`
	MealType      string          `json:"meal_type"`
	Instructions  []string        `json:"instructions"`
	Ingredients   []llmIngredient `json:"ingredients"`
	NewFoodGroups []llmFoodGroup  `json:"new_food_groups"`
}

type llmIngredient struct {
	OriginalText       string   `json:"original_text"`
	Name               string   `json:"name"`
	ExistingIngredient string   `json:"existing_ingredient"`
	Quantity           float64  `json:"quantity"`
	Unit               string   `json:"unit"`
	Description        string   `json:"description"`
	MeasurementUnits   []string `json:"measurement_units"`
	FoodGroups         []string `json:"food_groups"`
}

type llmFoodGroup struct {
	Name         string `json:"name"`
	ShoppingList string `json:"shopping_list"`
}

// parseCatalog is the data the model chooses from and results are resolved against.
type parseCatalog struct {
	ingredients   []store.Ingredient
	units         []store.MeasurementUnit
	foodGroups    []store.FoodGroup
	shoppingLists []store.ShoppingList
	mealTypes     []store.MealType
}

func (s *OpenAIService) loadCatalog() (*parseCatalog, error) {
	var c parseCatalog
	var err error
	if c.ingredients, err = s.storage.Ingredients.GetAll(); err != nil {
		return nil, err
	}
	if c.units, err = s.storage.MeasurementUnits.GetAll(); err != nil {
		return nil, err
	}
	if c.foodGroups, err = s.storage.FoodGroups.GetAll(); err != nil {
		return nil, err
	}
	if c.shoppingLists, err = s.storage.ShoppingLists.GetAll(); err != nil {
		return nil, err
	}
	if c.mealTypes, err = s.storage.MealTypes.GetAll(); err != nil {
		return nil, err
	}
	if len(c.units) == 0 || len(c.mealTypes) == 0 {
		return nil, fmt.Errorf("measurement units and meal types must be configured before parsing recipes")
	}
	return &c, nil
}

func names[T any](items []T, name func(T) string) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = name(item)
	}
	return out
}

// enumOrString restricts a string to values, unless there are none to choose from.
func enumOrString(description string, values []string) jsonschema.Definition {
	d := jsonschema.Definition{Type: jsonschema.String, Description: description}
	if len(values) > 0 {
		d.Enum = values
	}
	return d
}

// object builds a strict-mode object: every property required, nothing extra.
func object(props map[string]jsonschema.Definition) jsonschema.Definition {
	required := make([]string, 0, len(props))
	for k := range props {
		required = append(required, k)
	}
	return jsonschema.Definition{Type: jsonschema.Object, Properties: props, Required: required, AdditionalProperties: false}
}

func recipeSchema(c *parseCatalog) *jsonschema.Definition {
	unitNames := names(c.units, func(u store.MeasurementUnit) string { return u.Name })
	str := func(description string) jsonschema.Definition {
		return jsonschema.Definition{Type: jsonschema.String, Description: description}
	}
	ingredient := object(map[string]jsonschema.Definition{
		"original_text":       str("the ingredient name exactly as written in the recipe, without quantity and unit"),
		"name":                str("canonical English ingredient name: singular, Title Case, base product only"),
		"existing_ingredient": str("exact name from the known ingredient list if it is the same ingredient, otherwise empty string"),
		"quantity":            {Type: jsonschema.Number},
		"unit":                enumOrString("unit of quantity", unitNames),
		"description":         str("one short sentence describing the ingredient"),
		"measurement_units": {
			Type:        jsonschema.Array,
			Description: "all units this ingredient is commonly measured in",
			Items:       &jsonschema.Definition{Type: jsonschema.String, Enum: unitNames},
		},
		"food_groups": {
			Type:        jsonschema.Array,
			Description: "food groups of the ingredient, preferring known food groups",
			Items:       &jsonschema.Definition{Type: jsonschema.String},
		},
	})
	foodGroup := object(map[string]jsonschema.Definition{
		"name":          str("name of a food group that is not in the known list"),
		"shopping_list": enumOrString("shopping list this food group is bought from", names(c.shoppingLists, func(l store.ShoppingList) string { return l.Name })),
	})
	schema := object(map[string]jsonschema.Definition{
		"title":           str("recipe title in English"),
		"description":     str("one or two sentence description in English"),
		"cook_time":       {Type: jsonschema.Integer, Description: "total time in minutes, 0 if unknown"},
		"servings":        {Type: jsonschema.Integer, Description: "number of servings, 0 if unknown"},
		"meal_type":       enumOrString("meal type", names(c.mealTypes, func(m store.MealType) string { return m.Name })),
		"instructions":    {Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.String, Description: "one step, without its number"}},
		"ingredients":     {Type: jsonschema.Array, Items: &ingredient},
		"new_food_groups": {Type: jsonschema.Array, Items: &foodGroup},
	})
	return &schema
}

const parseSystemPrompt = `You turn recipes into structured data for a shopping list app.
The recipe may be written in German, English or a mix of both. Write every text field in English.

Ingredient names:
- "name" is the canonical English name: singular, Title Case, only the product you would buy.
  Drop preparation and size words (chopped, gehackt, fresh, large, fein, gewürfelt).
  Keep a variety when it changes what you buy (Brown Sugar, Red Onion, Greek Yogurt).
- If the ingredient is the same product as one in the known ingredient list, put that exact
  name in "existing_ingredient". Plurals, translations, sizes and preparations count as the
  same product ("Knoblauchzehen" -> "Garlic", "chopped onions" -> "Onion").
  A different product you would buy separately is NOT the same, even if the names overlap:
  Coconut Milk is not Whole Milk, Sweet Potato is not Potato, Tomato Paste is not Tomato.
  When unsure, use "".
- "original_text" is the name as written in the recipe.

Quantities and units (unit must be one of the allowed units):
- g, kg, ml, l map directly. EL, tbsp -> tablespoon. TL, tsp -> teaspoon. Tasse -> cup.
- Zehe, clove, Stück, Bund, Dose, can, Scheibe, slice, or no unit -> piece.
- Prise or pinch -> 0.1 teaspoon. "nach Geschmack", "etwas", "to taste" -> quantity 1, piece.
- Ranges become their midpoint ("1-2" -> 1.5). Fractions become decimals ("1/2" -> 0.5).
- If a weight or volume is given next to a count ("1 Dose (400 g)"), use the weight or volume.
- List every ingredient once per line in the recipe; do not skip any.

Food groups: prefer the known food groups. Only list a food group in "new_food_groups" if it is
not known, together with the shopping list it is bought from.

Instructions: one entry per step, in order, without numbering.`

func parseUserPrompt(c *parseCatalog, recipeText string) string {
	return fmt.Sprintf("Known ingredients:\n%s\n\nKnown food groups:\n%s\n\nRecipe:\n<<<\n%s\n>>>",
		strings.Join(names(c.ingredients, func(i store.Ingredient) string { return i.Name }), "\n"),
		strings.Join(names(c.foodGroups, func(f store.FoodGroup) string { return f.Name }), "\n"),
		recipeText)
}

func (s *OpenAIService) ParseRecipe(ctx context.Context, recipeText string) (*dto.ParsedRecipe, error) {
	catalog, err := s.loadCatalog()
	if err != nil {
		return nil, err
	}
	content, err := s.llm.completeJSON(ctx, parseSystemPrompt, parseUserPrompt(catalog, recipeText), recipeSchema(catalog))
	if err != nil {
		return nil, err
	}
	var parsed llmRecipe
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		logger.Log.Errorw("unexpected recipe format from openai", "error", err, "response", content)
		return nil, fmt.Errorf("unexpected recipe format from openai: %w", err)
	}
	return s.buildParsedRecipe(catalog, &parsed)
}

// buildParsedRecipe resolves the model's answer against the database: units,
// meal type and ingredients become IDs, duplicates are merged and new food
// groups are checked against existing ones.
func (s *OpenAIService) buildParsedRecipe(c *parseCatalog, parsed *llmRecipe) (*dto.ParsedRecipe, error) {
	unitsByName := make(map[string]store.MeasurementUnit, len(c.units))
	for _, u := range c.units {
		unitsByName[strings.ToLower(u.Name)] = u
	}
	unitFor := func(name string) (store.MeasurementUnit, error) {
		if u, ok := unitsByName[strings.ToLower(strings.TrimSpace(name))]; ok {
			return u, nil
		}
		if u, ok := unitsByName["piece"]; ok {
			return u, nil
		}
		return store.MeasurementUnit{}, fmt.Errorf("unknown measurement unit %q", name)
	}

	result := &dto.ParsedRecipe{
		Recipe: dto.RecipeRequest{
			Title:       parsed.Title,
			Description: parsed.Description,
			CookTime:    parsed.CookTime,
			Servings:    parsed.Servings,
		},
		Ingredients:    []dto.ParsedRecipeIngredient{},
		NewIngredients: []dto.ParsedIngredient{},
		NewFoodGrops:   []dto.ParsedFoodGroup{},
		Warnings:       []string{},
	}
	for _, m := range c.mealTypes {
		if strings.EqualFold(m.Name, parsed.MealType) {
			result.Recipe.MealType = int(m.MealTypeID)
		}
	}
	if result.Recipe.MealType == 0 {
		result.Recipe.MealType = int(c.mealTypes[0].MealTypeID)
	}
	for _, step := range parsed.Instructions {
		if step = strings.TrimSpace(step); step != "" {
			result.Instructions = append(result.Instructions, dto.InstructionRequest{Numbering: len(result.Instructions) + 1, Content: step})
		}
	}

	// food groups: known ones keep their stored name, unknown ones get a validated shopping list
	knownGroups := make(map[string]string, len(c.foodGroups))
	for _, fg := range c.foodGroups {
		knownGroups[strings.ToLower(fg.Name)] = fg.Name
	}
	listNames := make(map[string]string, len(c.shoppingLists))
	for _, l := range c.shoppingLists {
		listNames[strings.ToLower(l.Name)] = l.Name
	}
	newGroups := map[string]dto.ParsedFoodGroup{}
	var newGroupOrder []string
	addNewGroup := func(key string, group dto.ParsedFoodGroup) {
		if _, ok := newGroups[key]; !ok {
			newGroupOrder = append(newGroupOrder, key)
		}
		newGroups[key] = group
	}
	for _, fg := range parsed.NewFoodGroups {
		key := strings.ToLower(strings.TrimSpace(fg.Name))
		if key == "" || knownGroups[key] != "" {
			continue
		}
		addNewGroup(key, dto.ParsedFoodGroup{Name: strings.TrimSpace(fg.Name), ShoppingLists: listNames[strings.ToLower(fg.ShoppingList)]})
	}
	resolveGroup := func(name string) (dto.ParsedFoodGroup, bool) {
		key := strings.ToLower(strings.TrimSpace(name))
		if key == "" {
			return dto.ParsedFoodGroup{}, false
		}
		if known := knownGroups[key]; known != "" {
			return dto.ParsedFoodGroup{Name: known}, true
		}
		group, ok := newGroups[key]
		if !ok {
			// referenced but not declared as new: the user picks a shopping list during review
			group = dto.ParsedFoodGroup{Name: strings.TrimSpace(name)}
			addNewGroup(key, group)
		}
		return group, true
	}

	matcher := newIngredientMatcher(s.storage.Aliases.WithTransaction(s.storage.DB), c.ingredients)
	unitNames := make(map[uint]string, len(c.units))
	for _, u := range c.units {
		unitNames[u.MeasurementUnitID] = u.Name
	}
	// a recipe holds each ingredient once, so repeated lines are merged into the
	// first one, converting units of the same dimension
	merge := func(name string, quantity *float64, unitID uint, addQuantity float64, addUnit store.MeasurementUnit) {
		if addUnit.MeasurementUnitID == unitID {
			*quantity += addQuantity
			return
		}
		if converted, ok := convertQuantity(addQuantity, addUnit.Name, unitNames[unitID]); ok {
			*quantity += converted
			return
		}
		result.Warnings = append(result.Warnings, fmt.Sprintf(
			"%s is listed more than once in units that cannot be combined (%s and %s); only %g %s was kept",
			name, unitNames[unitID], addUnit.Name, *quantity, unitNames[unitID]))
	}
	matchedIdx := map[uint]int{}
	newIdx := map[string]int{}
	for _, ing := range parsed.Ingredients {
		unit, err := unitFor(ing.Unit)
		if err != nil {
			return nil, err
		}
		existing, matchedBy, err := matcher.resolve(ing.ExistingIngredient, ing.OriginalText, ing.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			if i, ok := matchedIdx[existing.IngredientID]; ok {
				line := &result.Ingredients[i]
				merge(existing.Name, &line.Quantity, line.MeasurementUnitID, ing.Quantity, unit)
				continue
			}
			matchedIdx[existing.IngredientID] = len(result.Ingredients)
			result.Ingredients = append(result.Ingredients, dto.ParsedRecipeIngredient{
				IngredientID:      existing.IngredientID,
				IngredientName:    existing.Name,
				OriginalName:      strings.TrimSpace(ing.OriginalText),
				Quantity:          ing.Quantity,
				MeasurementUnitID: unit.MeasurementUnitID,
				MatchedBy:         matchedBy,
			})
			continue
		}

		key := normalizeName(ing.Name)
		if i, ok := newIdx[key]; ok {
			line := &result.NewIngredients[i]
			merge(line.Name, &line.Quantity, line.MeasurementUnitID, ing.Quantity, unit)
			continue
		}
		groups := []dto.ParsedFoodGroup{}
		for _, name := range ing.FoodGroups {
			if group, ok := resolveGroup(name); ok {
				groups = append(groups, group)
			}
		}
		units := ing.MeasurementUnits
		if !containsFold(units, unit.Name) {
			units = append(units, unit.Name)
		}
		newIdx[key] = len(result.NewIngredients)
		result.NewIngredients = append(result.NewIngredients, dto.ParsedIngredient{
			Name:                 strings.TrimSpace(ing.Name),
			OriginalName:         strings.TrimSpace(ing.OriginalText),
			Description:          ing.Description,
			Quantity:             ing.Quantity,
			MeasurementUnitNames: units,
			MeasurementUnitID:    unit.MeasurementUnitID,
			ParsedFoodGroups:     groups,
		})
	}

	// only report new food groups that a new ingredient actually uses
	used := map[string]bool{}
	for _, ing := range result.NewIngredients {
		for _, g := range ing.ParsedFoodGroups {
			used[strings.ToLower(g.Name)] = true
		}
	}
	for _, key := range newGroupOrder {
		if used[key] {
			result.NewFoodGrops = append(result.NewFoodGrops, newGroups[key])
		}
	}
	return result, nil
}

func containsFold(values []string, value string) bool {
	for _, v := range values {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}
