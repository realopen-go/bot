package store

import "github.com/sluggishhackers/realopen.go/models"

type IStore interface {
	GetBill(ID string) *models.Bill
	GetBills() map[string]*models.Bill
	SaveBill(b models.Bill)
}

type Store struct {
	bills map[string]*models.Bill
}

func New() IStore {
	return &Store{
		bills: make(map[string]*models.Bill),
	}
}
