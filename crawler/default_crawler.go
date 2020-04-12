package crawler

import (
	"log"
	"os"

	"github.com/gocolly/colly"
)

func newDefaultCrawler() *colly.Collector {
	mberID := os.Getenv("REALOPEN_MEMBER_ID")
	PWD := os.Getenv("REALOPEN_PASSWORD")
	c := colly.NewCollector()

	err := c.Post("https://www.open.go.kr/pa/member/openLogin/ajaxSessionCheck.ajax", map[string]string{"redirectUrl": "/pa/member/openLogin/memberLogin.ajax"})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Post("https://www.open.go.kr/pa/member/openLogin/memberLogin.ajax", map[string]string{"mberId": mberID, "pwd": PWD, "agent": "PC"})
	if err != nil {
		log.Fatal(err)
	}

	c.Visit("https://www.open.go.kr")
	return c
}
