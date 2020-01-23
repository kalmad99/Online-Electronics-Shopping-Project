package service

import (
	//"../../../entity"
	//"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank"
)

// BankService implements menu.OrderService interface
type BankService struct {
	bankRepo bank.PayRepository
}

// NewOrderService returns new OrderService object
func NewBankService(bankRepository bank.PayRepository) bank.PayService {
	return &BankService{bankRepo: bankRepository}
}

func (bs *BankService) MakePayment(accno string) []error {
	panic("implement me")
}

// Orders returns all stored food orders
//func (os *OrderService) Orders() ([]entity.Order, []error) {
//	ords, errs := os.orderRepo.Orders()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return ords, errs
//}
//public boolean checkBankId(int bankNo){
//        try {
//            ResultSet rs=statement.executeQuery("select 1 from bank where bankAccount='"+bankNo+"'");
//            if(rs.first()==true){
//               return true;
//            }else{
//                return false;
//            }
//
//        } catch (Exception e) {
//        }
//        return false;
//    }
