package main

import "testing"

// Benchmark functions
func BenchmarkProcessTotalWithIter(b *testing.B) {
	m := CustomMapType{
		"2":   500,
		"1":   1000,
		"3":   400,
		"6":   300,
		"0.1": 200,
	}

	b.ReportAllocs() // Include memory allocation statistics

	for i := 0; i < b.N; i++ {
		ProcessTotalWithIter(m)
	}
}

func BenchmarkProcessTotalWithSepDataStructure(b *testing.B) {
	m := CustomMapType{
		"2":   500,
		"1":   1000,
		"3":   400,
		"6":   300,
		"0.1": 200,
	}

	b.ReportAllocs() // Include memory allocation statistics

	for i := 0; i < b.N; i++ {
		ProcessTotalWithSepDataStructure(m)
	}
}
