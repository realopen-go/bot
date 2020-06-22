package store

import "github.com/sluggishhackers/go-realopen/models"

type IStore interface {
	ClearBills()
	GetBill(ID string) *models.Bill
	GetBills() map[string]*models.Bill
	GetFiles() map[string][]models.File
	SaveBill(b models.Bill)
	SaveFiles(billID string, files []models.File)
}

type Store struct {
	bills map[string]*models.Bill
	files map[string][]models.File
}

func New() IStore {
	return &Store{
		bills: make(map[string]*models.Bill),
		files: make(map[string][]models.File),
	}
}
