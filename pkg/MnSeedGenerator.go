package minuit

import (
	"errors"
	"log"
	"math"
)

type MnSeedGenerator struct {
}

func NewMnSeedGenerator() *MnSeedGenerator {
	return &MnSeedGenerator{}
}

func (this *MnSeedGenerator) Generate(fcn MnFcnInterface, gc GradientCalculator, st *MnUserParameterState, stra *MnStrategy) (*MinimumSeed, error) {
	var n int = st.VariableParameters()
	var prec *MnMachinePrecision = st.Precision()

	// initial starting values
	var x *MnAlgebraicVector = NewMnAlgebraicVector(n)
	for i := 0; i < n; i++ {
		x.set(i, st.intParameters()[i])
	}
	var fcnmin float64 = fcn.valueOf(x)
	var pa *MinimumParameters = NewMinimumParameters(x, fcnmin)

	var dgrad *FunctionGradient
	if _, ok := interface{}(gc).(AnalyticalGradientCalculator); ok {
		var igc *InitialGradientCalculator = NewInitialGradientCalculator(fcn, st.trafo(), stra)
		tmp, err := igc.gradient(pa)
		if err != nil {
			return nil, err
		}
		grd, err := gc.Gradient(pa)
		if err != nil {
			return nil, err
		}
		dgrad = NewFunctionGradientFromMnAlgebraicVectors(grd.grad(), tmp.g2(), tmp.gstep())

		if gc.(*AnalyticalGradientCalculator).checkGradient() {
			var good bool = true
			var hgc *HessianGradientCalculator = NewHessianGradientCalculator(fcn, st.trafo(), NewMnStrategyWithStra(2))
			hgrd, fnErr := hgc.deltaGradient(pa, dgrad)
			if fnErr != nil {
				return nil, fnErr
			}
			for i := 0; i < n; i++ {
				if math.Abs(hgrd.First.grad().get(i)-grd.grad().get(i)) > hgrd.Second.get(i) {
					log.Printf("gradient discrepancy of external parameter %d (internal parameter %d) too large.\n", st.trafo().extOfInt(i), i)
					good = false
				}
			}
			if !good {
				log.Println("Minuit does not accept user specified gradient. To force acceptance, override 'virtual bool checkGradient() const' of FCNGradientBase.h in the derived class.")
				return nil, errors.New("assertion violation: asserted 'good' to be true but was false")
			}
		}
	} else {
		var err error
		dgrad, err = gc.Gradient(pa)
		if err != nil {
			return nil, err
		}
	}
	mat, err := NewMnAlgebraicSymMatrix(n)
	if err != nil {
		return nil, err
	}
	var dcovar float64 = 1.0
	if st.HasCovariance() {
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				err = mat.set(i, j, st.intCovariance().Get(i, j))
				if err != nil {
					return nil, err
				}
			}
		}
		dcovar = 0.0
	} else {
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
	}

	var minErr *MinimumError = NewMinimumError(mat, dcovar)
	edm, err := NewVariableMetricEDMEstimator().estimate(dgrad, minErr)
	if err != nil {
		return nil, err
	}
	var state *MinimumState = NewMinimumStateWithGrad(pa, minErr, dgrad, edm, fcn.numOfCalls())

	if NegativeG2LineSearch.hasNegativeG2(dgrad, prec) {
		if _, ok := interface{}(gc).(AnalyticalGradientCalculator); ok {
			var ngc *Numerical2PGradientCalculator = NewNumerical2PGradientCalculator(fcn, st.trafo(), stra)
			state, err = NegativeG2LineSearch.search(fcn, state, ngc, prec)
			if err != nil {
				return nil, err
			}
		} else {
			state, err = NegativeG2LineSearch.search(fcn, state, gc, prec)
			if err != nil {
				return nil, err
			}
		}
	}

	if stra.Strategy() == 2 && !st.HasCovariance() {
		//calculate full 2nd derivative
		tmp, err := NewMnHesseWithStrategy(stra).CalculateWithMnfcnStTrafoMaxcalls(fcn, state, st.trafo(), 0)
		if err != nil {
			return nil, err
		}
		return NewMinimumSeed(tmp, st.trafo()), nil
	}

	return NewMinimumSeed(state, st.trafo()), nil
}
