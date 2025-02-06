package minuit

// MnSimplex
/*
 * SIMPLEX is a function minimization method using the simplex method of Nelder and
 * Mead. MnSimplex provides minimization of the function by the method of SIMPLEX
 * and the functionality for parameters interaction. It also retains the result from the
 * last minimization in case the user may want to do subsequent minimization steps with
 * parameter interactions in between the minimization requests. As SIMPLEX is a
 * stepping method it does not produce a covariance matrix.
 */
type MnSimplex struct {
	*MnApplication
	theMinimizer *SimplexMinimizer
}

/* construct from FCNBase + double[] for parameters and errors with default strategy */
func NewMnSimplex(fcn FCNBase, par, err []float64) *MnSimplex {
	return NewMnSimplexWithParErrStra(fcn, par, err, DEFAULT_STRATEGY)
}

/* construct from FCNBase + double[] for parameters and errors */
func NewMnSimplexWithParErrStra(fcn FCNBase, par, err []float64, stra int) *MnSimplex {
	return NewMnSimplexWithParameterStateStrategy(fcn, NewMnUserParameterState(par, err), NewMnStrategyWithStra(stra))
}

/* construct from FCNBase + double[] for parameters and MnUserCovariance with default strategy */
func NewMnSimplexWithParCov(fcn FCNBase, par []float64, cov *MnUserCovariance) *MnSimplex {
	return NewMnSimplexWithParCovStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/* construct from FCNBase + double[] for parameters and MnUserCovariance */
func NewMnSimplexWithParCovStra(fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) *MnSimplex {
	return NewMnSimplexWithParameterStateStrategy(fcn, NewMnUserParameterState(par, cov), NewMnStrategyWithStra(stra))
}

/* construct from FCNBase + MnUserParameters with default strategy */
func NewMnSimplexWithParameters(fcn FCNBase, par *MnUserParameters) *MnSimplex {
	return NewMnSimplexWithParametersStra(fcn, par, DEFAULT_STRATEGY)
}

/* construct from FCNBase + MnUserParameters */
func NewMnSimplexWithParametersStra(fcn FCNBase, par *MnUserParameters, stra int) *MnSimplex {
	return NewMnSimplexWithParameterStateStrategy(fcn, NewMnUserParameterState(par), NewMnStrategyWithStra(stra))
}

/* construct from FCNBase + MnUserParameters + MnUserCovariance with default strategy */
func NewMnSimplexWithParametersCovariance(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) *MnSimplex {
	return NewMnSimplexWithParametersCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/* construct from FCNBase + MnUserParameters + MnUserCovariance */
func NewMnSimplexWithParametersCovarianceStra(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance, stra int) *MnSimplex {
	return NewMnSimplexWithParameterStateStrategy(fcn, NewMnUserParameterState(par, cov), NewMnStrategyWithStra(stra))
}

/* construct from FCNBase + MnUserParameterState + MnStrategy */
func NewMnSimplexWithParameterStateStrategy(fcn FCNBase, par *MnUserParameterState, str *MnStrategy) *MnSimplex {
	return &MnSimplex{
		MnApplication: NewMnApplicationWithFcnStateStra(fcn, par, str),
		theMinimizer:  NewSimplexMinimizer(),
	}
}

func (this *MnSimplex) Minimizer() *ModularFunctionMinimizer {
	return this.theMinimizer.ModularFunctionMinimizer
}
