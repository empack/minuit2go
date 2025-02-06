package minuit

import "math"

type SimplexSeedGenerator struct {
}

func NewSimplexSeedGenerator() *SimplexSeedGenerator {
	return &SimplexSeedGenerator{}
}

func (_ *SimplexSeedGenerator) Generate(fcn *MnFcn, gc GradientCalculator, st *MnUserParameterState, stra *MnStrategy) (*MinimumSeed, error) {
	var n int = st.VariableParameters()
	var prec *MnMachinePrecision = st.Precision()

	// initial starting values
	var x *MnAlgebraicVector = NewMnAlgebraicVector(n)
	for i := 0; i < n; i++ {
		x.set(i, st.intParameters()[i])
	}
	var fcnmin float64 = fcn.valueOf(x)
	var pa *MinimumParameters = NewMinimumParameters(x, fcnmin)
	var igc *InitialGradientCalculator = NewInitialGradientCalculator(fcn, st.trafo(), stra)
	dgrad, fnErr := igc.gradient(pa)
	if fnErr != nil {
		return nil, fnErr
	}
	mat, fnErr := NewMnAlgebraicSymMatrix(n)
	if fnErr != nil {
		return nil, fnErr
	}
	var dcovar float64 = 1.0
	for i := 0; i < n; i++ {
		if math.Abs(dgrad.g2().get(i)) > prec.eps2() {
			fnErr = mat.set(i, i, 1./dgrad.g2().get(i))
			if fnErr != nil {
				return nil, fnErr
			}
		} else {
			fnErr = mat.set(i, i, 1.0)
			if fnErr != nil {
				return nil, fnErr
			}
		}
	}

	var err *MinimumError = NewMinimumError(mat, dcovar)
	edm, fnErr := NewVariableMetricEDMEstimator().estimate(dgrad, err)
	if fnErr != nil {
		return nil, fnErr
	}
	var state *MinimumState = NewMinimumStateWithGrad(pa, err, dgrad, edm, fcn.numOfCalls())

	return NewMinimumSeed(state, st.trafo()), nil
}
