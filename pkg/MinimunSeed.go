/*
*

	 *
	 * @version $Id: MinimumSeed.java 8584 2006-08-10 23:06:37Z duns $

	 class MinimumSeed
	 {
		MinimumSeed(MinimumState state, MnUserTransformation trafo)
		{
		   theState = state;
		   theTrafo = trafo;
		   theValid = true;
		}

		MinimumState state()
		{
		   return theState;
		}
		MinimumParameters parameters()
		{
		   return state().parameters();
		}
		MinimumError error()
		{
		   return state().error();
		}
		FunctionGradient gradient()
		{
		   return state().gradient();
		}
		MnUserTransformation trafo()
		{
		   return theTrafo;
		}
		MnMachinePrecision precision()
		{
		   return theTrafo.precision();
		}
		double fval()
		{
		   return state().fval();
		}
		double edm()
		{
		   return state().edm();
		}
		int nfcn()
		{
		   return state().nfcn();
		}
		boolean isValid()
		{
		   return theValid;
		}

		private MinimumState theState;
		private MnUserTransformation theTrafo;
		private boolean theValid;
	 }
*/
package minuit

type MinimumSeed struct {
	theState MinimumState
	theTrafo MnUserTransformation
	theValid bool
}

func NewMinimumSeed(state MinimumState, trafo MnUserTransformation) *MinimumSeed {
	return &MinimumSeed{
		theState: state,
		theTrafo: trafo,
		theValid: true,
	}
}

func (this *MinimumSeed) state() *MinimumState {
	return &this.theState
}

func (this *MinimumSeed) parameters() *MinimumParameters {
	return this.theState.parameters()
}

func (this *MinimumSeed) error() *MinimumError {
	return this.theState.error()
}

func (this *MinimumSeed) gradient() *FunctionGradient {
	return this.theState.gradient()
}

func (this *MinimumSeed) trafo() *MnUserTransformation {
	return &this.theTrafo
}

func (this *MinimumSeed) precision() *MnMachinePrecision {
	return this.theTrafo.precision()
}

func (this *MinimumSeed) fval() float64 {
	return this.theState.fval()
}

func (this *MinimumSeed) edm() float64 {
	return this.theState.edm()
}

func (this *MinimumSeed) nfcn() int {
	return this.theState.nfcn()
}

func (this *MinimumSeed) isValid() bool {
	return this.theValid
}
