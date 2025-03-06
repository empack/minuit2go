package minuit

import "context"

type ScanMinimizer struct {
	theSeedGenerator *SimplexSeedGenerator
	theBuilder       *ScanBuilder
	baseImpl         *ModularFunctionMinimizer
}

func (this *ScanMinimizer) minimizeWithError(ctx context.Context, fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	return this.baseImpl.minimizeWithError(ctx, fcn, st, strategy, maxfcn, toler, errorDef, useAnalyticalGradient, checkGradient)
}

func (this *ScanMinimizer) minimize(ctx context.Context, mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.minimize(ctx, mfcn, gc, seed, strategy, maxfcn, toler)
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
