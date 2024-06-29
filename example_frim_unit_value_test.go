package johnny_test

import (
	"fmt"

	"github.com/profe-ajedrez/gyro"
	"github.com/profe-ajedrez/johnny"
)

func ExampleFromUnitValue() {
	// Define the values which will be used in the calculations

	unitValue := udfs("1044.543103448276")
	qty := udfs("35157")
	percDiscount := udfs("10")
	amountLineDiscount := udfs("100")
	percTax := udfs("16")
	amountLineTax := qty.Div(gyro.NewHundred()).Round(0).Mul(udfs("0.04"))

	// instance the calculator as a new FromUnitValue
	calc := johnny.NewFromUnitValue(unitValue)

	// define the visitors to be used in the calculations
	qtyVisitor := johnny.WithQTY(qty)
	percDiscVisitor := johnny.NewPercentualDiscount(percDiscount)
	amountDiscVisitor := johnny.NewAmountDiscount(amountLineDiscount)
	percTaxVisitor := johnny.NewUnbufferedPercTax(percTax)
	amountTaxVisitor := johnny.NewUnbufferedAmountTax(amountLineTax)

	// Receive the visitors to the calculator
	calc.Receive(qtyVisitor)
	calc.Receive(percDiscVisitor)
	calc.Receive(amountDiscVisitor)
	calc.Receive(percTaxVisitor)
	calc.Receive(amountTaxVisitor)

	// get the net value from the snapshot visitor
	net := calc.Snapshot()

	calc.Add(percTaxVisitor.Amount())
	calc.Add(amountTaxVisitor.Amount())

	// get the brute value from the snapshot visitor
	brute := calc.Snapshot()

	// get the total taxes amount from the percTaxVisitor and amountTaxVisitor visitors
	totalTaxes := percTaxVisitor.Amount().Add(amountTaxVisitor.Amount())

	// get the total discount amount from the percDiscVisitor and amountDiscVisitor visitors
	totalDiscounts := percDiscVisitor.Amount().Add(amountDiscVisitor.Amount())

	fmt.Println("net: ", net.String())
	fmt.Println("brute: ", brute.String())
	fmt.Println("total discounts: ", totalDiscounts.String())
	fmt.Println("total taxes: ", totalTaxes.String())

	// Output:
	// net:  33050601.699137935398800
	// brute:  38338712.510000050626080
	// total discounts:  3672400.188793103933200
	// total taxes:  5288110.3518620696638080
}
