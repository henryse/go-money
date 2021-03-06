package money

import (
	"errors"
	"html/template"
	"math"
	"strconv"
)

// Amount is a data structure that stores the amount being used for calculations
type Amount struct {
	val int64
}

// Money represents monetary value information, stores
// currency and amount value
type Money struct {
	amount   *Amount
	currency *Currency
}

// New creates and returns new instance of Money
func New(amount int64, code string) *Money {
	return &Money{
		amount:   &Amount{val: amount},
		currency: newCurrency(code).get(),
	}
}

// New creates and returns new instance of Money
func NewMoney(money *Money) *Money {
	amount := *money.amount
	currency := *money.currency

	return &Money{
		amount:   &amount,
		currency: &currency,
	}
}

// Currency returns the currency used by Money
func NewFloat64(amount float64, code string) *Money {
	amount = math.Round(amount * math.Pow10(currencies[code].Fraction))
	value := int64(amount)

	return New(value, code)
}

// Currency returns the currency used by Money
func NewString(amount string, code string) *Money {
	value, _ := strconv.ParseFloat(amount, 64)

	return NewFloat64(value, code)
}

// Currency returns the currency used by Money
func (m *Money) Currency() *Currency {
	return m.currency
}

// Amount returns a copy of the internal monetary value as an int64
func (m *Money) Amount() int64 {
	return m.amount.val
}

// Amount returns a copy of the internal monetary value as an float64
func (m *Money) AmountFloat64() float64 {
	return float64(m.amount.val) / math.Pow10(m.currency.Fraction)
}

// SameCurrency check if given Money is equals by currency
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return errors.New("currencies don't match")
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	switch {
	case m.amount.val > om.amount.val:
		return 1
	case m.amount.val < om.amount.val:
		return -1
	}

	return 0
}

// Equals checks equality between two Money types
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero
func (m *Money) IsZero() bool {
	return m.amount.val == 0
}

// IsPositive returns boolean of whether the value of Money is positive
func (m *Money) IsPositive() bool {
	return m.amount.val > 0
}

// IsNegative returns boolean of whether the value of Money is negative
func (m *Money) IsNegative() bool {
	return m.amount.val < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value
func (m *Money) Absolute() *Money {
	return &Money{amount: mutate.calc.absolute(m.amount), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value
func (m *Money) Negative() *Money {
	return &Money{amount: mutate.calc.negative(m.amount), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: mutate.calc.add(m.amount, om.amount), currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money
func (m *Money) Subtract(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: mutate.calc.subtract(m.amount, om.amount), currency: m.currency}, nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier
func (m *Money) Multiply(mul int64) *Money {
	return &Money{amount: mutate.calc.multiply(m.amount, mul), currency: m.currency}
}

// Divide returns new Money struct with value representing Self division value by given divider
func (m *Money) Divide(div int64) *Money {
	return &Money{amount: mutate.calc.divide(m.amount, div), currency: m.currency}
}

// Round returns new Money struct with value rounded to nearest zero
func (m *Money) Round() *Money {
	return &Money{amount: mutate.calc.round(m.amount), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("Split must be higher than zero")
	}

	a := mutate.calc.divide(m.amount, int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	l := mutate.calc.modulus(m.amount, int64(n)).val

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, &Amount{1})
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("No ratios specified")
	}

	// Calculate sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int64
	var ms []*Money
	for _, r := range rs {
		party := &Money{
			amount:   mutate.calc.allocate(m.amount, r, sum),
			currency: m.currency,
		}

		ms = append(ms, party)
		total += party.amount.val
	}

	// Calculate leftover value and divide to first parties
	lo := m.amount.val - total
	sub := int64(1)
	if lo < 0 {
		sub = -sub
	}

	if ms != nil {
		for p := 0; lo != 0; p++ {
			ms[p].amount = mutate.calc.add(ms[p].amount, &Amount{sub})
			lo -= sub
		}
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value
func (m *Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.amount.val)
}

// Display lets represent Money struct as string Accounting Format.
func (m *Money) DisplayAccounting() string {
	c := m.currency.get()
	return c.Formatter().FormatAccounting(m.amount.val)
}

// Display Money struct as string Accounting Format.
func (m *Money) DisplayDIV(positive string, negative string) template.HTML {
	c := m.currency.get()
	return c.Formatter().FormatDIV(m.amount.val, positive, negative)
}

// Display Money struct as string Accounting Format.
func (m *Money) DisplayAccountingDIV(positive string, negative string) template.HTML {
	c := m.currency.get()
	return c.Formatter().FormatAccountingDIV(m.amount.val, positive, negative)
}

// Display Money struct as string Accounting Format, but if it is Zero return nothing.
func (m *Money) DisplayAccountingNoZeroDIV(positive string, negative string) template.HTML {
	c := m.currency.get()
	if m.IsZero() {
		return ""
	}
	return c.Formatter().FormatAccountingDIV(m.amount.val, positive, negative)
}
