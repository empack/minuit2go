package minuit

type FunctionGradient struct {
	theGradient       *MnAlgebraicVector
	theG2ndDerivative *MnAlgebraicVector
	theGStepSize      *MnAlgebraicVector
	theValid          bool
	theAnalytical     bool
}

func NewFunctionGradientFromNumber(n int) *FunctionGradient {
	return &FunctionGradient{
		theGradient:       NewMnAlgebraicVector(n),
		theG2ndDerivative: NewMnAlgebraicVector(n),
		theGStepSize:      NewMnAlgebraicVector(n),
	}
}

func NewFunctionGradient(grd *MnAlgebraicVector) *FunctionGradient {
	return &FunctionGradient{
		theGradient:       grd,
		theG2ndDerivative: NewMnAlgebraicVector(grd.size()),
		theGStepSize:      NewMnAlgebraicVector(grd.size()),
		theValid:          true,
		theAnalytical:     true,
	}
}

func NewFunctionGradientFromMnAlgebraicVectors(grd, g2, gstep *MnAlgebraicVector) *FunctionGradient {
	return &FunctionGradient{
		theGradient:       grd,
		theG2ndDerivative: g2,
		theGStepSize:      gstep,
		theValid:          true,
		theAnalytical:     false,
	}
}

func (this *FunctionGradient) grad() *MnAlgebraicVector {
	return this.theGradient
}
func (this *FunctionGradient) vec() *MnAlgebraicVector {
	return this.theGradient
}
func (this *FunctionGradient) isValid() bool {
	return this.theValid
}

func (this *FunctionGradient) isAnalytical() bool {
	return this.theAnalytical
}
func (this *FunctionGradient) g2() *MnAlgebraicVector {
	return this.theG2ndDerivative
}
func (this *FunctionGradient) gstep() *MnAlgebraicVector {
	return this.theGStepSize
}
