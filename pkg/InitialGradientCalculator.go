package minuit

import (
	"errors"
	"math"
)

type InitialGradientCalculator struct {
	theFcn            MnFcnInterface
	theTransformation *MnUserTransformation
	theStrategy       *MnStrategy
}

func NewInitialGradientCalculator(fcn MnFcnInterface, par *MnUserTransformation, stra *MnStrategy) *InitialGradientCalculator {
	return &InitialGradientCalculator{
		theFcn:            fcn,
		theTransformation: par,
		theStrategy:       stra,
	}
}

func (this *InitialGradientCalculator) gradient(par *MinimumParameters) (*FunctionGradient, error) {
	if !par.isValid() {
		return nil, errors.New("parameters are invalid")
	} else {
		n := this.trafo().variableParameters()
		if n != par.vec().size() {
			return nil, errors.New("parameters have invalid size")
		} else {
			gr := NewMnAlgebraicVector(n)
			gr2 := NewMnAlgebraicVector(n)
			gst := NewMnAlgebraicVector(n)

			for i := 0; i < n; i++ {
				exOfIn := this.trafo().extOfInt(i)
				varValue := par.vec().get(i)
				werr := this.trafo().parameter(exOfIn).Error()
				sav := this.trafo().int2ext(i, varValue)
				sav2 := sav + werr
				if this.trafo().parameter(exOfIn).HasLimits() && this.trafo().parameter(exOfIn).HasUpperLimit() && sav2 > this.trafo().parameter(exOfIn).UpperLimit() {
					sav2 = this.trafo().parameter(exOfIn).UpperLimit()
				}

				var2 := this.trafo().ext2int(exOfIn, sav2)
				vplu := var2 - varValue
				sav2 = sav - werr
				if this.trafo().parameter(exOfIn).HasLimits() && this.trafo().parameter(exOfIn).HasLowerLimit() && sav2 < this.trafo().parameter(exOfIn).LowerLimit() {
					sav2 = this.trafo().parameter(exOfIn).LowerLimit()
				}

				var2 = this.trafo().ext2int(exOfIn, sav2)
				vmin := var2 - varValue
				dirin := 0.5 * (math.Abs(vplu) + math.Abs(vmin))
				g2 := 2.0 * this.theFcn.errorDef() / (dirin * dirin)
				gsmin := 8.0 * this.precision().eps2() * (math.Abs(varValue) + this.precision().eps2())
				gstep := math.Max(gsmin, 0.1*dirin)
				grd := g2 * dirin
				if this.trafo().parameter(exOfIn).HasLimits() && gstep > 0.5 {
					gstep = 0.5
				}

				gr.set(i, grd)
				gr2.set(i, g2)
				gst.set(i, gstep)
			}

			return NewFunctionGradientFromMnAlgebraicVectors(gr, gr2, gst), nil
		}
	}
}

func (this *InitialGradientCalculator) gradientWithGradient(par *MinimumParameters,
	gra *FunctionGradient) (*FunctionGradient, error) {
	g, err := this.gradient(par)
	return g, err
}

func (this *InitialGradientCalculator) trafo() *MnUserTransformation {
	return this.theTransformation
}

func (this *InitialGradientCalculator) fcn() MnFcnInterface {
	return this.theFcn
}

func (this *InitialGradientCalculator) precision() *MnMachinePrecision {
	return this.theTransformation.precision()
}

func (this *InitialGradientCalculator) strategy() *MnStrategy {
	return this.theStrategy
}

func (this *InitialGradientCalculator) ncycle() int {
	return this.strategy().GradientNCycles()
}

func (this *InitialGradientCalculator) stepTolerance() float64 {
	return this.strategy().GradientStepTolerance()
}

func (this *InitialGradientCalculator) gradTolerance() float64 {
	return this.strategy().GradientTolerance()
}
