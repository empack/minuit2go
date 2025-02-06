package minuit

type MinimumState struct {
	theParameters *MinimumParameters
	theError      *MinimumError
	theGradient   *FunctionGradient
	theEDM        float64
	theNFcn       int
}

func NewMinimumStateFromNumber(n int) (*MinimumState, error) {
	e, err := NewMinimumErrorFromNumber(n)
	if err != nil {
		return nil, err
	}
	return &MinimumState{
		theParameters: NewMinimumParametersFromNumber(n),
		theError:      e,
		theGradient:   NewFunctionGradientFromNumber(n),
	}, nil
}

func NewMinimumStateWithGrad(states *MinimumParameters, err *MinimumError, grad *FunctionGradient, edm float64, nfcn int) *MinimumState {
	return &MinimumState{
		theParameters: states,
		theError:      err,
		theGradient:   grad,
		theEDM:        edm,
		theNFcn:       nfcn,
	}
}

func NewMinimumState(states *MinimumParameters, edm float64, nfcn int) (*MinimumState, error) {
	e, err := NewMinimumErrorFromNumber(states.vec().size())
	if err != nil {
		return nil, err
	}
	return &MinimumState{
		theParameters: states,
		theError:      e,
		theGradient:   NewFunctionGradientFromNumber(states.vec().size()),
		theEDM:        edm,
		theNFcn:       nfcn,
	}, nil
}

func (this *MinimumState) parameters() *MinimumParameters {
	return this.theParameters
}

func (this *MinimumState) vec() *MnAlgebraicVector {
	return this.theParameters.vec()
}

func (this *MinimumState) size() int {
	return this.theParameters.vec().size()
}

func (this *MinimumState) error() *MinimumError {
	return this.theError
}

func (this *MinimumState) gradient() *FunctionGradient {
	return this.theGradient
}

func (this *MinimumState) fval() float64 {
	return this.theParameters.fval()
}

func (this *MinimumState) edm() float64 {
	return this.theEDM
}

func (this *MinimumState) nfcn() int {
	return this.theNFcn
}

func (this *MinimumState) isValid() bool {
	if this.hasParameters() && this.hasCovariance() {
		return this.parameters().isValid() && this.error().isValid()
	} else if this.hasParameters() {
		return this.parameters().isValid()
	} else {
		return false
	}
}

func (this *MinimumState) hasParameters() bool {
	return this.theParameters.isValid()
}

func (this *MinimumState) hasCovariance() bool {
	return this.theError.isAvailable()
}
