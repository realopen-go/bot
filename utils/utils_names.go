package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sluggishhackers/go-realopen/models"
)

func MakeFileDir(bill *models.Bill) string {
	return fmt.Sprintf("%s_%s", bill.RqestPot, strings.TrimSpace(strings.ReplaceAll(bill.RqestSj, " ", "_")))
}

func MakeFileName(bill *models.Bill, file models.File) string {
	symbol := regexp.MustCompile("['/\"]")
	fileName := fmt.Sprintf("%s_%s_%s_%s", bill.IfrmpPrcsRstrNo, bill.PrcsNstNm, strings.TrimSpace(strings.ReplaceAll(bill.PrcsDeptNm, " ", "_")), strings.TrimSpace(strings.ReplaceAll(file.UploadFileOrglNm, " ", "_")))
	return symbol.ReplaceAllString(fileName, "_")
}
