package main

// 1. does reading the entire request body into memory scale for large uploads?
// 2. how should we surface errors when the CSV contains non-numeric cells?
// 3. if a csv contains invalid values, save the partial and return an error
// 4. implement a way to resume the calculation from where it was interrupted
// 5. can we add a benchmark or integration test that parses a sample 10 MB CSV?

// city,jan,feb,mar
// Berlin,1.2,2.3,5
// Paris,3.4,4.5,6.7
// Madrid,5.O,6.1,7.2
// Milan,2,3,4

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type summary struct {
	Column string  `json:"column"`
	Sum    float64 `json:"sum"`
}

func summarizeCSV(data []byte) ([]summary, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	sums := make([]float64, len(header))

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		for i, v := range rec {
			f, _ := strconv.ParseFloat(v, 64)
			sums[i] += f
		}
	}

	out := make([]summary, len(header))
	for i, col := range header {
		out[i] = summary{Column: col, Sum: sums[i]}
	}
	return out, nil
}

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	list, err := summarizeCSV(data)
	if err != nil {
		http.Error(w, "csv error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func main() {
	http.HandleFunc("/summarize", summarizeHandler)
	log.Println("CSV Summarizer listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
