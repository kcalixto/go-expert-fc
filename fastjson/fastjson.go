package main

import "github.com/valyala/fastjson"

func GetCoordinatesLenWithFastJSON(jsonStr string) int {
	var p fastjson.Parser
	v, err := p.Parse(jsonStr)
	if err != nil {
		panic(err)
	}

	features := v.GetArray("features")
	if features == nil {
		panic("features is nil")
	}
	geometry := features[0].Get("geometry")
	if geometry == nil {
		panic("geometry is nil")
	}
	coordinates := geometry.GetArray("coordinates")
	if coordinates == nil {
		panic("coordinates is nil")
	}

	return len(coordinates)
}
