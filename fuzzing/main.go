package main

func CalcTax(in float64) (out float64) {
	if in <= 0 {
		return 0
	}
	if in >= 1000 && in < 20000 {
		return 10.0
	}
	if in >= 20000 {
		return 20.0
	}
	return 5.0
}
