package minuit

// MinimumSeedGenerator
/* base class for seed generators (starting values); the seed generator
 * prepares initial starting values from the input (MnUserParameterState)
 * for the minimization;
 */
type MinimumSeedGenerator interface {
	generate(fcn *MnFcn, calc GradientCalculator, user *MnUserParameterState, stra *MnStrategy) *MinimumSeed
}
