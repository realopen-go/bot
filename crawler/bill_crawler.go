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

func makeFileName(filePath string, bill *BillResultFormat, file *models.File) string {
	return fmt.Sprintf("%s/%s_%s_%s_%s", filePath, bill.DtlVo.IfrmpPrcsRstrNo, bill.DtlVo.PrcsNstNm, strings.Trim(strings.ReplaceAll(bill.DtlVo.PrcsDeptNm, " ", "_"), " "), strings.Trim(strings.ReplaceAll(file.UploadFileOrglNm, " ", "_"), " "))
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
	processingCrawler := c.defaultCrawler.Clone()
	processingCrawler.OnResponse(func(r2 *colly.Response) {
		body := string(r2.Body)
		startIndex := strings.Index(body, "var result")
		endIndex := strings.Index(body, "//var naviInfo")
		result := body[startIndex:endIndex]
		data := strings.TrimRight(strings.TrimSpace(result[strings.Index(result, "{"):]), ";")

		billResultFormat := &ProcessingBillResultFormat{}

		err := json.Unmarshal([]byte(data), billResultFormat)
		if err != nil {
			fmt.Println("Error to Unmarshall Bill Result Format")
			log.Fatal(err)
		}

		// TODO: formatter
		fmt.Println("Processing Bill")
		// fmt.Printf("%+v", billResultFormat.DtlVo)
		billResultFormat.DtlVo.RqestPot = strings.ReplaceAll(billResultFormat.DtlVo.RqestPot, ".", "-")

		c.store.SaveBill(billResultFormat.DtlVo)
	})

	crawler.OnResponse(func(r *colly.Response) {
		billID := r.Ctx.Get("billId")
		ifrmpPrcsRstrNo := r.Ctx.Get("ifrmpPrcsRstrNo")
		prcsStsCd := r.Ctx.Get("prcsStsCd")

		// fmt.Printf("Fetched: %s\n", billID)

		body := string(r.Body)
		startIndex := strings.Index(body, "var result")
		endIndex := strings.Index(body, "//var naviInfo")

		// 아직 처리중 단계에서는 청구건 상세페이지가 존재하지 않음
		if startIndex == -1 || endIndex == -1 {
			fmt.Println("처리중...")
			processingCrawler.Post(PROCESSING_HOST, NewParamsPostBill(billID, ifrmpPrcsRstrNo, prcsStsCd))
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
						// c.statusmanager.SetFileStatus(billID, true)
						close(ch)
					}
				}
			} else {
				// c.statusmanager.SetFileStatus(billID, false)
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

		filePath := fmt.Sprintf("%s/%s/%s_%s", wd, rmtstor.REALOPEN_DATA_DIR, bill.DtlVo.RqestPot, strings.Trim(strings.ReplaceAll(bill.DtlVo.RqestSj, " ", "_"), " "))

		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Errorf("😡 Error to create a new directory")
		}

		fileName := makeFileName(filePath, bill, &file)
		err = r.Save(fileName)
		if err != nil {
			fmt.Errorf("😡 Error to save a file to download", err)
			ch <- fmt.Sprintf("Failed: %s", file.UploadFileOrglNm)
		} else {
			ch <- fmt.Sprintf("Succeeded: %s", file.UploadFileOrglNm)
		}
	})

	downloader.Visit(fmt.Sprintf("https://www.open.go.kr/util/FileDownload.down?atchFileUploadNo=%s&fileSn=%s", file.AtchFileUploadNo, file.FileSn))
}
