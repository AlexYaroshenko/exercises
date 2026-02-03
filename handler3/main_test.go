package main

import (
	"log"
	"strings"
	"testing"
)

func generateCSV(rows int) string {
	var b strings.Builder
	b.WriteString("city,jan,feb,mar\n")
	for i := 0; i < rows; i++ {
		b.WriteString("X,1.1,2.2,3.3\n")
	}
	return b.String()
}

func BenchmarkSummarize_ReadAll_10MB(b *testing.B) {
	csvData := generateCSV(1_000_000) // Approx 10MB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := summarizeCSV(strings.NewReader(csvData))
		if err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}
}
