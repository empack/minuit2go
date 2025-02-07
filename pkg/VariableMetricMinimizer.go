package minuit

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
func (this *VariableMetricMinimizer) minimizeWithError(fcn FCNBase, st *MnUserParameterState, strategy *MnStrategy, maxfcn int, toler, errorDef float64, useAnalyticalGradient, checkGradient bool) (*FunctionMinimum, error) {
	return this.baseImpl.minimizeWithError(fcn, st, strategy, maxfcn, toler, errorDef, useAnalyticalGradient, checkGradient)
}

func (this *VariableMetricMinimizer) minimize(mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	return this.baseImpl.minimize(mfcn, gc, seed, strategy, maxfcn, toler)
}

func (this *VariableMetricMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theMinSeedGen
}
func (this *VariableMetricMinimizer) Builder() MinimumBuilder {
	return this.theMinBuilder
}
