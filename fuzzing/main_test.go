package main

import "testing"

func FuzzCalcTax(f *testing.F) {
	seed := []float64{0, 1, 999, 1000, 1001} // proximity to boundaries

	// we cannot add seed values directly to f, since we need to iterate over them, not the entire slice
	for _, amount := range seed {
		f.Add(amount)
	}

	f.Fuzz(func(t *testing.T, amount float64) {
		result := CalcTax(amount)
		if amount <= 0 && result != 0 {
			t.Errorf("CalcTax(%.2f) = %.2f; want %.2f", amount, result, 0.0)
		}
		if amount >= 20000 && result != 20.0 {
			t.Errorf("CalcTax(%.2f) = %.2f; want %.2f", amount, result, 20.0)
		}
	})
}
