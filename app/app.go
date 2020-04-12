package app

import (
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/sluggishhackers/realopen.go/crawler"
	"github.com/sluggishhackers/realopen.go/rmtstor"
	"github.com/sluggishhackers/realopen.go/statmanager"
	"github.com/sluggishhackers/realopen.go/store"
	"github.com/sluggishhackers/realopen.go/utils/date"
)

type IApp interface {
	crawl()
	DownloadFiles()
	Initialize()
	RunDailyCrawler()
	SyncIndexing(dateFrom string, dateTo string)
}

type App struct {
	crawler       crawler.ICrawler
	remoteStorage rmtstor.IRemoteStorage
	store         store.IStore
	statmanager   statmanager.Istatmanager
}

// 정보공개플랫폼 최초 날짜
var initialDateFrom = "2003-01-01"
var initialDateTo = "2003-01-01"

// 1. 어제자 청구 목록 인덱싱
// 2. 지금껏 인덱싱한 청구 목록 중 아직 공개되지 않은 청구 조회
func (app *App) crawl() {
	dateFrom := time.Now().AddDate(0, 0, -2).Format(date.DEFAULT_FORMAT)
	dateTo := time.Now().AddDate(0, 0, -1).Format(date.DEFAULT_FORMAT)

	// 1. Fetch all bills & Indexing
	app.SyncIndexing(dateFrom, dateTo)

	// 2. Download Files
	app.DownloadFiles()

	// 3. Push Git History
	app.remoteStorage.UploadIndex(false)
	app.remoteStorage.UploadFiles(false)
}

func (app *App) DownloadFiles() {
	app.statmanager.Load()

	bills := app.store.GetBills()
	for _, b := range bills {
		app.crawler.FetchBill(b)
	}

	app.statmanager.Update()
}

func (app *App) Initialize() {
	app.remoteStorage.Initialize()
	app.statmanager.Initialize()

	dateTo := date.Now().Format(date.DEFAULT_FORMAT)

	// 1. Fetch all bills & Indexing
	app.SyncIndexing(initialDateFrom, dateTo)

	// 2. Download Files
	app.DownloadFiles()

	// 3. Push Git History
	app.remoteStorage.UploadIndex(true)
	app.remoteStorage.UploadFiles(true)
}

func (app *App) RunDailyCrawler() {
	gocron.Every(1).Day().Do(app.crawl)
}

func (app *App) SyncIndexing(dateFrom string, dateTo string) {
	// 💡 정보공개청구플랫폼에서 지원하는 검색시작일자: "2003-01-01"
	// app.crawler.FetchBills("2003-01-01", date.Now().Format(date.DEFAULT_FORMAT))
	app.crawler.FetchBills(dateFrom, dateTo)
	app.statmanager.Indexing(app.store.GetBills())
}

func New(c crawler.ICrawler, rs rmtstor.IRemoteStorage, s store.IStore, sm statmanager.Istatmanager) IApp {
	return &App{
		crawler:       c,
		remoteStorage: rs,
		store:         s,
		statmanager:   sm,
	}
}
