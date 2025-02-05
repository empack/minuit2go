package minuit

type MnScan struct {
	theMinimizer *ScanMinimizer
	MnApplication
}

func NewMnScan(fcn FCNBase, par []float64, err []float64) *MnScan {
	return NewMnScanWithStrategy(fcn, par, err, DEFAULT_STRATEGY)
}

func NewMnScanWithStrategy(fcn FCNBase, par []float64, err []float64, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(fcn, NewMnUserParameterState(par, err), NewMnStrategyWithStra(stra))
}

func NewMnScanWithCovariance(fcn FCNBase, par []float64, cov MnUserCovariance) *MnScan {
	return NewMnScanWithCovarianceAndStrategy(fcn, par, cov, DEFAULT_STRATEGY)
}

func NewMnScanWithCovarianceAndStrategy(fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(fcn, NewMnUserParameterStateWithCovariance(par, cov), NewMnStrategyWithStra(stra))
}

func NewMnScanWithParameters(fcn FCNBase, par *MnUserParameters) *MnScan {
	return NewMnScanWithParametersAndStrategy(fcn, par, DEFAULT_STRATEGY)
}

func NewMnScanWithParametersAndStrategy(fcn FCNBase, par *MnUserParameters, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(fcn, NewMnUserParameterStateWithParameters(par), NewMnStrategyWithStra(stra))
}

func NewMnScanWithParametersAndCovariance(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) *MnScan {
	return NewMnScanWithParametersAndCovarianceAndStrategy(fcn, par, cov, DEFAULT_STRATEGY)
}

func NewMnScanWithParametersAndCovarianceAndStrategy(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance,
	stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(fcn, NewMnUserParameterStateWithParametersAndCovariance(par, cov), NewMnStrategyWithStra(stra))
}

func NewMnScanWithStateAndStrategy(fcn FCNBase, state *MnUserParameterState, str *MnStrategy) *MnScan {
	return &MnScan{
		theMinimizer: NewScanMinimizer(),
	}
}

func (this *MnScan) Minimizer() *ModularFunctionMinimizer {
	return this.theMinimizer
}

func (this *MnScan) Scan(par int) []*Point {
	return this.ScanWithMaxsteps(par, 41)
}

func (this *MnScan) ScanWithMaxsteps(par, maxsteps int) []*Point {
	return this.ScanWithMaxstepsRange(par, maxsteps, 0.0, 0.0)
}

func (this *MnScan) ScanWithMaxstepsRange(par, maxsteps int, low, high float64) []*Point {
	scan := NewMnParameterScan(this.theFCN, this.theState.parameters())
	amin := scan.fval()
	result := scan.scan(par, maxsteps, low, high)
	if scan.fval() < amin {
		this.theState.SetValue(par, scan.parameters().value(par))
		amin = scan.fval()
	}

	return result
}
