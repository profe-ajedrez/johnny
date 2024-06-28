package main

import (
	"fmt"

	"github.com/profe-ajedrez/gyro"
	"github.com/profe-ajedrez/johnny"
)

func udsf(d string) gyro.Gyro {
	return unsafeDecFromStr(d)
}

func unsafeDecFromStr(d string) gyro.Gyro {
	dec, _ := gyro.NewFromString(d)
	return dec
}

func main() {
	// Define the values which will be used in the calculations

	unitValue := udsf("1044.543103448276")
	qty := udsf("35157")
	percDiscount := udsf("10")
	amountLineDiscount := udsf("100")
	percTax := udsf("16")
	amountLineTax := qty.Div(gyro.NewHundred()).Round(0).Mul(udsf("0.04"))

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
}
