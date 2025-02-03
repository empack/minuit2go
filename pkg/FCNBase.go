package minuit

// FCNBase
/*
 * User function base class, has to be implemented by the user.
 */
type FCNBase interface {
	/*
	 * Returns the value of the function with the given parameters.
	 */
	ValueOf(par []float64) float64
}
