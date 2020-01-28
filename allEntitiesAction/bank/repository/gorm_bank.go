package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"log"
)

// BankGormRepo implements the bank.BankRepository interface
type BankGormRepo struct {
	conn *gorm.DB
}

// NewBankGormRepo returns new object of BankGormRepo
func NewBankGormRepo(db *gorm.DB) bank.PayRepository {
	return &BankGormRepo{conn: db}
}

// BankExists check if a given account is found
func (bankRepo *BankGormRepo) BankExists(acc string) bool {
	bank := entity.Bank{}
	errs := bankRepo.conn.Find(&bank, "account_no=?", acc).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// MakePayment accepts account number and amount to make payment
func (bankRepo *BankGormRepo) MakePayment(accno string, amount float64) (*entity.Bank, []error) {
	bank1 := []entity.Bank{}
	bal := entity.Bank{}
	row := bankRepo.conn.Select("balance").First(&bank1).Where("account_no = ?", accno).Scan(&bal)
	log.Println("Balance", bal.Balance)
	if row.RecordNotFound() {
		panic(row.Error)
	}
	if bal.Balance < amount{
		return nil, []error{errors.New("not Sufficient Balance")}
	} else{
		row = bankRepo.conn.Exec("UPDATE banks SET balance=? where account_no=?", bal.Balance - amount, accno)
		if row.RowsAffected < 1 {
			return nil, []error{errors.New("error")}
		}
		return &bal, nil
	}
}


