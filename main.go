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

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getNextLine() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recs := getContentFromCSV()
		if accessCount < len(recs) {
			line := strings.Join(recs[accessCount], ",")
			accessCount++

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(line))

			fmt.Println(line)
		} else {

			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("End of CSV reached"))

			fmt.Println("Hit end of CSV file")

		}
	})
}

func getContentFromCSV() [][]string {
	f, err := os.Open("tmp/example.csv")
	if err != nil {
		panic("CSV File not found")
	}
	defer f.Close()

	r := csv.NewReader(f)

	recs, err := r.ReadAll()
	recs = recs[1:]
	if err != nil {
		panic("Failed")
	}
	fmt.Println(recs[0])
	return recs
}
