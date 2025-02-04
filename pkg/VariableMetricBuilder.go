package minuit

import (
	"errors"
	"fmt"
	"math"
)

type VariableMetricBuilder struct {
	theEstimator    *VariableMetricEDMEstimator
	theErrorUpdator *DavidonErrorUpdator
}

func NewVariableMetricBuilder() *VariableMetricBuilder {
	return &VariableMetricBuilder{
		theEstimator:    NewVariableMetricEDMEstimator(),
		theErrorUpdator: NewDavidonErrorUpdator(),
	}
}

func (this *VariableMetricBuilder) Minimum(fcn *MnFcn, gc GradientCalculator, seed *MinimumSeed, strategy *MnStrategy,
	maxfcn int,
	edmval float64) (*FunctionMinimum, error) {
	fmin, err := this.minimum(fcn, gc, seed, maxfcn, edmval)
	if err != nil {
		return nil, err
	}

	if strategy.strategy() == 2 || strategy.strategy() == 1 && fmin.Error().dcovar() > 0.05 {
		st := NewMnHesse(strategy).calculate(fcn, fmin.state(), fmin.seed().trafo(), 0)
		fmin.add(st)
	}

	if !fmin.IsValid() {
		fmt.Println("FunctionMinimum is invalid.")
	}

	return fmin, nil
}

func (this *VariableMetricBuilder) minimum(fcn *MnFcn, gc GradientCalculator, seed *MinimumSeed, maxfcn int,
	edmval float64) (*FunctionMinimum, error) {
	edmval *= 1.0e-4
	if seed.parameters().vec().size() == 0 {
		return NewFunctionMinimumWithSeedUp(seed, fcn.errorDef()), nil
	} else {
		prec := seed.precision()
		result := make([]*MinimumState, 8)
		edm := seed.state().edm()
		if edm < 0.0 {
			fmt.Println("VariableMetricBuilder: initial matrix not pos.def")
			if seed.error().isPosDef() {
				return nil, errors.New("something is wrong")
			} else {
				return NewFunctionMinimumWithSeedUp(seed, fcn.errorDef()), nil
			}
		} else {
			result = append(result, seed.state())
			edm *= 1.0 + 3.0*seed.error().dcovar()

			for {
				s0 := result[(len(result) - 1)]
				mvsm, err := MnUtils.MulVSM(s0.error().invHessian(), s0.gradient().vec())
				if err != nil {
					return nil, err
				}
				step := MnUtils.MulV(mvsm, -1.0)
				gdel, err := MnUtils.InnerProduct(step, s0.gradient().grad())
				if err != nil {
					return nil, err
				}
				if gdel > 0.0 {
					fmt.Println("VariableMetricBuilder: matrix not pos.def")
					fmt.Printf("gdel > 0: %f\n", gdel)
					s0 = MnPosDef.TestState(s0, prec)
					mvsm, err = MnUtils.MulVSM(s0.error().invHessian(), s0.gradient().vec())
					if err != nil {
						return nil, err
					}
					step = MnUtils.MulV(mvsm, -1.0)
					gdel, err = MnUtils.InnerProduct(step, s0.gradient().grad())
					if err != nil {
						return nil, err
					}
					fmt.Printf("gdel: %f", gdel)
					if gdel > 0.0 {
						result = append(result, s0)
						return NewFunctionMinimumWithSeedStatesUp(seed, result, fcn.errorDef()), nil
					}
				}

				pp := MnLineSearch.search(fcn, s0.parameters(), step, gdel, prec)
				if math.Abs(pp.y()-s0.fval()) < prec.eps() {
					fmt.Println("VariableMetricBuilder: no improvement")
					break
				}

				added, err := MnUtils.AddV(s0.vec(), MnUtils.MulV(step, pp.x()))
				if err != nil {
					return nil, err
				}
				p := NewMinimumParameters(added, pp.y())
				g := gc.GradientWithGrad(p, s0.gradient())
				edm, err = this.estimator().estimate(*g, *s0.error())
				if err != nil {
					return nil, err
				}
				if edm < 0.0 {
					fmt.Println("VariableMetricBuilder: matrix not pos.def")
					fmt.Println("edm < 0")
					s0 = MnPosDef.TestState(s0, prec)
					edm, err = this.estimator().estimate(*g, *s0.error())
					if err != nil {
						return nil, err
					}
					if edm < 0.0 {
						result = append(result, s0)
						return NewFunctionMinimumWithSeedStatesUp(seed, result, fcn.errorDef()), nil
					}
				}

				e, err := this.errorUpdator().Update(s0, p, g)
				if err != nil {
					return nil, err
				}
				result = append(result, NewMinimumStateWithGrad(p, e, g, edm, fcn.numOfCalls()))
				edm *= 1.0 + 3.0*e.dcovar()
				if edm > edmval && fcn.numOfCalls() < maxfcn {
					break
				}
			}

			if fcn.numOfCalls() >= maxfcn {
				fmt.Println("VariableMetricBuilder: call limit exceeded")
				return NewFunctionMinimumWithSeedStatesUpReachedCallLimit(seed, result, fcn.errorDef()), nil
			} else if edm > edmval {
				if edm < math.Abs(prec.eps2()*result[len(result)-1].fval()) {
					fmt.Println("VariableMetricBuilder: machine accuracy limits further improvement.")
					return NewFunctionMinimumWithSeedStatesUp(seed, result, fcn.errorDef()), nil
				} else if edm < 10.0*edmval {
					return NewFunctionMinimumWithSeedStatesUp(seed, result, fcn.errorDef()), nil
				} else {
					fmt.Println("VariableMetricBuilder: finishes without convergence.")
					fmt.Printf("VariableMetricBuilder: edm= %f requested: %f\n", edm, edmval)
					return NewFunctionMinimumWithSeedStatesUp(seed, result, fcn.errorDef()), nil
				}
			} else {
				return NewFunctionMinimumWithSeedStatesUp(seed, result, fcn.errorDef()), nil
			}
		}
	}
}

func (this *VariableMetricBuilder) estimator() *VariableMetricEDMEstimator {
	return this.theEstimator
}

func (this *VariableMetricBuilder) errorUpdator() *DavidonErrorUpdator {
	return this.theErrorUpdator
}
