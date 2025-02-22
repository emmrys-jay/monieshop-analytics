package transaction

import "time"

type Transaction struct {
	SalesStaffId    int
	TransactionTime time.Time
	Products        []Product
	SaleAmount      float64
}

type Product struct {
	Id       string
	Quantity int
}
