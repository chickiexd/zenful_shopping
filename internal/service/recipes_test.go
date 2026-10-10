package service

import (
	"os"
	"strings"
	"testing"

	"github.com/chickiexd/zenful_shopping/internal/dto"
	"github.com/chickiexd/zenful_shopping/internal/logger"
	"github.com/chickiexd/zenful_shopping/internal/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	unitGram     = 1
	unitTeaspoon = 5
)

// newTestStorage returns storage bound to a transaction that is rolled back
// when the test ends. Requires TEST_DB_DSN pointing at a migrated database.
func newTestStorage(t *testing.T) *store.Storage {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set")
	}
	logger.Init()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormlogger.Discard})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })
	storage := store.NewStorage(tx)
	return &storage
}

type fixture struct {
	listID       uint
	ingredientID uint
	recipeID     uint
}

// newFixture creates a shopping list linked to a fresh food group and ingredient,
// plus a recipe using 200g of that ingredient.
func newFixture(t *testing.T, s *store.Storage) fixture {
	t.Helper()
	list := &store.ShoppingList{Name: "Test List", Color: "WHITE"}
	fg := &store.FoodGroup{Name: "Test Food Group"}
	ing := &store.Ingredient{Name: "Test Ingredient"}
	for _, v := range []any{list, fg, ing} {
		if err := s.DB.Create(v).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
	must(t, s.DB.Exec("INSERT INTO shopping_list_food_groups (shopping_list_id, food_group_id) VALUES (?, ?)", list.ShoppingListID, fg.FoodGroupID).Error)
	must(t, s.DB.Exec("INSERT INTO ingredient_food_groups (ingredient_id, food_group_id) VALUES (?, ?)", ing.IngredientID, fg.FoodGroupID).Error)

	recipe := &store.Recipe{Title: "Test Recipe", MealTypeID: 1}
	must(t, s.DB.Create(recipe).Error)
	must(t, s.DB.Create(&store.RecipeIngredient{
		RecipeID: recipe.RecipeID, IngredientID: ing.IngredientID, MeasurementUnitID: unitGram, Quantity: 200,
	}).Error)
	return fixture{list.ShoppingListID, ing.IngredientID, recipe.RecipeID}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func listItems(t *testing.T, s *store.Storage, listID uint) []store.ShoppingListItem {
	t.Helper()
	items, err := s.ShoppingLists.GetItemsByShoppingListID(listID)
	must(t, err)
	return items
}

func TestAddRecipeTwiceSumsQuantity(t *testing.T) {
	s := newTestStorage(t)
	f := newFixture(t, s)
	svc := &recipeService{s}

	must(t, svc.AddToShoppingList(f.recipeID))
	must(t, svc.AddToShoppingList(f.recipeID))

	items := listItems(t, s, f.listID)
	if len(items) != 1 || items[0].Quantity != 400 {
		t.Fatalf("want one item with 400, got %+v", items)
	}
}

func TestAddRecipeSkipsPantryIngredients(t *testing.T) {
	s := newTestStorage(t)
	f := newFixture(t, s)
	must(t, (&pantryService{s}).Create(f.ingredientID))

	must(t, (&recipeService{s}).AddToShoppingList(f.recipeID))

	if items := listItems(t, s, f.listID); len(items) != 0 {
		t.Fatalf("want no items, got %+v", items)
	}
}

func TestAddRecipeKeepsDifferentUnitsSeparate(t *testing.T) {
	s := newTestStorage(t)
	f := newFixture(t, s)
	_, err := (&ingredientService{s}).AddToShoppingList(&dto.AddIngredientToShoppingListRequest{
		ShoppingListID: f.listID, IngredientID: f.ingredientID, MeasurementUnitID: unitTeaspoon, Quantity: 2,
	})
	must(t, err)

	must(t, (&recipeService{s}).AddToShoppingList(f.recipeID))

	items := listItems(t, s, f.listID)
	if len(items) != 2 {
		t.Fatalf("want 2 items (grams and teaspoons), got %+v", items)
	}
}

func TestRemoveRecipeRestoresQuantity(t *testing.T) {
	s := newTestStorage(t)
	f := newFixture(t, s)
	svc := &recipeService{s}

	must(t, svc.AddToShoppingList(f.recipeID))
	must(t, svc.AddToShoppingList(f.recipeID))
	must(t, svc.RemoveFromShoppingList(f.recipeID))

	items := listItems(t, s, f.listID)
	if len(items) != 1 || items[0].Quantity != 200 {
		t.Fatalf("want one item with 200, got %+v", items)
	}

	must(t, svc.RemoveFromShoppingList(f.recipeID))
	if items := listItems(t, s, f.listID); len(items) != 0 {
		t.Fatalf("want item deleted, got %+v", items)
	}
}

func TestAddUnknownRecipeReturnsNotFound(t *testing.T) {
	s := newTestStorage(t)
	err := (&recipeService{s}).AddToShoppingList(999999)
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("want ErrRecordNotFound, got %v", err)
	}
}

func TestPrepareSyncData(t *testing.T) {
	s := newTestStorage(t)
	f := newFixture(t, s)
	must(t, (&recipeService{s}).AddToShoppingList(f.recipeID))

	data, err := (&keepSyncService{s}).prepareData()
	must(t, err)
	for _, l := range data {
		if l.Title == "Test List" {
			if len(l.Ingredients) != 1 || l.Ingredients[0].Name != "Test Ingredient" || !strings.HasPrefix(l.Ingredients[0].Quantity, "200 g") {
				t.Fatalf("unexpected sync data: %+v", l.Ingredients)
			}
			return
		}
	}
	t.Fatal("test list missing from sync data")
}

func TestRecipePreloadsMatchForeignKeys(t *testing.T) {
	s := newTestStorage(t)
	f := newFixture(t, s)

	recipe, err := s.Recipes.GetByID(f.recipeID)
	must(t, err)
	if recipe.MealType.MealTypeID != recipe.MealTypeID {
		t.Errorf("meal type preload: want %d, got %d", recipe.MealTypeID, recipe.MealType.MealTypeID)
	}
	ri := recipe.RecipeIngredients[0]
	if ri.Ingredient.IngredientID != f.ingredientID {
		t.Errorf("ingredient preload: want %d, got %d", f.ingredientID, ri.Ingredient.IngredientID)
	}
	if ri.MeasurementUnit.MeasurementUnitID != unitGram {
		t.Errorf("unit preload: want %d, got %d", unitGram, ri.MeasurementUnit.MeasurementUnitID)
	}
}
