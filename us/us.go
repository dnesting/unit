package us

import (
	"math"

	"github.com/dnesting/unit"
	"github.com/dnesting/unit/si"
)

var Registry = unit.NewRegistry("us", nil)
var Survey = unit.NewRegistry("survey", Registry)
var Fluid = unit.NewRegistry("fluid", Registry)
var Deg = unit.NewRegistry("degree", Registry)

var (
	Inch  = Registry.Derive("in", si.Meter(0.0254), "\"", "″")
	Pica  = Registry.Derive("P̸", Inch(1.0/6))
	Point = Registry.Derive("p", Inch(1.0/72))
	Foot  = Registry.Derive("ft", Inch(12), "'", "′")
	Yard  = Registry.Derive("yd", Foot(3))
	Mile  = Registry.Derive("mi", Foot(5280))

	SurveyFoot = Survey.Derive("ft", si.Meter(1200.0/3937))
	SurveyMile = Survey.Derive("mi", si.Meter(6336000.0/3937))

	Fathom       = Registry.Derive("ftm", Yard(2))
	NauticalMile = Registry.Derive("NM", si.Meter(1852), "nmi")

	Acre = Registry.Derive("acre", SurveyFoot.Pow(2)(43560))

	Teaspoon   = Registry.Derive("tsp", si.Liter(4.92892159375/1000))
	Tablespoon = Registry.Derive("Tbsp", Teaspoon(3))
	FluidOunce = Fluid.Derive("oz", Tablespoon(2))
	Shot       = Registry.Derive("jig", Tablespoon(3))
	Cup        = Registry.Derive("cp", FluidOunce(8))
	Pint       = Registry.Derive("pt", Cup(2))
	Quart      = Registry.Derive("qt", Pint(2))
	Gallon     = Registry.Derive("gal", Quart(4))

	Dram  = Registry.Derive("dr", si.Kilogram(0.0017718451953125))
	Ounce = Registry.Derive("oz", Dram(16))
	Pound = Registry.Derive("lb", Ounce(16))
	Ton   = Registry.Derive("ton", Pound(2000))

	DegFahrenheit = Registry.Derive("°F", si.DegCelsius.Mul(unit.Scalar(float64(5)/9)), "℉", "degF")
	Calorie       = Registry.Derive("cal", si.Joule(4.184))
	KiloCalorie   = Registry.Derive("kcal", si.Joule(4184), "Cal")
	PoundForce    = Registry.Derive("lbf", Pound.Mul(si.Meter.Div(si.Second.Pow(2)))(9.80665))

	Second = si.Second
	Minute = Registry.Derive("m", Second(60))
	Hour   = Registry.Derive("h", Minute(60))
	Day    = Registry.Derive("d", Hour(24))

	Degree    = Registry.Derive("°", si.Radian(1).DivN(2*math.Pi).MulN(360), "deg")
	DegMinute = Deg.Derive("'", Degree(1.0/60))
	DegSecond = Deg.Derive("\"", DegMinute(1.0/60))
)

func init() {
	unit.DefaultFormatter.Config(unit.WithNoGapFor(Degree, DegMinute, DegSecond))
}

func ToDMS(deg unit.Value) (d, m, s unit.Value, ok bool) {
	var remain unit.Units
	d, remain = deg.Convert(Degree)
	if !remain.Empty() {
		return
	}
	di := math.Floor(d.S)

	mf := (d.S - di) * 60
	mi := math.Floor(mf)

	sf := (mf - mi) * 60

	d.S = di
	m = DegMinute(mi)
	s = DegSecond(sf)
	ok = true
	return
}

func ToCelsius(degF unit.Value) (c unit.Value, ok bool) {
	fvalue, remain := degF.Convert(DegFahrenheit)
	if !remain.Empty() {
		return
	}
	fvalue.S -= 32
	fvalue, _ = fvalue.Convert(si.DegCelsius)
	return fvalue, true
}

func ToFahrenheit(degC unit.Value) (f unit.Value, ok bool) {
	cvalue, remain := degC.Convert(si.DegCelsius)
	if !remain.Empty() {
		return
	}
	cvalue, _ = cvalue.Convert(DegFahrenheit)
	cvalue.S += 32
	return cvalue, true
}
