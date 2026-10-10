package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/chickiexd/zenful_shopping/internal/store"
	"gorm.io/gorm"
)

// normalizeName reduces an ingredient name to a comparable key: lowercase,
// punctuation stripped, whitespace collapsed and the last word singularized
// with simple English rules. German plurals are left alone; aliases cover those.
func normalizeName(name string) string {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, name)
	words := strings.Fields(cleaned)
	if len(words) == 0 {
		return ""
	}
	words[len(words)-1] = singularize(words[len(words)-1])
	return strings.Join(words, " ")
}

func singularize(word string) string {
	if len(word) <= 3 {
		return word
	}
	switch {
	case strings.HasSuffix(word, "ies"):
		return strings.TrimSuffix(word, "ies") + "y"
	case strings.HasSuffix(word, "oes"),
		strings.HasSuffix(word, "ches"),
		strings.HasSuffix(word, "shes"),
		strings.HasSuffix(word, "xes"):
		return strings.TrimSuffix(word, "es")
	case strings.HasSuffix(word, "ss"), strings.HasSuffix(word, "us"), strings.HasSuffix(word, "is"):
		return word
	case strings.HasSuffix(word, "s"):
		return strings.TrimSuffix(word, "s")
	}
	return word
}

// ingredientMatcher resolves names coming from the recipe parser to existing ingredients.
type ingredientMatcher struct {
	aliases   *store.AliasRepository
	byName    map[string]store.Ingredient // lowercase exact name
	byNormKey map[string]store.Ingredient // normalizeName(name)
}

func newIngredientMatcher(aliases *store.AliasRepository, catalog []store.Ingredient) *ingredientMatcher {
	m := &ingredientMatcher{
		aliases:   aliases,
		byName:    make(map[string]store.Ingredient, len(catalog)),
		byNormKey: make(map[string]store.Ingredient, len(catalog)),
	}
	for _, ing := range catalog {
		m.byName[strings.ToLower(ing.Name)] = ing
		if key := normalizeName(ing.Name); key != "" {
			if _, taken := m.byNormKey[key]; !taken {
				m.byNormKey[key] = ing
			}
		}
	}
	return m
}

// how an ingredient was matched, reported to the import review
const (
	matchedByModel = "model"
	matchedByAlias = "alias"
	matchedByName  = "name"
)

// resolve returns the matching ingredient and how it was matched, or nil when
// the ingredient is new. Order: the model's pick from the catalog, remembered
// aliases for the original and canonical name, then the normalized canonical name.
func (m *ingredientMatcher) resolve(existing, original, name string) (*store.Ingredient, string, error) {
	if ing, ok := m.byName[strings.ToLower(strings.TrimSpace(existing))]; ok && existing != "" {
		return &ing, matchedByModel, nil
	}
	for _, candidate := range []string{original, name} {
		key := normalizeName(candidate)
		if key == "" {
			continue
		}
		ing, err := m.aliases.GetIngredientByAlias(key)
		if err == nil {
			return ing, matchedByAlias, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
	}
	if ing, ok := m.byNormKey[normalizeName(name)]; ok {
		return &ing, matchedByName, nil
	}
	return nil, "", nil
}

// unit sizes in a common base (gram or milliliter) for merging quantities
var unitSizes = map[string]struct {
	dimension string
	factor    float64
}{
	"gram":       {"mass", 1},
	"kilogram":   {"mass", 1000},
	"milliliter": {"volume", 1},
	"liter":      {"volume", 1000},
	"teaspoon":   {"volume", 5},
	"tablespoon": {"volume", 15},
	"cup":        {"volume", 240},
}

// convertQuantity converts quantity between two units of the same dimension.
func convertQuantity(quantity float64, from, to string) (float64, bool) {
	f, okFrom := unitSizes[strings.ToLower(from)]
	t, okTo := unitSizes[strings.ToLower(to)]
	if !okFrom || !okTo || f.dimension != t.dimension {
		return 0, false
	}
	return quantity * f.factor / t.factor, true
}
