package minuit

import "context"

// MnMigrad
/*
 * MnMigrad provides minimization of the function by the method of MIGRAD, the most
 * efficient and complete single method, recommended for general functions,
 * and the functionality for parameters interaction. It also retains the result from
 * the last minimization in case the user may want to do subsequent minimization steps
 * with parameter interactions in between the minimization requests.
 * The minimization produces as a by-product the error matrix of the parameters, which
 * is usually reliable unless warning messages are produced.
 */
type MnMigrad struct {
	baseImpl     *MnApplication
	theMinimizer *VariableMetricMinimizer
}

func (this *MnMigrad) Minimize(ctx context.Context) (*FunctionMinimum, error) {
	return this.baseImpl.Minimize(ctx)
}

func (this *MnMigrad) MinimizeWithMaxfcn(ctx context.Context, maxfcn int) (*FunctionMinimum, error) {
	return this.baseImpl.MinimizeWithMaxfcn(ctx, maxfcn)
}

func (this *MnMigrad) MinimizeWithMaxfcnToler(ctx context.Context, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.MinimizeWithMaxfcnToler(ctx, maxfcn, toler)
}

func (this *MnMigrad) Precision() *MnMachinePrecision {
	return this.baseImpl.Precision()
}

func (this *MnMigrad) State() *MnUserParameterState {
	return this.baseImpl.State()
}

func (this *MnMigrad) Parameters() *MnUserParameters {
	return this.baseImpl.Parameters()
}

func (this *MnMigrad) Covariance() *MnUserCovariance {
	return this.baseImpl.Covariance()
}

func (this *MnMigrad) Fcnbase() FCNBase {
	return this.baseImpl.Fcnbase()
}

func (this *MnMigrad) Strategy() *MnStrategy {
	return this.baseImpl.Strategy()
}

func (this *MnMigrad) NumOfCalls() int {
	return this.baseImpl.NumOfCalls()
}

func (this *MnMigrad) minuitParameters() []*MinuitParameter {
	return this.baseImpl.minuitParameters()
}

func (this *MnMigrad) Params() []float64 {
	return this.baseImpl.Params()
}

func (this *MnMigrad) Errors() []float64 {
	return this.baseImpl.Errors()
}

func (this *MnMigrad) parameter(i int) *MinuitParameter {
	return this.baseImpl.parameter(i)
}

func (this *MnMigrad) AddWithErr(name string, val, err float64) {
	this.baseImpl.AddWithErr(name, val, err)
}

func (this *MnMigrad) AddWithErrLowUp(name string, val, err, low, up float64) {
	this.baseImpl.AddWithErrLowUp(name, val, err, low, up)
}

func (this *MnMigrad) Add(name string, val float64) {
	this.baseImpl.Add(name, val)
}

func (this *MnMigrad) Fix(index int) {
	this.baseImpl.Fix(index)
}

func (this *MnMigrad) Release(index int) {
	this.baseImpl.Release(index)
}

func (this *MnMigrad) SetValue(index int, val float64) {
	this.baseImpl.SetValue(index, val)
}

func (this *MnMigrad) SetError(index int, err float64) {
	this.baseImpl.SetError(index, err)
}

func (this *MnMigrad) SetLimits(index int, low, up float64) {
	this.baseImpl.SetLimits(index, low, up)
}

func (this *MnMigrad) RemoveLimits(index int) {
	this.baseImpl.RemoveLimits(index)
}

func (this *MnMigrad) Value(index int) float64 {
	return this.baseImpl.Value(index)
}

func (this *MnMigrad) Error(index int) float64 {
	return this.baseImpl.Error(index)
}

func (this *MnMigrad) FixWithName(name string) {
	this.baseImpl.FixWithName(name)
}

func (this *MnMigrad) ReleaseWithName(name string) {
	this.baseImpl.ReleaseWithName(name)
}

func (this *MnMigrad) SetValueWithName(name string, val float64) {
	this.baseImpl.SetValueWithName(name, val)
}

func (this *MnMigrad) SetErrorWithName(name string, err float64) {
	this.baseImpl.SetErrorWithName(name, err)
}

func (this *MnMigrad) SetLimitsWithName(name string, low, up float64) {
	this.baseImpl.SetLimitsWithName(name, low, up)
}

func (this *MnMigrad) RemoveLimitsWithName(name string) {
	this.baseImpl.RemoveLimitsWithName(name)
}

func (this *MnMigrad) SetPrecision(prec float64) {
	this.baseImpl.SetPrecision(prec)
}

func (this *MnMigrad) ValueWithName(name string) float64 {
	return this.baseImpl.ValueWithName(name)
}

func (this *MnMigrad) ErrorWithName(name string) float64 {
	return this.baseImpl.ErrorWithName(name)
}

func (this *MnMigrad) Index(name string) int {
	return this.baseImpl.Index(name)
}

func (this *MnMigrad) Name(index int) string {
	return this.baseImpl.Name(index)
}

func (this *MnMigrad) int2ext(i int, value float64) float64 {
	return this.baseImpl.int2ext(i, value)
}

func (this *MnMigrad) ext2int(i int, value float64) float64 {
	return this.baseImpl.ext2int(i, value)
}

func (this *MnMigrad) intOfExt(i int) (int, error) {
	return this.baseImpl.intOfExt(i)
}

func (this *MnMigrad) extOfInt(i int) int {
	return this.baseImpl.extOfInt(i)
}

func (this *MnMigrad) VariableParameters() int {
	return this.baseImpl.VariableParameters()
}

func (this *MnMigrad) SetUseAnalyticalDerivatives(use bool) {
	this.baseImpl.SetUseAnalyticalDerivatives(use)
}

func (this *MnMigrad) UseAnalyticalDerivatives() bool {
	return this.baseImpl.UseAnalyticalDerivatives()
}

func (this *MnMigrad) SetCheckAnalyticalDerivatives(check bool) {
	this.baseImpl.SetCheckAnalyticalDerivatives(check)
}

func (this *MnMigrad) CheckAnalyticalDerivatives() bool {
	return this.baseImpl.CheckAnalyticalDerivatives()
}

func (this *MnMigrad) SetErrorDef(errorDef float64) {
	this.baseImpl.SetErrorDef(errorDef)
}

func (this *MnMigrad) ErrorDef() float64 {
	return this.baseImpl.ErrorDef()
}

/** construct from FCNBase + double[] for parameters and errors with default strategy */
func NewMnMigradWithParErr(fcn FCNBase, par, err []float64) *MnMigrad {
	return NewMnMigradWithParErrStrategy(fcn, par, err, DEFAULT_STRATEGY)
}

/** construct from FCNBase + double[] for parameters and errors */
func NewMnMigradWithParErrStrategy(fcn FCNBase, par, err []float64, stra int) *MnMigrad {
	return NewMnMigradWithParameterStateStrategy(fcn, NewUserParamStateFromParamAndErrValues(par, err), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + double[] for parameters and MnUserCovariance with default strategy */
func NewMnMigradWithParCovariance(fcn FCNBase, par []float64, cov *MnUserCovariance) *MnMigrad {
	return NewMnMigradWithParCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/** construct from FCNBase + double[] for parameters and MnUserCovariance */
func NewMnMigradWithParCovarianceStra(fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) *MnMigrad {
	ups, _ := NewMnUserParameterStateFlUc(par, cov)
	return NewMnMigradWithParameterStateStrategy(fcn, ups, NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameters with default strategy */
func NewMnMigradWithParameters(fcn FCNBase, par *MnUserParameters) *MnMigrad {
	return NewMnMigradWithParametersStra(fcn, par, DEFAULT_STRATEGY)
}

/** construct from FCNBase + MnUserParameters */
func NewMnMigradWithParametersStra(fcn FCNBase, par *MnUserParameters, stra int) *MnMigrad {
	return NewMnMigradWithParameterStateStrategy(fcn, NewUserParameterStateFromUserParameter(par), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameters + MnUserCovariance with default strategy */
func NewMnMigradWithParametersCovariance(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) *MnMigrad {
	return NewMnMigradWithParametersCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/** construct from FCNBase + MnUserParameters + MnUserCovariance */
func NewMnMigradWithParametersCovarianceStra(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance, stra int) *MnMigrad {
	ups, _ := NewUserParamStateFromUserParamCovariance(par, cov)
	return NewMnMigradWithParameterStateStrategy(fcn, ups, NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameterState + MnStrategy */
func NewMnMigradWithParameterStateStrategy(fcn FCNBase, par *MnUserParameterState, str *MnStrategy) *MnMigrad {
	mingrad := &MnMigrad{
		baseImpl:     NewMnApplicationWithFcnStateStra(fcn, par, str),
		theMinimizer: NewVariableMetricMinimizer(),
	}
	mingrad.baseImpl.super = mingrad
	return mingrad
}

func (this *MnMigrad) Minimizer() ModularFunctionMinimizerInterface {
	return this.theMinimizer
}
