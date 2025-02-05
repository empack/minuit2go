package minuit

type MnCross struct {
	theValue  float64
	theState  *MnUserParameterState
	theNFcn   int
	theValid  bool
	theLimset bool
	theMaxFcn bool
	theNewMin bool
}

func NewMnCross() *MnCross {
	return &MnCross{
		theValue:  0,
		theState:  NewMnUserParameterState(),
		theNFcn:   0,
		theValid:  false,
		theLimset: false,
		theMaxFcn: false,
		theNewMin: false,
	}
}

func NewMnCrossWithNfcn(nfcn int) *MnCross {
	return &MnCross{
		theValue:  0,
		theState:  NewMnUserParameterState(),
		theNFcn:   nfcn,
		theValid:  false,
		theLimset: false,
		theMaxFcn: false,
		theNewMin: false,
	}
}

func NewMnCrossWithValueStateNfcn(value float64, state *MnUserParameterState, nfcn int) *MnCross {
	return &MnCross{
		theValue:  value,
		theState:  state,
		theNFcn:   nfcn,
		theValid:  false,
		theLimset: false,
		theMaxFcn: false,
		theNewMin: false,
	}
}

func NewMnCrossWithStateNfcnLimset(state *MnUserParameterState, nfcn int) *MnCross {
	return &MnCross{
		theValue:  0,
		theState:  state,
		theNFcn:   nfcn,
		theValid:  false,
		theLimset: true,
		theMaxFcn: false,
		theNewMin: false,
	}
}

func NewMnCrossWithStateNfcnMaxNfcn(state *MnUserParameterState, nfcn int) *MnCross {
	return &MnCross{
		theValue:  0,
		theState:  state,
		theNFcn:   nfcn,
		theValid:  false,
		theLimset: false,
		theMaxFcn: true,
		theNewMin: false,
	}
}

func NewMnCrossWithStateNfcnNewMin(state *MnUserParameterState, nfcn int) *MnCross {
	return &MnCross{
		theValue:  0,
		theState:  state,
		theNFcn:   nfcn,
		theValid:  false,
		theLimset: false,
		theMaxFcn: false,
		theNewMin: true,
	}
}

func (this *MnCross) value() float64 {
	return this.theValue
}

func (this *MnCross) state() *MnUserParameterState {
	return this.theState
}

func (this *MnCross) isValid() bool {
	return this.theValid
}

func (this *MnCross) atLimit() bool {
	return this.theLimset
}

func (this *MnCross) atMaxFcn() bool {
	return this.theMaxFcn
}

func (this *MnCross) newMinimum() bool {
	return this.theNewMin
}

func (this *MnCross) nfcn() int {
	return this.theNFcn
}
