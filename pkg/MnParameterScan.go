package minuit

import "math"

// MnParameterScan
/* Scans the values of FCN as a function of one parameter and retains the
 * best function and parameter values found
 */
type MnParameterScan struct {
	theFCN        FCNBase
	theParameters *MnUserParameters
	theAmin       float64
}

func NewMnParameterScan(fcn FCNBase, par *MnUserParameters) *MnParameterScan {
	return &MnParameterScan{
		theFCN:        fcn,
		theParameters: par,
		theAmin:       fcn.ValueOf(par.params()),
	}
}

func NewMnParameterScanWithFval(fcn FCNBase, par *MnUserParameters, fval float64) *MnParameterScan {
	return &MnParameterScan{
		theFCN:        fcn,
		theParameters: par,
		theAmin:       fval,
	}
}

func (this *MnParameterScan) scan(par int) []*Point {
	return this.scanWithMaxSteps(par, 41)
}
func (this *MnParameterScan) scanWithMaxSteps(par, maxsteps int) []*Point {
	return this.scanWithMaxStepsLowHigh(par, maxsteps, 0, 0)
}

/* returns pairs of (x,y) points, x=parameter value, y=function value of FCN */
func (this *MnParameterScan) scanWithMaxStepsLowHigh(par, maxsteps int, low, high float64) []*Point {
	if maxsteps > 101 {
		maxsteps = 101
	}
	var result []*Point = make([]*Point, 0, maxsteps+1)
	var params []float64 = this.theParameters.params()
	result = append(result, NewPoint(params[par], this.theAmin))

	if low > high {
		return result
	}
	if maxsteps < 2 {
		return result
	}

	if low == 0.0 && high == 0.0 {
		low = params[par] - 2.*this.theParameters.error(par)
		high = params[par] + 2.*this.theParameters.error(par)
	}

	if low == 0. && high == 0. && this.theParameters.parameter(par).hasLimits() {
		if this.theParameters.parameter(par).hasLowerLimit() {
			low = this.theParameters.parameter(par).lowerLimit()
		}
		if this.theParameters.parameter(par).hasUpperLimit() {
			high = this.theParameters.parameter(par).upperLimit()
		}
	}

	if theParameters.parameter(par).hasLimits() {
		if this.theParameters.parameter(par).hasLowerLimit() {
			low = math.Max(low, this.theParameters.parameter(par).lowerLimit())
		}
		if this.theParameters.parameter(par).hasUpperLimit() {
			high = math.Min(high, this.theParameters.parameter(par).upperLimit())
		}
	}

	var x0 float64 = low
	var stp float64 = (high - low) / (float64(maxsteps) - 1.0)
	for i := 0; i < maxsteps; i++ {
		params[par] = x0 + (float64(i) * stp)
		var fval float64 = this.theFCN.ValueOf(params)
		if fval < this.theAmin {
			this.theParameters.setValue(par, params[par])
			this.theAmin = fval
		}
		result = append(result, NewPoint(params[par], fval))
	}

	return result
}

func (this *MnParameterScan) parameters() *MnUserParameters {
	return this.theParameters
}
func (this *MnParameterScan) fval() float64 {
	return this.theAmin
}
