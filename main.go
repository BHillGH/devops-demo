package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
)

var accessCount int = 0

type CSVLine struct {
	ID   string
	Name string
	Hash string
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/get", getNextLine())
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready to go"))
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		accessCount = 0
		w.Write([]byte("Reset count"))
	})

	log.Fatal(http.ListenAndServe(":88", mux))
}

func getNextLine() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dataLines, err := GetContentFromCSV()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Count: %v\n", accessCount)
		if accessCount == len(dataLines) {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("End of CSV reached"))

			fmt.Println("Hit end of CSV file")
			return
		}

		line := ConvertToCSVLineString(dataLines[accessCount])
		accessCount++

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(line))
	})
}

func ConvertToCSVLineString(line CSVLine) string {
	return fmt.Sprintf("%s,%s,%s", line.ID, line.Name, line.Hash)
}

func GetContentFromCSV() ([]CSVLine, error) {
	f, err := os.Open("data/example.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	recs, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip header row
	recs = recs[1:]

	var dataLine []CSVLine

	for _, record := range recs {
		if len(record) != 3 {
			return nil, fmt.Errorf("unexpected record format: %v", record)
		}
		dataLine = append(dataLine, CSVLine{
			ID:   record[0],
			Name: record[1],
			Hash: record[2],
		})
	}

	return dataLine, nil
}
