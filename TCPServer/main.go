package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
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

type queryStringMap struct {
	Map map[string]string `json:"query"`
}

type encoded struct {
	Response []CovidData `json:"response"`
}

func main() {

	var addr string
	var network string

	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	lis, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}
	defer lis.Close()

	log.Println("**** Global Currency Service ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to closse listener:", err)
			}
			continue
		}
		log.Println("connected to", conn.RemoteAddr())

		go returnData(conn)

	}
}

func returnData(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("error closing the connection", err)
		}
	}()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	table := Load()

	for {
		var reqAsMap queryStringMap
		if err := dec.Decode(&reqAsMap); err != nil {
			fmt.Println("failed to unmarshal: ", err)
		}

		for k, v := range reqAsMap.Map {
			if k == "region" {
				result := Find(table, v)

				var encode encoded
				encode.Response = result
				if err := enc.Encode(&encode); err != nil {
					fmt.Println("failed to encode data:", err)
					return
				}
			} else if k == "date" {
				result := Find(table, v)

				var encode encoded
				encode.Response = result
				if err := enc.Encode(&encode); err != nil {
					fmt.Println("failed to encode", err)
					return
				}
			} else {
				fmt.Println("Enter correct Json")
			}
		}

	}

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
