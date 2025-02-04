package minuit

import "errors"

type AnalyticalGradientCalculator struct {
	theGradCalc       FCNGradientBase
	theTransformation *MnUserTransformation
	theCheckGradient  bool
}

func NewAnalyticalGradientCalculator(fcn FCNGradientBase, state *MnUserTransformation, checkGradient bool) *AnalyticalGradientCalculator {
	return &AnalyticalGradientCalculator{
		theGradCalc:       fcn,
		theTransformation: state,
		theCheckGradient:  checkGradient,
	}
}

func (this *AnalyticalGradientCalculator) Gradient(par *MinimumParameters) (*FunctionGradient, error) {
	var grad []float64 = this.theGradCalc.Gradient(this.theTransformation.transform(par.vec()).data)
	if len(grad) != len(this.theTransformation.parameters()) {
		return nil, errors.New("IllegalArgumentException: Invalid parameter size")
	}

	var v *MnAlgebraicVector = NewMnAlgebraicVector(par.vec().size())
	for i := 0; i < par.vec().size(); i++ {
		var ext int = this.theTransformation.extOfInt(i)
		if this.theTransformation.parameter(ext).HasLimits() {
			var dd float64 = this.theTransformation.dInt2Ext(i, par.vec().get(i))
			v.set(i, dd*grad[ext])
		} else {
			v.set(i, grad[ext])
		}
	}

	return NewFunctionGradient(v), nil
}

func (this *AnalyticalGradientCalculator) checkGradient() bool {
	return this.theCheckGradient
}
