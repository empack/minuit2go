package minuit

type MinimumParameters struct {
	theParameters *MnAlgebraicVector
	theStepSize   *MnAlgebraicVector
	theFVal       float64
	theValid      bool
	theHasStep    bool
}

func NewMinimumParametersByNumber(n int) *MinimumParameters {
	return &MinimumParameters{
		theParameters: NewMnAlgebraicVector(n),
		theStepSize:   NewMnAlgebraicVector(n),
	}
}

func NewMinimumParameters(avec *MnAlgebraicVector, fval float64) *MinimumParameters {
	return &MinimumParameters{
		theParameters: avec,
		theStepSize:   NewMnAlgebraicVector(avec.size()),
		theFVal:       fval,
		theValid:      true,
	}

}

func NewMinimumParametersByMnAlgebraicVectors(avec, dirin *MnAlgebraicVector, fval float64) *MinimumParameters {
	return &MinimumParameters{
		theParameters: avec,
		theStepSize:   dirin,
		theFVal:       fval,
		theValid:      true,
		theHasStep:    true,
	}
}

func (this *MinimumParameters) vec() *MnAlgebraicVector {
	return this.theParameters
}
func (this *MinimumParameters) dirin() *MnAlgebraicVector {
	return this.theStepSize
}
func (this *MinimumParameters) fval() float64 {
	return this.theFVal
}
func (this *MinimumParameters) isValid() bool {
	return this.theValid
}
func (this *MinimumParameters) hasStepSize() bool {
	return this.theHasStep
}
