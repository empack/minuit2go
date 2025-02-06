package minuit

import "math"

/* Performs a minimization using the simplex method of Nelder and Mead
 * (ref. Comp. J. 7, 308 (1965)).
 */
type ScanBuilder struct {
}

func (s ScanBuilder) Minimum(mfcn *MnFcn, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, toler float64) (*FunctionMinimum, error) {
	var x *MnAlgebraicVector = seed.parameters().vec().Clone()
	var upst *MnUserParameterState = NewMnUserParameterState(seed.state(), mfcn.errorDef(), seed.trafo())
	var scan *MnParameterScan = NewMnParameterScan(mfcn.fcn(), upst.parameters(), seed.fval())
	var amin float64 = scan.fval()
	var n int = seed.trafo().variableParameters()
	var dirin *MnAlgebraicVector = NewMnAlgebraicVector(n)
	for i := 0; i < n; i++ {
		var ext int = seed.trafo().extOfInt(i)
		scan.scan(ext)
		if scan.fval() < amin {
			amin = scan.fval()
			x.set(i, seed.trafo().ext2int(ext, scan.parameters().value(ext)))
		}
		dirin.set(i, math.Sqrt(2.*mfcn.errorDef()*seed.error().invHessian().get(i, i)))
	}

	var mp *MinimumParameters = NewMinimumParametersFromMnAlgebraicVectors(x, dirin, amin)
	st, err := NewMinimumState(mp, 0., mfcn.numOfCalls())
	if err != nil {
		return nil, err
	}

	var states []*MinimumState = make([]*MinimumState, 0, 1)
	states = append(states, st)
	return NewFunctionMinimumWithSeedStatesUp(seed, states, mfcn.errorDef()), nil
}
