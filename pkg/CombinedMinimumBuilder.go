package minuit

import "log"

type CombinedMinimumBuilder struct {
	theVMMinimizer      *VariableMetricMinimizer
	theSimplexMinimizer *SimplexMinimizer
}

func NewCombinedMinimumBuilder() *CombinedMinimumBuilder {
	return &CombinedMinimumBuilder{
		theVMMinimizer:      NewVariableMetricMinimizer(),
		theSimplexMinimizer: NewSimplexMinimizer(),
	}
}

func (this *CombinedMinimumBuilder) Minimum(fcn *MnFcn, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	var min *FunctionMinimum = this.theVMMinimizer.minimize(fcn, gc, seed, strategy, maxfcn, toler)

	if !min.isValid() {
		log.Println("CombinedMinimumBuilder: migrad method fails, will try with simplex method first.")
		var str *MnStrategy = NewMnStrategy(2)
		var min1 *FunctionMinimum = this.theSimplexMinimizer.minimize(fcn, gc, seed, str, maxfcn, toler)
		if !min1.isValid() {
			log.Println("CombinedMinimumBuilder: both migrad and simplex method fail.")
			return min1, nil
		}
		var seed1 *MinimumSeed = this.theVMMinimizer.SeedGenerator().generate(fcn, gc, min1.userState(), str)

		var min2 *FunctionMinimum = this.theVMMinimizer.minimize(fcn, gc, seed1, str, maxfcn, toler)
		if !min2.isValid() {
			log.Println("CombinedMinimumBuilder: both migrad and method fails also at 2nd attempt.")
			log.Println("CombinedMinimumBuilder: return simplex minimum.")
			return min1, nil
		}
		return min2, nil
	}
	return min, nil
}
