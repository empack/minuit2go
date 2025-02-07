package minuit

type ScanMinimizer struct {
	theSeedGenerator *SimplexSeedGenerator
	theBuilder       *ScanBuilder
	baseImpl         *ModularFunctionMinimizer
}

func (this *ScanMinimizer) minimizeWithError(fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	return this.baseImpl.minimizeWithError(fcn, st, strategy, maxfcn, toler, errorDef, useAnalyticalGradient, checkGradient)
}

func (this *ScanMinimizer) minimize(mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.minimize(mfcn, gc, seed, strategy, maxfcn, toler)
}

func NewScanMinimizer() *ScanMinimizer {
	mini := &ScanMinimizer{
		theSeedGenerator: NewSimplexSeedGenerator(),
		theBuilder:       NewScanBuilder(),
		baseImpl:         NewModularFunctionMinimizer(),
	}
	mini.baseImpl.super = mini
	return mini
}

func (this *ScanMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theSeedGenerator
}

func (this *ScanMinimizer) Builder() MinimumBuilder {
	return this.theBuilder
}
