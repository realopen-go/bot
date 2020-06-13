package store

import "github.com/sluggishhackers/go-realopen/models"

type IStore interface {
	ClearBills()
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
