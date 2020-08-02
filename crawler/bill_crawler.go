package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sluggishhackers/go-realopen/models"
	"github.com/sluggishhackers/go-realopen/rmtstor"
	"github.com/sluggishhackers/go-realopen/utils"
)

var PROCESSED_HOST = "https://www.open.go.kr/pa/billing/openBilling/openBillingDntcDtl.do"
var PROCESSING_HOST = "https://www.open.go.kr/pa/billing/openBilling/openBillingDtl.do"

type BillResultFormat struct {
	FileList []models.File `json:"dntcFileList"`
	DtlVo    models.Bill   `json:"dtlVo"`
}

type ProcessingBillResultFormat struct {
	DtlVo models.Bill `json:"dtlVo"`
}

func NewParamsPostBill(ID string, IfrmpPrcsRstrNo string, PrcsStsCd string) map[string]string {
	return map[string]string{
		"keyword":         "",
		"rqestNo":         ID,
		"ifrmpPrcsRstrNo": IfrmpPrcsRstrNo,
		"prcsRstrNo":      IfrmpPrcsRstrNo,
		"prcsStsCd":       PrcsStsCd,
		"hash":            "true",
	}

}

func (c *Crawler) NewBillCrawler() *colly.Collector {
	crawler := c.defaultCrawler.Clone()

	crawler.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		startIndex := strings.Index(body, "var result")
		endIndex := strings.Index(body, "//var naviInfo")

		// 아직 처리중 단계에서는 청구건 상세페이지가 존재하지 않음
		if startIndex == -1 || endIndex == -1 {
			fmt.Println("처리중...")
		} else {
			result := body[startIndex:endIndex]
			data := strings.TrimRight(strings.TrimSpace(result[strings.Index(result, "{"):]), ";")

			billResultFormat := &BillResultFormat{}

			err := json.Unmarshal([]byte(data), billResultFormat)
			if err != nil {
				fmt.Println("Error to Unmarshall Bill Result Format")
				log.Fatal(err)
			}

			// TODO: formatter
			// fmt.Printf("%+v", billResultFormat.DtlVo)
			billResultFormat.DtlVo.RqestPot = strings.ReplaceAll(billResultFormat.DtlVo.RqestPot, ".", "-")

			c.store.SaveBill(billResultFormat.DtlVo)
			c.store.SaveFiles(billResultFormat.DtlVo.ID, billResultFormat.FileList)

			fileCount := len(billResultFormat.FileList)
			ch := make(chan string, fileCount)

			if len(billResultFormat.FileList) > 0 {
				for _, f := range billResultFormat.FileList {
					fmt.Println(fmt.Sprintf("Download : %s", f.UploadFileOrglNm))
					go c.DownloadFile(billResultFormat, f, ch)
				}

				downloadFinishedCount := 0
				for channel := range ch {
					downloadFinishedCount++

					// Download Message
					fmt.Println(channel)

					if downloadFinishedCount == fileCount {
						close(ch)
					}
				}
			} else {
				close(ch)
			}

		}
	})

	return crawler
}

func (c *Crawler) FetchBill(billID string, ifrmpPrcsRstrNo string, prcsStsCd string) {
	fmt.Println(fmt.Sprintf("Start to fetch a bill: %s", billID))

	// the key of "url" into the context of the request
	c.billCrawler.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("billId", billID)
		r.Ctx.Put("ifrmpPrcsRstrNo", ifrmpPrcsRstrNo)
		r.Ctx.Put("prcsStsCd", prcsStsCd)
	})

	err := c.billCrawler.Post("https://www.open.go.kr/pa/billing/openBilling/openBillingDntcDtl.do", NewParamsPostBill(billID, ifrmpPrcsRstrNo, prcsStsCd))
	if err != nil {
		fmt.Println("😡 Error to fetch a bill - 1")
		log.Fatal(err)
	}
}

// TODO: FileDownloader 분리
func (c *Crawler) DownloadFile(bill *BillResultFormat, file models.File, ch chan string) {
	downloader := c.defaultCrawler.Clone()

	downloader.OnResponse(func(r *colly.Response) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		fileDir := fmt.Sprintf("%s/%s/%s", wd, rmtstor.REALOPEN_DATA_DIR, utils.MakeFileDir(&bill.DtlVo))

		err = os.Mkdir(fileDir, os.ModePerm)
		if err != nil {
			fmt.Errorf("😡 Error to create a new directory")
		}

		filePath := fmt.Sprintf("%s/%s", fileDir, utils.MakeFileName(&bill.DtlVo, file))
		err = r.Save(filePath)
		if err != nil {
			fmt.Errorf("😡 Error to save a file to download", err)
			ch <- fmt.Sprintf("Failed: %s", file.UploadFileOrglNm)
		} else {
			ch <- fmt.Sprintf("Succeeded: %s", file.UploadFileOrglNm)
		}
	})

	downloader.Visit(fmt.Sprintf("https://www.open.go.kr/util/FileDownload.down?atchFileUploadNo=%s&fileSn=%s", file.AtchFileUploadNo, file.FileSn))
}
