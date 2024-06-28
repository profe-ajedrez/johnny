package johnny_test

import (
	"fmt"

	"github.com/profe-ajedrez/gyro"
	"github.com/profe-ajedrez/johnny"
)

// udfs stands for unsafe decimal from string.
// helps to have a decimal value from string ignoring errors.
func udfs(s string) gyro.Gyro {
	d, _ := gyro.NewFromString(s)
	return d
}

func ExampleFromBrute() {
	const maxScale = 12
	const scaleForNet = 6

	// we define the inmstance of johnny to use
	// this time we will calculate the detail values from the brute value
	bg := johnny.NewFromBruteDefault().WithBrute(udfs("1619.1"))

	// We define the visitors to use to calculate the values
	brute := &johnny.SnapshotVisitor{}
	net := &johnny.SnapshotVisitor{}
	netWD := &johnny.SnapshotVisitor{}
	unitValue := johnny.NewUnitValue(udfs("3"))

	// using the snapshot visitor we preserve the value of the brute
	bg.Receive(brute)
	// We apply a percentualUntac visitor to remove the taxes from the brute,
	// means we get the net value
	bg.Receive(johnny.NewPercentualUnTax(udfs("16")))
	// We apply a snapshot visitor to preserve the net value
	bg.Receive(net)
	// We apply a PercentualUnDiscount visitor to remove the discount from the net,
	// means we get the net value without discount
	bg.Receive(johnny.NewPercentualUnDiscount(udfs("0")))
	// We apply a snapshot visitor to preserve the net value without discount
	bg.Receive(netWD)
	// We apply a unit value visitor to get the unit value
	bg.Receive(unitValue)
	// the Round visitor is used to round the values to a given scale
	bg.Receive(johnny.NewRound(maxScale))

	netRounded := net.Get().Round(scaleForNet)
	unitValue.Round(maxScale)

	fmt.Printf("Brute value: %v\nNet value: %v\nNet rounded: %v\nNet value with discount: %v\nUnit value: %v\nBuffer value: %v",
		brute.Get().String(), net.Get().String(), netRounded.String(), netWD.Get().String(), unitValue.Get().String(), bg.Value().String())

	// Output:
	// Brute value: 1619.1
	// Net value: 0.139578
	// Net rounded: 0.0
	// Net value with discount: 0.139578
	// Unit value: 0.47
	// Buffer value: 0.47
}
