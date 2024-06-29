// Package johnny provides a custom implementation of the Visitor pattern
// to calculate sales subtotals. It defines several types, including Johnny,
// Handler, DefaultJohnny, FromUnitValue, and FromBrute, which work together
// to perform calculations and transformations on sales data.
package johnny

import (
	"strings"

	"github.com/profe-ajedrez/gyro"
)

var _ Johnny = &DefaultJohnny{}

type Handler interface {
	Snapshot() gyro.Gyro
	Restore(gyro.Gyro)
}

// Johnny is a thing that represents the core operations for working with sales data.
type Johnny interface {
	Receive(Visitor)
	Value() gyro.Gyro

	Add(gyro.Gyro)
	Sub(gyro.Gyro)
	Mul(gyro.Gyro)
	Div(gyro.Gyro)
	set(gyro.Gyro)
	String() string

	Handler
}

// DefaultJohnny is a concrete implementation of the [Johnny] interface.
// It holds the current value of the Johnny as an [gyro.Gyro].
// Is used as a common default implementation of the Johnny interface.
// Also, you could implement your own Johnny type by embedding this struct, to get the basic functionality
type DefaultJohnny struct {
	v gyro.Gyro
}

// Value returns the current value of the Johnny.
func (b *DefaultJohnny) Value() gyro.Gyro {
	return b.v
}

// Add adds the given decimal value to the Johnny.
func (b *DefaultJohnny) Add(v gyro.Gyro) {
	b.v = b.v.Add(v)
}

// Sub subtracts the given gyro.Gyro value from the Johnny.
func (b *DefaultJohnny) Sub(v gyro.Gyro) {
	b.v = b.v.Sub(v)
}

// Mul multiplies the Johnny by the given decimal value.
func (b *DefaultJohnny) Mul(v gyro.Gyro) {
	b.v = b.v.Mul(v)
}

// Div divides the Johnny by the given decimal value.
// This could trigger a division by zero panic because this implementation
// Visitesn't check if the given value is zero or not.
func (b *DefaultJohnny) Div(v gyro.Gyro) {
	b.v = b.v.Div(v)
}

// String returns a string representation of the Johnny value.
func (b *DefaultJohnny) String() string {
	w := strings.Builder{}

	w.WriteString("buffer: ")
	w.WriteString(b.v.String())

	return w.String()
}

// Receive binds the given Visitor to the defaultJohnny instance.
// The Visitor will be invoked with the defaultJohnny instance
// when the Visit method is called on the Visitor.
func (b *DefaultJohnny) Receive(e Visitor) {
	e.Visit(b)
}

// Snapshot returns the current value of the Johnny.
func (b *DefaultJohnny) Snapshot() gyro.Gyro {
	return b.Value()
}

// Restore sets the value of the Johnny instance to the provided decimal value.
func (b *DefaultJohnny) Restore(s gyro.Gyro) {
	b.set(s)
}

func (b *DefaultJohnny) set(s gyro.Gyro) {
	b.v = s
}

type FromUnitValue struct {
	*DefaultJohnny
}

// NewFromUnitValueDefault returns a new instance of FromUnitValue with a zero-valued [gyro.Gyro]
func NewFromUnitValueDefault() FromUnitValue {
	return FromUnitValue{
		DefaultJohnny: &DefaultJohnny{},
	}
}

// NewFromUnitValue returns a new instance of FromUnitValue with the provided entry value
func NewFromUnitValue(entry gyro.Gyro) FromUnitValue {
	return FromUnitValue{
		DefaultJohnny: &DefaultJohnny{
			v: entry,
		},
	}
}

func (f FromUnitValue) Visit(v Visitor) {
	v.Visit(f)
}

// FromBrute is a thing able to be converted from the brute subtotal removing
// elements as discounts and taxes through defined binded visitors
type FromBrute struct {
	*DefaultJohnny
}

// NewFromBruteDefault returns a new instance of FromBrute with a zero-valued.
func NewFromBruteDefault() FromBrute {
	return FromBrute{
		DefaultJohnny: &DefaultJohnny{},
	}
}

// NewFromBrute returns a new instance of FromBrute with the provided brute value set as the Johnny.
func NewFromBrute(brute gyro.Gyro) FromBrute {
	return FromBrute{
		DefaultJohnny: &DefaultJohnny{
			v: brute,
		},
	}
}

// WithBrute sets the value of the FromBrute instance to the provided brute value and returns the updated instance.
func (f FromBrute) WithBrute(brute gyro.Gyro) FromBrute {
	f.v = brute
	return f
}

func (f FromBrute) Receive(v Visitor) {
	v.Visit(f)
}
