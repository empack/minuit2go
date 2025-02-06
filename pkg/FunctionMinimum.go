package minuit

type FunctionMinimum struct {
	theSeed             *MinimumSeed
	theStates           []*MinimumState
	theErrorDef         float64
	theAboveMaxEdm      bool
	theReachedCallLimit bool
	theUserState        *MnUserParameterState
}

func NewFunctionMinimumWithSeedUp(seed *MinimumSeed, up float64) *FunctionMinimum {
	localTheStates := make([]*MinimumState, 0)
	localTheStates = append(localTheStates, NewMinimumStateWithGrad(seed.parameters(), seed.error(), seed.gradient(),
		seed.parameters().fval(), seed.nfcn()))
	return &FunctionMinimum{
		theSeed:             seed,
		theStates:           localTheStates,
		theErrorDef:         up,
		theAboveMaxEdm:      false,
		theReachedCallLimit: false,
		theUserState:        NewMnUserParameterState(),
	}
}

func NewFunctionMinimumWithSeedStatesUp(seed *MinimumSeed, states []*MinimumState, up float64) *FunctionMinimum {
	return &FunctionMinimum{
		theSeed:             seed,
		theStates:           states,
		theErrorDef:         up,
		theAboveMaxEdm:      false,
		theReachedCallLimit: false,
		theUserState:        NewMnUserParameterState(),
	}
}

func NewFunctionMinimumWithSeedStatesUpReachedCallLimit(seed *MinimumSeed, states []*MinimumState,
	up float64) *FunctionMinimum {
	return &FunctionMinimum{
		theSeed:             seed,
		theStates:           states,
		theErrorDef:         up,
		theAboveMaxEdm:      false,
		theReachedCallLimit: true,
		theUserState:        NewMnUserParameterState(),
	}
}

func NewFunctionMinimumWithSeedStatesUpAboveMaxEdm(seed *MinimumSeed, states []*MinimumState,
	up float64) *FunctionMinimum {
	return &FunctionMinimum{
		theSeed:             seed,
		theStates:           states,
		theErrorDef:         up,
		theAboveMaxEdm:      true,
		theReachedCallLimit: false,
		theUserState:        NewMnUserParameterState(),
	}
}

func (this *FunctionMinimum) add(state *MinimumState) {
	this.theStates = append(this.theStates, state)
}

func (this *FunctionMinimum) seed() *MinimumSeed {
	return this.theSeed
}

func (this *FunctionMinimum) states() []*MinimumState {
	return this.theStates
}

func (this *FunctionMinimum) UserState() *MnUserParameterState {
	if !this.theUserState.IsValid() {
		this.theUserState = NewUserParameterStateWith(this.state(), this.ErrorDef(),
			this.seed().trafo())
	}

	return this.theUserState
}

func (this *FunctionMinimum) UserParameters() *MnUserParameters {
	if !this.theUserState.IsValid() {
		this.theUserState = NewMnUserParameterStateWith(this.state(), this.ErrorDef(), this.seed().Trafo())
	}

	return this.theUserState.Parameters()
}

func (this *FunctionMinimum) UserCovariance() *MnUserCovariance {
	if !this.theUserState.IsValid() {
		this.theUserState = NewMnUserParameterStateWith(this.state(), this.ErrorDef(), this.seed().Trafo())
	}

	return this.theUserState.Covariance()
}

func (this *FunctionMinimum) lastState() *MinimumState {
	return this.theStates[len(this.theStates)-1]
}

func (this *FunctionMinimum) state() *MinimumState {
	return this.lastState()
}

func (this *FunctionMinimum) parameters() *MinimumParameters {
	return this.lastState().parameters()
}

func (this *FunctionMinimum) error() *MinimumError {
	return this.lastState().error()
}

func (this *FunctionMinimum) grad() *FunctionGradient {
	return this.lastState().gradient()
}

func (this *FunctionMinimum) Fval() float64 {
	return this.lastState().fval()
}

func (this *FunctionMinimum) Edm() float64 {
	return this.lastState().edm()
}

func (this *FunctionMinimum) Nfcn() int {
	return this.lastState().nfcn()
}

func (this *FunctionMinimum) ErrorDef() float64 {
	return this.theErrorDef
}

func (this *FunctionMinimum) IsValid() bool {
	return this.state().isValid() && !this.isAboveMaxEdm() && !this.hasReachedCallLimit()
}

func (this *FunctionMinimum) hasValidParameters() bool {
	return this.state().parameters().isValid()
}

func (this *FunctionMinimum) hasValidCovariance() bool {
	return this.state().error().isValid()
}

func (this *FunctionMinimum) hasAccurateCovar() bool {
	return this.state().error().isAccurate()
}

func (this *FunctionMinimum) hasPosDefCovar() bool {
	return this.state().error().isPosDef()
}

func (this *FunctionMinimum) hasMadePosDefCovar() bool {
	return this.state().error().isMadePosDef()
}

func (this *FunctionMinimum) hesseFailed() bool {
	return this.state().error().hesseFailed()
}

func (this *FunctionMinimum) hasCovariance() bool {
	return this.state().error().isAvailable()
}

func (this *FunctionMinimum) isAboveMaxEdm() bool {
	return this.theAboveMaxEdm
}

func (this *FunctionMinimum) hasReachedCallLimit() bool {
	return this.theReachedCallLimit
}
