package main

import "os"

func GetJSONStr() string {
	jsonBytes, err := os.ReadFile("canada.json")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
