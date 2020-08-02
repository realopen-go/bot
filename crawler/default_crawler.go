package crawler

import (
	"log"
	"os"
	"time"

	b64 "encoding/base64"

	"github.com/gocolly/colly"
)

func newDefaultCrawler() *colly.Collector {
	decodedMberID := os.Getenv("REALOPEN_MEMBER_ID")
	decodedPWD := os.Getenv("REALOPEN_PASSWORD")
	encodedMberId := b64.StdEncoding.EncodeToString([]byte(decodedMberID))
	encodedPWD := b64.StdEncoding.EncodeToString([]byte(decodedPWD))
	c := colly.NewCollector()

	timeout, _ := time.ParseDuration("1m")
	c.SetRequestTimeout(timeout)
	err := c.Post("https://www.open.go.kr/pa/member/openLogin/ajaxSessionCheck.ajax", map[string]string{"redirectUrl": "/pa/member/openLogin/memberLogin.ajax"})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Post("https://www.open.go.kr/pa/member/openLogin/memberLogin.ajax", map[string]string{"mberId": encodedMberId, "pwd": encodedPWD, "agent": "PC"})
	if err != nil {
		log.Fatal(err)
	}

	c.Visit("https://www.open.go.kr")
	return c
}
