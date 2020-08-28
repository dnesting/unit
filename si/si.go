// Package si defines common SI and Metric units in terms of natural fundamental
// constants from the natural package.
package si

import (
	"time"

	"github.com/dnesting/unit"
	"github.com/dnesting/unit/natural"
)

var Registry unit.Registry

var (
	Yotta = Registry.Prefix("Y", 1e24)
	Zetta = Registry.Prefix("Z", 1e21)
	Exa   = Registry.Prefix("E", 1e18)
	Peta  = Registry.Prefix("P", 1e15)
	Tera  = Registry.Prefix("T", 1e12)
	Giga  = Registry.Prefix("G", 1e9)
	Mega  = Registry.Prefix("M", 1e6)
	Kilo  = Registry.Prefix("k", 1e3)
	Hecto = Registry.Prefix("h", 1e2)
	Deka  = Registry.Prefix("da", 1e1)

	Deci  = Registry.Prefix("d", 1e-1)
	Centi = Registry.Prefix("c", 1e-2)
	Milli = Registry.Prefix("m", 1e-3)
	Micro = Registry.Prefix("µ", 1e-6, "u")
	Nano  = Registry.Prefix("n", 1e-9)
	Pico  = Registry.Prefix("p", 1e-12)
	Femto = Registry.Prefix("f", 1e-15)
	Atto  = Registry.Prefix("a", 1e-18)
	Zepto = Registry.Prefix("z", 1e-21)
	Yocto = Registry.Prefix("y", 1e-24)

	// Hz is defined in terms of the radiation produced by the
	// transition between the two hyperfine ground states of caesium,
	// which has a frequency, ΔνCs (natural.Caesium), of exactly
	// 9192631770 Hz.
	Hertz   = Registry.Derive("Hz", natural.Caesium(1.0/9192631770))
	Second  = Registry.Derive("s", Hertz(1).Recip())
	Caesium = unit.MustConvert(natural.Caesium, Hertz)

	// Metre is defined to be 1/299 792 458 of the distance the
	// speed of light (natural.C) travels in 1 second.
	Metre = Registry.Derive("m", natural.C.Mul(Second)(1.0/299792458))
	Meter = Metre
	C     = unit.MustConvert(natural.C, Metre.Div(Second))

	// Gram is derived from the Planck constant (h), which is defined
	// to be 6.62607015×10−34 kg m^2/s.
	Gram     = Registry.Derive("g", natural.H.Mul(Second).Div(Metre.Pow(2))(1.0/6.62607015e-38))
	Kilogram = Kilo(Gram)
	Newton   = Registry.Derive("N", Kilogram.Mul(Metre).Div(Second.Pow(2)))
	Joule    = Registry.Derive("J", Metre.Mul(Newton))
	H        = unit.MustConvert(natural.H, Joule.Mul(Second))

	// Ampere is derived from the elementary charge e, defined to be 1.602176634×10^−19 A s.
	Ampere  = Registry.Derive("A", natural.E.Div(Second)(1.0/1.602176634e-19))
	Coulomb = Registry.Derive("C", Ampere.Mul(Second))
	E       = unit.MustConvert(natural.E, Coulomb)

	// Kelvin is derived from the Boltzmann constant k, defined to be 1.380649×10^−23 J/K.
	Kelvin = Registry.Derive("K", Joule.Div(natural.KB)(1.380649e-23))
	K      = unit.MustConvert(natural.KB, Joule.Div(Kelvin))

	Mole = Registry.Primitive("mol")
	NA   = unit.Scalar(6.02214076e23).Div(Mole)(1)

	Watt      = Registry.Derive("W", Joule.Div(Second))
	Steradian = Registry.Derive("sr", Metre.Pow(2).Div(Metre.Pow(2)))

	// Candela is derived from the luminous efficacy of monochromatic radiation of frequency
	// 540×10^12 Hz, Kcd, defined to be 683 cd sr/W.
	Candela = Registry.Derive("cd", natural.Kcd.Mul(Watt).Div(Steradian)(1.0/683))
	Kcd     = unit.MustConvert(natural.Kcd, Candela.Mul(Steradian).Div(Watt))

	Becquerel = Registry.Derive("Bq", unit.Scalar(1).Div(Second))
	Farad     = Registry.Derive("F", Coulomb.Div(Volt))
	Gray      = Registry.Derive("Gy", Joule.Div(Kilogram))
	Henry     = Registry.Derive("H", Volt.Mul(Second).Div(Ampere))
	Katal     = Registry.Derive("kat", Mole.Div(Second))
	Liter     = Registry.Derive("L", Centi(Metre).Pow(3)(1000))
	Lumen     = Registry.Derive("lm", Candela.Mul(Steradian))
	Lux       = Registry.Derive("lx", Lumen.Div(Metre.Pow(2)))
	Ohm       = Registry.Derive("Ω", Volt.Div(Ampere), "Ohm")
	Pascal    = Registry.Derive("Pa", Newton.Div(Metre.Pow(2)))
	Radian    = Registry.Derive("rad", Metre.Div(Metre))
	Siemens   = Registry.Derive("S", Ampere.Div(Volt))
	Sievert   = Registry.Derive("Sv", Joule.Div(Kilogram))
	Tesla     = Registry.Derive("T", Volt.Mul(Second).Div(Metre.Pow(2)))
	Volt      = Registry.Derive("V", Watt.Div(Ampere))
	Weber     = Registry.Derive("Wb", Joule.Div(Ampere))

	// DegCelsius shares the same underlying unit as Kelvin, but on a
	// different scale.  No conversions are therefore needed for values
	// that represent changes in temperature.  To convert between values
	// on their respective scales, use CToK and KToC.
	DegCelsius = Registry.Derive("°C", Kelvin, "℃", "degC")
)

// FromDuration converts a time.Duration to a unit.Value with unit Second.
func FromDuration(d time.Duration) unit.Value {
	return Second(d.Seconds())
}

// ToDuration converts a unit.Value to a time.Duration.  If v's units
// cannot be converted to Second, returns 0, false.
func ToDuration(v unit.Value) (d time.Duration, ok bool) {
	v, remain := v.Convert(Second)
	if !remain.Empty() {
		return 0, false
	}
	return time.Duration(v.S * float64(time.Second)), true
}

// FromCelsius converts a DegCelsius value on the Celsius scale to its
// corresponding value on the Kelvin scale.  Returns false if v cannot
// be converted to DegCelsius.
func FromCelsius(v unit.Value) (unit.Value, bool) {
	c, remain := v.Convert(DegCelsius)
	if !remain.Empty() {
		return unit.Value{}, false
	}
	return Kelvin(c.S + 273.15), true
}

// ToCelsius converts a Kelvin value on the Kelvin scale to its corresponding
// DegCelsius value on the Celsius scale.  Returns false if v cannot be
// converted to Kelvin.
func ToCelsius(v unit.Value) (unit.Value, bool) {
	k, remain := v.Convert(Kelvin)
	if !remain.Empty() {
		return unit.Value{}, false
	}
	return DegCelsius(k.S - 273.15), true
}
