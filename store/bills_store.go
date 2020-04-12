package store

import (
	"github.com/sluggishhackers/realopen.go/models"
)

func (store *Store) GetBill(ID string) *models.Bill {
	return store.bills[ID]
}

func (store *Store) GetBills() map[string]*models.Bill {
	return store.bills
}

func (store *Store) SaveBill(b models.Bill) {
	store.bills[b.ID] = &b
}
