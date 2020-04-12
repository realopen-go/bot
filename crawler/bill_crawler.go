package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sluggishhackers/realopen.go/models"
	"github.com/sluggishhackers/realopen.go/rmtstor"
)

type BillResultFormat struct {
	FileList []models.File `json:"dntcFileList"`
}

func NewParamsPostBill(ID string, IfrmpPrcsRstrNo string) map[string]string {
	return map[string]string{
		"rqestNo":         ID,
		"ifrmpPrcsRstrNo": IfrmpPrcsRstrNo,
	}
}

func (c *Crawler) NewBillCrawler() *colly.Collector {
	crawler := c.defaultCrawler.Clone()
	crawler.Limit(&colly.LimitRule{Parallelism: 2})

	crawler.OnResponse(func(r *colly.Response) {
		billID := r.Ctx.Get("billId")
		body := string(r.Body)
		result := body[strings.Index(body, "var result"):strings.Index(body, "//var naviInfo")]
		data := strings.TrimRight(strings.TrimSpace(result[strings.Index(result, "{"):]), ";")

		billResultFormat := &BillResultFormat{}

		err := json.Unmarshal([]byte(data), billResultFormat)
		if err != nil {
			log.Fatal(err)
		}

		fileCount := len(billResultFormat.FileList)
		ch := make(chan string, fileCount)
		// go func(ch chan string) {

		if len(billResultFormat.FileList) > 0 {
			for _, f := range billResultFormat.FileList {
				fmt.Println(fmt.Sprintf("Download : %s", f.UploadFileOrglNm))
				go c.DownloadFile(billID, f, ch)
			}
			// }(ch)

			downloadFinishedCount := 0
			for channel := range ch {
				downloadFinishedCount++

				// Download Message
				fmt.Println(channel)

				if downloadFinishedCount == fileCount {
					c.statmanager.SetFileStatus(billID, true)
					close(ch)
				}
			}
		} else {
			c.statmanager.SetFileStatus(billID, false)
			close(ch)
		}
	})

	return crawler
}

func (c *Crawler) FetchBill(bill *models.Bill) {
	fmt.Println(fmt.Sprintf("Start to fetch a bill: %s", bill.ID))
	// the key of "url" into the context of the request
	c.billCrawler.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("billId", bill.ID)
	})
	err := c.billCrawler.Post("https://www.open.go.kr/pa/billing/openBilling/openBillingDntcDtl.do", NewParamsPostBill(bill.ID, bill.IfrmpPrcsRstrNo))
	if err != nil {
		fmt.Println("😡 Error to fetch a bill")
		log.Fatal(err)
	}
}

func (c *Crawler) DownloadFile(billID string, file models.File, ch chan string) {
	downloader := c.defaultCrawler.Clone()

	downloader.OnResponse(func(r *colly.Response) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		filePath := fmt.Sprintf("%s/%s/%s", wd, rmtstor.REALOPEN_DATA_DIR, billID)

		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Errorf("😡 Error to create a new directory")
		}

		err = r.Save(fmt.Sprintf("%s/%s", filePath, file.UploadFileOrglNm))
		if err != nil {
			fmt.Errorf("😡 Error to save a file to download", err)
			ch <- fmt.Sprintf("Failed: %s", file.UploadFileOrglNm)
		} else {
			ch <- fmt.Sprintf("Succeeded: %s", file.UploadFileOrglNm)
		}
	})

	downloader.Visit(fmt.Sprintf("https://www.open.go.kr/util/FileDownload.down?atchFileUploadNo=%s&fileSn=%s", file.AtchFileUploadNo, file.FileSn))
}
