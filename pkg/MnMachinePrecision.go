package minuit

import "math"

// MnMachinePrecision
/*
 * Determines the relative floating point arithmetic precision. The
 * setPrecision() method can be used to override Minuit's own determination,
 * when the user knows that the {FCN} function value is not calculated to
 * the nominal machine accuracy.
 */
type MnMachinePrecision struct {
	theEpsMac float64
	theEpsMa2 float64
}

func NewMnMachinePrecision() *MnMachinePrecision {
	res := &MnMachinePrecision{
		theEpsMac: 4.0e-7,
		theEpsMa2: 2. * math.Sqrt(4.0e-7),
	}

	var epstry float64 = 0.5
	var one float64 = 1.0
	for i := 0; i < 100; i++ {
		epstry *= 0.5
		var epsp1 float64 = one + epstry
		var epsbak float64 = epsp1 - one
		if epsbak < epstry {
			res.SetPrecision(8.0 * epstry)
			break
		}
	}
	return res
}

/** eps returns the smallest possible number so that 1.+eps > 1. */
func (this *MnMachinePrecision) eps() float64 {
	return this.theEpsMac
}

/** eps2 returns 2*sqrt(eps) */
func (this *MnMachinePrecision) eps2() float64 {
	return this.theEpsMa2
}

/** override Minuit's own determination */
func (this *MnMachinePrecision) SetPrecision(prec float64) {
	this.theEpsMac = prec
	this.theEpsMa2 = 2. * math.Sqrt(this.theEpsMac)
}
