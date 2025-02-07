package minuit

type SimplexMinimizer struct {
	baseImpl         *ModularFunctionMinimizer
	theSeedGenerator *SimplexSeedGenerator
	theBuilder       *SimplexBuilder
}

func NewSimplexMinimizer() *SimplexMinimizer {
	mini := &SimplexMinimizer{
		baseImpl:         NewModularFunctionMinimizer(),
		theSeedGenerator: NewSimplexSeedGenerator(),
		theBuilder:       NewSimplexBuilder(),
	}
	mini.baseImpl.super = mini
	return mini
}
func (this *SimplexMinimizer) minimizeWithError(fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	return this.baseImpl.minimizeWithError(fcn, st, strategy, maxfcn, toler, errorDef, useAnalyticalGradient, checkGradient)
}

func (this *SimplexMinimizer) minimize(mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.minimize(mfcn, gc, seed, strategy, maxfcn, toler)
}

func (this *SimplexMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theSeedGenerator
}
func (this *SimplexMinimizer) Builder() MinimumBuilder {
	return this.theBuilder
}
