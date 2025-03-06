package minuit

import "context"

type VariableMetricMinimizer struct {
	baseImpl      *ModularFunctionMinimizer
	theMinSeedGen *MnSeedGenerator
	theMinBuilder *VariableMetricBuilder
}

func NewVariableMetricMinimizer() *VariableMetricMinimizer {
	mini := &VariableMetricMinimizer{
		baseImpl:      NewModularFunctionMinimizer(),
		theMinSeedGen: NewMnSeedGenerator(),
		theMinBuilder: NewVariableMetricBuilder(),
	}
	mini.baseImpl.super = mini
	return mini
}
func (this *VariableMetricMinimizer) minimizeWithError(ctx context.Context, fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	return this.baseImpl.minimizeWithError(ctx, fcn, st, strategy, maxfcn, toler, errorDef, useAnalyticalGradient, checkGradient)
}

func (this *VariableMetricMinimizer) minimize(ctx context.Context, mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.minimize(ctx, mfcn, gc, seed, strategy, maxfcn, toler)
}

func (this *VariableMetricMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theMinSeedGen
}
func (this *VariableMetricMinimizer) Builder() MinimumBuilder {
	return this.theMinBuilder
}
