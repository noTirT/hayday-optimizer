package api

import (
	"log"

	"github.com/noTirT/hayday-optimizer/base"
	"github.com/noTirT/hayday-optimizer/models"
)

type GoodsRepository struct {
	goods models.HayDayGoodList
}

func NewGoodsRepository(fileManager *base.FileManager[models.HayDayGoodList]) *GoodsRepository {
	goods, err := fileManager.Read("goods.json")
	if err != nil {
		log.Fatal("Error reading in goods")
	}
	return &GoodsRepository{
		goods: goods,
	}
}

func (repo *GoodsRepository) GetAllGoods() models.HayDayGoodList {
	return repo.goods
}

func (repo *GoodsRepository) GetGoodByName(name string) (*models.HayDayGood, error) {
	for _, good := range repo.goods {
		if good.Name == name {
			return &good, nil
		}
	}

	return nil, base.ErrNoGoodByNameFound
}

func (repo *GoodsRepository) GetGoodsByLevel(level int) models.HayDayGoodList {
	var result models.HayDayGoodList
	for _, good := range repo.goods {
		if good.RequiredLevel <= level {
			result = append(result, good)
		}
	}
	return result
}
