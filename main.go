package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/emmrys-jay/monieshop/analytics"
	logger "github.com/emmrys-jay/monieshop/log"
	"github.com/emmrys-jay/monieshop/parser"
)

var directory string
var help bool

func init() {
	flag.StringVar(&directory, "dir", "", "Folder where txt files are stored")
	flag.BoolVar(&help, "help", false, "Display help information")
	flag.Parse()

	if help {
		fmt.Println("Usage of monieshop analytics software:")
		fmt.Println("  -dir string")
		fmt.Println("        directory to find txt files")
		fmt.Println("  -help")
		fmt.Println("        Display help information")
		os.Exit(0)
	}
}

func main() {
	worker := analytics.NewWorker()
	logger := logger.NewLogger("errors.log")
	parser := parser.NewParser(directory, logger)

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !isValidFile(file.Name()) {
			continue
		}

		transactions, errors := parser.ParseTransactions(file.Name())
		if errors {
			log.Println("There were errors while parsing csv files. Check 'errors.log' for details")
		}

		if len(transactions) > 0 {
			worker.Analyze(transactions)
		}
	}

	rs := worker.GetResult()
	PrintResult(rs)
}

func PrintResult(result *analytics.Analytics) {
	fmt.Println("Highest Sales Volume In a Day: ")
	fmt.Printf("Day: %v, Sales Volume: %v\n\n", result.HighestDaySalesVolume.Day, result.HighestDaySalesVolume.Volume)

	fmt.Println("Highest Sales Value In a Day: ")
	fmt.Printf("Day: %v, Sales Value: %v\n\n", result.HighestDaySalesValue.Day, result.HighestDaySalesValue.Value)

	fmt.Println("Most Sold Product ID By Volume: ")
	fmt.Printf("Product ID: %v\n\n", result.MostSoldProductID)

	length := 9
	fmt.Println("Highest Sales Staff ID for each month: ")
	fmt.Println()
	sort.Slice(result.HighestSalesStaffID, func(i, j int) bool {
		monthI, _ := time.Parse("January", result.HighestSalesStaffID[i].Month)
		monthJ, _ := time.Parse("January", result.HighestSalesStaffID[j].Month)
		return monthI.Month() < monthJ.Month()
	})

	for _, v := range result.HighestSalesStaffID {
		fmt.Printf("Month: %v StaffID: %v Volume: %v\n", v.Month+strings.Repeat(" ", length-len(v.Month)), v.StaffId, v.Sales)
	}
	fmt.Println()

	fmt.Println("Highest Hour of The Day By Average Transaction Volume: ")
	fmt.Printf("Hour: %v\n", result.HighestHourOfTheDay)

}
