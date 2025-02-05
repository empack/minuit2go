package minuit

type ModularFunctionMinimizer struct {
}

func (this *ModularFunctionMinimizer) minimizeWithError(fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	var mfcn *MnUserFcn = NewMnUserFcn(fcn, errorDef, st.trafo())

	var gc GradientCalculator
	if _, ok := interface{}(fcn).(FCNGradientBase); ok && useAnalyticalGradient {
		gc = NewAnalyticalGradientCalculator(fcn.(FCNGradientBase), st.trafo(), checkGradient)
	} else {
		gc = NewNumerical2PGradientCalculator(mfcn, st.trafo(), strategy)
	}

	var npar int = st.variableParameters()
	if maxfcn == 0 {
		maxfcn = 200 + 100*npar + 5*npar*npar
	}
	var mnseeds *MinimumSeed = this.SeedGenerator().Generate(mfcn.ParentClass, gc, st, strategy)

	return this.minimize(mfcn.ParentClass, gc, mnseeds, strategy, maxfcn, toler)
}

func (this *ModularFunctionMinimizer) minimize(mfcn *MnFcn, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.Builder().Minimum(mfcn, gc, seed, strategy, maxfcn, toler*mfcn.errorDef())
}

func (this *ModularFunctionMinimizer) SeedGenerator() MinimumSeedGenerator {
	panic("this SeedGenerator base implementation should never be called and overwritten by super stuct")
}
func (this *ModularFunctionMinimizer) Builder() MinimumBuilder {
	panic("this Builder base implementation should never be called and overwritten by super stuct")
}
