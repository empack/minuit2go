package minuit

import "errors"

var (
	DEFAULT_STRATEGY = 1
	DEFAULT_MAXFCN   = 0
	DEFAULT_TOLER    = 0.1
)

type MnApplication struct {
	useAnalyticalDerivatives   bool
	checkAnalyticalDerivatives bool
	theFCN                     FCNBase
	theState                   *MnUserParameterState
	theStrategy                *MnStrategy
	theNumCall                 int
	theErrorDef                float64
}

func NewMnApplicationWithFcnStateStra(fcn FCNBase, state *MnUserParameterState, stra *MnStrategy) *MnApplication {
	return &MnApplication{
		useAnalyticalDerivatives:   true,
		checkAnalyticalDerivatives: true,
		theFCN:                     fcn,
		theState:                   state,
		theStrategy:                stra,
		theNumCall:                 0,
		theErrorDef:                1.0,
	}
}

func NewMnApplicationWithFcnStateStraNfcn(fcn FCNBase, state *MnUserParameterState, stra *MnStrategy,
	nfcn int) *MnApplication {
	return &MnApplication{
		useAnalyticalDerivatives:   true,
		checkAnalyticalDerivatives: true,
		theFCN:                     fcn,
		theState:                   state,
		theStrategy:                stra,
		theNumCall:                 0,
		theErrorDef:                1.0,
	}
}

func (this *MnApplication) Minimize() (*FunctionMinimum, error) {
	return this.MinimizeWithMaxfcn(DEFAULT_MAXFCN)
}

func (this *MnApplication) MinimizeWithMaxfcn(maxfcn int) (*FunctionMinimum, error) {
	return this.MinimizeWithMaxfcnToler(maxfcn, DEFAULT_TOLER)
}

func (this *MnApplication) MinimizeWithMaxfcnToler(maxfcn int, toler float64) (*FunctionMinimum, error) {
	if !this.theState.IsValid() {
		return nil, errors.New("invalid state")
	} else {
		npar := this.VariableParameters()
		if maxfcn == 0 {
			maxfcn = 200 + 100*npar + 5*npar*npar
		}

		min, _ := this.Minimizer().minimize(this.theFCN, this.theState, this.theStrategy, maxfcn, toler,
			this.theErrorDef, this.useAnalyticalDerivatives, this.checkAnalyticalDerivatives)
		this.theNumCall += min.Nfcn()
		this.theState = min.UserState()
		return min, nil
	}
}

func (this *MnApplication) Minimizer() *ModularFunctionMinimizer {
	// TODO: MnApplication.class:63 this is an abstract function so..
	return nil
}

func (this *MnApplication) Precision() *MnMachinePrecision {
	return this.theState.Precision()
}

func (this *MnApplication) State() *MnUserParameterState {
	return this.theState
}

func (this *MnApplication) Parameters() *MnUserParameters {
	return this.theState.parameters()
}

func (this *MnApplication) Covariance() *MnUserCovariance {
	return this.theState.covariance()
}

func (this *MnApplication) Fcnbase() FCNBase {
	return this.theFCN
}

func (this *MnApplication) Strategy() *MnStrategy {
	return this.theStrategy
}

func (this *MnApplication) NumOfCalls() int {
	return this.theNumCall
}

func (this *MnApplication) minuitParameters() []*MinuitParameter {
	return this.theState.MinuitParameters()
}

func (this *MnApplication) Params() []float64 {
	return this.theState.Params()
}

func (this *MnApplication) Errors() []float64 {
	return this.theState.errors()
}

func (this *MnApplication) parameter(i int) *MinuitParameter {
	return this.theState.parameter(i)
}

func (this *MnApplication) AddWithErr(name string, val, err float64) {
	this.theState.AddStFlFl(name, val, err)
}

func (this *MnApplication) AddWithErrLowUp(name string, val, err, low, up float64) {
	this.theState.AddStFlFlFlFl(name, val, err, low, up)
}

func (this *MnApplication) Add(name string, val float64) {
	this.theState.AddStFl(name, val)
}

func (this *MnApplication) Fix(index int) {
	this.theState.Fix(index)
}

func (this *MnApplication) Release(index int) {
	this.theState.Release(index)
}

func (this *MnApplication) SetValue(index int, val float64) {
	this.theState.SetValue(index, val)
}

func (this *MnApplication) SetError(index int, err float64) {
	this.theState.SetError(index, err)
}

func (this *MnApplication) SetLimits(index int, low, up float64) {
	this.theState.SetLimits(index, low, up)
}

func (this *MnApplication) RemoveLimits(index int) {
	this.theState.RemoveLimits(index)
}

func (this *MnApplication) Value(index int) float64 {
	return this.theState.Value(index)
}

func (this *MnApplication) Error(index int) float64 {
	return this.theState.Error(index)
}

func (this *MnApplication) FixWithName(name string) {
	this.theState.FixSt(name)
}

func (this *MnApplication) ReleaseWithName(name string) {
	this.theState.ReleaseSt(name)
}

func (this *MnApplication) SetValueWithName(name string, val float64) {
	this.theState.SetValueStFl(name, val)
}

func (this *MnApplication) SetErrorWithName(name string, err float64) {
	this.theState.SetErrorStFl(name, err)
}

func (this *MnApplication) SetLimitsWithName(name string, low, up float64) {
	this.theState.SetLimitsStFlFl(name, low, up)
}

func (this *MnApplication) RemoveLimitsWithName(name string) {
	this.theState.RemoveLimitsSt(name)
}

func (this *MnApplication) SetPrecision(prec float64) {
	this.theState.SetPrecision(prec)
}

func (this *MnApplication) ValueWithName(name string) float64 {
	return this.theState.ValueSt(name)
}

func (this *MnApplication) ErrorWithName(name string) float64 {
	return this.theState.ErrorSt(name)
}

func (this *MnApplication) Index(name string) int {
	return this.theState.Index(name)
}

func (this *MnApplication) Name(index int) string {
	return this.theState.Name(index)
}

func (this *MnApplication) int2ext(i int, value float64) float64 {
	return this.theState.int2ext(i, value)
}

func (this *MnApplication) ext2int(i int, value float64) float64 {
	return this.theState.ext2int(i, value)
}

func (this *MnApplication) intOfExt(i int) int {
	return this.theState.intOfExt(i)
}

func (this *MnApplication) extOfInt(i int) int {
	return this.theState.ExtOfInt(i)
}

func (this *MnApplication) VariableParameters() int {
	return this.theState.VariableParameters()
}

func (this *MnApplication) SetUseAnalyticalDerivatives(use bool) {
	this.useAnalyticalDerivatives = use
}

func (this *MnApplication) UseAnalyticalDerivatives() bool {
	return this.useAnalyticalDerivatives
}

func (this *MnApplication) SetCheckAnalyticalDerivatives(check bool) {
	this.checkAnalyticalDerivatives = check
}

func (this *MnApplication) CheckAnalyticalDerivatives() bool {
	return this.checkAnalyticalDerivatives
}

func (this *MnApplication) SetErrorDef(errorDef float64) {
	this.theErrorDef = errorDef
}

func (this *MnApplication) ErrorDef() float64 {
	return this.theErrorDef
}
