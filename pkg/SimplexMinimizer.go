package minuit

type SimplexMinimizer struct {
	ModularFunctionMinimizer
	theSeedGenerator *SimplexSeedGenerator
	theBuilder       *SimplexBuilder
}

func NewSimplexMinimizer() *SimplexMinimizer {
	return &SimplexMinimizer{
		theSeedGenerator: NewSimplexSeedGenerator(),
		theBuilder:       NewSimplexBuilder(),
	}
}

func (this *SimplexMinimizer) SeedGenerator() MinimumSeedGenerator {
	return this.theSeedGenerator
}
func (this *SimplexMinimizer) Builder() MinimumBuilder {
	return this.theBuilder
}
