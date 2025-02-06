package minuit

import (
	"errors"
	"fmt"
	"math"
)

// TODO: MnHesse.java needs full rewrite
type MnHesse struct {
	theStrategy *MnStrategy
}

func NewMnHesse() *MnHesse {
	return &MnHesse{theStrategy: NewMnStrategyWithStra(1)}
}

func NewMnHesseWithIntStrategy(stra int) *MnHesse {
	return &MnHesse{theStrategy: NewMnStrategyWithStra(stra)}
}

func NewMnHesseWithStrategy(stra *MnStrategy) *MnHesse {
	return &MnHesse{theStrategy: stra}
}

func (this *MnHesse) CalculateWithFcnParErr(fcn FCNBase, par, err []float64) *MnUserParameterState {
	return this.CalculateWithFcnParErrMaxcalls(fcn, NewMnUserParameterStateWithParErr(par, err), 0)
}

func (this *MnHesse) CalculateWithFcnParErrMaxcalls(fcn FCNBase, par, err []float64,
	maxcalls int) *MnUserParameterState {
	return this.CalculateWithFcnStateMaxcalls(fcn, NewMnUserParameterStateWithParErr(par, err), maxcalls)
}

func (this *MnHesse) CalculateWithFcnParCovariance(fcn FCNBase, par []float64,
	cov *MnUserCovariance) *MnUserParameterState {
	return this.CalculateWithFcnParCovarianceMaxcalls(fcn, par, cov, 0)
}

func (this *MnHesse) CalculateWithFcnParCovarianceMaxcalls(fcn FCNBase, par []float64, cov *MnUserCovariance,
	maxcalls int) *MnUserParameterState {
	return this.CalculateWithFcnStateMaxcalls(fcn, NewUserParameterStateWithParCov(par, cov), maxcalls)
}

func (this *MnHesse) CalculateWithFcnPar(fcn FCNBase, par *MnUserParameters) *MnUserParameterState {
	return this.CalculateWithFcnParMaxcalls(fcn, par, 0)
}

func (this *MnHesse) CalculateWithFcnParMaxcalls(fcn FCNBase, par *MnuserParameters,
	maxcalls int) *MnUserParameterState {
	return this.CalculateWithFcnStateMaxcalls(fcn, NewMnUserParameterStateWithPar(par), maxcalls)
}

func (this *MnHesse) CalculateWithFcnParCovMaxcalls(fcn FCNBase, par *MnUserParameters,
	cov *MnUserCovariance, maxcalls int) *MnUserParameterState {
	return this.CalculateWithFcnStateMaxcalls(fcn, NewMnUserParameterStateWithParCov(par, cov), maxcalls)
}

func (this *MnHesse) CalculateWithFcnStateMaxcalls(fcn FCNBase, state *MnUserParameterState,
	maxcalls int) (*MnUserParameterState, error) {
	errDef := 1.0
	n := state.VariableParameters()
	mfcn := NewMnUserFcn(fcn, errDef, state.trafo())
	x := NewMnAlgebraicVector(n)

	for i := 0; i < n; i++ {
		x.set(i, state.intParameters()[i])
	}

	amin := mfcn.valueOf(x)
	gc := NewNumerical2PGradientCalculator(mfcn, state.trafo(), this.theStrategy)
	par := NewMinimumParameters(x, amin)
	gra := gc.GradientWithPar(par)
	symmatrix, err := NewMnAlgebraicSymMatrix(n)
	if err != nil {
		return nil, err
	}
	tmp := this.CalculateWithMnfcnStTrafoMaxcalls(mfcn, NewMinimumStateWithGrad(par, NewMinimumError(symmatrix,
		1.0), gra, state.Edm(), state.Nfcn()), state.trafo(), maxcalls)

	return NewMnuserParamterStateWithStateErrdefTrafo(tmp, errDef, state.trafo()), nil
}

// TODO: MnHesse.java:76 needs full rewrite here
func (this *MnHesse) CalculateWithMnfcnStTrafoMaxcalls(mfcn *MnFcn, st *MinimumState, trafo *MnUserTransformation,
	maxcalls int) (*MinimumState, error) {
	prec := trafo.precision()
	amin := mfcn.valueOf(st.vec())
	aimsag := math.Sqrt(prec.eps2()) * (math.Abs(amin) + mfcn.errorDef())
	n := st.parameters().vec().size()
	if maxcalls == 0 {
		maxcalls = 200 + 100 * n + 5 * n * n
	}

	vhmat, err := NewMnAlgebraicSymMatrix(n)
	if err != nil {
		return nil, err
	}
	g2 := st.gradient().g2().Clone()
	gst := st.gradient().gstep().Clone()
	grd := st.gradient().grad().Clone()
	dirin := st.gradient().gstep().Clone()
	yy := NewMnAlgebraicVector(n)
	if st.gradient().isAnalytical() {
		ifc := NewInitialGradientCalculator(fcn, trafo, this.theStrategy)
		tmp := igc.gradient(st.parameters())
		gst = tmp.gstep().Clone()
		dirin = tmp.gstep().Clone()
		g2 = tmp.g2().Clone()
	}

	// try catch block
	var tryToCatch bool
	x := st.parameters().vec().Clone()

	for i := 0; i < n; i++ {
		xtf := x.get(i)
		dmin := 8.0 * prec.eps2() * (math.Abs(xtf) + prec.eps2())
		d := math.Abs(gst.get(i))
		if d < dmin {
			d = dmin
		}

		for icyc := 0; icyc < this.ncycles(); icyc++ {
			sag, fs1, fs2 := 0.0, 0.0, 0.0

			var multpy int
			for multpy = 0; multpy < 5; multpy++ {
				x.set(i, xtf + d)
				fs1 = mfcn.valueOf(x)
				x.set(i, xtf - d)
				fs2 = mfcn.valueOf(x)
				x.set(i, xtf)
				sag = 0.5 * (fs1 + fs2 - 2.0 * amin)
				if sag > prec.eps2() {
					break
				}

				if trafo.parameter(i).HasLimits() {
					if d > 0.5 {
						// error catch handling
						// throw new MnHesseFailed("MnHesse: 2nd derivative zero for parameter");
					}

					d *= 10.0
					if d > 0.5 {
						d = 0.51
					}
				} else {
					d *= 10.0
				}
			}

			if multpy >= 5 {
				// error catch handling
				// throw new MnHesseFailed("MnHesse: 2nd derivative zero for parameter");
			}

			g2bfor := g2.get(i)
			g2.set(i, 2.0 * sag / (d * d))
			grd.set(i, (fs1 - fs2) / (2.0 * d))
			gst.set(i, d)
			dirin.set(i, d)
			yy.set(i, fs1)
			d = math.Sqrt(2.0 * aimsag / math.Abs(g2.get(i)))
			if trafo.parameter(i).HasLimits() {
				d = math.Min(0.5, d)
			}

			if d < dmin {
				d = dmin
			}

			if math.Abs((d - d) / d) < this.tolerstp() ||math.Abs((g2.get(i) - g2bfor) / g2.get(i)) < this.tolerg2() {
				break
			}

			var53 := math.Min(d, 10.5 * d)
			d = math.Max(var53, 0.1 * d)
		}

		_ = vhmat.set(i, i, g2.get(i))
		if mfcn.numOfCalls() - st.nfcn() > maxcalls {
			// error catch handling
			// throw new MnHesseFailed("MnHesse: maximum number of allowed function calls exhausted.");
		}
	}

	if this.theStrategy.strategy() > 0 {
		hgc := NewHessianGradientCalculator(mfcn, trafo, this.theStrategy)
		gr := hgc.GradientWithGradient(st.parameters(), NewFunctionGradientFromMnAlgebraicVectors(grd, g2, gst))
		grd = gr.grad()
	}

	for i := 0; i < n; i++ {
		x.set(i, x.get(i) + dirin.get(i))

		for j := i + 1; j < n; j++ {
			x.set(j, x.get(j) + dirin.get(j))
			fs1 := mfcn.valueOf(x)
			elem := (fs1 + amin - yy.get(i) - yy.get(j)) / (dirin.get(i) * dirin.get(j))
			_ = vhmat.set(i, j, elem)
			x.set(j, x.get(j) - dirin.get(i))
		}

		x.set(i, x.get(i) - dirin.get(j))
	}

	tmp := MnPosDef.TestError(NewMinimumError(vhmat, 1.0), prec)
	vhmat = tmp.invHessian()

	// try catch block
	failedHesse := vhmat.invert()
	if failedHesse != nil {
		// error catch handling
		// throw new MnHesseFailed("MnHesse: matrix inversion fails!");
		fmt.Println("MnHesse: matrix inversion fails!")
	}

	gr := NewFunctionGradientFromMnAlgebraicVectors(grd, g2, gst)
	if tmp.isMadePosDef() {
		fmt.Println("MnHesse: matrix is invalid!")
		fmt.Println("MnHesse: matrix is not pos. def.!")
		fmt.Println("MnHesse: matrix was forced pos. def.")
		// TODO: MnHesse.java:208 this is tricky
		return NewMinimumStateWithGrad(st.parameters(), NewMinimumErrorFromMnMadePosDef(vhmat,
			THIS IS UNKNOWN), gr, st.edm(), mfcn.numOfCalls()), nil
	} else {
		e = NewMinimumError(vhmat, 0.0)
		edm, err := NewVariableMetricEDMEstimator().estimate(gr, err)
		if err != nil {
			return nil, err
		}
		// TODO: MnHesse.java:223 this is tricky
		return NewMinimumStateWithGrad(st.parameters(), e, gr, edm, mfcn.numOfCalls()), nil
	}

	if tryToCatch {
		fmt.Println(tryToCatch)
		fmt.Println("MnHesse fails and will return diagonal matrix")

		for j := 0; j < n; j++ {
			var tmp float64
			if g2.get(j) < prec.eps2() {
				tmp = 1.0
			} else {
				tmp = 1.0 / g2.get(j)
			}
			if tmp < prec.eps2() {
				vhmat.set(j, j, 1.0)
			} else {
				vhmat.set(j, j, tmp)
			}

			// TODO: MnHesse.java:223 this is tricky
			return NewMinimumStateWithGrad(st.parameters(), NewMinimumErrorFromHesse(vham, THIS IS UNKNOWN),
				st.gradient(), st.edm(), st.nfcn() + mfcn.numOfCalls()), nil
		}
	}

	return nil, errors.New("this should not happen")
}

func (this *MnHesse) ncycles() int {
	return this.theStrategy.HessianNCycles()
}

func (this *MnHesse) tolerstp() float64 {
	return this.theStrategy.HessianStepTolerance()
}

func (this *MnHesse) tolerg2() float64 {
	return this.theStrategy.hessianG2Tolerance()
}