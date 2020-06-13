package store

import (
	"github.com/sluggishhackers/go-realopen/models"
)

func (store *Store) ClearBills() {
	store.bills = make(map[string]*models.Bill)
}

func (store *Store) GetBill(ID string) *models.Bill {
	return store.bills[ID]
}

func (store *Store) GetBills() map[string]*models.Bill {
	return store.bills
}

func (store *Store) SaveBill(b models.Bill) {
	store.bills[b.ID] = &b
}
