package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/*
POST /sum

	{
	  "numbers": [1, 2, 3, 4]
	}

	{
	  "sum": 10
	}
*/

func main() {
	http.HandleFunc("/sum", sumHandler)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatalf("failed to start server : %w", err)
	}
}

type SumRequest struct {
	Numbers []int `json:"numbers"`
}

type SumResponse struct {
	Sum int `json:"sum"`
}

func sumHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var s SumRequest
	if err := json.Unmarshal(body, &s); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(s.Numbers) == 0 {
		http.Error(w, "numbers must not be empty", http.StatusBadRequest)
		return
	}

	r, err := json.Marshal(SumResponse{Sum: sum(s.Numbers)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(r)
}

func sum(s []int) int {
	var total int
	for _, n := range s {
		total += n
	}
	return total
}
