package minuit

type CombinedMinimizer struct {
	*ModularFunctionMinimizer
	theMinSeedGen *MnSeedGenerator
	theMinBuilder *CombinedMinimumBuilder
}

func NewCombinedMinimizer() *CombinedMinimizer {
	return &CombinedMinimizer{
		theMinSeedGen: NewMnSeedGenerator(),
		theMinBuilder: NewCombinedMinimumBuilder(),
	}
}

func (this *CombinedMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theMinSeedGen
}
func (this *CombinedMinimizer) Builder() MinimumBuilder {
	return this.theMinBuilder
}
