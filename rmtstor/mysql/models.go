package mysql

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

type Bill struct {
	MultiID string `gorm:"type:varchar(255);not null`

	BillID                    string `gorm:"type:varchar(32);primary_key;unique;not null"` // 접수번호 - RqestNo
	Content                   string `gorm:"type:text"`                                    // 공개내용 - OppCn
	MultiPrcsYn               bool   `gorm:"type:boolean"`
	OpenType                  string `gorm:"type:varchar(100)"`                      // 공개방법 - OppStleSeNm
	OpenStatus                string `gorm:"type:varchar(10);index:open_status"`     // 공개여부 - OpenSeNm
	PublicDate                string `gorm:"index:public_date"`                      // 리얼오픈 공개일자
	ProcessorCode             string `gorm:"type:varchar(100);index:processor_code"` // 처리기관/코드 - ChrgDeptCd
	ProcessorDepartmentName   string `gorm:"type:varchar(100)"`                      // 처리기관/처리과명 - PrcsDeptNm
	ProcessorDrafterName      string `gorm:"type:varchar(100)"`                      // 처리기관/기안자명 - DrftrNmpn
	ProcessorDrafterPosition  string `gorm:"type:varchar(100)"`                      // 처리기관/기안자직위직급 - DrftrClsfNm
	ProcessorName             string `gorm:"type:varchar(100)"`                      // 처리기관명 - PrcsNstNm
	ProcessorRstrNumber       string `gorm:"type:varchar(100)"`                      // 처리기관명 - PrcsRstrNumber
	ProcessorReviewerName     string `gorm:"type:varchar(100)"`                      // 처리기관/검토자이름 - ChkrNmpn
	ProcessorReviewerPosition string `gorm:"type:varchar(100)"`                      // 처리기관/검토자 직위직급 - ChkrClsfNm
	ProcessorStsCd            string `gorm:"type:varchar(100)"`                      // - PrcsStsCd
	RequestContent            string `gorm:"type:text"`                              // 청구내용 - RqestInfoDtls
	RequestDate               string `gorm:"index:request_date"`                     // 접수일자 - RqestPot
	BillTitle                 string `gorm:"type:varchar(100)"`                      // 접수명 - RqestSj
	UserID                    string
	User                      User
}

func (b Bill) CreateMultiID() string {
	if !b.MultiPrcsYn {
		return b.BillID
	}

	sha := sha256.New()
	sha.Write([]byte(strings.TrimSpace(b.RequestContent)))
	hash := base64.URLEncoding.EncodeToString(sha.Sum(nil))
	return fmt.Sprintf("%s_%s", b.UserID, hash)
}

type File struct {
	gorm.Model
	FileName string `gorm:"type:varchar(255);not null"`
	BillID   string
	Bill     Bill
}

type User struct {
	ID          string        `gorm:"type:varchar(128);primary_key;unique;not null"`
	EmbagoMonth sql.NullInt64 `gorm:"type:(255);default:NULL`
	Password    string        `gorm:"type:varchar(255);`
	Username    string        `gorm:"type:varchar(128);unique;not null"`
}

func (Bill) TableName() string {
	return "bills"
}

func (File) TableName() string {
	return "files"
}

func (User) TableName() string {
	return "users"
}
