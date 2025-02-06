package minuit

import "math"

type SimplexSeedGenerator struct {
}

func NewSimplexSeedGenerator() *SimplexSeedGenerator {
	return &SimplexSeedGenerator{}
}

func (_ *SimplexSeedGenerator) Generate(fcn *MnFcn, gc GradientCalculator, st *MnUserParameterState, stra *MnStrategy) (*MinimumSeed, error) {
	var n int = st.variableParameters()
	var prec *MnMachinePrecision = st.precision()

	// initial starting values
	var x *MnAlgebraicVector = NewMnAlgebraicVector(n)
	for i := 0; i < n; i++ {
		x.set(i, st.intParameters().get(i))
	}
	var fcnmin float64 = fcn.valueOf(x)
	var pa *MinimumParameters = NewMinimumParameters(x, fcnmin)
	var igc *InitialGradientCalculator = NewInitialGradientCalculator(fcn, st.trafo(), stra)
	var dgrad *FunctionGradient = igc.Gradient(pa)
	var mat *MnAlgebraicSymMatrix = NewMnAlgebraicSymMatrix(n)
	var dcovar float64 = 1.0
	for i := 0; i < n; i++ {
		if math.Abs(dgrad.g2().get(i)) > prec.eps2() {
			mat.set(i, i, 1./dgrad.g2().get(i))
		} else {
			mat.set(i, i, 1.0)
		}
	}

	var err *MinimumError = NewMinimumError(mat, dcovar)
	var edm float64 = NewVariableMetricEDMEstimator().estimate(dgrad, err)
	var state *MinimumState = NewMinimumStateWithGrad(pa, err, dgrad, edm, fcn.numOfCalls())

	return NewMinimumSeed(state, st.trafo())
}
