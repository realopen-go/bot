package crawler

import (
	"github.com/gocolly/colly"
	"github.com/sluggishhackers/realopen.go/models"
	"github.com/sluggishhackers/realopen.go/statmanager"
	"github.com/sluggishhackers/realopen.go/store"
)

type ICrawler interface {
	FetchBill(bill *models.Bill)
	FetchBills(dateFrom string, dateTo string)
	NewBillsCrawler() *colly.Collector
	NewBillCrawler() *colly.Collector
}

type Crawler struct {
	defaultCrawler *colly.Collector
	billCrawler    *colly.Collector
	billsCrawler   *colly.Collector
	downloader     *colly.Collector
	store          store.IStore
	statmanager    statmanager.Istatmanager
}

func New(s store.IStore, sm statmanager.Istatmanager) ICrawler {
	c := &Crawler{
		defaultCrawler: newDefaultCrawler(),
		store:          s,
		statmanager:    sm,
	}

	c.billsCrawler = c.NewBillsCrawler()
	c.billCrawler = c.NewBillCrawler()

	return c
}
