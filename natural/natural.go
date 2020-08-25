package natural

import (
	"math"

	"github.com/dnesting/unit"
)

var (
	// Planck Units

	// C is the speed of light in a vacuum.
	C = unit.Primitive("c") // 299792458 m/s
	// Hr is the reduced Planck constant.
	Hr = unit.Derive("ħ", H(float64(1)/2*math.Pi))
	// G is the gravitational constant.
	G = unit.Primitive("G")
	// KB is the Boltzmann constant.
	KB = unit.Primitive("k_B") // 1.380649e-23 J/K

	// Caesium is one period of the radiation corresponding to the
	// transition between the two hyperfine levels of the ground
	// state of the caesium 133 atom.
	Caesium = unit.Primitive("∆νCs") // 9192631770 Hz

	// H is the Planck constant.
	H = unit.Primitive("h") // 6.62607015e-34 kg m2 s−1

	// E is the elementary charge.
	E = unit.Primitive("e") // 1.602176634e-19 A s

	// Kcd is the luminous efficacy of monochromatic radiation of frequency 540e12 Hz.
	Kcd = unit.Primitive("Kcd") // 683 cd sr/W
)
