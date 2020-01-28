package service

import (
	//"../../../entity"
	//"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

// BankService implements bank.BankService interface
type BankService struct {
	bankRepo bank.PayRepository
}

// NewBankService returns new BankService object
func NewBankService(bankRepository bank.PayRepository) bank.PayService {
	return &BankService{bankRepo: bankRepository}
}

func (bs *BankService) BankExists(acc string) bool {
	exists := bs.bankRepo.BankExists(acc)
	return exists
}

func (bs *BankService) MakePayment(accno string, amount float64) (*entity.Bank, []error) {
	ban, errs := bs.bankRepo.MakePayment(accno, amount)
	if len(errs) > 0 {
		return ban, errs
	}
	return ban, nil
}
