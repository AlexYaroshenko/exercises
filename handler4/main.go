package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

// We want to build an API endpoint to analyze web server logs.
//
// 1.  **Implement the `analyzeHandler` function.**
//     - It should decode the incoming JSON request body into the `LogRequest` struct.
//     - Handle potential JSON decoding errors gracefully.
//
// 2.  **Process the logs to find the top 3 most visited pages.**
//     - Count the occurrences of each `path`.
//     - Determine the top 3 paths based on their counts.
//
// 3.  **Send the response.**
//     - The response should be a JSON object matching the `TopPagesResponse` struct.
//     - The `top_pages` slice should be sorted in descending order of `count`.
//     - Set the correct `Content-Type` header to `application/json`.
//
// Request Example:
// {
//   "logs": [
//     {"path": "/home", "status_code": 200},
//     {"path": "/products", "status_code": 200},
//     {"path": "/home", "status_code": 200},
//     {"path": "/about", "status_code": 200},
//     {"path": "/products", "status_code": 200},
//     {"path": "/home", "status_code": 500}
//   ]
// }
//
// Response Example:
// {
//   "top_pages": [
//     {"path": "/home", "count": 3},
//     {"path": "/products", "count": 2},
//     {"path": "/about", "count": 1}
//   ]
// }

type LogRequest struct {
	Logs []LogEntry `json:"logs"`
}

type LogEntry struct {
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
}

type TopPagesResponse struct {
	TopPages []PageCount `json:"top_pages"`
}

type PageCount struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// GetTop returns the top n pages by visit count
func (l *LogRequest) GetTop(n int) []PageCount {
	counts := make(map[string]int)
	for _, log := range l.Logs {
		counts[log.Path]++
	}

	pageCounts := make([]PageCount, 0, len(counts))
	for path, count := range counts {
		pageCounts = append(pageCounts, PageCount{Path: path, Count: count})
	}

	sort.Slice(pageCounts, func(i, j int) bool {
		return pageCounts[i].Count > pageCounts[j].Count
	})

	if len(pageCounts) > n {
		return pageCounts[:n]
	}
	return pageCounts
}

// TopN defines how many top pages to return
const TopN = 3

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := TopPagesResponse{req.GetTop(TopN)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/analyze/top-pages", analyzeHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
