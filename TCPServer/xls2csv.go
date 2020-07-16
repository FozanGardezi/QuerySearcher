package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	excel2Csv("covid_final_data.xlsx")
}

func excel2Csv(fileName string) {
	file, err := excelize.OpenFile("covid_final_data.xlsx")

	//if error while opening file occurs exit
	if err != nil {
		log.Fatalln("Couldn't open the xlsx file", err)
		return
	}

	// Get value from cell by given worksheet name and axis.
	cell, err := file.GetCellValue("TimeSeries_KeyIndicators", "B2")
	if err != nil {
		println(err.Error())
		return
	}
	println(cell)

	// Get all the rows in the sheet
	rows, err := file.GetRows("TimeSeries_KeyIndicators")
	for _, row := range rows {
		for _, colCell := range row {
			print(colCell, "\t")
		}
		println()
	}

	csvfile, errr := os.Create("TimeSeries_KeyIndicators.csv")

	if errr != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)

	for _, row := range rows {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()

	csvfile.Close()
}
