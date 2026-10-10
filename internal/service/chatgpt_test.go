package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/sashabaranov/go-openai/jsonschema"
	"gorm.io/gorm"
)

// seed ingredient IDs from migration 000005
const (
	seedSalt   = 1
	seedTomato = 7
	seedGarlic = 8
	seedOnion  = 9
)

type fakeCompleter struct {
	content string
	err     error
	schema  *jsonschema.Definition
}

func (f *fakeCompleter) completeJSON(_ context.Context, _, _ string, schema *jsonschema.Definition) (string, error) {
	f.schema = schema
	return f.content, f.err
}

func TestNormalizeName(t *testing.T) {
	cases := map[string]string{
		"Tomatoes":            "tomato",
		"  Red   Onions ":     "red onion",
		"Eggs":                "egg",
		"Berries":             "berry",
		"Peaches":             "peach",
		"Grass":               "grass",
		"Hummus":              "hummus",
		"Knoblauchzehen":      "knoblauchzehen",
		"Tuna (Canned)":       "tuna canned",
		"Crème fraîche":       "crème fraîche",
		"Green Beans":         "green bean",
		"Chopped, fresh Dill": "chopped fresh dill",
	}
	for in, want := range cases {
		if got := normalizeName(in); got != want {
			t.Errorf("normalizeName(%q) = %q, want %q", in, got, want)
		}
	}
}

const germanRecipe = `{
  "title": "Chickpea Curry",
  "description": "A quick curry.",
  "cook_time": 30,
  "servings": 4,
  "meal_type": "main",
  "instructions": ["Chop the onion.", "  ", "Cook everything."],
  "ingredients": [
    {"original_text": "Knoblauchzehen", "name": "Garlic", "existing_ingredient": "Garlic", "quantity": 2, "unit": "piece", "description": "", "measurement_units": ["piece"], "food_groups": ["Vegetables"]},
    {"original_text": "Zwiebeln", "name": "Onions", "existing_ingredient": "", "quantity": 1, "unit": "piece", "description": "", "measurement_units": ["piece"], "food_groups": ["Vegetables"]},
    {"original_text": "Salz", "name": "Salt", "existing_ingredient": "salt", "quantity": 0.1, "unit": "Teaspoon", "description": "", "measurement_units": ["gram"], "food_groups": ["Spices"]},
    {"original_text": "Knoblauch", "name": "Garlic", "existing_ingredient": "Garlic", "quantity": 1, "unit": "piece", "description": "", "measurement_units": ["piece"], "food_groups": ["Vegetables"]},
    {"original_text": "Meersalz", "name": "Salt", "existing_ingredient": "Salt", "quantity": 1, "unit": "tablespoon", "description": "", "measurement_units": ["gram"], "food_groups": ["Spices"]},
    {"original_text": "Knoblauchpulver", "name": "Garlic", "existing_ingredient": "Garlic", "quantity": 10, "unit": "gram", "description": "", "measurement_units": ["gram"], "food_groups": ["Spices"]},
    {"original_text": "Kichererbsen", "name": "Chickpea", "existing_ingredient": "", "quantity": 400, "unit": "gram", "description": "A legume.", "measurement_units": ["piece"], "food_groups": ["Legumes", "vegetables", "Grains"]},
    {"original_text": "Kichererbsen (Dose)", "name": "Chickpeas", "existing_ingredient": "", "quantity": 200, "unit": "gram", "description": "A legume.", "measurement_units": ["gram"], "food_groups": ["Legumes"]}
  ],
  "new_food_groups": [
    {"name": "Legumes", "shopping_list": "Other"},
    {"name": "Vegetables", "shopping_list": "Produce"},
    {"name": "Unused Group", "shopping_list": "Other"}
  ]
}`

func TestParseRecipeResolvesAgainstCatalog(t *testing.T) {
	s := newTestStorage(t)
	llm := &fakeCompleter{content: germanRecipe}
	svc := &OpenAIService{llm: llm, storage: s}

	got, err := svc.ParseRecipe(context.Background(), "irrelevant")
	must(t, err)

	if got.Recipe.Title != "Chickpea Curry" || got.Recipe.MealType != 2 {
		t.Errorf("recipe: got %+v, want title and meal type main (2)", got.Recipe)
	}
	if len(got.Instructions) != 2 || got.Instructions[1].Numbering != 2 || got.Instructions[1].Content != "Cook everything." {
		t.Errorf("instructions not renumbered: %+v", got.Instructions)
	}

	byID := map[uint]dto.ParsedRecipeIngredient{}
	for _, ing := range got.Ingredients {
		byID[ing.IngredientID] = ing
	}
	if len(got.Ingredients) != 3 {
		t.Fatalf("want garlic, onion and salt matched, got %+v", got.Ingredients)
	}
	if g := byID[seedGarlic]; g.Quantity != 3 {
		t.Errorf("garlic lines should merge to 3, got %+v", g)
	}
	if _, ok := byID[seedOnion]; !ok {
		t.Error("plural 'Onions' should match Onion")
	}
	if salt := byID[seedSalt]; salt.MeasurementUnitID != unitTeaspoon || salt.Quantity < 3.099 || salt.Quantity > 3.101 {
		t.Errorf("salt: 'Teaspoon' should resolve case-insensitively and 1 tablespoon merge as 3 teaspoons, got %+v", salt)
	}
	if g := byID[seedGarlic]; g.MatchedBy != matchedByModel || g.OriginalName != "Knoblauchzehen" || g.IngredientName != "Garlic" {
		t.Errorf("garlic match details: got %+v", g)
	}
	if o := byID[seedOnion]; o.MatchedBy != matchedByName {
		t.Errorf("onion should be matched by name, got %+v", o)
	}
	if len(got.Warnings) != 1 || !strings.Contains(got.Warnings[0], "Garlic") {
		t.Errorf("garlic in piece and gram cannot be combined and should warn, got %v", got.Warnings)
	}

	if len(got.NewIngredients) != 1 {
		t.Fatalf("both chickpea lines should merge into one new ingredient, got %+v", got.NewIngredients)
	}
	chickpea := got.NewIngredients[0]
	if chickpea.Quantity != 600 || chickpea.OriginalName != "Kichererbsen" || chickpea.MeasurementUnitID != unitGram {
		t.Errorf("unexpected new ingredient %+v", chickpea)
	}
	if !containsFold(chickpea.MeasurementUnitNames, "gram") {
		t.Errorf("the unit used in the recipe should be offered, got %v", chickpea.MeasurementUnitNames)
	}
	wantGroups := []dto.ParsedFoodGroup{{Name: "Legumes", ShoppingLists: "Other"}, {Name: "Vegetables"}, {Name: "Grains"}}
	if len(chickpea.ParsedFoodGroups) != len(wantGroups) {
		t.Fatalf("food groups: got %+v", chickpea.ParsedFoodGroups)
	}
	for i, want := range wantGroups {
		if chickpea.ParsedFoodGroups[i] != want {
			t.Errorf("food group %d: got %+v, want %+v", i, chickpea.ParsedFoodGroups[i], want)
		}
	}
	// known "Vegetables" and the unused group are not new; undeclared "Grains" is, without a list
	wantNew := []dto.ParsedFoodGroup{{Name: "Legumes", ShoppingLists: "Other"}, {Name: "Grains"}}
	if len(got.NewFoodGrops) != 2 || got.NewFoodGrops[0] != wantNew[0] || got.NewFoodGrops[1] != wantNew[1] {
		t.Errorf("new food groups: got %+v, want %+v", got.NewFoodGrops, wantNew)
	}

	if unit := llm.schema.Properties["ingredients"].Items.Properties["unit"]; len(unit.Enum) == 0 {
		t.Error("unit should be restricted to the stored measurement units")
	}
}

func TestParseRecipeUsesRememberedAlias(t *testing.T) {
	s := newTestStorage(t)
	must(t, (&ingredientService{s}).AddAlias(&dto.CreateIngredientAliasRequest{Alias: "Passierte Tomaten", IngredientID: seedTomato}))

	svc := &OpenAIService{storage: s, llm: &fakeCompleter{content: `{
	  "title": "Sauce", "description": "", "cook_time": 0, "servings": 0, "meal_type": "unknown",
	  "instructions": [], "new_food_groups": [],
	  "ingredients": [{"original_text": "passierte Tomaten", "name": "Tomato Passata", "existing_ingredient": "",
	    "quantity": 500, "unit": "gram", "description": "", "measurement_units": [], "food_groups": []}]
	}`}}

	got, err := svc.ParseRecipe(context.Background(), "irrelevant")
	must(t, err)
	if len(got.Ingredients) != 1 || got.Ingredients[0].IngredientID != seedTomato || got.Ingredients[0].MatchedBy != matchedByAlias {
		t.Fatalf("alias should resolve to Tomato, got %+v / new %+v", got.Ingredients, got.NewIngredients)
	}
	if got.Recipe.MealType == 0 {
		t.Error("unknown meal type should fall back to a valid one")
	}
}

func TestParseRecipeErrors(t *testing.T) {
	s := newTestStorage(t)
	refused := &OpenAIService{storage: s, llm: &fakeCompleter{err: ErrRecipeRefused}}
	if _, err := refused.ParseRecipe(context.Background(), "x"); !errors.Is(err, ErrRecipeRefused) {
		t.Errorf("want ErrRecipeRefused, got %v", err)
	}
	garbage := &OpenAIService{storage: s, llm: &fakeCompleter{content: "```json\n{}"}}
	if _, err := garbage.ParseRecipe(context.Background(), "x"); err == nil {
		t.Error("invalid JSON should return an error")
	}
}

func TestAddAlias(t *testing.T) {
	s := newTestStorage(t)
	svc := &ingredientService{s}

	must(t, svc.AddAlias(&dto.CreateIngredientAliasRequest{Alias: "Garlics", IngredientID: seedGarlic}))
	if _, err := s.Aliases.GetIngredientByAlias("garlic"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("alias equal to the ingredient's own name should be skipped, got %v", err)
	}

	must(t, svc.AddAlias(&dto.CreateIngredientAliasRequest{Alias: "Knoblauchzehen", IngredientID: seedOnion}))
	must(t, svc.AddAlias(&dto.CreateIngredientAliasRequest{Alias: "knoblauchzehen", IngredientID: seedGarlic}))
	ing, err := s.Aliases.GetIngredientByAlias("knoblauchzehen")
	must(t, err)
	if ing.IngredientID != seedGarlic {
		t.Errorf("re-adding an alias should repoint it, got ingredient %d", ing.IngredientID)
	}

	err = svc.AddAlias(&dto.CreateIngredientAliasRequest{Alias: "x", IngredientID: 999999})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("unknown ingredient should be not found, got %v", err)
	}
}

func TestConvertQuantity(t *testing.T) {
	if q, ok := convertQuantity(2, "tablespoon", "teaspoon"); !ok || q != 6 {
		t.Errorf("2 tablespoon -> teaspoon: got %v %v", q, ok)
	}
	if q, ok := convertQuantity(500, "gram", "kilogram"); !ok || q != 0.5 {
		t.Errorf("500 gram -> kilogram: got %v %v", q, ok)
	}
	if _, ok := convertQuantity(1, "gram", "cup"); ok {
		t.Error("mass and volume must not convert")
	}
	if _, ok := convertQuantity(1, "piece", "gram"); ok {
		t.Error("piece must not convert")
	}
}

func TestRecipeCreateRejectsInvalidIngredients(t *testing.T) {
	s := newTestStorage(t)
	svc := &recipeService{s}
	recipe := dto.RecipeRequest{Title: "Test", MealType: 1}
	cases := map[string][]dto.IngredientRequest{
		"duplicate":  {{IngredientID: seedGarlic, MeasurementUnitID: unitGram, Quantity: 1}, {IngredientID: seedGarlic, MeasurementUnitID: unitTeaspoon, Quantity: 1}},
		"missing id": {{IngredientID: 0, MeasurementUnitID: unitGram, Quantity: 1}},
		"unknown":    {{IngredientID: 999999, MeasurementUnitID: unitGram, Quantity: 1}},
		"negative":   {{IngredientID: seedGarlic, MeasurementUnitID: unitGram, Quantity: -1}},
	}
	for name, ingredients := range cases {
		_, err := svc.Create(&dto.CreateRecipeRequest{Recipe: recipe, Ingredients: ingredients})
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("%s: want a validation error, got %v", name, err)
		}
	}
	_, err := svc.Create(&dto.CreateRecipeRequest{Recipe: dto.RecipeRequest{Title: " "}})
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("empty title: want a validation error, got %v", err)
	}
}

func TestIngredientCreateRejectsUnknownFoodGroupAndDuplicates(t *testing.T) {
	s := newTestStorage(t)
	svc := &ingredientService{s}
	var validationErr *ValidationError

	_, err := svc.Create(&dto.CreateIngredientRequest{Name: "Chickpea", FoodGroups: []string{"No Such Group"}})
	if !errors.As(err, &validationErr) {
		t.Errorf("unknown food group: want a validation error, got %v", err)
	}
	_, err = svc.Create(&dto.CreateIngredientRequest{Name: "garlic"})
	if !errors.As(err, &validationErr) {
		t.Errorf("duplicate name (case-insensitive): want a validation error, got %v", err)
	}

	created, err := svc.Create(&dto.CreateIngredientRequest{Name: "Chickpea", Description: "A legume.", FoodGroups: []string{"vegetables"}})
	must(t, err)
	stored, err := s.Ingredients.GetByID(created.IngredientID)
	must(t, err)
	if stored.Description != "A legume." || len(created.FoodGroups) != 1 {
		t.Errorf("description and food group should be saved, got %+v / %+v", stored, created)
	}
}
