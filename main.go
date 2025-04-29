package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
	"github.com/noTirT/hayday-optimizer/api"
	"github.com/noTirT/hayday-optimizer/base"
	"github.com/noTirT/hayday-optimizer/models"
	"github.com/noTirT/hayday-optimizer/scraping"
)

func main() {
	fetch := flag.Bool("fetch", false, "Re-Fetch data from the website")
	flag.Parse()

	hayDayFilemanager, err := base.NewJsonFileManager[models.HayDayGoodList]("./data")
	if err != nil {
		log.Fatalf("Failed to create file manager: %v", err)
	}

	if *fetch {
		fetchGoods(hayDayFilemanager)
	}

	apiIP := "localhost"
	apiPort := 5000
	serverAddr := fmt.Sprintf("%s:%d", apiIP, apiPort)

	log.Printf("Starting API server at: %s\n", serverAddr)

	r := http.NewServeMux()

	goodsRepository := api.NewGoodsRepository(hayDayFilemanager)
	goodsController := api.NewGoodsController(goodsRepository)
	goodsController.Init(r)

	log.Println("API server started")

	log.Fatal(http.ListenAndServe(serverAddr, r))
}

func fetchGoods(hayDayFilemanager *base.FileManager[models.HayDayGoodList]) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.ExecPath("/usr/bin/chromium-browser"))
	url := "https://hayday.fandom.com/wiki/Goods_List"

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	scraper := scraping.NewHayDayScraper(ctx, url)
	goods, err := scraper.Scrape()
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	log.Printf("Rows scraped: %d\n", len(goods))

	goodsFileName := "goods.json"

	if hayDayFilemanager.Exists(goodsFileName) {
		hayDayFilemanager.Delete(goodsFileName)
	}

	hayDayFilemanager.Write(goodsFileName, goods)

}
