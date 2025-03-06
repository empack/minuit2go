package minuit

import "context"

type MnScan struct {
	theMinimizer *ScanMinimizer
	baseImpl     *MnApplication
}

func (this *MnScan) Minimize(ctx context.Context) (*FunctionMinimum, error) {
	return this.baseImpl.Minimize(ctx)
}

func (this *MnScan) MinimizeWithMaxfcn(ctx context.Context, maxfcn int) (*FunctionMinimum, error) {
	return this.baseImpl.MinimizeWithMaxfcn(ctx, maxfcn)
}

func (this *MnScan) MinimizeWithMaxfcnToler(ctx context.Context, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.MinimizeWithMaxfcnToler(ctx, maxfcn, toler)
}

func (this *MnScan) Precision() *MnMachinePrecision {
	return this.baseImpl.Precision()
}

func (this *MnScan) State() *MnUserParameterState {
	return this.baseImpl.State()
}

func (this *MnScan) Parameters() *MnUserParameters {
	return this.baseImpl.Parameters()
}

func (this *MnScan) Covariance() *MnUserCovariance {
	return this.baseImpl.Covariance()
}

func (this *MnScan) Fcnbase() FCNBase {
	return this.baseImpl.Fcnbase()
}

func (this *MnScan) Strategy() *MnStrategy {
	return this.baseImpl.Strategy()
}

func (this *MnScan) NumOfCalls() int {
	return this.baseImpl.NumOfCalls()
}

func (this *MnScan) minuitParameters() []*MinuitParameter {
	return this.baseImpl.minuitParameters()
}

func (this *MnScan) Params() []float64 {
	return this.baseImpl.Params()
}

func (this *MnScan) Errors() []float64 {
	return this.baseImpl.Errors()
}

func (this *MnScan) parameter(i int) *MinuitParameter {
	return this.baseImpl.parameter(i)
}

func (this *MnScan) AddWithErr(name string, val, err float64) {
	this.baseImpl.AddWithErr(name, val, err)
}

func (this *MnScan) AddWithErrLowUp(name string, val, err, low, up float64) {
	this.baseImpl.AddWithErrLowUp(name, val, err, low, up)
}

func (this *MnScan) Add(name string, val float64) {
	this.baseImpl.Add(name, val)
}

func (this *MnScan) Fix(index int) {
	this.baseImpl.Fix(index)
}

func (this *MnScan) Release(index int) {
	this.baseImpl.Release(index)
}

func (this *MnScan) SetValue(index int, val float64) {
	this.baseImpl.SetValue(index, val)
}

func (this *MnScan) SetError(index int, err float64) {
	this.baseImpl.SetError(index, err)
}

func (this *MnScan) SetLimits(index int, low, up float64) {
	this.baseImpl.SetLimits(index, low, up)
}

func (this *MnScan) RemoveLimits(index int) {
	this.baseImpl.RemoveLimits(index)
}

func (this *MnScan) Value(index int) float64 {
	return this.baseImpl.Value(index)
}

func (this *MnScan) Error(index int) float64 {
	return this.baseImpl.Error(index)
}

func (this *MnScan) FixWithName(name string) {
	this.baseImpl.FixWithName(name)
}

func (this *MnScan) ReleaseWithName(name string) {
	this.baseImpl.ReleaseWithName(name)
}

func (this *MnScan) SetValueWithName(name string, val float64) {
	this.baseImpl.SetValueWithName(name, val)
}

func (this *MnScan) SetErrorWithName(name string, err float64) {
	this.baseImpl.SetErrorWithName(name, err)
}

func (this *MnScan) SetLimitsWithName(name string, low, up float64) {
	this.baseImpl.SetLimitsWithName(name, low, up)
}

func (this *MnScan) RemoveLimitsWithName(name string) {
	this.baseImpl.RemoveLimitsWithName(name)
}

func (this *MnScan) SetPrecision(prec float64) {
	this.baseImpl.SetPrecision(prec)
}

func (this *MnScan) ValueWithName(name string) float64 {
	return this.baseImpl.ValueWithName(name)
}

func (this *MnScan) ErrorWithName(name string) float64 {
	return this.baseImpl.ErrorWithName(name)
}

func (this *MnScan) Index(name string) int {
	return this.baseImpl.Index(name)
}

func (this *MnScan) Name(index int) string {
	return this.baseImpl.Name(index)
}

func (this *MnScan) int2ext(i int, value float64) float64 {
	return this.baseImpl.int2ext(i, value)
}

func (this *MnScan) ext2int(i int, value float64) float64 {
	return this.baseImpl.ext2int(i, value)
}

func (this *MnScan) intOfExt(i int) (int, error) {
	return this.baseImpl.intOfExt(i)
}

func (this *MnScan) extOfInt(i int) int {
	return this.baseImpl.extOfInt(i)
}

func (this *MnScan) VariableParameters() int {
	return this.baseImpl.VariableParameters()
}

func (this *MnScan) SetUseAnalyticalDerivatives(use bool) {
	this.baseImpl.SetUseAnalyticalDerivatives(use)
}

func (this *MnScan) UseAnalyticalDerivatives() bool {
	return this.baseImpl.UseAnalyticalDerivatives()
}

func (this *MnScan) SetCheckAnalyticalDerivatives(check bool) {
	this.baseImpl.SetCheckAnalyticalDerivatives(check)
}

func (this *MnScan) CheckAnalyticalDerivatives() bool {
	return this.baseImpl.CheckAnalyticalDerivatives()
}

func (this *MnScan) SetErrorDef(errorDef float64) {
	this.baseImpl.SetErrorDef(errorDef)
}

func (this *MnScan) ErrorDef() float64 {
	return this.baseImpl.ErrorDef()
}

func NewMnScan(ctx context.Context, fcn FCNBase, par []float64, err []float64) *MnScan {
	return NewMnScanWithStrategy(ctx, fcn, par, err, DEFAULT_STRATEGY)
}

func NewMnScanWithStrategy(ctx context.Context, fcn FCNBase, par []float64, err []float64, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(ctx, fcn, NewUserParamStateFromParamAndErrValues(par, err), NewMnStrategyWithStra(stra))
}

func NewMnScanWithCovariance(ctx context.Context, fcn FCNBase, par []float64, cov *MnUserCovariance) (*MnScan, error) {
	return NewMnScanWithCovarianceAndStrategy(ctx, fcn, par, cov, DEFAULT_STRATEGY)
}

func NewMnScanWithCovarianceAndStrategy(ctx context.Context, fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) (*MnScan, error) {
	state, fnErr := NewMnUserParameterStateFlUc(par, cov)
	if fnErr != nil {
		return nil, fnErr
	}
	return NewMnScanWithStateAndStrategy(ctx, fcn, state, NewMnStrategyWithStra(stra)), nil
}

func NewMnScanWithParameters(ctx context.Context, fcn FCNBase, par *MnUserParameters) *MnScan {
	return NewMnScanWithParametersAndStrategy(ctx, fcn, par, DEFAULT_STRATEGY)
}

func NewMnScanWithParametersAndStrategy(ctx context.Context, fcn FCNBase, par *MnUserParameters, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(ctx, fcn, NewUserParameterStateFromUserParameter(par), NewMnStrategyWithStra(stra))
}

func NewMnScanWithParametersAndCovariance(ctx context.Context, fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) (*MnScan, error) {
	return NewMnScanWithParametersAndCovarianceAndStrategy(ctx, fcn, par, cov, DEFAULT_STRATEGY)
}

func NewMnScanWithParametersAndCovarianceAndStrategy(ctx context.Context, fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance,
	stra int) (*MnScan, error) {
	state, fnErr := NewUserParamStateFromUserParamCovariance(par, cov)
	if fnErr != nil {
		return nil, fnErr
	}
	return NewMnScanWithStateAndStrategy(ctx, fcn, state, NewMnStrategyWithStra(stra)), nil
}

func NewMnScanWithStateAndStrategy(ctx context.Context, fcn FCNBase, state *MnUserParameterState, str *MnStrategy) *MnScan {
	ret := &MnScan{
		baseImpl:     NewMnApplicationWithFcnStateStra(fcn, state, str),
		theMinimizer: NewScanMinimizer(),
	}
	ret.baseImpl.super = ret
	return ret
}

func (this *MnScan) Minimizer() ModularFunctionMinimizerInterface {
	return this.theMinimizer
}

func (this *MnScan) Scan(par int) ([]*Point, error) {
	return this.ScanWithMaxsteps(par, 41)
}

func (this *MnScan) ScanWithMaxsteps(par, maxsteps int) ([]*Point, error) {
	return this.ScanWithMaxstepsRange(par, maxsteps, 0.0, 0.0)
}

func (this *MnScan) ScanWithMaxstepsRange(par, maxsteps int, low, high float64) ([]*Point, error) {
	scan := NewMnParameterScan(this.baseImpl.theFCN, this.baseImpl.theState.parameters())
	amin := scan.fval()
	result := scan.scanWithMaxStepsLowHigh(par, maxsteps, low, high)
	if scan.fval() < amin {
		fnError := this.baseImpl.theState.SetValue(par, scan.parameters().Value(par))
		if fnError != nil {
			return nil, fnError
		}
		amin = scan.fval()
	}

	return result, nil
}
