package minuit

type MinimumBuilder interface {
	Minimum(fcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error)
}
