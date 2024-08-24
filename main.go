package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var accessCount int = 0

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
		recs, err := GetContentFromCSV()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Count: %v\n", accessCount)
		fmt.Printf("length of recs: %v\n", len(recs))
		if accessCount == len(recs) {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("End of CSV reached"))

			fmt.Println("Hit end of CSV file")
			return
		}

		line := strings.Join(recs[accessCount], ",")
		accessCount++

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(line))
	})
}

func GetContentFromCSV() ([][]string, error) {
	f, err := os.Open("data/example.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	recs, err := r.ReadAll()
	recs = recs[1:]
	if err != nil {
		return nil, err
	}

	return recs, nil
}
