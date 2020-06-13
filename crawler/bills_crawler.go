package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/sluggishhackers/go-realopen/models"
)

// https://www.open.go.kr/pa/billing/openBilling/openBillingList.do

type BillListResultFormat struct {
	Result struct {
		List []models.Bill `json:"list"`
	} `json:"result"`
}

type BillsCountResultFormat struct {
	Result struct {
		Vo struct {
			TotalCount int `json:"totalPage"`
		}
	} `json:"result"`
}

type Params map[string]string

func NewParamsPostBills(
	dateFrom string, // 접수일자 시작일
	dateTo string, // 접수일자 종료일
	page int,
) map[string]string {
	return map[string]string{
		"stRceptPot": dateFrom,
		"edRceptPot": dateTo,
		"searchYn":   "Y",
		"chkDate":    "nonClass",
		"rowPage":    "10",
		"totalPage":  string(page),
		"moveStatus": "L",
		"viewPage":   "1",
	}
}

func (c *Crawler) NewBillsCrawler() *colly.Collector {
	crawler := c.defaultCrawler.Clone()

	crawler.OnResponse(func(r *colly.Response) {
		billListResultFormat := &BillListResultFormat{}

		err := json.Unmarshal(r.Body, billListResultFormat)
		if err != nil {
			log.Fatal(err)
		}

		for _, b := range billListResultFormat.Result.List {
			fmt.Printf("Save the bill from the list: %s\n", b.ID)
			c.store.SaveBill(b)
		}
	})

	return crawler
}

func (c *Crawler) FetchBills(
	dateFrom string,
	dateTo string,
) {
	fmt.Printf("Start to crawling from %s to %s\n", dateFrom, dateTo)

	// 1. 전체 TotalPage 를 먼저 가져온 후
	// 2. TotalPage를 기반으로 전체 청구 목록을 크롤링한다
	totalCountCrawler := c.billsCrawler.Clone()
	totalCountCrawler.OnResponse(func(r *colly.Response) {
		billsCountResultFormat := &BillsCountResultFormat{}

		err := json.Unmarshal(r.Body, billsCountResultFormat)
		if err != nil {
			log.Fatal(err)
		}

		// totalCount := billsCountResultFormat.Result.Vo.TotalCount
		for i := 0; i < 1; i++ {
			err = c.billsCrawler.Post("https://www.open.go.kr/pa/billing/openBilling/openBillingSrchList.ajax", NewParamsPostBills(dateFrom, dateTo, i))
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	err := totalCountCrawler.Post("https://www.open.go.kr/pa/billing/openBilling/openBillingSrchList.ajax", NewParamsPostBills(dateFrom, dateTo, 0))
	if err != nil {
		log.Fatal(err)
	}
}
