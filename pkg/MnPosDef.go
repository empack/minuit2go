package minuit

import (
	"log"
	"math"
)

var MnPosDef = &mnPosDefStruct{}

type mnPosDefStruct struct {
}

func (this *mnPosDefStruct) TestState(st *MinimumState, prec *MnMachinePrecision) (*MinimumState, error) {
	err, fnErr := this.TestError(st.error(), prec)
	if fnErr != nil {
		return nil, fnErr
	}
	return NewMinimumStateWithGrad(st.parameters(), err, st.gradient(), st.edm(), st.nfcn()), nil
}

func (this *mnPosDefStruct) TestError(e *MinimumError, prec *MnMachinePrecision) (*MinimumError, error) {
	var err *MnAlgebraicSymMatrix = e.invHessian().Clone()
	v_, fnErr := err.get(0, 0)
	if fnErr != nil {
		return nil, fnErr
	}
	if err.size() == 1 && v_ < prec.eps() {
		fnErr := err.set(0, 0, 1.)
		if fnErr != nil {
			return nil, fnErr
		}
		return NewMinimumErrorFromMnMadePosDef(err, &MnMadePosDef{}), nil
	}
	v_, fnErr = err.get(0, 0)
	if fnErr != nil {
		return nil, fnErr
	}
	if err.size() == 1 && v_ > prec.eps() {
		return e, nil
	}
	//   std::cout<<"MnPosDef init matrix= "<<err<<std::endl;
	var epspdf float64 = math.Max(1.e-6, prec.eps2())
	v_, fnErr = err.get(0, 0)
	if fnErr != nil {
		return nil, fnErr
	}
	var dgmin float64 = v_

	for i := 0; i < err.nrow(); i++ {
		v_, fnErr = err.get(i, i)
		if fnErr != nil {
			return nil, fnErr
		}
		if v_ < prec.eps2() {
			log.Printf("negative or zero diagonal element %d in covariance matrix\n", i)
		}
		if v_ < dgmin {
			dgmin = v_
		}
	}

	var dg float64 = 0.0
	if dgmin < prec.eps2() {
		dg = 1. + epspdf - dgmin
		//     dg = 0.5*(1. + epspdf - dgmin);
		log.Printf("added %f to diagonal of error matrix \n", dg)
	}

	var s *MnAlgebraicVector = NewMnAlgebraicVector(err.nrow())
	p, fnErr := NewMnAlgebraicSymMatrix(err.nrow())
	if fnErr != nil {
		return nil, fnErr
	}
	for i := 0; i < err.nrow(); i++ {
		v_, fnErr = err.get(i, i)
		if fnErr != nil {
			return nil, fnErr
		}
		fnErr = err.set(i, i, v_+dg)
		if fnErr != nil {
			return nil, fnErr
		}
		if v_ < 0. {
			fnErr = err.set(i, i, 1.)
			if fnErr != nil {
				return nil, fnErr
			}
		}
		s.set(i, 1./math.Sqrt(v_))
		for j := 0; j <= i; j++ {
			fnErr = p.set(i, j, v_*s.get(i)*s.get(j))
			if fnErr != nil {
				return nil, fnErr
			}
		}
	}

	//   std::cout<<"MnPosDef p: "<<p<<std::endl
	eval, fnErr := p.eigenvalues()
	if fnErr != nil {
		return nil, fnErr
	}
	var pmin float64 = eval.get(0)
	var pmax float64 = eval.get(eval.size() - 1)
	//   std::cout<<"pmin= "<<pmin<<" pmax= "<<pmax<<std::endl;
	pmax = math.Max(math.Abs(pmax), 1.)
	if pmin > epspdf*pmax {
		return e, nil
	}

	var padd float64 = 0.001*pmax - pmin
	log.Println("eigenvalues: ")
	for i := 0; i < err.nrow(); i++ {
		v_, fnErr = err.get(i, i)
		if fnErr != nil {
			return nil, fnErr
		}
		fnErr = err.set(i, i, v_*(1.+padd))
		if fnErr != nil {
			return nil, fnErr
		}
		log.Printf("%5g\n", eval.get(i))
	}
	//   std::cout<<"MnPosDef final matrix: "<<err<<std::endl;
	log.Printf("matrix forced pos-def by adding %f to diagonal\n", padd)
	//   std::cout<<"eigenvalues: "<<eval<<std::endl;
	return NewMinimumErrorFromMnMadePosDef(err, &MnMadePosDef{}), nil
}
