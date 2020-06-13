package models

type Bill struct {
	ID        string `json:"rqestNo"`   // "6474270"
	Result    string `json:"oppSeNm"`   // "비공개"
	Status    string `json:"prcsStsNm"` // "통지완료"
	Requester string `json:"bllrNm"`    // "투명사회를위한정보공개센터"

	ChargeYn        string `json:"chargeYn"` // "N"
	ChkrNmpn        string `json:"chkrNmpn"`
	ChkrClsfNm      string `json:"chkrClsfNm"`
	ChrgDeptCd      string `json:"chrgDeptCd'`     // "B550021"
	ChrgDeptFullNm  string `json:"chrgDeptFullNm"` // "금융감독원"
	ChrgDeptNm      string `json:"chrgDeptNm"`     // "금융감독원"
	CprGrpBizrno    string `json:"cprGrpBizrno"`   // "101-80-05483"
	DrftrNmpn       string `json:"drftrNmpn"`
	DrftrClsfNm     string `json:"drftrClsfNm"`
	EndRow          int    `json:"endRow"`          // 0
	FeeFpyYn        string `json:"feeFpyYn"`        // "Y"
	FeeSumAmt       string `json:"feeSumAmt"`       // 0
	IfrmpPrcsRstrNo string `json:"ifrmpPrcsRstrNo"` // "8564760"
	IsOpened        string `json:""`                // "Y"
	MberId          string `json:"mberId"`          // "opengirok"
	MberSeCd        string `json:"mberSeCd"`        // "401"
	MultiPrcsYn     string `json:"multiPrcsYn"`     // "N"
	NeedCertYn      string `json:"needCertYn"`      // "N"
	OpenPot         string `json:"openPot"`         // "20.03.03"
	OpetrNmpn       string `json:"opetrNmpn"`       // "김다혜"
	OppCn           string `json:"oppCn"`           // ""
	OppSeCd         string `json:"oppSeCd"`         // "3"
	OppSeNm         string `json:"oppSeNm"`
	OppStleSeNm     string `json:"oppStleSeNm"`
	PrcsDeptCd      string `json:"prcsDeptCd"`      // "B550021"
	PrcsDeptNm      string `json:"prcsDeptNm"`      // "금융감독원"
	PrcsDsnTmlmtYmd string `json:"prcsDsnTmlmtYmd"` // "2020.03.04"
	PrcsFullNstNm   string `json:"prcsFullNstNm"`   // "금융감독원"
	PrcsNstCd       string `json:"prcsNstCd"`       // "B550021"
	PrcsNstNm       string `json:"prcsNstNm"`       // "금융감독원"
	PrcsPot         string `json:"prcsPot"`         // "2020.03.03"
	PrcsRstrNo      string `json:"prcsRstrNo"`      // "8564760"
	PrcsStsCd       string `json:"prcsStsCd"`       // "143"
	RceptDate       string `json:"rceptDate"`       // "2020.02.20"
	RltvRqestNoCnt  string `json:"rltvRqestNoCnt"`  // "1"
	RowPage         int    `json:"rowPage"`         // 10
	RqestFullNstNm  string `json:"mberId"`          // "금융감독원"
	RqestNstCd      string `json:"rqestNstCd"`      // "B550021"
	RqestNstNm      string `json:"rqestNstNm"`      // "금융감독원"
	RqestPot        string `json:"rqestPot"`        // "2020.02.20"
	RqestSj         string `json:"rqestSj"`         // "금융 기관의 정보 보안 사고 내역 정보공개 청구"
	RqestInfoDtls   string `json:"rqestInfoDtls"`
	StartRow        int    `json:"startRow"`   // 0
	TotalCount      int    `json:"totalCount"` // 0
	TotalPage       int    `json:"totalPage"`  // 0
	ViewPage        int    `json:"viewPage"`   // 1
}

type File struct {
	Addr1              string `json:"addr1"`              // ""
	Addr2              string `json:"addr2"`              // ""
	AtchFilePrsrvNm    string `json:"atchFilePrsrvNm"`    // "202003021735501750000.zip"
	AtchFileUploadNo   string `json:"atchFileUploadNo"`   // "ZkExMS84ODgzZUVWRTdoMlg3aXd3dz09"
	Atch_fileByteNum   string `json:"atch_fileByteNum"`   // ""
	CsdCnvrId          string `json:"csdCnvrId"`          // ""
	CsdCnvrStsCd       string `json:"csdCnvrStsCd"`       // "020"
	CsdDocNm           string `json:"csdDocNm"`           // ""
	CsdFileCoursDtls   string `json:"csdFileCoursDtls"`   // ""
	Error_code         string `json:"error_code"`         // ""
	Error_msg          string `json:"error_msg"`          // ""
	FileAbsltCoursDtls string `json:"fileAbsltCoursDtls"` // "/pidfiles/uploads/pb/dlsrinfo/"
	FileSn             string `json:"fileSn"`             // "1"
	FrstRgstPot        string `json:"frstRgstPot"`        // "2020-03-02 17:35:50.0"
	FrstRgstrId        string `json:"frstRgstrId"`        // ""
	FrstRgstrNm        string `json:"frstRgstrNm"`        // ""
	LastUpdtPot        string `json:"lastUpdtPot"`        // "2020-03-02 17:35:50.0"
	LastUpdtrId        string `json:"lastUpdtrId"`        // ""
	LastUpdtrNm        string `json:"lastUpdtrNm"`        // ""
	RefrnFileUploadNo  string `json:"refrnFileUploadNo"`  // ""
	RowPage            int    `json:"rowPage"`            // 10
	TmstmpPrcsKeyVal   string `json:"tmstmpPrcsKeyVal"`   // ""
	TmstmpRsulCd       string `json:"tmstmpRsulCd"`       // ""
	TotalPage          int    `json:""`                   // 0
	UploadFileOrglNm   string `json:"uploadFileOrglNm"`   // "2015.zip"
	ViewPage           int    `json:""`                   // 1
	ZipCode            string `json:""`                   // ""
}
