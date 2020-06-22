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

type BillResultFormat struct {
	FileList []models.File `json:"dntcFileList"`
	DtlVo    models.Bill   `json:"dtlVo"`
	//DtlVo    struct {
	//ChkrClsfNm      string `json:"chkrClsfNm"`  // 처리기관 - 검토자 직위/직급
	//ChkrNmpn        string `json:"chkrNmpn"`    // 처리기관 - 검토자 이름
	//ChrgDeptCd      string `json:"chrgDeptCd"`  // 처리기관 - 코드
	//ChrgDeptNm      string `json:"chrgDeptNm"`  // 처리기관 - 처리과명
	//DrftrClsfNm     string `json:"drftrClsfNm"` // 처리기관 -기안자 직위/직급
	//DrftrNmpn       string `json:"drftrNmpn"`   // 처리기관 -기안자 이름
	//IfrmpPrcsRstrNo string `json:"ifrmpPrcsRstrNo"`
	//InfoOppPrcsCn   string `json:"infoOppPrcsCn"`
	//OppCn           string `json:"oppCn"`       // 공개내용
	//OppOprtnPot     string `json:"oppOprtnPot"` // 통지일자
	//OppSeNm         string `json:"oppSeNm"`
	//OppStleSeNm     string `json:"oppStleSeNm"`   // 공개방법 - 교부형태
	//PrcsDeptNm      string `json:"prcsDeptNm"`    // 처리기관 - 처리과명
	//PrcsNstNm       string `json:"prcsNstNm"`     // 처리기관명
	//PrcsStsCd       string `json:"prcsStsCd"`     //
	//RqestInfoDtls   string `json:"rqestInfoDtls"` // 청구내용
	//RqestNo         string `json:"rqestNo"`       // 접수번호
	//RqestPot        string `json:"rqestPot"`      // 접수일자
	//RqestSj         string `json:"rqestSj"`       // 접수명
	//SnctrlClsfNm    string `json:"snctrlClsfNm"`  // 처리기관 -기안자 직위/직급
	//SnctrlNmpn      string `json:"snctrlNmpn"`    // 처리기관 -결재권자 직위/직급
	//} `json:"dtlVo"`
}

func makeFileName(filePath string, bill *BillResultFormat, file *models.File) string {
	return fmt.Sprintf("%s/%s_%s", filePath, bill.DtlVo.IfrmpPrcsRstrNo, file.UploadFileOrglNm)
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
		billID := r.Ctx.Get("billId")
		fmt.Printf("Fetched: %s\n", billID)

		body := string(r.Body)
		startIndex := strings.Index(body, "var result")
		endIndex := strings.Index(body, "//var naviInfo")

		// TODO: 왜 찾을 수 없는 페이지가 뜨는거지?
		if startIndex == -1 || endIndex == -1 {
			return
		}

		result := body[startIndex:endIndex]
		data := strings.TrimRight(strings.TrimSpace(result[strings.Index(result, "{"):]), ";")

		billResultFormat := &BillResultFormat{}

		err := json.Unmarshal([]byte(data), billResultFormat)
		if err != nil {
			fmt.Println("Error to Unmarshall Bill Result Format")
			log.Fatal(err)
		}

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
	})

	return crawler
}

func (c *Crawler) FetchBill(billID string, ifrmpPrcsRstrNo string, prcsStsCd string) {
	fmt.Println(fmt.Sprintf("Start to fetch a bill: %s", billID))

	// the key of "url" into the context of the request
	c.billCrawler.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("billId", billID)
	})

	err := c.billCrawler.Post("https://www.open.go.kr/pa/billing/openBilling/openBillingDntcDtl.do", NewParamsPostBill(billID, ifrmpPrcsRstrNo, prcsStsCd))
	if err != nil {
		fmt.Println("😡 Error to fetch a bill - 1")
		log.Fatal(err)
	}
}

func (c *Crawler) DownloadFile(bill *BillResultFormat, file models.File, ch chan string) {
	downloader := c.defaultCrawler.Clone()

	downloader.OnResponse(func(r *colly.Response) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		filePath := fmt.Sprintf("%s/%s/%s", wd, rmtstor.REALOPEN_DATA_DIR, bill.DtlVo.RqestSj)

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
