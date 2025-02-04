package minuit

import "math"

type SqrtLowParameterTransformation struct{}

func NewSqrtLowParameterTransformation() *SqrtLowParameterTransformation {
	return &SqrtLowParameterTransformation{}
}

// transformation from internal to external
func (this *SqrtLowParameterTransformation) int2ext(value, lower float64) float64 {
	return lower - 1. + math.Sqrt(value*value+1.)
}

// transformation from external to internal
func (this *SqrtLowParameterTransformation) ext2int(value, lower float64, prec *MnMachinePrecision) float64 {
	var yy float64 = value - lower + 1.0
	var yy2 float64 = yy * yy
	if yy2 < (1. + prec.eps2()) {
		return 8 * math.Sqrt(prec.eps2())
	} else {
		return math.Sqrt(yy2 - 1)
	}
}

// derivative of transformation from internal to external
func (this *SqrtLowParameterTransformation) dInt2Ext(value, lower float64) float64 {
	return value / (math.Sqrt(value*value + 1.))
}
