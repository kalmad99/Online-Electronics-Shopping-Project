package entity

// ProductMock mocks product lists
var ProductMock = Product{
	ID:         1,
	Name:       "Mock Item 01",
	CategoryID: 2,
	//Category: []Category{},
	Quantity: 10,
	Price:    50.5,
	//Seller: "Mock Seller 1",
	Description: "Mock Item Description",
	Image:       "mock_item.jpg",
	Rating:      4,
	RatersCount: 2,
}
