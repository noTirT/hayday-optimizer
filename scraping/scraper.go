package scraping

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/noTirT/hayday-optimizer/base"
	"github.com/noTirT/hayday-optimizer/models"
)

type HayDayScraper struct {
	ctx      context.Context
	url      string
	goods    models.HayDayGoodList
	nameToID map[string]uuid.UUID
}

func NewHayDayScraper(ctx context.Context, url string) *HayDayScraper {
	return &HayDayScraper{
		ctx:      ctx,
		url:      url,
		nameToID: make(map[string]uuid.UUID),
	}
}

func (s *HayDayScraper) Scrape() (models.HayDayGoodList, error) {
	if err := chromedp.Run(s.ctx, chromedp.Navigate(s.url)); err != nil {
		log.Fatalf("Failed to navigate to page: %v", err)
	}

	rawRows, err := s.extractTableData()
	if err != nil {
		return nil, err
	}

	if err := s.processRawData(rawRows); err != nil {
		return nil, err
	}

	if err := s.processIngredients(); err != nil {
		return nil, err
	}

	return s.goods, nil
}

func (s *HayDayScraper) extractTableData() ([]models.RawTableRow, error) {
	log.Println("Extracting table data...")
	var rawRows []models.RawTableRow

	if err := chromedp.Run(s.ctx, chromedp.Evaluate(`
		Array.from(document.querySelectorAll('table tbody tr')).map(row => {
			const cells = row.querySelectorAll('td');
			return {
				name: cells[0].querySelector('a')?.textContent.trim() || '',
				level: parseInt(cells[1].textContent.trim()) || 0,
				price: parseInt(cells[2].textContent.trim()) || 0,
				timeStr: cells[3].textContent.trim(),
				xp: parseInt(cells[4].textContent.trim()) || 0,
				rawNeeds: cells[5].textContent.trim(),
				source: cells[6].textContent.trim()
			};
		})
	`, &rawRows)); err != nil {
		return nil, base.ErrExtractingTableData
	}

	log.Printf("Scanned %d raw rows:\n", len(rawRows))
	return rawRows, nil
}

func (s *HayDayScraper) processRawData(rawRows []models.RawTableRow) error {
	log.Println("Processing raw table rows...")
	s.goods = make([]models.HayDayGood, 0, len(rawRows))

	for _, row := range rawRows {

		if row.Name == "" {
			continue
		}

		good := models.HayDayGood{
			ID:             uuid.New(),
			Name:           base.CapializeWordsOfString(row.Name),
			RequiredLevel:  row.Level,
			MaxPrice:       row.Price,
			GainedXP:       row.XP,
			RawIngredients: row.RawNeeds,
		}

		duration, err := base.ParseDurationString(row.TimeStr)
		if err != nil {
			log.Printf("Error parsing time for %s: %v\n", row.Name, err)
		}
		good.ProductionTime = duration

		source := strings.SplitN(row.Source, "(", 2)[0]
		good.Source = strings.TrimSpace(source)

		s.goods = append(s.goods, good)
		s.nameToID[good.Name] = good.ID
	}

	return nil
}

func (s *HayDayScraper) processIngredients() error {
	log.Println("Processing ingredients...")
	ingredientParser := NewIngredientParser(s.nameToID)

	for i := range s.goods {
		ingredients, err := ingredientParser.Parse(s.goods[i].RawIngredients)
		if err != nil {
			log.Printf("Error parsing ingredients for %s: %v\n", s.goods[i].Name, err)
			continue
		}
		if len(ingredients) > 0 && ingredients[0].ProductID != s.goods[i].ID {
			s.goods[i].Ingredients = ingredients
		}

		s.goods[i].RawIngredients = ""
	}
	return nil
}
