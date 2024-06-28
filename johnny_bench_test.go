package johnny

import (
	"testing"

	"github.com/profe-ajedrez/gyro"
)

func BenchmarkDiscount(b *testing.B) {
	qty := udfs("10")

	benchCases := []struct {
		name   string
		entry  gyro.Gyro
		ratios []gyro.Gyro
	}{
		{
			name:   "discountable 100.123 perc discount 10 and 15",
			entry:  udfs("100.123"),
			ratios: func() []gyro.Gyro { g := make([]gyro.Gyro, 2); g[0] = udfs("15"); g[1] = udfs("10"); return g }(),
		},
		{
			name:  "discountable 172780372728901.12323223 perc discount 10.343244, 15 and 34.5654664",
			entry: udfs("172780372728901.12323223"),
			ratios: func() []gyro.Gyro {
				g := make([]gyro.Gyro, 3)
				g[0] = udfs("10.343244")
				g[1] = udfs("10")
				g[2] = udfs("34.5654664")
				return g
			}(),
		},
	}

	var bg FromUnitValue

	for _, tc := range benchCases {
		discountOverEntryValue := NewPercentualDiscount(tc.ratios[0])
		discountConsideringQty := NewPercentualDiscount(tc.ratios[1])

		b.ResetTimer()

		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i <= b.N; i++ {
				bg = NewFromUnitValue(tc.entry)
				bg.Visit(discountOverEntryValue)
				bg.Visit(WithQTY(qty))
				bg.Visit(discountConsideringQty)
			}
		})
	}
}

func BenchmarkCreateTaxe(b *testing.B) {
	benchCases := []struct {
		name  string
		entry gyro.Gyro
		ratio gyro.Gyro
	}{
		{
			name:  "taxable 100.123 tax 10%",
			entry: udfs("100.123"),
			ratio: udfs("10"),
		},
		{
			name:  "taxable 10.123 tax 1%",
			entry: udfs("10.123"),
			ratio: udfs("1"),
		},
		{
			name:  "taxable 100000.123001 tax 7.566%",
			entry: udfs("100000.123001"),
			ratio: udfs("7.566"),
		},
		{
			name:  "taxable 111111111111111100000.123001 tax 7.11111111566%",
			entry: udfs("111111111111111100000.123001"),
			ratio: udfs("7.11111111566"),
		},
	}

	b.ResetTimer()
	for _, tc := range benchCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i <= b.N; i++ {
				bg := NewFromUnitValue(tc.entry)
				taxOverEntryValue := NewPercTax(tc.ratio)
				bg.Receive(taxOverEntryValue)
			}
		})
	}
}

func BenchmarkTax(b *testing.B) {
	entry := udfs("100.123")
	ratio := udfs("10")
	ratio2 := udfs("15")
	qty := udfs("10")

	taxOverEntryValue := NewPercTax(ratio)
	taxConsideringQty := NewPercTax(ratio2)

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		bg := NewFromUnitValue(entry)

		bg.Receive(taxOverEntryValue)
		bg.Receive(WithQTY(qty))
		bg.Receive(taxConsideringQty)
	}
}

// BenchmarkPercentualUndiscount: Benchmarks the Visit method of the PercentualUndiscount struct with positive, negative, zero, and fractional ratios.
func BenchmarkPercentualUndiscount(b *testing.B) {
	testCases := []struct {
		name    string
		initial gyro.Gyro
		ratio   gyro.Gyro
	}{
		{
			name:    "Positive ratio",
			initial: udfs("100"),
			ratio:   udfs("10"),
		},
		{
			name:    "Zero ratio",
			initial: udfs("100"),
			ratio:   udfs("0"),
		},
		{
			name:    "Fractional ratio",
			initial: udfs("100"),
			ratio:   udfs("5.5"),
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := NewPercentualUnDiscount(tc.ratio)
				bg := &DefaultJohnny{v: tc.initial}
				u.Visit(bg)
			}
		})
	}
}

// BenchmarkAmountUndiscount: Benchmarks the Visit method of the AmountUndiscount struct with positive, negative, zero, and fractional amounts.
func BenchmarkAmountUndiscount(b *testing.B) {
	testCases := []struct {
		name    string
		initial gyro.Gyro
		amount  gyro.Gyro
	}{
		{
			name:    "Positive amount",
			initial: udfs("100"),
			amount:  udfs("10"),
		},
		{
			name:    "Zero amount",
			initial: udfs("100"),
			amount:  udfs("0"),
		},
		{
			name:    "Fractional amount",
			initial: udfs("100"),
			amount:  udfs("5.5"),
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			bg := &DefaultJohnny{v: tc.initial}
			u := NewAmountUnDiscount(tc.amount)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				u.Visit(bg)
			}
		})
	}
}

func BenchmarkTaxHandlerFromUnitValue(b *testing.B) {
	entry := udfs("232.5")
	qty := udfs("3")
	tax := udfs("16")
	th := NewTaxHandlerFromUnitValue()

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		bg := NewFromUnitValue(entry)
		th.WithPercentualTax(tax)

		qtyEv := WithQTY(qty)

		bg.Visit(qtyEv)
		bg.Snapshot() // snapshot of the net value
		th.Visit(bg)
		bg.Snapshot() // snapshot of the brute value
	}
}

func BenchmarkFromBrute(b *testing.B) {
	const maxScale = 12
	const scaleForNet = 6

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		// we define the inmstance of johnny to use
		// this time we will calculate the detail values from the brute value
		bg := NewFromBruteDefault().WithBrute(udfs("1619.1"))

		// We define the visitors to use to calculate the values
		brute := &SnapshotVisitor{}
		net := &SnapshotVisitor{}
		netWD := &SnapshotVisitor{}
		unitValue := NewUnitValue(udfs("3"))

		// using the snapshot visitor we preserve the value of the brute
		bg.Receive(brute)
		// We apply a percentualUntac visitor to remove the taxes from the brute,
		// means we get the net value
		bg.Receive(NewPercentualUnTax(udfs("16")))
		// We apply a snapshot visitor to preserve the net value
		bg.Receive(net)
		// We apply a PercentualUnDiscount visitor to remove the discount from the net,
		// means we get the net value without discount
		bg.Receive(NewPercentualUnDiscount(udfs("0")))
		// We apply a snapshot visitor to preserve the net value without discount
		bg.Receive(netWD)
		// We apply a unit value visitor to get the unit value
		bg.Receive(unitValue)
		// the Round visitor is used to round the values to a given scale
		bg.Receive(NewRound(maxScale))

		_ = net.Get().Round(scaleForNet)
		unitValue.Round(maxScale)
	}
}

func BenchmarkFromUnitValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		unitValue := udfs("1044.543103448276")
		qty := udfs("35157")
		percDiscount := udfs("10")
		amountLineDiscount := udfs("100")
		percTax := udfs("16")
		amountLineTax := qty.Div(gyro.NewHundred()).Round(0).Mul(udfs("0.04"))

		// instance the calculator as a new FromUnitValue
		calc := NewFromUnitValue(unitValue)

		// define the visitors to be used in the calculations
		qtyVisitor := WithQTY(qty)
		percDiscVisitor := NewPercentualDiscount(percDiscount)
		amountDiscVisitor := NewAmountDiscount(amountLineDiscount)
		percTaxVisitor := NewUnbufferedPercTax(percTax)
		amountTaxVisitor := NewUnbufferedAmountTax(amountLineTax)

		// Receive the visitors to the calculator
		calc.Receive(qtyVisitor)
		calc.Receive(percDiscVisitor)
		calc.Receive(amountDiscVisitor)
		calc.Receive(percTaxVisitor)
		calc.Receive(amountTaxVisitor)

		// get the net value from the snapshot visitor
		_ = calc.Snapshot()

		calc.Add(percTaxVisitor.Amount())
		calc.Add(amountTaxVisitor.Amount())

		// get the brute value from the snapshot visitor
		_ = calc.Snapshot()

		// get the total taxes amount from the percTaxVisitor and amountTaxVisitor visitors
		_ = percTaxVisitor.Amount().Add(amountTaxVisitor.Amount())

		// get the total discount amount from the percDiscVisitor and amountDiscVisitor visitors
		_ = percDiscVisitor.Amount().Add(amountDiscVisitor.Amount())
	}
}
