package statusmanager

import (
	"fmt"
	"os"

	"github.com/sluggishhackers/go-realopen/rmtstor"
)

var statusFilePath string

type Istatusmanager interface {
	Initialize()
	GetBillsToUpdate() []*BillStatus
	SetFileStatus(string, bool)
}

type BillStatus struct {
	BillID          string `yaml:"BILL_ID"`
	Status          string `yaml:"STATUS"`
	FileStatus      string `yaml:"FILE_STATUS"`
	Link            string `yaml:"LINK"`
	IfrmpPrcsRstrNo string `yaml:"ifrmpPrcsRstrNo"`
	PrcsDeptNm      string `yaml:"prcsDeptNm"`
	RqestSj         string `yaml:"rqestSj"`
}

type statusmanager struct {
	statuses map[string]*BillStatus
}

var MEMBER_NAME string

func init() {
	MEMBER_NAME = os.Getenv("REALOPEN_MEMBER_NAME")
}

func (sm *statusmanager) Initialize() {
}

func (sm *statusmanager) GetBillsToUpdate() []*BillStatus {
	var bills []*BillStatus
	for _, s := range sm.statuses {
		if s.FileStatus == "" {
			bills = append(bills, s)
		}
	}
	return bills
}

func (sm *statusmanager) SetFileStatus(billID string, downloaded bool) {
	if downloaded {
		sm.statuses[billID].FileStatus = "DOWNLOADED"
		sm.statuses[billID].Link = fmt.Sprintf("%s/%s", rmtstor.REALOPEN_DATA_REPOSITORY, billID)
	} else {
		sm.statuses[billID].FileStatus = "NO_FILE"
	}
}

func New() Istatusmanager {
	return &statusmanager{
		statuses: make(map[string]*BillStatus),
	}
}

func ParseStatus(koStatus string, koResult string) string {
	switch koStatus {
	case "통지완료":
		switch koResult {
		case "공개":
			return "DECIDED_OPEN"
		case "부분공개":
			return "DECIDED_OPEN_PARTIAL"
		case "비공개":
			return "DECIDED_NO"
		default:
			return ""
		}
	default:
		return ""
	}
}
