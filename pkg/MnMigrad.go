package minuit

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
	*MnApplication
	theMinimizer *VariableMetricMinimizer
}

/** construct from FCNBase + double[] for parameters and errors with default strategy */
func NewMnMigradWithParErr(fcn FCNBase, par, err []float64) *MnMigrad {
	return NewMnMigradWithParErrStrategy(fcn, par, err, DEFAULT_STRATEGY)
}

/** construct from FCNBase + double[] for parameters and errors */
func NewMnMigradWithParErrStrategy(fcn FCNBase, par, err []float64, stra int) *MnMigrad {
	return NewMnMigradWithParameterStateStrategy(fcn, NewMnUserParameterState(par, err), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + double[] for parameters and MnUserCovariance with default strategy */
func NewMnMigradWithParCovariance(fcn FCNBase, par []float64, cov *MnUserCovariance) *MnMigrad {
	return NewMnMigradWithParCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/** construct from FCNBase + double[] for parameters and MnUserCovariance */
func NewMnMigradWithParCovarianceStra(fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) *MnMigrad {
	return NewMnMigradWithParameterStateStrategy(fcn, NewMnUserParameterState(par, cov), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameters with default strategy */
func NewMnMigradWithParameters(fcn FCNBase, par *MnUserParameters) *MnMigrad {
	return NewMnMigradWithParametersStra(fcn, par, DEFAULT_STRATEGY)
}

/** construct from FCNBase + MnUserParameters */
func NewMnMigradWithParametersStra(fcn FCNBase, par *MnUserParameters, stra int) *MnMigrad {
	return NewMnMigradWithParametersCovariance(fcn, NewMnUserParameterState(par), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameters + MnUserCovariance with default strategy */
func NewMnMigradWithParametersCovariance(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) *MnMigrad {
	return NewMnMigradWithParametersCovarianceStra(fcn, par, cov, DEFAULT_STRATEGY)
}

/** construct from FCNBase + MnUserParameters + MnUserCovariance */
func NewMnMigradWithParametersCovarianceStra(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance, stra int) *MnMigrad {
	return NewMnMigradWithParameterStateStrategy(fcn, NewMnUserParameterState(par, cov), NewMnStrategyWithStra(stra))
}

/** construct from FCNBase + MnUserParameterState + MnStrategy */
func NewMnMigradWithParameterStateStrategy(fcn FCNBase, par *MnUserParameterState, str *MnStrategy) *MnMigrad {
	return &MnMigrad{
		MnApplication: NewMnApplicationWithFcnStateStra(fcn, par, str),
		theMinimizer:  NewVariableMetricMinimizer(),
	}
}

func (this *MnMigrad) minimizer() ModularFunctionMinimizer {
	return this.theMinimizer.ModularFunctionMinimizer //TODO fix when getting panic from calling SeedGenerator()
}
