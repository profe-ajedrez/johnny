package main

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

func main() {
	const maxScale = 12
	const scaleForNet = 6

	bg := johnny.NewFromBruteDefault().WithBrute(udfs("1619.1"))

	brute := &johnny.SnapshotVisitor{}
	net := &johnny.SnapshotVisitor{}
	netWD := &johnny.SnapshotVisitor{}
	unitValue := johnny.NewUnitValue(udfs("3"))

	bg.Receive(brute)
	bg.Receive(johnny.NewPercentualUnTax(udfs("16")))
	bg.Receive(net)
	bg.Receive(johnny.NewPercentualUnDiscount(udfs("0")))
	bg.Receive(netWD)
	bg.Receive(unitValue)
	bg.Receive(johnny.NewRound(maxScale))

	netRounded := net.Get().Round(scaleForNet)
	unitValue.Round(maxScale)

	fmt.Printf("Brute value: %v\nNet value: %v\nNet rounded: %v\nNet value with discount: %v\nUnit value: %v\nBuffer value: %v",
		brute.Get().String(), net.Get().String(), netRounded.String(), netWD.Get().String(), unitValue.Get().String(), bg.Value().String())

	//	Output:
	//	Brute value: 1619.1
	//	Net value: 0.139578
	//	Net rounded: 0.0
	//	Net value with discount: 0.139578
	//	Unit value: 0.47
	//	Buffer value: 0.47
}
