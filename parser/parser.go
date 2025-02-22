package parser

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/emmrys-jay/monieshop/log"
	"github.com/emmrys-jay/monieshop/transaction"
)

type Parser struct {
	directory string
	logger    *log.Logger
}

func NewParser(dir string, l *log.Logger) *Parser {
	return &Parser{
		directory: dir,
		logger:    l,
	}
}

func (pa *Parser) ParseTransactions(filename string) (ts []transaction.Transaction, errors bool) {
	f, err := os.Open(filepath.Join(pa.directory, filename))
	if err != nil {
		pa.logger.Printf("Error opening file: %v\n", err)
		return nil, true
	}
	defer f.Close()

	return pa.parse(f)
}

func (pa *Parser) parse(file *os.File) (ts []transaction.Transaction, errors bool) {

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		pa.logger.Println("Error reading file: ", err)
		return nil, true
	}

	var transactions = make([]transaction.Transaction, 0)

	for idx, record := range records {
		staffIdStr := record[0]
		transactionTimeStr := record[1]
		productsStr := record[2]
		salesAmountStr := record[3]

		staffId, err := strconv.Atoi(staffIdStr)
		if err != nil {
			errors = true
			pa.logger.Printf("Invalid staff id '%v' found in line %v, file '%v'\n", staffIdStr, idx+1, file.Name())
			continue
		}

		if len(strings.Split(transactionTimeStr, ":")) <= 2 {
			transactionTimeStr += ":00"
		}
		transactionTime, err := time.Parse("2006-01-02T15:04:05", transactionTimeStr)
		if err != nil {
			errors = true
			pa.logger.Printf("Invalid transaction time '%v' found in line %v, file '%v'\n", transactionTimeStr, idx+1, file.Name())
			continue
		}

		p := strings.TrimSuffix(strings.TrimPrefix(productsStr, "["), "]")
		prods := strings.Split(p, "|")

		products := make([]transaction.Product, 0, len(prods))
		for _, v := range prods {
			pi := strings.Split(v, ":")

			r, _ := strconv.Atoi(pi[1])
			products = append(products, transaction.Product{
				Id:       pi[0],
				Quantity: r,
			})
		}

		salesAmount, err := strconv.ParseFloat(salesAmountStr, 64)
		if err != nil {
			errors = true
			pa.logger.Printf("Invalid sales amount format '%v' found on line %v, file '%v'\n", salesAmountStr, idx+1, file.Name())
			continue
		}

		t := transaction.Transaction{
			SalesStaffId:    staffId,
			TransactionTime: transactionTime,
			Products:        products,
			SaleAmount:      salesAmount,
		}

		transactions = append(transactions, t)
	}

	return transactions, errors
}
