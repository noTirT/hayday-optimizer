package api

import (
	"sort"

	"github.com/noTirT/hayday-optimizer/models"
)

func groupGoodsBySource(goods models.HayDayGoodList) map[string]models.HayDayGoodList {
	sourceMap := make(map[string]models.HayDayGoodList)
	for _, good := range goods {
		sourceMap[good.Source] = append(sourceMap[good.Source], good)
	}
	return sourceMap
}

func sortGoodsByPriceDescending(goods models.HayDayGoodList) {
	sort.Slice(goods, func(i, j int) bool {
		return (goods)[i].MaxPrice > (goods)[j].MaxPrice
	})
}

func isBaseProduct(good models.HayDayGood) bool {
	return len(good.Ingredients) == 0 || good.Ingredients == nil
}
