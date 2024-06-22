package main

import "encoding/json"

func GetCoordinatesLenWithEncodingJSON(jsonStr string) int {
	var canada Canada
	err := json.Unmarshal([]byte(jsonStr), &canada)
	if err != nil {
		panic(err)
	}

	return len(canada.Features[0].Geometry.Coordinates)
}
