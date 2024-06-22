package main

import "testing"

func BenchmarkGetCoordinatesLenWithFastJSON(b *testing.B) {
	jsonStr := GetJSONStr()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetCoordinatesLenWithFastJSON(jsonStr)
	}
}

func BenchmarkGetCoordinatesLenWithEncodingJSON(b *testing.B) {
	jsonStr := GetJSONStr()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetCoordinatesLenWithEncodingJSON(jsonStr)
	}
}
