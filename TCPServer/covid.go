package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CovidData struct {
	Date           string `json:"date"`
	Positive_Test  int    `json:"positive_date"`
	Test_Performed int    `json:"test_performed"`
	Expired        int    `json:"expired"`
	Admitted       int    `json:"admitted"`
	Discharged     int    `json:"discharged"`
	Region         string `json:"region"`
}

func Find(table []CovidData, filter string) []CovidData {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]CovidData, 0)
	filter = strings.ToUpper(filter)
	for _, covid := range table {
		if covid.Date == filter ||
			covid.Region == filter ||
			strings.Contains(strings.ToUpper(covid.Date), filter) ||
			strings.Contains(strings.ToUpper(covid.Region), filter) {
			result = append(result, covid)
		}
	}
	return result
}

func Load() []CovidData {
	lines, err := ReadCsv("TimeSeries_KeyIndicators.csv")
	if err != nil {
		fmt.Println("Error Reading CSV:", err)
	}
	table := make([]CovidData, 0)
	for _, line := range lines {
		line2, _ := strconv.Atoi(line[2])
		line3, _ := strconv.Atoi(line[3])
		line6, _ := strconv.Atoi(line[6])
		line10, _ := strconv.Atoi(line[10])
		line5, _ := strconv.Atoi(line[5])
		data := CovidData{
			Date:           line[4],
			Positive_Test:  line2,
			Test_Performed: line3,
			Expired:        line6,
			Admitted:       line10,
			Discharged:     line5,
			Region:         line[9],
		}
		table = append(table, data)
	}
	return table

}

func ReadCsv(filename string) ([][]string, error) {
	csvfile, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer csvfile.Close()

	lines, err := csv.NewReader(csvfile).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
