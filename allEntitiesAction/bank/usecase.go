package bank

//OrderService specifies customer menu order related services
type PayService interface {
	MakePayment(accno string) []error
}
