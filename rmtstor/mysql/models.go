package mysql

type Bill struct {
	BillID                    string `gorm:"type:varchar(32);primary_key;unique;not null"` // 접수번호 - RqestNo
	Content                   string `gorm:"type:text"`                                    // 공개내용 - OppCn
	OpenType                  string `gorm:"type:varchar(100)"`                            // 공개방법 - OppStleSeNm
	OpenStatus                string `gorm:"type:varchar(10);index:open_status"`           // 공개여부 - OpenSeNm
	ProcessorCode             string `gorm:"type:varchar(100);index:processor_code"`       // 처리기관/코드 - ChrgDeptCd
	ProcessorDepartmentName   string `gorm:"type:varchar(100)"`                            // 처리기관/처리과명 - PrcsDeptNm
	ProcessorDrafterName      string `gorm:"type:varchar(100)"`                            // 처리기관/기안자명 - DrftrNmpn
	ProcessorDrafterPosition  string `gorm:"type:varchar(100)"`                            // 처리기관/기안자직위직급 - DrftrClsfNm
	ProcessorName             string `gorm:"type:varchar(100)"`                            // 처리기관명 - PrcsNstNm
	ProcessorRstrNumber       string `gorm:"type:varchar(100)"`                            // 처리기관명 - PrcsRstrNumber
	ProcessorReviewerName     string `gorm:"type:varchar(100)"`                            // 처리기관/검토자이름 - ChkrNmpn
	ProcessorReviewerPosition string `gorm:"type:varchar(100)"`                            // 처리기관/검토자 직위직급 - ChkrClsfNm
	ProcessorStsCd            string `gorm:"type:varchar(100)"`                            // - PrcsStsCd
	RequestContent            string `gorm:"type:text"`                                    // 청구내용 - RqestInfoDtls
	RequestDate               string `gorm:"index:request_date"`                           // 접수일자 - RqestPot
	BillTitle                 string `gorm:"type:varchar(100)"`                            // 접수명 - RqestSj
	UserID                    string
	User                      User
}

type User struct {
	ID       string `gorm:"type:varchar(128);primary_key;unique;not null"`
	Username string `gorm:"type:varchar(128);unique;not null"`
	Password string `gorm:"type:varchar(255);`
}

func (Bill) TableName() string {
	return "bills"
}

func (User) TableName() string {
	return "users"
}
