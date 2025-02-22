package analytics

import (
	"runtime"
	"sync"

	"github.com/emmrys-jay/monieshop/transaction"
)

type Worker struct {
	wg *sync.WaitGroup
	mu *sync.Mutex
}

func NewWorker() *Worker {
	return &Worker{
		wg: &sync.WaitGroup{},
		mu: &sync.Mutex{},
	}
}

func (w *Worker) Analyze(transactions []transaction.Transaction) {
	numCPU := runtime.NumCPU()
	resultsChan := make(chan *Result)

	w.wg.Add(numCPU)
	partition := len(transactions) / numCPU
	for i := 1; i <= numCPU; i++ {
		if i == 1 {
			go w.start(transactions[:partition], resultsChan)
		} else if i < numCPU {
			go w.start(transactions[partition*(i-1):partition*i], resultsChan)
		} else {
			go w.start(transactions[partition*i:], resultsChan)
		}
	}

	day := transactions[0].TransactionTime.Format("2006-01-02")
	month := transactions[0].TransactionTime.Month().String()

	go func() {
		w.wg.Wait()
		close(resultsChan)
	}()

	tStore := w.populateTemporaryStore(month, resultsChan)

	// Populate global store value
	// Get from your temporary store
	w.mu.Lock()
	defer w.mu.Unlock()
	storeValue := store.Load().(*TemporaryStore)

	if tStore.HighestDaySales.Volume > storeValue.HighestDaySales.Volume {
		storeValue.HighestDaySales.Volume = tStore.HighestDaySales.Volume
		storeValue.HighestDaySales.Day = day
	}

	if tStore.HighestDaySalesValue.Value > storeValue.HighestDaySalesValue.Value {
		storeValue.HighestDaySalesValue.Value = tStore.HighestDaySalesValue.Value
		storeValue.HighestDaySalesValue.Day = day
	}

	for k, v := range tStore.MostSoldProductID {
		storeValue.MostSoldProductID[k] += v
	}

	salesStaffNumbers := storeValue.HighestSalesStaff[month]
	if salesStaffNumbers == nil {
		salesStaffNumbers = make(map[int]int)
	}
	for k, v := range tStore.HighestSalesStaff[month] {
		salesStaffNumbers[k] = v
	}
	storeValue.HighestSalesStaff[month] = salesStaffNumbers

	for k, v := range tStore.HighestHourOfTheDay {
		storeValue.HighestHourOfTheDay[k] += v
	}

	storeValue.NumberOfDays++
	store.Store(storeValue)
}

func (w *Worker) start(ts []transaction.Transaction, resultChan chan *Result) {
	defer w.wg.Done()

	totalSalesVolume := 0
	totalSalesValue := 0.0
	mostSoldProductID := make(map[string]int)
	salesStaffId := make(map[int]int)      // staff Id to numoftransactions
	volumeForEachHour := make(map[int]int) // hour to total transaction volume

	for _, transaction := range ts {

		volumes := 0
		for _, product := range transaction.Products {
			// Total Sales Volume in a day
			volumes += product.Quantity

			mostSoldProductID[product.Id] += product.Quantity

		}

		totalSalesVolume += volumes                       // find total sales volume
		totalSalesValue += transaction.SaleAmount         // Find total sales value
		salesStaffId[transaction.SalesStaffId] += volumes // Add one to each staff id you see

		hour := transaction.TransactionTime.Hour()
		volumeForEachHour[hour]++ // Add one to each hour for each transaction
	}

	resultChan <- &Result{
		SalesVolume:         totalSalesVolume,
		SalesValue:          totalSalesValue,
		SoldProductID:       mostSoldProductID,
		SalesStaffPerMonth:  salesStaffId,
		HighestHourOfTheDay: volumeForEachHour,
	}
}

// Only call after all the workers are done
func (w *Worker) GetResult() *Analytics {
	storeValue := store.Load().(*TemporaryStore)

	maxSoldVolume := 0
	productId := ""
	for k, v := range storeValue.MostSoldProductID {
		if v > maxSoldVolume {
			productId = k
			maxSoldVolume = v
		}
	}

	monthlyHighestStaffId := make([]StaffIDPerMonth, 0, 12)
	for k, v := range storeValue.HighestSalesStaff {
		maxStaffIdVolume, maxStaffId := 0, 0
		for k, v := range v {
			if v > maxStaffIdVolume {
				maxStaffId = k
				maxStaffIdVolume = v
			}
		}

		monthlyHighestStaffId = append(monthlyHighestStaffId, StaffIDPerMonth{Sales: maxStaffIdVolume, StaffId: maxStaffId, Month: k})
	}

	average := 0.0
	maxHour := 0
	for k, v := range storeValue.HighestHourOfTheDay {
		avg := float64(v) / float64(storeValue.NumberOfDays)
		if avg > average {
			average = avg
			maxHour = k
		}
	}

	return &Analytics{
		HighestDaySalesVolume: storeValue.HighestDaySales,
		HighestDaySalesValue:  storeValue.HighestDaySalesValue,
		MostSoldProductID:     productId,
		HighestSalesStaffID:   monthlyHighestStaffId,
		HighestHourOfTheDay:   maxHour,
	}
}

func (w *Worker) populateTemporaryStore(month string, resultsChan chan *Result) *TemporaryStore {
	var tstore = TemporaryStore{
		HighestDaySales:      new(DaySalesVolume),
		HighestDaySalesValue: new(DaySalesValue),
		MostSoldProductID:    make(map[string]int),
		HighestSalesStaff:    make(map[string]map[int]int),
		HighestHourOfTheDay:  make(map[int]int),
	}

	mapOfStaffIdToVolume := make(map[int]int)

	for result := range resultsChan {
		tstore.HighestDaySales.Volume += result.SalesVolume
		tstore.HighestDaySalesValue.Value += result.SalesValue

		for k, v := range result.SoldProductID {
			tstore.MostSoldProductID[k] += v
		}

		for k, v := range result.SalesStaffPerMonth {
			mapOfStaffIdToVolume[k] += v
		}

		for k, v := range result.HighestHourOfTheDay {
			tstore.HighestHourOfTheDay[k] += v
		}
	}

	tstore.HighestSalesStaff[month] = mapOfStaffIdToVolume

	return &tstore
}
