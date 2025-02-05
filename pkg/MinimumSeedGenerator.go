package minuit

// MinimumSeedGenerator
/* base class for seed generators (starting values); the seed generator
 * prepares initial starting values from the input (MnUserParameterState)
 * for the minimization;
 */
type MinimumSeedGenerator interface {
	Generate(fcn *MnFcn, calc GradientCalculator, user *MnUserParameterState, stra *MnStrategy) (*MinimumSeed, error)
}
