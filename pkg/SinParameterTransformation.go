package minuit

import "math"

type SinParameterTransformation struct{}

func (this *SinParameterTransformation) int2ext(value, upper, lower float64) float64 {
	return lower + 0.5*(upper-lower)*(math.Sin(value)+1.0)
}

func (this *SinParameterTransformation) ext2int(value, upper, lower float64, prec *MnMachinePrecision) float64 {
	var piby2 float64 = 2. * math.Atan(1.)
	var distnn float64 = 8. * math.Sqrt(prec.eps2())
	var vlimhi float64 = piby2 - distnn
	var vlimlo float64 = -piby2 + distnn

	var yy float64 = 2.*(value-lower)/(upper-lower) - 1.
	var yy2 float64 = yy * yy
	if yy2 > (1. - prec.eps2()) {
		if yy < 0.0 {
			return vlimlo
		} else {
			return vlimhi
		}
	} else {
		return math.Asin(yy)
	}
}
func (this *SinParameterTransformation) dInt2Ext(value, upper, lower float64) float64 {
	return 0.5 * math.Abs((upper-lower)*math.Cos(value))
}
