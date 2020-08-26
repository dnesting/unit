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

	Hertz   = Registry.Derive("Hz", natural.Caesium(1).DivN(9192631770))
	Second  = Registry.Derive("s", Hertz(1).Recip())
	Caesium = unit.MustConvert(natural.Caesium, Hertz)

	Metre = Registry.Derive("m", natural.C.Mul(Second)(1).DivN(299792458))
	Meter = Metre
	C     = unit.MustConvert(natural.C, Metre.Div(Second))

	Gram     = Registry.Derive("g", natural.H.Mul(Second).Div(Metre.Pow(2))(float64(1)/6.62607015e-38))
	Kilogram = Kilo(Gram)
	Newton   = Registry.Derive("N", Kilogram.Mul(Metre).Div(Second.Pow(2)))
	Joule    = Registry.Derive("J", Metre.Mul(Newton))
	H        = unit.MustConvert(natural.H, Joule.Mul(Second))

	Ampere  = Registry.Derive("A", natural.E.Div(Second)(float64(1)/1.602176634e-19))
	Coulomb = Registry.Derive("C", Ampere.Mul(Second))
	E       = unit.MustConvert(natural.E, Coulomb)

	Kelvin = Registry.Derive("K", Joule.Div(natural.KB)(1.380649e-23))
	K      = unit.MustConvert(natural.KB, Joule.Div(Kelvin))

	Mole = Registry.Primitive("mol")
	NA   = unit.Scalar(6.02214076e23).Div(Mole)(1)

	Watt      = Registry.Derive("W", Joule.Div(Second))
	Steradian = Registry.Derive("sr", Metre.Pow(2).Div(Metre.Pow(2)))
	Candela   = Registry.Derive("cd", natural.Kcd.Mul(Watt).Div(Steradian)(float64(1)/683))
	Kcd       = unit.MustConvert(natural.Kcd, Candela.Mul(Steradian).Div(Watt))

	Becquerel  = Registry.Derive("Bq", unit.Scalar(1).Div(Second))
	DegCelsius = Registry.Derive("°C", Kelvin, "℃", "degC") // Use CToK or KToC to covert temperature
	Farad      = Registry.Derive("F", Coulomb.Div(Volt))
	Gray       = Registry.Derive("Gy", Joule.Div(Kilogram))
	Henry      = Registry.Derive("H", Volt.Mul(Second).Div(Ampere))
	Katal      = Registry.Derive("kat", Mole.Div(Second))
	Liter      = Registry.Derive("L", Centi(Metre).Pow(3)(1000))
	Lumen      = Registry.Derive("lm", Candela.Mul(Steradian))
	Lux        = Registry.Derive("lx", Lumen.Div(Metre.Pow(2)))
	Ohm        = Registry.Derive("Ω", Volt.Div(Ampere), "Ohm")
	Pascal     = Registry.Derive("Pa", Newton.Div(Metre.Pow(2)))
	Radian     = Registry.Derive("rad", Metre.Div(Metre))
	Siemens    = Registry.Derive("S", Ampere.Div(Volt))
	Sievert    = Registry.Derive("Sv", Joule.Div(Kilogram))
	Tesla      = Registry.Derive("T", Volt.Mul(Second).Div(Metre.Pow(2)))
	Volt       = Registry.Derive("V", Watt.Div(Ampere))
	Weber      = Registry.Derive("Wb", Joule.Div(Ampere))
)

func FromDuration(d time.Duration) unit.Value {
	return Second(d.Seconds())
}

func ToDuration(v unit.Value) (d time.Duration, ok bool) {
	v, remain := v.Convert(Second)
	if !remain.IsEmpty() {
		return 0, false
	}
	return time.Duration(v.S * float64(time.Second)), true
}

func CToK(v unit.Value) (unit.Value, bool) {
	c, remain := v.Convert(DegCelsius)
	if c.U.IsEmpty() {
		return v, false
	}
	return Kelvin(c.S + 273.15).Mul(remain), true
}

func KToC(v unit.Value) (unit.Value, bool) {
	k, remain := v.Convert(Kelvin)
	if k.U.IsEmpty() {
		return v, false
	}
	return DegCelsius(k.S - 273.15).Mul(remain), true
}
