package bank

// OrderRepository specifies customer menu order related database operations
type PayRepository interface {
	MakePayment(accno string) []error
}
