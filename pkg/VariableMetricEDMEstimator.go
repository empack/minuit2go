package minuit

type VariableMetricEDMEstimator struct {
}

func NewVariableMetricEDMEstimator() *VariableMetricEDMEstimator {
	return &VariableMetricEDMEstimator{}
}

func (vme *VariableMetricEDMEstimator) estimate(g *FunctionGradient, e *MinimumError) (float64, error) {
	if e.invHessian().size() == 1 {
		ih, err := e.invHessian().get(0, 0)
		if err != nil {
			return 0, err
		}
		return 0.5 * g.grad().get(0) * g.grad().get(0) * ih, nil
	} else {
		rho, err := MnUtils.Similarity(g.grad(), e.invHessian())
		if err != nil {
			return 0, err
		}
		return 0.5 * rho, nil
	}
}
