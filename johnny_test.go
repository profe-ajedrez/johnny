package johnny

import (
	"testing"

	"github.com/profe-ajedrez/gyro"
)

func TestDiscount(t *testing.T) {
	for i, tc := range testCasesDiscounts {
		b, d, shouldFail, whatErrorShouldBe, err := tc.tester()

		if !shouldFail && err != nil {
			t.Logf("[FAIL test case %d] %v", i, err)
			t.FailNow()
		}

		if shouldFail && err == nil {
			t.Logf("[FAIL test case %d %s] error expected: %v", i, tc.name, whatErrorShouldBe)
			t.FailNow()
		}

		if !shouldFail && !b.v.Equal(tc.expecteds.v) {
			t.Logf("[FAIL test case %d %s] got johnny.buffer %v. Expected %v", i, tc.name, b.v, tc.expecteds.v)
			t.FailNow()
		}

		if !shouldFail && !d.ratio.Equal(tc.expecteds.ratio) {
			t.Logf("[FAIL test case %d %s] got discount.ratio %v. Expected %v", i, tc.name, d.ratio, tc.expecteds.ratio)
			t.FailNow()
		}

		if !shouldFail && !d.amount.Equal(tc.expecteds.amount) {
			t.Logf("[FAIL test case %d] got discount.amount %v. Expected %v", i, d.amount, tc.expecteds.amount)
			t.FailNow()
		}
	}
}

func udfs(s string) gyro.Gyro {
	g, _ := gyro.NewFromString(s)
	return g
}

// testCasesDiscounts test cases list
var testCasesDiscounts = []struct {
	name string
	// tester is a function that implements test cases
	// should return the *johnny instance constructed by the case, the one from Discount,
	// a bool indicating whether the case should end with an error,
	// the possible error or, failing that, nil
	// and a string explaining why it should end in an error if this should happen
	tester func() (FromUnitValue, Discount, bool, string, error)

	// expecteds contiene un struct con los datos que la función tester debería producir
	expecteds struct {
		// FromUnitv debe contener los valores que el FromUnitv devuelto por tester debería contener
		FromUnitValue

		// Discount debe contener los valores que el Discount devuelto por tester debería contener
		Discount
	}
}{
	{
		name: "Discount over entry unit v",
		tester: func() (FromUnitValue, Discount, bool, string, error) {
			// Se crea un entry unit v de tipo decimal
			entry := udfs("3.453561112")

			// Se crea un johnny de scala 12 con el entry v
			b := NewFromUnitValue(entry)

			// Se define un valor de descuento porcentual de 10%
			ratio := gyro.NewFromInt64(10)

			// Se crea el Evr con el ratio de descuento indicaVisit
			d1 := NewPercentualDiscount(ratio)

			// Se Visitea para evaluación al Evr d1, que es un descuento porcentual
			b.Visit(d1)

			// En este punto, b.buffer debería contener el valor de b.entryv - 10% del valor de d1.ratio
			// y d1.amount debería contener el valor equivalente al 10% de b.entryv
			return b, d1.Discount, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Discount
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					// expected buffer should be 90% of entryv, because discount is 10%
					v: udfs("3.453561112").Mul(udfs("0.9")),
				},
			},
			Discount: Discount{
				ratio: udfs("10"),
				// expected amount should be 10% of entryv
				amount: udfs("3.453561112").Mul(udfs("0.1")),
			},
		},
	},
	{
		name: "Amount Discount over entry unit v",
		tester: func() (FromUnitValue, Discount, bool, string, error) {
			entry := udfs("3.453561112")
			b := NewFromUnitValue(entry)
			amount := udfs("1.834566333")
			d1 := NewAmountDiscount(amount)
			b.Visit(d1)

			return b, d1.Discount, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Discount
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("3.453561112").Sub(udfs("1.834566333")),
				},
			},
			Discount: Discount{
				ratio:  gyro.NewHundred().Mul(udfs("1.834566333")).Div(udfs("3.453561112")),
				amount: udfs("1.834566333"),
			},
		},
	},
	{
		name: "Combo Discount percentual over entry unit v and other considering quantity",
		tester: func() (FromUnitValue, Discount, bool, string, error) {
			entry := udfs("100.123")
			b := NewFromUnitValue(entry)

			ratio := udfs("10")
			discountOverEntryv := NewPercentualDiscount(ratio)

			b.Visit(discountOverEntryv)

			ratio = udfs("15")
			discountConsideringQty := NewPercentualDiscount(ratio)

			qty := udfs("10")

			b.Visit(WithQTY(qty))
			b.Visit(discountConsideringQty)

			totalDiscountApplied := Discount{
				ratio:  udfs("1"),
				amount: udfs("100"),
			}

			return b, totalDiscountApplied, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Discount
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("100.123").Mul(udfs("0.9")).Mul(udfs("10")).Mul(udfs("0.85")),
				},
			},
			Discount: Discount{
				ratio:  udfs("1"),
				amount: udfs("100"),
			},
		},
	},
}

func TestTax(t *testing.T) {
	for i, tc := range testCasesTaxes {
		b, d, shouldFail, whatErrorShouldBe, err := tc.tester()

		if !shouldFail && err != nil {
			t.Logf("[FAIL test case %d] %v", i, err)
			t.FailNow()
		}

		if shouldFail && err == nil {
			t.Logf("[FAIL test case %d %s] error expected: %v", i, tc.name, whatErrorShouldBe)
			t.FailNow()
		}

		if !shouldFail && !b.v.Equal(tc.expecteds.v) {
			t.Logf("[FAIL test case %d %s] got johnny.buffer %v. Expected %v", i, tc.name, b.v, tc.expecteds.v)
			t.FailNow()
		}

		if !shouldFail && !d.ratio.Equal(tc.expecteds.ratio) {
			t.Logf("[FAIL test case %d %s] got tax.ratio %v. Expected %v", i, tc.name, d.ratio, tc.expecteds.ratio)
			t.FailNow()
		}

		if !shouldFail && !d.amount.Equal(tc.expecteds.amount) {
			t.Logf("[FAIL test case %d] got tax.amount %v. Expected %v", i, d.amount, tc.expecteds.amount)
			t.FailNow()
		}

		if !shouldFail && !d.taxable.Equal(tc.expecteds.taxable) {
			t.Logf("[FAIL test case %d] got tax.taxable %v. Expected %v", i, d.taxable, tc.expecteds.taxable)
			t.FailNow()
		}
	}
}

// testCasesTaxes test cases list for taxes
var testCasesTaxes = []struct {
	name string
	// tester is a function that implements test cases
	// should return the *johnny instance constructed by the case, the one from Tax,
	// a bool indicating whether the case should end with an error,
	// the possible error or, failing that, nil
	// and a string explaining why it should end in an error if this should happen
	tester func() (FromUnitValue, Tax, bool, string, error)

	// expecteds contiene un struct con los datos que la función tester debería producir
	expecteds struct {
		// FromUnitValue debe contener los valores que el FromUnitv devuelto por tester debería contener
		FromUnitValue

		// Tax debe contener los valores que el Discount devuelto por tester debería contener
		Tax
	}
}{
	{
		name: "Tax over entry unit v",
		tester: func() (FromUnitValue, Tax, bool, string, error) {
			// Se crea un entry unit v de tipo decimal
			entry := udfs("17.3475345")

			// Se crea un johnny de scala 12 con el entry v
			b := NewFromUnitValue(entry)

			// Se define un valor de impuesto porcentual de 10%
			ratio := gyro.NewFromInt64(10)

			// Se crea el Evr con el ratio de impuesto indicaVisit
			t1 := NewPercTax(ratio)

			// Se Receiveea para evaluación al Evr d1, que es un descuento porcentual
			b.Receive(t1)

			// En este punto, b.buffer debería contener el valor de b.entryv - 10% del valor de d1.ratio
			// y d1.amount debería contener el valor equivalente al 10% de b.entryv
			return b, t1.Tax, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Tax
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					// expected buffer should be 90% of entryv, because discount is 10%
					v: udfs("17.3475345").Mul(udfs("1.1")),
				},
			},
			Tax: Tax{
				ratio: udfs("10"),
				// expected amount should be 10% of entryv
				amount: udfs("17.3475345").Mul(udfs("0.1")),

				taxable: udfs("17.3475345"),
			},
		},
	},
	{
		name: "Amount tax over entry unit v",
		tester: func() (FromUnitValue, Tax, bool, string, error) {
			entry := udfs("100")
			b := NewFromUnitValue(entry)
			amount := udfs("9.8")
			t1 := NewAmountTax(amount)
			b.Receive(t1)

			return b, t1.Tax, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Tax
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("109.8"),
				},
			},
			Tax: Tax{
				ratio:   udfs("9.8"),
				amount:  udfs("9.8"),
				taxable: udfs("100"),
			},
		},
	},
	{
		name: "Combo Tax percentual over entry unit v and other considering quantity",
		tester: func() (FromUnitValue, Tax, bool, string, error) {
			entry := udfs("100.123")
			b := NewFromUnitValue(entry)

			ratio := udfs("10")
			taxOverEntryv := NewPercTax(ratio)

			b.Receive(taxOverEntryv)

			ratio = udfs("15")
			taxConsideringQty := NewPercTax(ratio)

			qty := udfs("10")

			b.Receive(WithQTY(qty))
			b.Receive(taxConsideringQty)

			// Nos interesa validar el valor de b.buffer, pues no se deberían mezclar
			// impuestos calculaVisits en distintos pasos
			totalTaxApplied := Tax{
				ratio:   udfs("1"),
				amount:  udfs("100"),
				taxable: udfs("1"),
			}

			return b, totalTaxApplied, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Tax
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("1266.55595"),
				},
			},
			Tax: Tax{
				ratio:   udfs("1"),
				amount:  udfs("100"),
				taxable: udfs("1"),
			},
		},
	},
	{
		name: "Tax over entry unit v with quantity",
		tester: func() (FromUnitValue, Tax, bool, string, error) {
			entry := udfs("17.3475345")
			b := NewFromUnitValue(entry)
			ratio := gyro.NewFromInt64(10)
			t1 := NewPercTax(ratio)

			qty := udfs("5")
			b.Receive(WithQTY(qty))
			b.Receive(t1)

			return b, t1.Tax, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Tax
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("17.3475345").Mul(udfs("5")).Mul(udfs("1.1")),
				},
			},
			Tax: Tax{
				ratio:   udfs("10"),
				amount:  udfs("17.3475345").Mul(udfs("0.1")).Mul(udfs("5")),
				taxable: udfs("17.3475345").Mul(udfs("5")),
			},
		},
	},
	{
		name: "Tax over entry unit v with multiple taxes",
		tester: func() (FromUnitValue, Tax, bool, string, error) {
			entry := udfs("100")
			b := NewFromUnitValue(entry)
			ratio1 := udfs("10")
			t1 := NewPercTax(ratio1)

			ratio2 := udfs("5")
			t2 := NewPercTax(ratio2)

			b.Receive(t1)
			// this tax will be applied over the previous tax t1
			b.Receive(t2)

			return b, t2.Tax, false, "", nil
		},
		expecteds: struct {
			FromUnitValue
			Tax
		}{
			FromUnitValue: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("100").Mul(udfs("1.1")).Mul(udfs("1.05")),
				},
			},
			Tax: Tax{
				ratio:   udfs("5"),
				amount:  udfs("100").Mul(udfs("1.1")).Mul(udfs("0.05")),
				taxable: udfs("100").Mul(udfs("1.1")),
			},
		},
	},
}

func TestPercentualUndiscount(t *testing.T) {
	testCases := []struct {
		name     string
		initial  gyro.Gyro
		ratio    gyro.Gyro
		expected gyro.Gyro
	}{
		{
			name:     "Fractional ratio",
			initial:  udfs("100"),
			ratio:    udfs("5.5"),
			expected: udfs("105.820105820105800"),
		},
		{
			name:     "Positive ratio",
			initial:  udfs("100"),
			ratio:    udfs("10"),
			expected: udfs("111.111111111111100"),
		},
		{
			name:     "Negative ratio",
			initial:  udfs("100"),
			ratio:    udfs("-10"),
			expected: udfs("90.909090909090900"),
		},
		{
			name:     "Zero ratio",
			initial:  udfs("100"),
			ratio:    udfs("0"),
			expected: udfs("100"),
		},
	}

	for k, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &DefaultJohnny{v: tc.initial}
			u := NewPercentualUnDiscount(tc.ratio)
			u.Visit(b)
			if !b.v.Equal(tc.expected) {
				t.Errorf("[test case %d] Expected %v, got %v", tc.expected, k, b.v)
			}
		})
	}
}

func TestAmountUndiscount(t *testing.T) {
	testCases := []struct {
		name     string
		initial  gyro.Gyro
		amount   gyro.Gyro
		expected gyro.Gyro
	}{
		{
			name:     "Positive amount",
			initial:  udfs("100"),
			amount:   udfs("10"),
			expected: udfs("110"),
		},
		{
			name:     "Negative amount",
			initial:  udfs("100"),
			amount:   udfs("-10"),
			expected: udfs("90"),
		},
		{
			name:     "Zero amount",
			initial:  udfs("100"),
			amount:   udfs("0"),
			expected: udfs("100"),
		},
		{
			name:     "Fractional amount",
			initial:  udfs("100"),
			amount:   udfs("5.5"),
			expected: udfs("105.5"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &DefaultJohnny{v: tc.initial}
			u := NewAmountUnDiscount(tc.amount)
			u.Visit(b)
			if !b.v.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, b.v)
			}
		})
	}
}

func TestHandlerFromUnitv(t *testing.T) {
	for i, tc := range testCaseTaxHandlerFromUnitValue {
		b, th, shouldFail, err := tc.tester()

		if !shouldFail && err != nil {
			t.Logf("[test %d FAILED] should not fail. %v", i, err)
			t.FailNow()
		}

		if shouldFail && err == nil {
			t.Logf("[test %d FAILED] should has been failed.", i)
			t.FailNow()
		}

		if shouldFail {
			continue
		}

		if !b.Value().Equal(tc.expected.Value()) {
			t.Logf("[test %d FAILED] buffer. Got %v  Expected %v", i, b.Value(), tc.expected.Value())
			t.FailNow()
		}

		if !th.totalRatio.Equal(tc.expected.totalRatio) {
			t.Logf("[test %d FAILED] total ratio. Got %v  Expected %v", i, th.totalRatio, tc.expected.totalRatio)
			t.FailNow()
		}

		if !th.totalAmount.Equal(tc.expected.totalAmount) {
			t.Logf("[test %d FAILED] total amount. Got %v  Expected %v", i, th.totalAmount, tc.expected.totalAmount)
			t.FailNow()
		}

		if !th.taxable.Equal(tc.expected.taxable) {
			t.Logf("[test %d FAILED] taxable. Got %v  Expected %v", i, th.taxable, tc.expected.taxable)
			t.FailNow()
		}
	}
}

var testCaseTaxHandlerFromUnitValue = []struct {
	tester   func() (Johnny, *TaxHandlerFromUnitValue, bool, error)
	expected struct {
		Johnny
		TaxHandlerFromUnitValue
	}
}{
	{
		tester: func() (Johnny, *TaxHandlerFromUnitValue, bool, error) {
			entry := udfs("232.5")
			b := NewFromUnitValue(entry)

			b.Receive(WithQTY(udfs("3")))

			th := NewTaxHandlerFromUnitValue()
			th.WithPercentualTax(udfs("16"))
			net := SnapshotVisitor{}
			net.Visit(b)
			th.Visit(b)
			brute := SnapshotVisitor{}
			brute.Visit(b)

			return b, th, false, nil
		},
		expected: struct {
			Johnny
			TaxHandlerFromUnitValue
		}{
			Johnny: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("232.5").Mul(udfs("3")).Mul(udfs("1.16")),
				},
			},
			TaxHandlerFromUnitValue: TaxHandlerFromUnitValue{
				TaxHandler: &TaxHandler{
					totalRatio:  udfs("16"),
					totalAmount: udfs("232.5").Mul(udfs("0.16")).Mul(udfs("3")),
					taxable:     udfs("232.5").Mul(udfs("3")),
				},
			},
		},
	},
	{
		tester: func() (Johnny, *TaxHandlerFromUnitValue, bool, error) {
			entry := udfs("100")
			b := NewFromUnitValue(entry)

			b.Receive(WithQTY(udfs("2")))

			th := NewTaxHandlerFromUnitValue()
			th.WithPercentualTax(udfs("20"))

			net := SnapshotVisitor{}
			net.Visit(b)

			th.Visit(b)

			brute := SnapshotVisitor{}
			brute.Visit(b)

			return b, th, false, nil
		},
		expected: struct {
			Johnny
			TaxHandlerFromUnitValue
		}{
			Johnny: FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("100").Mul(udfs("2")).Mul(udfs("1.2")),
				},
			},
			TaxHandlerFromUnitValue: TaxHandlerFromUnitValue{
				TaxHandler: &TaxHandler{
					totalRatio:  udfs("20"),
					totalAmount: udfs("100").Mul(udfs("0.2")).Mul(udfs("2")),
					taxable:     udfs("100").Mul(udfs("2")),
				},
			},
		},
	},
	{
		tester: func() (Johnny, *TaxHandlerFromUnitValue, bool, error) {
			entry := udfs("50")
			b := NewFromUnitValue(entry)

			b.Receive(WithQTY(udfs("4")))

			th := NewTaxHandlerFromUnitValue()
			th.WithPercentualTax(udfs("8"))

			net := SnapshotVisitor{}
			net.Visit(b)

			th.Visit(b)

			brute := SnapshotVisitor{}
			brute.Visit(b)

			return b, th, false, nil
		},
		expected: struct {
			Johnny
			TaxHandlerFromUnitValue
		}{
			Johnny: &FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("50").Mul(udfs("4")).Mul(udfs("1.08")),
				},
			},
			TaxHandlerFromUnitValue: TaxHandlerFromUnitValue{
				TaxHandler: &TaxHandler{
					totalRatio:  udfs("8"),
					totalAmount: udfs("50").Mul(udfs("0.08")).Mul(udfs("4")),
					taxable:     udfs("50").Mul(udfs("4")),
				},
			},
		},
	},
	{
		tester: func() (Johnny, *TaxHandlerFromUnitValue, bool, error) {
			entry := udfs("75.25")
			b := NewFromUnitValue(entry)

			b.Receive(WithQTY(udfs("1.5")))

			th := NewTaxHandlerFromUnitValue()
			th.WithPercentualTax(udfs("12.5"))

			net := SnapshotVisitor{}
			net.Visit(b)

			th.Visit(b)

			brute := SnapshotVisitor{}
			brute.Visit(b)

			return b, th, false, nil
		},
		expected: struct {
			Johnny
			TaxHandlerFromUnitValue
		}{
			Johnny: &FromUnitValue{
				DefaultJohnny: &DefaultJohnny{
					v: udfs("75.25").Mul(udfs("1.5")).Mul(udfs("1.125")),
				},
			},
			TaxHandlerFromUnitValue: TaxHandlerFromUnitValue{
				TaxHandler: &TaxHandler{
					totalRatio:  udfs("12.5"),
					totalAmount: udfs("75.25").Mul(udfs("0.125")).Mul(udfs("1.5")),
					taxable:     udfs("75.25").Mul(udfs("1.5")),
				},
			},
		},
	},
}
