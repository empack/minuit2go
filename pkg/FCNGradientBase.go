package minuit

// FCNGradientBase
/*
 * Extension of FCNBase for providing the analytical gradient of the
 * function. The user-gradient is checked at the beginning of the
 * minimization against the Minuit internal numerical gradient in order to
 * spot problems in the analytical gradient calculation.
 * @version $Id: FCNGradientBase.java 8584 2006-08-10 23:06:37Z duns $
 * @see MnApplication#setUseAnalyticalDerivatives
 * @see MnApplication#setCheckAnalyticalDerivatives
 */
type FCNGradientBase interface {
	FCNBase
	/*
	 * Calculate the function gradient with respect to each parameter at the
	 * given point in parameter space. The size of the output gradient vector
	 * must be equal to the size of the input parameter vector.
	 */
	Gradient(par []float64) []float64
}
