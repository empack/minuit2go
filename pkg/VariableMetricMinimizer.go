package minuit

type VariableMetricMinimizer struct {
	theMinSeedGen *MnSeedGenerator
	theMinBuilder *VariableMetricBuilder
}

func NewVariableMetricMinimizer() *VariableMetricMinimizer {
	return &VariableMetricMinimizer{
		theMinSeedGen: NewMnSeedGenerator(),
		theMinBuilder: NewVariableMetricBuilder(),
	}
}

func (this *VariableMetricMinimizer) SeedGenerator() *MinimumSeedGenerator {
	return this.theMinSeedGen
}
func (this *VariableMetricMinimizer) Builder() *MinimumBuilder {
	return this.theMinBuilder
}
