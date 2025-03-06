package minuit

import (
	"context"
	"log"
	"math"
)

type SimplexBuilder struct {
}

func NewSimplexBuilder() *SimplexBuilder {
	return &SimplexBuilder{}
}

// TODO: ADD BREAKPOINT
func (this *SimplexBuilder) Minimum(ctx context.Context, mfcn MnFcnInterface, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy, maxfcn int, minedm float64) (*FunctionMinimum, error) {
	var prec *MnMachinePrecision = seed.precision()
	var x *MnAlgebraicVector = seed.parameters().vec().Clone()
	var step *MnAlgebraicVector = MnUtils.MulV(seed.gradient().gstep(), 10.0)

	var n int = x.size()
	var wg float64 = 1.0 / float64(n)
	var alpha float64 = 1.0
	var beta float64 = 0.5
	var gamma float64 = 2.0
	var rhomin float64 = 4.0
	var rhomax float64 = 8.0
	var rho1 float64 = 1.0 + alpha
	var rho2 float64 = 1.0 + alpha*gamma

	var simpl []*Pair[float64, *MnAlgebraicVector] = make([]*Pair[float64, *MnAlgebraicVector], 0, n+1)
	simpl = append(simpl, NewPair[float64, *MnAlgebraicVector](seed.fval(), x.Clone()))

	var jl int = 0
	var jh int = 0
	var amin float64 = seed.fval()
	var aming float64 = seed.fval()

	for i := 0; i < n; i++ {
		log.Println("SimplexBuilder.Minimum: i=", i)
		var dmin float64 = 8. * prec.eps2() * (math.Abs(x.get(i)) + prec.eps2())
		if step.get(i) < dmin {
			step.set(i, dmin)
		}
		x.set(i, x.get(i)+step.get(i))
		var tmp float64 = mfcn.valueOf(x)
		if tmp < amin {
			amin = tmp
			jl = i + 1
		}
		if tmp > aming {
			aming = tmp
			jh = i + 1
		}
		simpl = append(simpl, NewPair[float64, *MnAlgebraicVector](tmp, x.Clone()))
		x.set(i, x.get(i)-step.get(i))
	}
	var simplex *SimplexParameters = NewSimplexParameters(simpl, jh, jl)
	ok := true

	stopped := false

	for ok {
		select {
		case <-ctx.Done():
			stopped = true
			ok = false
		default:
			amin = simplex.get(jl).First
			jl = simplex.jl()
			jh = simplex.jh()
			var pbar *MnAlgebraicVector = NewMnAlgebraicVector(n)
			for i := 0; i < n+1; i++ {
				if i == jh {
					continue
				}
				if res, err := MnUtils.AddV(pbar, MnUtils.MulV(simplex.get(i).Second, wg)); err == nil {
					pbar = res
				} else {
					return nil, err
				}
			}

			pstar, err := MnUtils.SubV(MnUtils.MulV(pbar, 1.+alpha), MnUtils.MulV(simplex.get(jh).Second, alpha))
			if err != nil {
				return nil, err
			}
			var ystar float64 = mfcn.valueOf(pstar)

			if ystar > amin {
				if ystar < simplex.get(jh).First {
					simplex.update(ystar, pstar)
					if jh != simplex.jh() {
						continue
					}
				}
				pstst, err := MnUtils.AddV(MnUtils.MulV(simplex.get(jh).Second, beta), MnUtils.MulV(pbar, 1.-beta))
				if err != nil {
					return nil, err
				}
				var ystst float64 = mfcn.valueOf(pstst)
				if ystst > simplex.get(jh).First {
					break
				}
				simplex.update(ystst, pstst)
				continue
			}

			pstst, err := MnUtils.AddV(MnUtils.MulV(pstar, gamma), MnUtils.MulV(pbar, 1.-gamma))
			if err != nil {
				return nil, err
			}
			var ystst float64 = mfcn.valueOf(pstst)

			var y1 float64 = (ystar - simplex.get(jh).First) * rho2
			var y2 float64 = (ystst - simplex.get(jh).First) * rho1
			var rho float64 = 0.5 * (rho2*y1 - rho1*y2) / (y1 - y2)
			if rho < rhomin {
				if ystst < simplex.get(jl).First {
					simplex.update(ystst, pstst)
				} else {
					simplex.update(ystar, pstar)
				}
				continue
			}
			if rho > rhomax {
				rho = rhomax
			}
			prho, err := MnUtils.AddV(MnUtils.MulV(pbar, rho), MnUtils.MulV(simplex.get(jh).Second, 1.0-rho))
			if err != nil {
				return nil, err
			}
			var yrho float64 = mfcn.valueOf(prho)
			if yrho < simplex.get(jl).First && yrho < ystst {
				simplex.update(yrho, prho)
				continue
			}
			if ystst < simplex.get(jl).First {
				simplex.update(ystst, pstst)
				continue
			}
			if yrho > simplex.get(jl).First {
				if ystst < simplex.get(jl).First {
					simplex.update(ystst, pstst)
				} else {
					simplex.update(ystar, pstar)
				}
				continue
			}
			if ystar > simplex.get(jh).First {
				pstst, err = MnUtils.AddV(MnUtils.MulV(simplex.get(jh).Second, beta), MnUtils.MulV(pbar, 1-beta))
				if err != nil {
					return nil, err
				}
				ystst = mfcn.valueOf(pstst)
				if ystst > simplex.get(jh).First {
					break
				}
				simplex.update(ystst, pstst)
			}

			// replacement for do while loop
			ok = simplex.edm() > minedm && mfcn.numOfCalls() < maxfcn
		}
	}

	amin = simplex.get(jl).First
	jl = simplex.jl()
	jh = simplex.jh()

	var pbar *MnAlgebraicVector = NewMnAlgebraicVector(n)
	for i := 0; i < n+1; i++ {
		if i == jh {
			continue
		}
		var fnErr error
		pbar, fnErr = MnUtils.AddV(pbar, MnUtils.MulV(simplex.get(i).Second, wg))
		if fnErr != nil {
			return nil, fnErr
		}
	}
	var ybar float64 = mfcn.valueOf(pbar)
	if ybar < amin {
		simplex.update(ybar, pbar)
	} else {
		pbar = simplex.get(jl).Second
		ybar = simplex.get(jl).First
	}

	var dirin *MnAlgebraicVector = simplex.dirin()
	//   scale to sigmas on parameters werr^2 = dirin^2 * (up/edm)
	dirin = MnUtils.MulV(dirin, math.Sqrt(mfcn.errorDef()/simplex.edm()))

	st, fnErr := NewMinimumState(NewMinimumParametersFromMnAlgebraicVectors(pbar, dirin, ybar), simplex.edm(), mfcn.numOfCalls())
	if fnErr != nil {
		return nil, fnErr
	}
	var states []*MinimumState = make([]*MinimumState, 0, 1)
	states = append(states, st)

	if stopped {
		log.Println("SimplexBuilder: stopped")
		return NewFunctionMinimumWithSeedStatesUpStopped(seed, states, mfcn.errorDef()), nil
	}

	if mfcn.numOfCalls() > maxfcn {
		log.Println("Simplex did not converge, #fcn calls exhausted.")
		return NewFunctionMinimumWithSeedStatesUpReachedCallLimit(seed, states, mfcn.errorDef()), nil
	}
	if simplex.edm() > minedm {
		log.Println("Simplex did not converge, edm > minedm.")
		return NewFunctionMinimumWithSeedStatesUpAboveMaxEdm(seed, states, mfcn.errorDef()), nil
	}
	return NewFunctionMinimumWithSeedStatesUp(seed, states, mfcn.errorDef()), nil
}
