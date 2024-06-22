package main

import (
	"fmt"
)

func main() {
	jsonStr := GetJSONStr()
	fmt.Println(GetCoordinatesLenWithEncodingJSON(jsonStr))
	fmt.Println(GetCoordinatesLenWithFastJSON(jsonStr))
}
