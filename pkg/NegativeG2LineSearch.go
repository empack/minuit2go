package minuit

import "math"

// NegativeG2LineSearch
/* In case that one of the components of the second derivative g2 calculated
 * by the numerical gradient calculator is negative, a 1dim line search in
 * the direction of that component is done in order to find a better position
 * where g2 is again positive.
 */
var NegativeG2LineSearch *negativeG2LineSearchStruct = &negativeG2LineSearchStruct{}

type negativeG2LineSearchStruct struct {
}

func (this *negativeG2LineSearchStruct) search(fcn MnFcnInterface, st *MinimumState, gc GradientCalculator, prec *MnMachinePrecision) (*MinimumState, error) {
	var negG2 bool = this.hasNegativeG2(st.gradient(), prec)
	if !negG2 {
		return st, nil
	}

	var n int = st.parameters().vec().size()
	var dgrad *FunctionGradient = st.gradient()
	var pa *MinimumParameters = st.parameters()
	var iterate bool = false
	var iter int = 0
	for ok := true; ok; ok = ok {
		iterate = false
		for i := 0; i < n; i++ {
			if dgrad.g2().get(i) < prec.eps2() {
				// do line search if second derivative negative
				var step *MnAlgebraicVector = NewMnAlgebraicVector(n)
				step.set(i, dgrad.gstep().get(i)*dgrad.vec().get(i))
				if math.Abs(dgrad.vec().get(i)) > prec.eps2() {
					step.set(i, step.get(i)*(-1./math.Abs(dgrad.vec().get(i))))
				}
				var gdel float64 = step.get(i) * dgrad.vec().get(i)
				pp, fnErr := MnLineSearch.search(fcn, pa, step, gdel, prec)
				if fnErr != nil {
					return nil, fnErr
				}
				step = MnUtils.MulV(step, pp.x())
				v_, err := MnUtils.AddV(pa.vec(), step)
				if err != nil {
					return nil, err
				}
				pa = NewMinimumParameters(v_, pp.y())
				dgrad, err = gc.GradientWithGrad(pa, dgrad)
				if err != nil {
					return nil, err
				}
				iterate = true
				break
			}
		}
		ok = iter < 2*n && iterate
		iter += 1
	}

	mat, err := NewMnAlgebraicSymMatrix(n)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		if math.Abs(dgrad.g2().get(i)) > prec.eps2() {
			err = mat.set(i, i, 1.0/dgrad.g2().get(i))
			if err != nil {
				return nil, err
			}
		} else {
			err = mat.set(i, i, 1.0)
			if err != nil {
				return nil, err
			}
		}
	}

	var minErr *MinimumError = NewMinimumError(mat, 1.)
	edm, err := NewVariableMetricEDMEstimator().estimate(dgrad, minErr)
	if err != nil {
		return nil, err
	}
	return NewMinimumStateWithGrad(pa, minErr, dgrad, edm, fcn.numOfCalls()), nil
}

func (this *negativeG2LineSearchStruct) hasNegativeG2(grad *FunctionGradient, prec *MnMachinePrecision) bool {
	for i := 0; i < grad.vec().size(); i++ {
		if grad.g2().get(i) < prec.eps2() {
			return true
		}
	}
	return false
}
