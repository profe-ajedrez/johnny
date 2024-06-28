package johnny

import (
	"strings"

	"github.com/profe-ajedrez/gyro"
)

// Visitor is an interface that defines a method for performing an operation
// on a Johnny.
// The Visit method takes a Johnny as an argument and performs some operation on it.
type Visitor interface {
	Visit(Johnny)
}

var _ Visitor = &PercentualDiscount{}
var _ Visitor = &AmountDiscount{}

// Discount represents a discount that can be applied to a Johnny value.
// The ratio field represents the percentage discount, and the amount field
// represents the fixed amount discount.
type Discount struct {
	ratio  gyro.Gyro
	amount gyro.Gyro
}

// Ratio returns the ratio of the discount.
func (d *Discount) Ratio() gyro.Gyro {
	return d.ratio
}

// Amount returns the fixed amount discount.
func (d *Discount) Amount() gyro.Gyro {
	return d.amount
}

// String returns a string representation of the Discount, including the ratio and amount.
func (d *Discount) String() string {
	w := strings.Builder{}

	w.WriteString("ratio: ")
	w.WriteString(d.ratio.String())
	w.WriteString(" amount: ")
	w.WriteString(d.amount.String())
	return w.String()
}

// PercentualDiscount represents a discount that is applied as a percentage of the Johnny value.
// It embeds the Discount struct, which contains the ratio and amount fields.
type PercentualDiscount struct {
	Discount
}

// NewPercentualDiscount creates a new PercentualDiscount instance with the given ratio.
// The ratio represents the percentage discount to be applied.
func NewPercentualDiscount(ratio gyro.Gyro) *PercentualDiscount {
	return &PercentualDiscount{
		Discount: Discount{
			ratio: ratio,
		},
	}
}

// Visit applies the percentual discount to the given Johnny value.
// It calculates the discount amount by multiplying the Johnny value's buffer
// by the discount ratio, and then dividing by 100 to get the percentage.
// The calculated discount amount is then subtracted from the Johnny value.
// This implemenetation Visitesnt check for negative discounts
func (pd *PercentualDiscount) Visit(b Johnny) {
	pd.amount = b.Value().Mul(pd.ratio).Div(gyro.NewHundred())
	b.Sub(pd.amount)
}

// AmountDiscount represents a discount that is applied as a fixed amount.
// It embeds the Discount struct, which contains the ratio and amount fields.
type AmountDiscount struct {
	Discount
}

// NewAmountDiscount creates a new AmountDiscount instance with the given fixed amount.
// The amount represents the fixed discount to be applied.
func NewAmountDiscount(amount gyro.Gyro) *AmountDiscount {
	return &AmountDiscount{
		Discount: Discount{
			amount: amount,
		},
	}
}

// Visit applies the fixed amount discount to the given Johnny value.
// If the Johnny value's buffer is zero, the discount ratio is set to zero.
// This implemenetation Visitesnt check for negative discounts
func (pd *AmountDiscount) Visit(b Johnny) {
	if b.Value().Equal(gyro.NewZero()) {
		pd.ratio = gyro.NewZero()
		return
	}

	pd.ratio = gyro.NewHundred().Mul(pd.amount).Div(b.Value())
	b.Sub(pd.amount)
}

type Qty struct {
	qty gyro.Gyro
}

// Qty represents a quantity visitor that multiplies the Johnny instance by a given gyro.Gyro value.
func WithQTY(qty gyro.Gyro) Qty {
	return Qty{qty: qty}
}

// WithQTY returns a new Qty instance with the provided gyro.Gyro value.
func (q Qty) Visit(b Johnny) {
	b.Mul(q.qty)
}

// Visit multiplies the given Johnny instance by the Qty's gyro.Gyro value.
type UnitValue struct {
	qty       gyro.Gyro
	unitValue gyro.Gyro
}

func NewUnitValue(qty gyro.Gyro) *UnitValue {
	return &UnitValue{
		qty: qty,
	}
}

// NewUnitValue returns a new instance of UnitValue with the provided quantity value.
func (q *UnitValue) Visit(b Johnny) {
	if q.qty.Cmp(gyro.NewZero()) > 0 {
		q.unitValue = b.Value().Div(q.qty)
		b.set(q.unitValue)
	}
}

func (q *UnitValue) Get() gyro.Gyro {
	return q.unitValue
}

// Get returns the underlying [gyro.Gyro] value of the UnitValue instance.
func (q *UnitValue) Round(sc int32) {
	q.unitValue = q.unitValue.Round(sc)
}

// Tax struct holds the components necessary for tax calculation on a Johnny value.
// It includes the tax ratio, the tax amount, and the taxable base amount.
// This struct is typically used as a visitor to apply tax calculations to a Johnny value.
type Tax struct {
	ratio   gyro.Gyro
	amount  gyro.Gyro
	taxable gyro.Gyro
}

// Amount returns the amount of tax calculated for the Johnny value.
func (pt *Tax) Amount() gyro.Gyro {
	return pt.amount
}

// Ratio returns the tax ratio for the Johnny value.
func (pt *Tax) Ratio() gyro.Gyro {
	return pt.ratio
}

// Taxable returns the taxable value for the Johnny value that this Tax was applied to.
func (pt *Tax) Taxable() gyro.Gyro {
	return pt.taxable
}

// PercTax is a percentual based tax visitor.
// Wraps over Tax structure and implements the Visit method.
type PercTax struct {
	Tax
}

// NewPercTax creates a new PercTax instance with the given tax ratio.
// The PercTax struct wraps over the Tax struct and implements the Visit method
// to calculate the tax amount based on the given ratio.
func NewPercTax(ratio gyro.Gyro) *PercTax {
	return &PercTax{
		Tax: Tax{
			ratio: ratio,
		},
	}
}

// Visit applies the percentual tax to the Johnny instance's value and updates the Johnny instance's buffer directly.
// It calculates the tax amount based on the Johnny instance's value and the tax ratio,
// and adds the tax amount to the Johnny instance's buffer.
// It also stores the calculated tax amount and the taxable value in the PercTax struct.
// This implemenetation doesnt check for negative taxes
func (pt *PercTax) Visit(b Johnny) {
	pt.amount = b.Value().Mul(pt.ratio.Div(gyro.NewHundred()))
	pt.taxable = b.Value()
	b.Add(pt.amount)
}

// UnbufferedPercTax is a Visitor that applies a percentual tax to the Johnny instance's value.
// It does not modify the Johnny instance's buffer directly.
type UnbufferedPercTax struct {
	Tax
}

// NewUnbufferedPercTax returns a new instance of [UnbufferedPercTax] with the provided ratio value.
// The ratio value is used to initialize the [Tax] struct, which is embedded in the UnbufferedPercTax struct.
func NewUnbufferedPercTax(ratio gyro.Gyro) *UnbufferedPercTax {
	return &UnbufferedPercTax{
		Tax: Tax{
			ratio: ratio,
		},
	}
}

// Visit applies the percentual tax to the Johnny instance's value.
// It does not modify the Johnny instance's buffer directly.
// This implemenetation doesnt check for negative taxes
func (pt *UnbufferedPercTax) Visit(b Johnny) {
	pt.amount = b.Value().Mul(pt.ratio.Div(gyro.NewHundred()))
	pt.taxable = b.Value()
}

// AmountTax is a Tax that applies a fixed amount to the Johnny value.
// It wraps over the Tax struct and implements the Visit method to calculate the tax amount.
type AmountTax struct {
	Tax
}

// NewAmountTax returns a new instance of AmountTax with the provided amount value.
// The amount value is used to initialize the Tax struct, which is embedded in the AmountTax struct.
func NewAmountTax(amount gyro.Gyro) *AmountTax {
	return &AmountTax{
		Tax: Tax{
			amount: amount,
		},
	}
}

// Visit implements the Visitor pattern for AmountTax.
// It calculates the taxable value by adding the amount to the Johnny instance,
// and then calculates the ratio by dividing the amount by the taxable value.
func (pt *AmountTax) Visit(b Johnny) {
	pt.taxable = b.Value()
	b.Add(pt.amount)
	pt.ratio = pt.amount.Mul(gyro.NewHundred()).Div(pt.taxable)
}

// UnbufferedAmountTax is a Tax that applies a fixed amount to the Johnny value.
// It wraps over the Tax struct and implements the Visit method to calculate the tax amount,
// but does not modify the Johnny instance's buffer directly.
type UnbufferedAmountTax struct {
	Tax
}

// NewUnbufferedAmountTax returns a new instance of UnbufferedAmountTax with the provided amount value.
// The amount value is used to initialize the Tax struct, which is embedded in the UnbufferedAmountTax struct.
func NewUnbufferedAmountTax(amount gyro.Gyro) *UnbufferedAmountTax {
	return &UnbufferedAmountTax{
		Tax: Tax{
			amount: amount,
		},
	}
}

// Visit applies the fixed amount tax to the Johnny instance's value.
// It does not modify the Johnny instance's buffer directly.
// This implementation calculates the tax ratio based on the fixed amount and the Johnny instance's value.
// This implemenetation doesnt check for negative taxes
func (pt *UnbufferedAmountTax) Visit(b Johnny) {
	pt.taxable = b.Value()
	pt.ratio = pt.amount.Mul(gyro.NewHundred()).Div(pt.taxable)
}

// PercentualUndiscount represents a percentual undiscount.
type PercentualUndiscount struct {
	*Discount
}

// NewPercentualUnDiscount returns a new instance of PercentualUndiscount with the provided ratio value.
func NewPercentualUnDiscount(ratio gyro.Gyro) *PercentualUndiscount {
	return &PercentualUndiscount{
		Discount: &Discount{
			ratio: ratio,
		},
	}
}

// Visit implements the Visitor pattern for PercentualUndiscount.
// It calculates the undiscounted value by dividing the current value by the ratio,
// and then sets the result as the new value of the Johnny instance.
// It also calculates the amount of the undiscount by multiplying the new value by the ratio.
func (u *PercentualUndiscount) Visit(b Johnny) {
	if u.ratio.Equal(gyro.NewZero()) {
		return
	}

	d := gyro.NewHundred().Sub(u.ratio)
	v := b.Value().Div(d)
	v = v.Mul(gyro.NewHundred())
	b.set(v)
	u.amount = b.Value().Mul(u.ratio.Div(gyro.NewHundred()))
}

// AmountUndiscount represents an amount undiscount.
type AmountUndiscount struct {
	*Discount
}

// NewAmountUnDiscount returns a new instance of AmountUndiscount with the provided amount value.
func NewAmountUnDiscount(amount gyro.Gyro) *AmountUndiscount {
	return &AmountUndiscount{
		Discount: &Discount{
			amount: amount,
		},
	}
}

// Visit implements the Visitor pattern for AmountUndiscount.
// It adds the amount to the current value of the Johnny instance,
// and then calculates the ratio by dividing the amount by the new value.
func (u *AmountUndiscount) Visit(b Johnny) {
	b.Add(u.amount)
	u.ratio = u.amount.Mul(gyro.NewHundred()).Div(b.Value())
}

// Round is a visitor which performs a rounding operation with a specified scale.
// rounding usually implies a rescale operation, which is costly, use with care.
type Round struct {
	scale int32
}

// NewRound creates a new Round visitor with the specified scale.
// The Round visitor can be used to perform a rounding operation on a gyro.Gyro value.
// The scale parameter determines the number of decimal places to round to.
func NewRound(scale int32) Round {
	return Round{
		scale: scale,
	}
}

// Do applies a rounding operation to the given gyro.Gyro value, using the scale
// specified when the Round visitor was created. This effectively rescales the
// gyro.Gyro value to the desired number of decimal places.
func (r Round) Visit(b Johnny) {
	b.set(b.Value().Round(r.scale))
}

// SnapshotVisitor is a visitor that takes a snapshot of the current value of a gyro.Gyro.
type SnapshotVisitor struct {
	// buffer stores the snapshot of the gyro.Gyro value.
	buffer gyro.Gyro
}

// NewSnapshot returns a new instance of SnapshotVisitor.
func NewSnapshot() *SnapshotVisitor {
	return &SnapshotVisitor{}
}

// Visit sets the buffer to the current value of the Johnny object.
func (s *SnapshotVisitor) Visit(b Johnny) {
	s.buffer = b.Value()
}

// Get returns the snapshot of the gyro.Gyro value.
func (s *SnapshotVisitor) Get() gyro.Gyro {
	return s.buffer
}

// PercentualUntax is a tax calculator that calculates the tax as a percentage of the value.
type PercentualUntax struct {
	// Tax is the base tax structure.
	Tax
}

// NewPercentualUnTax returns a new instance of PercentualUntax with the given ratio.
func NewPercentualUnTax(ratio gyro.Gyro) *PercentualUntax {
	return &PercentualUntax{
		Tax: Tax{
			ratio: ratio,
		},
	}
}

// Visit calculates the tax amount based on the ratio and updates the Johnny object.
func (pu *PercentualUntax) Visit(b Johnny) {
	ratio := pu.ratio.Div(gyro.NewHundred())
	b.set(b.Value().Div(gyro.NewOne().Add(ratio)))
	pu.amount = b.Value().Mul(ratio)
}

// AmountUntax is a tax calculator that calculates the tax as a fixed amount.
type AmountUntax struct {
	// Tax is the base tax structure.
	Tax
}

// NewAmountUnTax returns a new instance of AmountUntax with the given amount.
func NewAmountUnTax(amount gyro.Gyro) *AmountUntax {
	return &AmountUntax{
		Tax: Tax{
			amount: amount,
		},
	}
}

// Visit calculates the ratio based on the amount and updates the Johnny object.
func (pu *AmountUntax) Visit(b Johnny) {
	b.Sub(pu.amount)
	pu.ratio = pu.amount.Mul(gyro.NewHundred()).Div(b.Value())
}

// TaxHandler is a handler that applies multiple taxes to a gyro.Gyro value.
type TaxHandler struct {
	// totalRatio is the total ratio of all taxes.
	totalRatio gyro.Gyro
	// totalAmount is the total amount of all taxes.
	totalAmount gyro.Gyro
	// taxable is the original value that taxes are applied to.
	taxable gyro.Gyro
}

// NewTaxHandler returns a new instance of TaxHandler.
func NewTaxHandler() *TaxHandler {
	return &TaxHandler{}
}

// WithPercentualTax adds a new percentual tax to the total ratio.
func (t *TaxHandler) WithPercentualTax(value gyro.Gyro) {
	t.totalRatio = t.totalRatio.Add(value)
}

// WithAmountTax adds a new amount tax to the total amount.
func (t *TaxHandler) WithAmountTax(value gyro.Gyro) {
	t.totalAmount = t.totalAmount.Add(value)
}

// TaxHandlerFromUnitValue is a handler that applies multiple taxes to a gyro.Gyro value, starting from a unit value.
type TaxHandlerFromUnitValue struct {
	*TaxHandler
}

// NewTaxHandlerFromUnitValue returns a new instance of TaxHandlerFromUnitValue.
func NewTaxHandlerFromUnitValue() *TaxHandlerFromUnitValue {
	return &TaxHandlerFromUnitValue{
		TaxHandler: NewTaxHandler(),
	}
}

// Visit applies the taxes to the given Johnny object.
func (t *TaxHandlerFromUnitValue) Visit(b Johnny) {
	t.taxable = b.Value()

	t1 := NewPercTax(t.totalRatio)
	t2 := NewAmountTax(t.totalAmount)

	do(b, t1, t2)

	t.totalRatio = t1.ratio.Add(t2.ratio)
	t.totalAmount = t1.amount.Add(t2.amount)
}

// Taxable returns the original value that taxes are applied to.
func (t *TaxHandlerFromUnitValue) Taxable() gyro.Gyro {
	return t.taxable
}

// TotalRatio returns the total ratio of all taxes.
func (t *TaxHandlerFromUnitValue) TotalRatio() gyro.Gyro {
	return t.totalRatio
}

// TotalAmount returns the total amount of all taxes.
func (t *TaxHandlerFromUnitValue) TotalAmount() gyro.Gyro {
	return t.totalAmount
}

// DiscountHandler is a handler that applies multiple discounts to a gyro.Gyro value.
type DiscountHandler struct {
	// totalRatio is the total ratio of all discounts.
	totalRatio gyro.Gyro
	// totalAmount is the total amount of all discounts.
	totalAmount gyro.Gyro
	// discountable is the original value that discounts are applied to.
	discountable gyro.Gyro
}

// NewDiscountHandler returns a new instance of DiscountHandler.
func NewDiscountHandler() *DiscountHandler {
	return &DiscountHandler{}
}

// WithPercentualDiscount adds a new percentual discount to the total ratio.
func (t *DiscountHandler) WithPercentualDiscount(value gyro.Gyro) {
	t.totalRatio = t.totalRatio.Add(value)
}

// WithAmountDiscount adds a new amount discount to the total amount.
func (t *DiscountHandler) WithAmountDiscount(value gyro.Gyro) {
	t.totalAmount = t.totalAmount.Add(value)
}

// DiscountHandlerFromUnitValue is a handler that applies multiple discounts to a gyro.Gyro value, starting from a unit value.
type DiscountHandlerFromUnitValue struct {
	*DiscountHandler
}

// NewDiscHandlerFromUnitValue returns a new instance of DiscountHandlerFromUnitValue.
func NewDiscHandlerFromUnitValue() *DiscountHandlerFromUnitValue {
	return &DiscountHandlerFromUnitValue{
		DiscountHandler: NewDiscountHandler(),
	}
}

// Visit applies the discounts to the given Johnny object.
func (t *DiscountHandlerFromUnitValue) Visit(b Johnny) {
	t.discountable = b.Value()

	t1 := NewPercentualDiscount(t.totalRatio)
	t2 := NewAmountDiscount(t.totalAmount)

	do(b, t1, t2)

	t.totalRatio = t1.ratio.Add(t2.ratio)
	t.totalAmount = t1.amount.Add(t2.amount)
}

// Discountable returns the original value that discounts are applied to.
func (t *DiscountHandlerFromUnitValue) Discountable() gyro.Gyro {
	return t.discountable
}

// TotalRatio returns the total ratio of all discounts.
func (t *DiscountHandlerFromUnitValue) TotalRatio() gyro.Gyro {
	return t.totalRatio
}

// TotalAmount returns the total amount of all discounts.
func (t *DiscountHandlerFromUnitValue) TotalAmount() gyro.Gyro {
	return t.totalAmount
}

// do applies the given visitors to the Johnny object.
func do(b Johnny, e1, e2 Visitor) {
	b.Receive(e1)
}
