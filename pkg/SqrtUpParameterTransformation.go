package minuit

import "math"

type SqrtUpParameterTransformation struct{}

func NewSqrtUpParameterTransformation() *SqrtUpParameterTransformation {
	return &SqrtUpParameterTransformation{}
}

// transformation from internal to external
func (this *SqrtUpParameterTransformation) int2ext(value, upper float64) float64 {
	return upper + 1. - math.Sqrt(value*value+1.)
}

// transformation from external to internal
func (this *SqrtUpParameterTransformation) ext2int(value, upper float64, prec *MnMachinePrecision) float64 {
	var yy float64 = upper - value + 1.
	var yy2 float64 = yy * yy
	if yy2 < (1. + prec.eps2()) {
		return 8 * math.Sqrt(prec.eps2())
	} else {
		return math.Sqrt(yy2 - 1)
	}
}

// derivative of transformation from internal to external
func (this *SqrtUpParameterTransformation) dInt2Ext(value, upper float64) float64 {
	return -value / (math.Sqrt(value*value + 1.))
}
