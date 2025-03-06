package minuit

import "context"

type MinimumBuilder interface {
	Minimum(ctx context.Context, fcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error)
}
