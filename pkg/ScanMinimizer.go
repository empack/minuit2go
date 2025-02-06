package minuit

type ScanMinimizer struct {
	theSeedGenerator *SimplexSeedGenerator
	theBuilder       *ScanBuilder
	ParentClass      *ModularFunctionMinimizer
}

func NewScanMinimizer() *ScanMinimizer {
	return &ScanMinimizer{
		theSeedGenerator: NewSimplexSeedGenerator(),
		theBuilder:       NewScanBuilder(),
		ParentClass:      NewModularFunctionMinimizer(),
	}
}

func (this *ScanMinimizer) seedGenerator() MinimumSeedGenerator {
	return this.theSeedGenerator
}

func (this *ScanMinimizer) builder() MinimumBuilder {
	return this.theBuilder
}
