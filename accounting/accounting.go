package accounting

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/LUSHDigital/core-lush/currency"
)

const (
	// QuadruplePrecision describes 128 bits of precision for IEEE 754 decimals
	// see: https://en.wikipedia.org/wiki/Quadruple-precision_floating-point_format
	QuadruplePrecision uint = 128

	// OctuplePrecision describes 256 bits of precision for for IEEE 754 decimals
	// see: https://en.wikipedia.org/wiki/Octuple-precision_floating-point_format
	OctuplePrecision uint = 256
)

var (
	// Minimum rate of tax must be 0.
	// As far as we know, there are as of writing, no markets with a negative sales tax.
	min = itof(0)

	// Baseline for creating a rate divisor.
	base = itof(1)
)

// ValidateManyFloatsArePrecise tests that the given float64 arguments
// have the desired precision, this is a convenience wrapper
// around ValidateFloatIsPrecise.
// NOTE: 2 digits past the dot is a business rule.
func ValidateManyFloatsArePrecise(args ...float64) error {
	for _, f := range args {
		if err := ValidateFloatIsPrecise(f); err != nil {
			return err
		}
	}
	return nil
}

// ValidateFloatIsPrecise ensures that a float64 value does not exceed
// a precision of 2 digits past the dot. This ensures we do not
// store incorrect currency data.
// NOTE: 2 digits past the dot is a business rule.
func ValidateFloatIsPrecise(f float64) error {
	// parse the float, with the smallest number of digits necessary
	parsed := strconv.FormatFloat(f, 'f', -1, 64)

	// split the parsed number on the dot
	parts := strings.Split(parsed, ".")

	// in case of exact number...
	if len(parts) == 1 {
		return nil
	}

	// check the float's precision
	if prec := len([]rune(parts[1])); prec > 2 {
		return ErrFloatPrecision{
			Value:     parsed,
			Precision: prec,
		}
	}
	return nil
}

// ToMinorUnit returns the currency data in it's minor unit, int64 format
func ToMinorUnit(c currency.Currency, value float64) int64 {
	return int64(math.Round(value * float64(c.Factor())))
}

// FromMinorUnit returns the currency data as a floating point from it's
// minor currency unit format.
func FromMinorUnit(c currency.Currency, value int64) float64 {
	// fast path for currencies like JPY with a factor of 1.
	if c.Factor() == 1 {
		return float64(value)
	}
	var (
		v = itof(value)
		f = ftof(c.FactorAsFloat64())
	)
	f64, _ := newf().Quo(v, f).Float64()
	return f64
}

// Exchange - Apply currency exchange rates to an amount.
//
// value - should always be given in the minor currency unit.
// exchange - should always be given from the approved finance list.
//
// Rounding to the nearest even is a defined business rule.
// Tills may round up to the nearest penny, but for reporting, the rule is
// always to use banker's rounding.
//
// If unclear, see: // http://wiki.c2.com/?BankersRounding.
func Exchange(c currency.Currency, value, exchange float64) float64 {
	// this can happen for example when dealing
	// with a GBP->GBP exchange, for example.
	if exchange == 0 {
		exchange = 1
	}

	var (
		v = ftof(value)
		f = ftof(c.FactorAsFloat64())
		e = ftof(exchange)
	)
	// Here we divide the value, by it's minor currency
	// unit factor, then divide it once more by the
	// exchange rate.
	// -> v / f / e
	f64, _ := v.Quo(v, f).Quo(v, e).Float64()

	// basic banker's rounding
	// http://wiki.c2.com/?BankersRounding
	return math.RoundToEven(f64*100) / 100
}

// RatNetAmount applies a VAT rate to a big.Rat value. This method returns a big.Float
// so it's accuracy can be checked, and it's value applied with .Rat(some.field)
func RatNetAmount(gross, rate *big.Rat) (*big.Float, error) {
	// Here we go for octuple precision as we are dealing with rational numbers.
	bf := func(rat *big.Rat) *big.Float {
		return big.NewFloat(0).SetRat(rat).SetPrec(OctuplePrecision).SetMode(big.ToNearestEven)
	}

	v := bf(gross)
	r := bf(rate)

	// Guard against impossible (negative) tax rates.
	switch r.Cmp(min) {
	case -1:
		return min, ErrSubZeroRate
	case 0:
		return v, nil
	}
	// Turn the rate into a divisor by making it superior to 1.
	divisor := newf().Add(base, r)

	// Here we divide the gross by it's vat:
	// -> val / vat
	// where vat is a gross superior to 1.
	return newf().Quo(v, divisor), nil
}

// NetAmount derives the net amount before tax is applied using the given rate.
func NetAmount(gross int64, rate float64) (int64, error) {
	// Coerce the amount and rate into big floats to perform accurate calculations.
	g := itof(gross)
	r := ftof(rate)
	// Guard against impossible (negative) tax rates.
	switch r.Cmp(min) {
	case -1:
		return 0, ErrSubZeroRate
	case 0:
		return gross, nil
	}
	// Turn the rate into a divisor by making it superior to 1.
	divisor := newf().Add(base, r)
	// The net amount must be:
	// amount / (rate + 1)
	netFloat := newf().Quo(g, divisor)
	// To avoid integer rounding errors we pass the float into a string first then cast to an integer.
	netStr := netFloat.Text('f', 0)
	net, _ := new(big.Int).SetString(netStr, 10)
	// Derive the integer value of the net amount
	return net.Int64(), nil
}

// TaxAmount returns the difference between the gross and the net amounts.
func TaxAmount(gross, net int64) (int64, error) {
	// Guard against values that are not allowed in this context.
	if gross < 0 {
		return 0, ErrSubZeroGross
	}
	if net < 0 {
		return 0, ErrSubZeroNet
	}
	if net > gross {
		return 0, ErrNetOverGrossAmount
	}
	// list amount - net amount = tax amount
	return gross - net, nil
}
