package minuit

type CombinedMinimizer struct {
	baseImpl      *ModularFunctionMinimizer
	theMinSeedGen *MnSeedGenerator
	theMinBuilder *CombinedMinimumBuilder
}

func (this *CombinedMinimizer) minimizeWithError(fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	return this.baseImpl.minimizeWithError(fcn, st, strategy, maxfcn, toler, errorDef, useAnalyticalGradient, checkGradient)
}

func (this *CombinedMinimizer) minimize(mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.minimize(mfcn, gc, seed, strategy, maxfcn, toler)
}

func NewCombinedMinimizer() *CombinedMinimizer {
	mini := &CombinedMinimizer{
		baseImpl:      NewModularFunctionMinimizer(),
		theMinSeedGen: NewMnSeedGenerator(),
		theMinBuilder: NewCombinedMinimumBuilder(),
	}
	mini.baseImpl.super = mini
	return mini
}

func (this *CombinedMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theMinSeedGen
}
func (this *CombinedMinimizer) Builder() MinimumBuilder {
	return this.theMinBuilder
}
