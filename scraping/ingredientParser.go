package scraping

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/noTirT/hayday-optimizer/base"
	"github.com/noTirT/hayday-optimizer/models"
)

type IngredientParser struct {
	NameToID map[string]uuid.UUID
}

func NewIngredientParser(nameToID map[string]uuid.UUID) *IngredientParser {
	return &IngredientParser{
		NameToID: nameToID,
	}
}

func (ip *IngredientParser) Parse(rawNeeds string) ([]models.Ingredient, error) {
	if rawNeeds == "N/A" || rawNeeds == "" {
		return nil, nil
	}

	var ingredients []models.Ingredient

	ingredientRegex := regexp.MustCompile(`([^()]+)\s*\((\d+)\)`)
	matches := ingredientRegex.FindAllStringSubmatch(rawNeeds, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		productName := base.CapializeWordsOfString(strings.TrimSpace(match[1]))
		amount, err := strconv.Atoi(match[2])
		if err != nil {
			log.Printf("Error parsing amount for %s: %v\n", productName, err)
			continue
		}
		productID, exists := ip.NameToID[productName]
		if !exists {
			log.Printf("Warning: Could not find product ID for '%s'", productName)
			continue
		}

		ingredients = append(ingredients, models.Ingredient{
			ProductID:   productID,
			ProductName: productName,
			Amount:      amount,
		})
	}

	return ingredients, nil
}
