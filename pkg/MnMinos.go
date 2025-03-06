package minuit

import (
	"context"
	"errors"
	"log"
	"math"
)

type MnMinos struct {
	theFCN      FCNBase
	theMinimum  *FunctionMinimum
	theStrategy *MnStrategy
}

/** construct from FCN + minimum */
func NewMnMinos(fcn FCNBase, min *FunctionMinimum) *MnMinos {
	return NewMnMinosWithFunctionMinimumStra(fcn, min, DEFAULT_STRATEGY)
}

/** construct from FCN + minimum + strategy */
func NewMnMinosWithFunctionMinimumStra(fcn FCNBase, min *FunctionMinimum, stra int) *MnMinos {
	return NewMnMinosWithFunctionMinimumStrategy(fcn, min, NewMnStrategyWithStra(stra))
}

/** construct from FCN + minimum + strategy */
func NewMnMinosWithFunctionMinimumStrategy(fcn FCNBase, min *FunctionMinimum, stra *MnStrategy) *MnMinos {
	return &MnMinos{
		theFCN:      fcn,
		theMinimum:  min,
		theStrategy: stra,
	}
}

func (this *MnMinos) minos(ctx context.Context, par int) (*MinosError, error) {
	return this.minosWithErrDef(ctx, par, 1.)
}

func (this *MnMinos) minosWithErrDef(ctx context.Context, par int, errDef float64) (*MinosError, error) {
	return this.minosWithErrDefMaxCalls(ctx, par, errDef, DEFAULT_MAXFCN)
}

/**
 * Causes a MINOS error analysis to be performed on the parameter whose number is
 * specified. MINOS errors may be expensive to calculate, but are very reliable since
 * they take account of non-linearities in the problem as well as parameter correlations,
 * and are in general asymmetric.
 * @param maxcalls Specifies the (approximate) maximum number of function calls per parameter
 * requested, after which the calculation will be stopped for that parameter.
 */
func (this *MnMinos) minosWithErrDefMaxCalls(ctx context.Context, par int, errDef float64, maxcalls int) (*MinosError, error) {
	if !this.theMinimum.IsValid() {
		return nil, errors.New("assertion violation: Minimum is Invalid")
	}
	if this.theMinimum.UserState().parameter(par).IsFixed() {
		return nil, errors.New("assertion violation: parameter is fixed")
	}
	if this.theMinimum.UserState().parameter(par).IsConst() {
		return nil, errors.New("assertion violation: parameter is constant")
	}

	up, err := this.UpvalWithErrDefMaxCalls(ctx, par, errDef, maxcalls)
	if err != nil {
		return nil, err
	}
	lo, err := this.LovalWithErrDefMaxCalls(ctx, par, errDef, maxcalls)
	if err != nil {
		return nil, err
	}

	return NewMinosErrorWithValues(par, this.theMinimum.UserState().Value(par), lo, up), nil
}

func (this *MnMinos) Range(ctx context.Context, par int) (*Point, error) {
	return this.RangeWithErrDef(ctx, par, 1)
}

func (this *MnMinos) RangeWithErrDef(ctx context.Context, par int, errDef float64) (*Point, error) {
	return this.RangeWithErrDefMaxCalls(ctx, par, errDef, DEFAULT_MAXFCN)
}

// RangeWithErrDefMaxCalls
/*
 * Causes a MINOS error analysis for external parameter n.
 * @return The lower and upper bounds of parameter
 */
func (this *MnMinos) RangeWithErrDefMaxCalls(ctx context.Context, par int, errDef float64, maxcalls int) (*Point, error) {
	mnerr, err := this.minosWithErrDefMaxCalls(ctx, par, errDef, maxcalls)
	if err != nil {
		return nil, err
	}
	return mnerr.Range(), nil
}

func (this *MnMinos) Lower(ctx context.Context, par int) (float64, error) {
	return this.LowerWithErrDef(ctx, par, 1)
}

func (this *MnMinos) LowerWithErrDef(ctx context.Context, par int, errDef float64) (float64, error) {
	return this.LowerWithErrDefMaxCalls(ctx, par, errDef, DEFAULT_MAXFCN)
}

// LowerWithErrDefMaxCalls
/** calculate one side (negative or positive error) of the parameter */
func (this *MnMinos) LowerWithErrDefMaxCalls(ctx context.Context, par int, errDef float64, maxcalls int) (float64, error) {
	var upar *MnUserParameterState = this.theMinimum.UserState()
	var err float64 = this.theMinimum.UserState().Error(par)
	aopt, fnErr := this.LovalWithErrDefMaxCalls(ctx, par, errDef, maxcalls)
	if fnErr != nil {
		return 0, fnErr
	}
	if aopt.isValid() {
		return -1. * err * (1. + aopt.value()), nil
	} else if aopt.atLimit() {
		return upar.parameter(par).LowerLimit(), nil
	} else {
		return upar.Value(par), nil
	}
}

func (this *MnMinos) Upper(ctx context.Context, par int) (float64, error) {
	return this.UpperWithErrDef(ctx, par, 1)
}

func (this *MnMinos) UpperWithErrDef(ctx context.Context, par int, errDef float64) (float64, error) {
	return this.UpperWithErrDefMaxCalls(ctx, par, errDef, DEFAULT_MAXFCN)
}

func (this *MnMinos) UpperWithErrDefMaxCalls(ctx context.Context, par int, errDef float64, maxcalls int) (float64, error) {
	var upar *MnUserParameterState = this.theMinimum.UserState()
	var err float64 = this.theMinimum.UserState().Error(par)
	aopt, fnErr := this.UpvalWithErrDefMaxCalls(ctx, par, errDef, maxcalls)
	if fnErr != nil {
		return 0, fnErr
	}
	if aopt.isValid() {
		return err * (1. + aopt.value()), nil
	} else if aopt.atLimit() {
		return upar.parameter(par).UpperLimit(), nil
	} else {
		return upar.Value(par), nil
	}
}

func (this *MnMinos) Loval(ctx context.Context, par int) (*MnCross, error) {
	return this.LovalWithErrDef(ctx, par, 1)
}

func (this *MnMinos) LovalWithErrDef(ctx context.Context, par int, errDef float64) (*MnCross, error) {
	return this.LovalWithErrDefMaxCalls(ctx, par, errDef, DEFAULT_MAXFCN)
}

func (this *MnMinos) LovalWithErrDefMaxCalls(ctx context.Context, par int, errDef float64, maxcalls int) (*MnCross, error) {
	errDef *= this.theMinimum.ErrorDef()
	if !this.theMinimum.IsValid() {
		return nil, errors.New("assertion violation: Minimum is Invalid")
	}
	if this.theMinimum.UserState().parameter(par).IsFixed() {
		return nil, errors.New("assertion violation: parameter is fixed")
	}
	if this.theMinimum.UserState().parameter(par).IsConst() {
		return nil, errors.New("assertion violation: parameter is constant")
	}

	if maxcalls == 0 {
		var nvar int = this.theMinimum.UserState().VariableParameters()
		maxcalls = 2 * (nvar + 1) * (200 + 100*nvar + 5*nvar*nvar)
	}

	var para []int = []int{par}

	var upar *MnUserParameterState = this.theMinimum.UserState().clone()
	var err float64 = upar.Error(par)
	var val float64 = upar.Value(par) - err
	var xmid []float64 = []float64{val}
	var xdir []float64 = []float64{-err}

	ind, fnErr := upar.intOfExt(par)
	if fnErr != nil {
		return nil, fnErr
	}
	var m *MnAlgebraicSymMatrix = this.theMinimum.error().matrix()
	var xunit float64 = math.Sqrt(errDef / err)
	for i := 0; i < m.nrow(); i++ {
		if i == ind {
			continue
		}
		v, fnErr := m.get(ind, i)
		if fnErr != nil {
			return nil, fnErr
		}
		var xdev float64 = xunit * v
		var ext int = upar.ExtOfInt(i)
		upar.SetValue(ext, upar.Value(ext)-xdev)
	}

	upar.Fix(par)
	upar.SetValue(par, val)

	var toler float64 = 0.1
	var cross *MnFunctionCross = NewMnFunctionCross(this.theFCN, upar, this.theMinimum.Fval(), this.theStrategy, errDef)
	aopt, fnErr := cross.cross(ctx, para, xmid, xdir, toler, maxcalls)
	if fnErr != nil {
		return nil, fnErr
	}

	if aopt.atLimit() {
		log.Printf("MnMinos parameter %d is at lower limit.\n", par)
	}
	if aopt.atMaxFcn() {
		log.Printf("MnMinos maximum number of function calls exceeded for parameter %d.\n", par)
	}
	if aopt.newMinimum() {
		log.Printf("MnMinos new minimum found while looking for parameter %d.\n", par)
	}
	if !aopt.isValid() {
		log.Printf("MnMinos could not find lower value for parameter %d.\n", par)
	}
	return aopt, nil
}

func (this *MnMinos) Upval(ctx context.Context, par int) (*MnCross, error) {
	return this.UpvalWithErrDef(ctx, par, 1)
}
func (this *MnMinos) UpvalWithErrDef(ctx context.Context, par int, errDef float64) (*MnCross, error) {
	return this.UpvalWithErrDefMaxCalls(ctx, par, errDef, DEFAULT_MAXFCN)
}
func (this *MnMinos) UpvalWithErrDefMaxCalls(ctx context.Context, par int, errDef float64, maxcalls int) (*MnCross, error) {
	errDef *= this.theMinimum.ErrorDef()
	if !this.theMinimum.IsValid() {
		return nil, errors.New("assertion violation: Minimum is Invalid")
	}
	if this.theMinimum.UserState().parameter(par).IsFixed() {
		return nil, errors.New("assertion violation: parameter is fixed")
	}
	if this.theMinimum.UserState().parameter(par).IsConst() {
		return nil, errors.New("assertion violation: parameter is constant")
	}
	if maxcalls == 0 {
		var nvar int = this.theMinimum.UserState().VariableParameters()
		maxcalls = 2 * (nvar + 1) * (200 + 100*nvar + 5*nvar*nvar)
	}

	var para []int = []int{par}

	var upar *MnUserParameterState = this.theMinimum.UserState().clone()
	var err float64 = upar.Error(par)
	var val float64 = upar.Value(par) + err
	var xmid []float64 = []float64{val}
	var xdir []float64 = []float64{err}

	ind, fnErr := upar.intOfExt(par)
	if fnErr != nil {
		return nil, fnErr
	}
	var m *MnAlgebraicSymMatrix = this.theMinimum.error().matrix()
	var xunit float64 = math.Sqrt(errDef / err)
	for i := 0; i < m.nrow(); i++ {
		if i == ind {
			continue
		}
		v, fnErr := m.get(ind, i)
		if fnErr != nil {
			return nil, fnErr
		}
		var xdev float64 = xunit * v
		var ext int = upar.ExtOfInt(i)
		upar.SetValue(ext, upar.Value(ext)+xdev)
	}

	upar.Fix(par)
	upar.SetValue(par, val)

	var toler float64 = 0.1
	var cross *MnFunctionCross = NewMnFunctionCross(this.theFCN, upar, this.theMinimum.Fval(), this.theStrategy, errDef)
	aopt, fnErr := cross.cross(ctx, para, xmid, xdir, toler, maxcalls)
	if fnErr != nil {
		return nil, fnErr
	}

	if aopt.atLimit() {
		log.Printf("MnMinos parameter %d is at upper limit.")
	}
	if aopt.atMaxFcn() {
		log.Printf("MnMinos maximum number of function calls exceeded for parameter %d.", par)
	}
	if aopt.newMinimum() {
		log.Printf("MnMinos new minimum found while looking for parameter %d.", par)
	}
	if !aopt.isValid() {
		log.Printf("MnMinos could not find upper value for parameter %d.", par)
	}
	return aopt, nil
}
