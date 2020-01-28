package repository

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"

	"github.com/jinzhu/gorm"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank"
)

// BankGormRepo implements the bank.BankRepository interface
type BankGormRepo struct {
	conn *gorm.DB
}

// NewBankGormRepo returns new object of BankGormRepo
func NewBankGormRepo(db *gorm.DB) bank.PayRepository {
	return &BankGormRepo{conn: db}
}

// Orders returns all customer orders stored in the database
func (bankRepo *BankGormRepo) MakePayment(accno string) []error {
	bank := []entity.Bank{}
	errs := bankRepo.conn.Find(&bank).GetErrors()
	if len(errs) > 0 {
		return errs
	}
	return nil
}
