package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type GoodsController struct {
	repo *GoodsRepository
}

func NewGoodsController(repo *GoodsRepository) *GoodsController {
	return &GoodsController{
		repo: repo,
	}
}

func (a *GoodsController) Init(router *http.ServeMux) {
	router.HandleFunc("GET /goods", a.getGoods)
	router.HandleFunc("GET /goods/{name}", a.getGoodByName)
	router.HandleFunc("GET /goods/level/{level}", a.getGoodsByLevel)
	router.HandleFunc("GET /goods/strategy/{level}", a.getMostProfitableGoods)
}

func (a *GoodsController) getGoods(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(a.repo.GetAllGoods())
}

func (a *GoodsController) getGoodByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	good, err := a.repo.GetGoodByName(name)
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(good)
}

func (a *GoodsController) getGoodsByLevel(w http.ResponseWriter, r *http.Request) {
	level := r.PathValue("level")
	parsedLevel, err := strconv.Atoi(level)
	if err != nil {
		return
	}

	goods := a.repo.GetGoodsByLevel(parsedLevel)
	json.NewEncoder(w).Encode(goods)
}

func (a *GoodsController) getMostProfitableGoods(w http.ResponseWriter, r *http.Request) {
	level := r.PathValue("level")

	parsedLevel, err := strconv.Atoi(level)
	if err != nil {
		return
	}

	optimizer := NewOptimizer(a.repo.GetAllGoods())

	availableGoods := a.repo.GetGoodsByLevel(parsedLevel)

	json.NewEncoder(w).Encode(optimizer.GetOptimizedPlan(availableGoods))
}
