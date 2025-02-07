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

func (this *CombinedMinimumBuilder) Minimum(fcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	min, err := this.theVMMinimizer.minimize(fcn, gc, seed, strategy, maxfcn, toler)
	if err != nil {
		return nil, err
	}

	if !min.IsValid() {
		log.Println("CombinedMinimumBuilder: migrad method fails, will try with simplex method first.")
		var str *MnStrategy = NewMnStrategyWithStra(2)
		min1, err := this.theSimplexMinimizer.minimize(fcn, gc, seed, str, maxfcn, toler)
		if err != nil {
			return nil, err
		}
		if !min1.IsValid() {
			log.Println("CombinedMinimumBuilder: both migrad and simplex method fail.")
			return min1, nil
		}
		seed1, err := this.theVMMinimizer.SeedGenerator().Generate(fcn, gc, min1.UserState(), str)
		if err != nil {
			return nil, err
		}

		min2, err := this.theVMMinimizer.minimize(fcn, gc, seed1, str, maxfcn, toler)
		if err != nil {
			return nil, err
		}
		if !min2.IsValid() {
			log.Println("CombinedMinimumBuilder: both migrad and method fails also at 2nd attempt.")
			log.Println("CombinedMinimumBuilder: return simplex minimum.")
			return min1, nil
		}
		return min2, nil
	}
	return min, nil
}
