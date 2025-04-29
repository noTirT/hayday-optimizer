package api

import (
	"github.com/google/uuid"
	"github.com/noTirT/hayday-optimizer/models"
)

type Optimizer struct {
	allGoods                   models.HayDayGoodList
	goodsMap                   map[uuid.UUID]models.HayDayGood
	currentMostProfitableGoods models.HayDayGoodList
}

func NewOptimizer(allGoods models.HayDayGoodList) *Optimizer {
	goodsMap := make(map[uuid.UUID]models.HayDayGood)
	for _, good := range allGoods {
		goodsMap[good.ID] = good
	}

	return &Optimizer{
		allGoods: allGoods,
		goodsMap: goodsMap,
		// For internal state management
		currentMostProfitableGoods: models.HayDayGoodList{},
	}
}

// Main process of optimization
func (o *Optimizer) GetOptimizedPlan(availableGoods models.HayDayGoodList) models.HayDayGoodList {
	o.selectMostProfitablePerSource(availableGoods)

	o.filterOutBaseProductsInIngredientChain()

	o.removeIngredientsOfHigherPricedProducts()

	o.filterGoods(func(good models.HayDayGood) bool {
		return good.Source != "Feed Mill"
	})

	o.removeProductsWithSourceConflicts()

	o.filterGoods(func(good models.HayDayGood) bool {
		return good.MaxPrice > 0
	})

	return o.currentMostProfitableGoods
}

// Return all the most profitable goods but only one per different source
func (o *Optimizer) selectMostProfitablePerSource(availableGoods models.HayDayGoodList) {
	goodsBySource := groupGoodsBySource(availableGoods)

	var mostProfitable models.HayDayGoodList

	for _, goods := range goodsBySource {
		if len(goods) == 0 {
			continue
		}
		mostProfitableGood := goods[0]
		maxProfit := mostProfitableGood.MaxPrice

		for _, good := range goods[1:] {
			profit := good.MaxPrice
			if profit > maxProfit {
				mostProfitableGood = good
				maxProfit = profit
			}
		}
		mostProfitable = append(mostProfitable, mostProfitableGood)
	}

	sortGoodsByPriceDescending(mostProfitable)

	o.currentMostProfitableGoods = mostProfitable
}

// Remove all base products from the List that are also ingredients of non-base products
func (o *Optimizer) filterOutBaseProductsInIngredientChain() {
	baseProductIDsToRemove := make(map[uuid.UUID]bool)

	for _, profitable := range o.currentMostProfitableGoods {
		if !isBaseProduct(profitable) {
			o.findBaseProductsInChain(profitable, baseProductIDsToRemove)
		}
	}

	o.filterGoods(func(good models.HayDayGood) bool {
		return !isBaseProduct(good) || !baseProductIDsToRemove[good.ID]
	})
}

// Recursively find nested ingredients that are base products
func (o *Optimizer) findBaseProductsInChain(good models.HayDayGood, baseProductIDs map[uuid.UUID]bool) {
	// Check direct ingredients
	for _, ingredient := range good.Ingredients {
		ingredientGood, exists := o.goodsMap[ingredient.ProductID]
		if !exists {
			continue
		}

		// If this ingredient is a base product, mark it
		if isBaseProduct(ingredientGood) {
			baseProductIDs[ingredientGood.ID] = true
		} else {
			// If it's not a base product, check its ingredients recursively
			o.findBaseProductsInChain(ingredientGood, baseProductIDs)
		}
	}
}

// Remove all products that are ingredients of products with higher price
func (o *Optimizer) removeIngredientsOfHigherPricedProducts() {
	// Track which products are ingredients of higher-priced products
	ingredientsToRemove := make(map[uuid.UUID]bool)

	sortGoodsByPriceDescending(o.currentMostProfitableGoods)

	// For each profitable good, check if any of the other profitable goods
	// are in its ingredient chain
	for _, profitable := range o.currentMostProfitableGoods {
		o.markIngredientsInChain(profitable, ingredientsToRemove)
	}

	// Filter out the products that are ingredients of higher-priced products
	o.filterGoods(func(good models.HayDayGood) bool {
		return !ingredientsToRemove[good.ID]
	})
}

func (o *Optimizer) markIngredientsInChain(good models.HayDayGood, ingredientsToRemove map[uuid.UUID]bool) {
	// Check direct ingredients
	for _, ingredient := range good.Ingredients {
		// If this ingredient is in our profitable goods, mark it for removal
		if _, exists := o.goodsMap[ingredient.ProductID]; exists {
			ingredientsToRemove[ingredient.ProductID] = true
		}

		// Continue checking the ingredient chain
		ingredientGood, exists := o.goodsMap[ingredient.ProductID]
		if exists && len(ingredientGood.Ingredients) > 0 {
			o.markIngredientsInChain(ingredientGood, ingredientsToRemove)
		}
	}
}

// Remove Products where their source is also the source of the ingredients of higher profitable products
func (o *Optimizer) removeProductsWithSourceConflicts() {
	sortGoodsByPriceDescending(o.currentMostProfitableGoods)

	// Keep track of sources required by ingredients of higher-priced products
	requiredSources := make(map[string]bool)

	// Products to remove due to source conflicts
	productsToRemove := make(map[uuid.UUID]bool)

	// Start with the highest priced product
	for i, highPricedProduct := range o.currentMostProfitableGoods {
		// Skip products already marked for removal
		if productsToRemove[highPricedProduct.ID] {
			continue
		}

		// Get all sources needed by this product's ingredient chain
		ingredientSources := o.getIngredientSources(highPricedProduct, make(map[uuid.UUID]bool))

		// Add these sources to our required sources
		for source := range ingredientSources {
			requiredSources[source] = true
		}

		// Check lower-priced products
		for j := i + 1; j < len(o.currentMostProfitableGoods); j++ {
			lowerPricedProduct := o.currentMostProfitableGoods[j]

			// If this product's source is needed by a higher-priced product's ingredients,
			// mark it for removal
			if requiredSources[lowerPricedProduct.Source] {
				productsToRemove[lowerPricedProduct.ID] = true
			}
		}
	}

	// Filter out products with source conflicts
	o.filterGoods(func(good models.HayDayGood) bool {
		return !productsToRemove[good.ID]
	})
}

func (o *Optimizer) getIngredientSources(good models.HayDayGood, visited map[uuid.UUID]bool) map[string]bool {
	// Prevent infinite recursion with cycles
	if visited[good.ID] {
		return make(map[string]bool)
	}
	visited[good.ID] = true

	sources := make(map[string]bool)

	// Check all ingredients
	for _, ingredient := range good.Ingredients {
		ingredientGood, exists := o.goodsMap[ingredient.ProductID]
		if !exists {
			continue
		}

		// Add this ingredient's source
		sources[ingredientGood.Source] = true

		// If this ingredient has its own ingredients, get their sources too
		if len(ingredientGood.Ingredients) > 0 {
			subSources := o.getIngredientSources(ingredientGood, visited)
			for source := range subSources {
				sources[source] = true
			}
		}
	}

	return sources
}

func (o *Optimizer) filterGoods(keepFn func(good models.HayDayGood) bool) {
	var result models.HayDayGoodList
	for _, good := range o.currentMostProfitableGoods {
		if keepFn(good) {
			result = append(result, good)
		}
	}
	o.currentMostProfitableGoods = result
}
