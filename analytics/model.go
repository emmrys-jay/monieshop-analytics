package analytics

import "sync/atomic"

type Analytics struct {
	HighestDaySalesVolume *DaySalesVolume
	HighestDaySalesValue  *DaySalesValue
	MostSoldProductID     string
	HighestSalesStaffID   []StaffIDPerMonth
	HighestHourOfTheDay   int
}

type StaffIDPerMonth struct {
	Month   string
	Sales   int
	StaffId int
}

type TemporaryStore struct {
	HighestDaySales      *DaySalesVolume        // Day to Sales volume (keep only the highest volume in a day)
	HighestDaySalesValue *DaySalesValue         // Day to sales Value (keep only the highest value in a day)
	MostSoldProductID    map[string]int         // Product Id to sales volume
	HighestSalesStaff    map[string]map[int]int // Month to map of staff id to volume
	HighestHourOfTheDay  map[int]int            // Hour of the day to transaction volume
	NumberOfDays         int                    // Total Number of days
}

type DaySalesVolume struct {
	Day    string
	Volume int
}

type DaySalesValue struct {
	Day   string
	Value float64
}

type Result struct {
	SalesVolume         int            // Total Sales volume
	SalesValue          float64        // TotalSalesValue
	SoldProductID       map[string]int // Map of ProductId to Volume
	SalesStaffPerMonth  map[int]int    // Map of staffID to numbers of sales transactions
	HighestHourOfTheDay map[int]int    // Map of hour of the Total number of transactions in that hour
}

type SalesStaffNumbers struct {
	Id     int
	Volume int
}

type TransactionVolume struct {
	SumOfTransactionVolume int
	NumberOfDays           int
}

var (
	store atomic.Value
)

func init() {
	store.Store(&TemporaryStore{
		HighestDaySales:      new(DaySalesVolume),
		HighestDaySalesValue: new(DaySalesValue),
		MostSoldProductID:    make(map[string]int),
		HighestHourOfTheDay:  make(map[int]int),
		HighestSalesStaff:    map[string]map[int]int{},
	})
}
