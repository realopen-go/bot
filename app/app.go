package app

import (
	"fmt"

	"github.com/sluggishhackers/go-realopen/models"

	"github.com/sluggishhackers/go-realopen/crawler"
	"github.com/sluggishhackers/go-realopen/rmtstor"
	"github.com/sluggishhackers/go-realopen/statusmanager"
	"github.com/sluggishhackers/go-realopen/store"
	"github.com/sluggishhackers/go-realopen/utils/date"
)

type IApp interface {
	Initialize()
	Install()
	RunDailyCrawler()
}

type App struct {
	crawler       crawler.ICrawler
	remoteStorage rmtstor.IRemoteStorage
	store         store.IStore
	statusmanager statusmanager.Istatusmanager
}

// 정보공개플랫폼 최초 날짜
var initialDateFrom = "2003-01-01"
var initialDateTo = "2003-01-01"

func (app *App) Initialize() {
	app.remoteStorage.Initialize()
	app.statusmanager.Initialize()
}

func (app *App) Install() {
	// 1. Fetch all bills before
	targetDate := date.Now().AddDate(0, 0, -1).Format(date.DEFAULT_FORMAT)
	app.crawler.FetchBills(initialDateFrom, targetDate)

	// 2. Fetch Bills Details
	for _, b := range app.store.GetBills() {
		app.crawler.FetchBill(b.ID, b.IfrmpPrcsRstrNo, b.PrcsStsCd)
	}

	// 8. Create Rows on Database
	app.remoteStorage.CreateBills(app.store.GetBills())

	// 9. Clear Bills in Store
	app.store.ClearBills()

	app.remoteStorage.UploadFiles(true)
}

func (app *App) RunDailyCrawler() {
	fmt.Println("Run Daily Crawler")

	// 1. Fetch Old Bills Not Opened
	oldBills := app.remoteStorage.FetchBillsNotOpened()

	// 2. Crawl Bills only for decided
	for _, b := range oldBills {
		app.crawler.FetchBill(b.BillID, b.ProcessorRstrNumber, b.ProcessorStsCd)
	}

	// 3. Update Bills' Status
	var updatedBills []*models.Bill
	for _, b := range app.store.GetBills() {
		updatedBills = append(updatedBills, b)
	}

	// 4. Update Rows on Database
	app.remoteStorage.UpdateBills(updatedBills)

	// 5. Clear Bills in Store
	app.store.ClearBills()

	// 6. Fetch new bills list
	targetDate := date.Now().AddDate(0, 0, -1).Format(date.DEFAULT_FORMAT)
	app.crawler.FetchBills(targetDate, targetDate)

	// 7. Fetch Bills Details
	for _, b := range app.store.GetBills() {
		app.crawler.FetchBill(b.ID, b.IfrmpPrcsRstrNo, b.PrcsStsCd)
	}

	// 8. Create Rows on Database
	app.remoteStorage.CreateBills(app.store.GetBills())

	// 9. Clear Bills in Store
	app.store.ClearBills()
}

func New(c crawler.ICrawler, rs rmtstor.IRemoteStorage, s store.IStore, sm statusmanager.Istatusmanager) IApp {
	return &App{
		crawler:       c,
		remoteStorage: rs,
		store:         s,
		statusmanager: sm,
	}
}
