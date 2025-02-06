package minuit

// MnMinimize
/*
 * Causes minimization of the function by the method of MIGRAD, as does the MnMigrad
 * class, but switches to the SIMPLEX method if MIGRAD fails to converge. Constructor
 * arguments, methods arguments and names of methods are the same as for MnMigrad
 * or MnSimplex.
 */
type MnMinimize struct {
	*MnApplication
	theMinimizer *CombinedMinimizer
}

/** construct from FCNBase + double[] for parameters and errors with default strategy */
func NewMnMinimize(fcn FCNBase, par, err []float64) *MnMinimize {
	return NewMnMinimizeWithParErrStra(fcn, par, err, DEFAULT_STRATEGY)
}

/** construct from FCNBase + double[] for parameters and errors */
func NewMnMinimizeWithParErrStra(fcn FCNBase, par, err []float64, stra int) *MnMinimize {
	return NewMnMinimizeWithParameterStateStrategy(fcn, NewMnUserParameterState(par, err), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + double[] for parameters and MnUserCovariance with default strategy */
func NewMnMinimizeWithParCovariance(fcn FCNBase, par []float64, cov *MnUserCovariance) *MnMinimize {
	return NewMnMinimizeWithParCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/** construct from FCNBase + double[] for parameters and MnUserCovariance */
func NewMnMinimizeWithParCovarianceStra(fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) *MnMinimize {
	return NewMnMinimizeWithParameterStateStrategy(fcn, NewMnUserParameterState(par, cov), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameters with default strategy */
func NewMnMinimizeWithParameters(fcn FCNBase, par *MnUserParameters) *MnMinimize {
	return NewMnMinimizeWithParametersStra(fcn, par, DEFAULT_STRATEGY)
}

/** construct from FCNBase + MnUserParameters */
func NewMnMinimizeWithParametersStra(fcn FCNBase, par *MnUserParameters, stra int) *MnMinimize {
	return NewMnMinimizeWithParameterStateStrategy(fcn, NewMnUserParameterState(par), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameters + MnUserCovariance with default strategy */
func NewMnMinimizeWithParametersCovariance(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) *MnMinimize {
	return NewMnMinimizeWithParametersCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/** construct from FCNBase + MnUserParameters + MnUserCovariance */
func NewMnMinimizeWithParametersCovarianceStra(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance, stra int) *MnMinimize {
	return NewMnMinimizeWithParameterStateStrategy(fcn, NewMnUserParameterState(par, cov), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameterState + MnStrategy */
func NewMnMinimizeWithParameterStateStrategy(fcn FCNBase, par *MnUserParameterState, str *MnStrategy) *MnMinimize {
	return &MnMinimize{
		MnApplication: NewMnApplicationWithFcnStateStra(fcn, par, str),
		theMinimizer:  NewCombinedMinimizer(),
	}
}

func (this *MnMinimize) minimizer() *ModularFunctionMinimizer {
	return this.theMinimizer.ModularFunctionMinimizer
}
