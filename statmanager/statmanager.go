package statmanager

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/sluggishhackers/realopen.go/models"
	"github.com/sluggishhackers/realopen.go/rmtstor"
	"github.com/sluggishhackers/realopen.go/utils/date"
	"gopkg.in/yaml.v3"
)

var statusFilePath string

type Istatmanager interface {
	Indexing(map[string]*models.Bill)
	Initialize()
	Load()
	Update()
	SetFileStatus(string, bool)
}

type BillStatus struct {
	BillID     string `yaml:"BILL_ID"`
	Status     string `yaml:"STATUS"`
	FileStatus string `yaml:"FILE_STATUS"`
	Link       string `yaml:"LINK"`
}

type IndexFileType struct {
	Updated  string        `yaml:"UPDATED"`
	Statuses []*BillStatus `yaml:"BILLS"`
}

type statmanager struct {
	statuses map[string]*BillStatus
}

var INDEXING_DIR string
var MEMBER_NAME string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	gitDir := rmtstor.REALOPEN_INDEX_DIR
	MEMBER_NAME = os.Getenv("REALOPEN_MEMBER_NAME")
	INDEXING_DIR = fmt.Sprintf("%s/%s/%s", wd, gitDir, MEMBER_NAME)
}

func initializeIndexFile() *IndexFileType {
	return &IndexFileType{
		Updated: date.Now().Format(date.DEFAULT_FORMAT),
	}
}

func (sm *statmanager) Initialize() {
	// clean
	cleanCmd := exec.Command("rm", "-rf", INDEXING_DIR)
	cleanCmd.Run()

	mkdirCmd := exec.Command("mkdir", INDEXING_DIR)
	mkdirCmd.Run()

	statusFilePath = fmt.Sprintf("%s/status.yml", INDEXING_DIR)
	touchCmd := exec.Command("touch", statusFilePath)
	touchCmd.Run()

}

func (sm *statmanager) Indexing(bills map[string]*models.Bill) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/status.yml", INDEXING_DIR))
	if err != nil {
		log.Fatal(err)
	}

	indexFile := initializeIndexFile()
	err = yaml.Unmarshal(data, indexFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, b := range bills {
		indexFile.Statuses = append(indexFile.Statuses, &BillStatus{
			BillID: b.ID,
			Status: ParseStatus(b.Status, b.Result),
		})
	}

	statusYml, err := yaml.Marshal(indexFile)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(statusFilePath, statusYml, 0644)
}

func (sm *statmanager) Load() {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/status.yml", INDEXING_DIR))
	if err != nil {
		log.Fatal(err)
	}

	indexFile := initializeIndexFile()
	err = yaml.Unmarshal(data, indexFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, b := range indexFile.Statuses {
		sm.statuses[b.BillID] = b
	}
}

func (sm *statmanager) Update() {
	indexFile := initializeIndexFile()

	for _, s := range sm.statuses {
		indexFile.Statuses = append(indexFile.Statuses, s)
	}

	changed, err := yaml.Marshal(indexFile)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(statusFilePath, changed, 0644)
}

func (sm *statmanager) SetFileStatus(billID string, downloaded bool) {
	if downloaded {
		sm.statuses[billID].FileStatus = "DOWNLOADED"
		sm.statuses[billID].Link = fmt.Sprintf("%s/%s", rmtstor.REALOPEN_DATA_REPOSITORY, billID)
	} else {
		sm.statuses[billID].FileStatus = "NO_FILE"
	}
}

func New() Istatmanager {
	return &statmanager{
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
