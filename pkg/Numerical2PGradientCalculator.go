package minuit

import (
	"errors"
	"math"
)

type Numerical2PGradientCalculator struct {
	theFcn            *MnFcn
	theTransformation *MnUserTransformation
	theStrategy       *MnStrategy
}

func NewNumerical2PGradientCalculator(fcn *MnFcn, par *MnUserTransformation,
	stra *MnStrategy) *Numerical2PGradientCalculator {
	return &Numerical2PGradientCalculator{
		theFcn:            fcn,
		theTransformation: par,
		theStrategy:       stra,
	}
}

func (this *Numerical2PGradientCalculator) Gradient(par *MinimumParameters) (*FunctionGradient, error) {
	gc := NewInitialGradientCalculator(this.theFcn, this.theTransformation, this.theStrategy)
	gra, err := gc.gradient(par)
	if err != nil {
		return nil, err
	}
	gwg, err := this.GradientWithGradient(par, gra)
	return gwg, err
}

func (this *Numerical2PGradientCalculator) GradientWithGradient(par *MinimumParameters,
	gradient *FunctionGradient) (*FunctionGradient, error) {
	if !par.isValid() {
		return nil, errors.New("parameters are invalid")
	} else {
		x := par.vec().Clone()
		fcnmin := par.fval()
		dfmin := 8.0 * this.precision().eps2() * (math.Abs(fcnmin) + this.theFcn.errorDef())
		vrysml := 8.0 * this.precision().eps() * this.precision().eps()
		n := x.size()
		grd := gradient.grad().Clone()
		g2 := gradient.g2().Clone()
		gstep := gradient.gstep().Clone()

		for i := 0; i < n; i++ {
			xtf := x.get(i)
			epspri := this.precision().eps2() + math.Abs(grd.get(i)*this.precision().eps2())
			stepb4 := 0.0

			for j := 0; j < this.ncycle(); j++ { // does this loop even make sense?
				optstp := math.Sqrt(dfmin / (math.Abs(g2.get(i)) + epspri))
				step := math.Max(optstp, math.Abs(0.1*gstep.get(i)))
				if this.trafo().parameter(this.trafo().extOfInt(i)).HasLimits() && step > 0.5 {
					step = 0.5
				}

				stpmax := 10.0 * math.Abs(gstep.get(i))
				if step > stpmax {
					step = stpmax
				}

				stpmin := math.Max(vrysml, 8.0*math.Abs(this.precision().eps2()*x.get(i)))
				if step < stpmin {
					step = stpmin
				}

				if math.Abs((step-stepb4)/step) < this.stepTolerance() {
					break
				}

				gstep.set(i, step)
				stepb4 = step
				x.set(i, xtf+step)
				fs1 := this.theFcn.valueOf(x)
				x.set(i, xtf-step)
				fs2 := this.theFcn.valueOf(x)
				x.set(i, xtf)
				grdb4 := grd.get(i)
				grd.set(i, 0.5*(fs1-fs2)/step)
				g2.set(i, (fs1+fs2-2.0*fcnmin)/step/step)
				if math.Abs(grdb4-grd.get(i))/math.Abs(grd.get(i))+dfmin/step < this.gradTolerance() {
					break
				}
			}
		}

		return NewFunctionGradientFromMnAlgebraicVectors(grd, g2, gstep), nil
	}

}

func (this *Numerical2PGradientCalculator) trafo() *MnUserTransformation {
	return this.theTransformation
}

func (this *Numerical2PGradientCalculator) fcn() *MnFcn {
	return this.theFcn
}

func (this *Numerical2PGradientCalculator) precision() *MnMachinePrecision {
	return this.theTransformation.precision()
}

func (this *Numerical2PGradientCalculator) strategy() *MnStrategy {
	return this.theStrategy
}

func (this *Numerical2PGradientCalculator) ncycle() int {
	return this.strategy().GradientNCycles()
}

func (this *Numerical2PGradientCalculator) stepTolerance() float64 {
	return this.strategy().GradientStepTolerance()
}

func (this *Numerical2PGradientCalculator) gradTolerance() float64 {
	return this.strategy().GradientTolerance()
}
