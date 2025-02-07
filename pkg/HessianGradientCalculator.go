package minuit

import (
	"errors"
	"math"
)

type HessianGradientCalculator struct {
	theFcn            MnFcnInterface
	theTransformation *MnUserTransformation
	theStrategy       *MnStrategy
}

func NewHessianGradientCalculator(fcn MnFcnInterface, par *MnUserTransformation, stra *MnStrategy) *HessianGradientCalculator {
	return &HessianGradientCalculator{
		theFcn:            fcn,
		theTransformation: par,
		theStrategy:       stra,
	}
}

func (this *HessianGradientCalculator) Gradient(par *MinimumParameters) (*FunctionGradient, error) {
	var gc *InitialGradientCalculator = NewInitialGradientCalculator(this.theFcn, this.theTransformation, this.theStrategy)
	gra, err := gc.gradient(par)
	if err != nil {
		return nil, err
	}
	return this.GradientWithGrad(par, gra)
}

func (this *HessianGradientCalculator) GradientWithGrad(par *MinimumParameters, grad *FunctionGradient) (*FunctionGradient, error) {
	res, err := this.deltaGradient(par, grad)
	return res.First, err
}

func (this *HessianGradientCalculator) deltaGradient(par *MinimumParameters, gradient *FunctionGradient) (*Pair[*FunctionGradient, *MnAlgebraicVector], error) {
	if !par.isValid() {
		return nil, errors.New("IllegalArgumentException: parameters are invalid")
	}

	var x *MnAlgebraicVector = par.vec().Clone()
	var grd *MnAlgebraicVector = gradient.grad().Clone()
	var g2 *MnAlgebraicVector = gradient.g2()
	var gstep *MnAlgebraicVector = gradient.gstep()

	var fcnmin float64 = par.fval()
	//   std::cout<<"fval: "<<fcnmin<<std::endl;

	var dfmin float64 = 4. * this.precision().eps2() * (math.Abs(fcnmin) + this.theFcn.errorDef())

	var n int = x.size()
	var dgrd *MnAlgebraicVector = NewMnAlgebraicVector(n)

	// initial starting values
	for i := 0; i < n; i++ {
		var xtf float64 = x.get(i)
		var dmin float64 = 4. * this.precision().eps2() * (xtf + this.precision().eps2())
		var epspri float64 = this.precision().eps2() + math.Abs(grd.get(i)*this.precision().eps2())
		var optstp float64 = math.Sqrt(dfmin / (math.Abs(g2.get(i)) + epspri))
		var d float64 = 0.2 * math.Abs(gstep.get(i))
		if d > optstp {
			d = optstp
		}
		if d < dmin {
			d = dmin
		}
		var chgold float64 = 10000.0
		var dgmin float64 = 0.0
		var grdold float64 = 0.0
		var grdnew float64 = 0.0
		for j := 0; j < this.ncycle(); j++ {
			x.set(i, xtf+d)
			var fs1 float64 = this.theFcn.valueOf(x)
			x.set(i, xtf-d)
			var fs2 float64 = this.theFcn.valueOf(x)
			x.set(i, xtf)
			//       double sag = 0.5*(fs1+fs2-2.*fcnmin);
			grdold = grd.get(i)
			grdnew = (fs1 - fs2) / (2. * d)
			dgmin = this.precision().eps() * (math.Abs(fs1) + math.Abs(fs2)) / d
			if math.Abs(grdnew) < this.precision().eps() {
				break
			}
			var change float64 = math.Abs((grdold - grdnew) / grdnew)
			if change > chgold && j > 1 {
				break
			}
			chgold = change
			grd.set(i, grdnew)
			if change < 0.05 {
				break
			}
			if math.Abs(grdold-grdnew) < dgmin {
				break
			}
			if d < dmin {
				break
			}
			d *= 0.2
		}
		dgrd.set(i, math.Max(dgmin, math.Abs(grdold-grdnew)))
	}
	return NewPair[*FunctionGradient, *MnAlgebraicVector](NewFunctionGradientFromMnAlgebraicVectors(grd, g2, gstep), dgrd), nil
}

func (this *HessianGradientCalculator) fcn() MnFcnInterface {
	return this.theFcn
}

func (this *HessianGradientCalculator) trafo() *MnUserTransformation {
	return this.theTransformation
}

func (this *HessianGradientCalculator) precision() *MnMachinePrecision {
	return this.theTransformation.precision()
}

func (this *HessianGradientCalculator) strategy() *MnStrategy {
	return this.theStrategy
}

func (this *HessianGradientCalculator) ncycle() int {
	return this.strategy().HessianGradientNCycles()
}

func (this *HessianGradientCalculator) stepTolerance() float64 {
	return this.strategy().GradientStepTolerance()
}

func (this *HessianGradientCalculator) gradTolerance() float64 {
	return this.strategy().GradientTolerance()
}
