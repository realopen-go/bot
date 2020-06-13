package crawler

import (
	"github.com/gocolly/colly"
	"github.com/sluggishhackers/go-realopen/statusmanager"
	"github.com/sluggishhackers/go-realopen/store"
)

type ICrawler interface {
	FetchBill(billID string, arg1 string, arg2 string)
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
	statusmanager  statusmanager.Istatusmanager
}

func New(s store.IStore, sm statusmanager.Istatusmanager) ICrawler {
	c := &Crawler{
		defaultCrawler: newDefaultCrawler(),
		store:          s,
		statusmanager:  sm,
	}
	c.billsCrawler = c.NewBillsCrawler()
	c.billCrawler = c.NewBillCrawler()
	return c
}
