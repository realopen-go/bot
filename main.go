package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/sluggishhackers/realopen.go/cmd"
)

//
//func init() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//}

func main() {
	cmd.Execute()
	// app.syncFiles(true)
	// app.uploadIndex()

	// app.indexing()

	// credential := url.Values{"mberId": {"ab3Blbmdpcm9r"}, "pwd": {"b3Blbmdpcm9r"}}
	// // Basic HTTP GET request
	// resp, err := http.PostForm("https://www.open.go.kr/pa/member/openLogin/memberLogin.ajax", credential)

	// if err != nil {
	// 	log.Fatal("Error getting response. ", err)
	// }
	// defer resp.Body.Close()

	// // Read body from response
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal("Error reading response. ", err)
	// }

	// fmt.Printf("%s\n", body)
}
