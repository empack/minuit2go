package minuit

type MnScan struct {
	theMinimizer *ScanMinimizer
	MnApplication
}

func NewMnScan(fcn FCNBase, par []float64, err []float64) *MnScan {
	return NewMnScanWithStrategy(fcn, par, err, DEFAULT_STRATEGY)
}

func NewMnScanWithStrategy(fcn FCNBase, par []float64, err []float64, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(fcn, NewUserParamStateFromParamAndErrValues(par, err), NewMnStrategyWithStra(stra))
}

func NewMnScanWithCovariance(fcn FCNBase, par []float64, cov *MnUserCovariance) (*MnScan, error) {
	return NewMnScanWithCovarianceAndStrategy(fcn, par, cov, DEFAULT_STRATEGY)
}

func NewMnScanWithCovarianceAndStrategy(fcn FCNBase, par []float64, cov *MnUserCovariance, stra int) (*MnScan, error) {
	state, fnErr := NewMnUserParameterStateFlUc(par, cov)
	if fnErr != nil {
		return nil, fnErr
	}
	return NewMnScanWithStateAndStrategy(fcn, state, NewMnStrategyWithStra(stra)), nil
}

func NewMnScanWithParameters(fcn FCNBase, par *MnUserParameters) *MnScan {
	return NewMnScanWithParametersAndStrategy(fcn, par, DEFAULT_STRATEGY)
}

func NewMnScanWithParametersAndStrategy(fcn FCNBase, par *MnUserParameters, stra int) *MnScan {
	return NewMnScanWithStateAndStrategy(fcn, NewUserParameterStateFromUserParameter(par), NewMnStrategyWithStra(stra))
}

func NewMnScanWithParametersAndCovariance(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance) (*MnScan, error) {
	return NewMnScanWithParametersAndCovarianceAndStrategy(fcn, par, cov, DEFAULT_STRATEGY)
}

func NewMnScanWithParametersAndCovarianceAndStrategy(fcn FCNBase, par *MnUserParameters, cov *MnUserCovariance,
	stra int) (*MnScan, error) {
	state, fnErr := NewUserParamStateFromUserParamCovariance(par, cov)
	if fnErr != nil {
		return nil, fnErr
	}
	return NewMnScanWithStateAndStrategy(fcn, state, NewMnStrategyWithStra(stra)), nil
}

func NewMnScanWithStateAndStrategy(fcn FCNBase, state *MnUserParameterState, str *MnStrategy) *MnScan {
	return &MnScan{
		theMinimizer: NewScanMinimizer(),
	}
}

func (this *MnScan) Minimizer() *ModularFunctionMinimizer {
	return this.theMinimizer.ParentClass
}

func (this *MnScan) Scan(par int) ([]*Point, error) {
	return this.ScanWithMaxsteps(par, 41)
}

func (this *MnScan) ScanWithMaxsteps(par, maxsteps int) ([]*Point, error) {
	return this.ScanWithMaxstepsRange(par, maxsteps, 0.0, 0.0)
}

func (this *MnScan) ScanWithMaxstepsRange(par, maxsteps int, low, high float64) ([]*Point, error) {
	scan := NewMnParameterScan(this.theFCN, this.theState.parameters())
	amin := scan.fval()
	result := scan.scanWithMaxStepsLowHigh(par, maxsteps, low, high)
	if scan.fval() < amin {
		fnError := this.theState.SetValue(par, scan.parameters().Value(par))
		if fnError != nil {
			return nil, fnError
		}
		amin = scan.fval()
	}

	return result, nil
}
