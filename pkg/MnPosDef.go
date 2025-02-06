package minuit

import (
	"log"
	"math"
)

var MnPosDef = &mnPosDefStruct{}

type mnPosDefStruct struct {
}

func (this *mnPosDefStruct) TestState(st *MinimumState, prec *MnMachinePrecision) *MinimumState {
	var err *MinimumError = this.TestError(st.error(), prec)
	return NewMinimumStateWithGrad(st.parameters(), err, st.gradient(), st.edm(), st.nfcn())
}

func (this *mnPosDefStruct) TestError(e *MinimumError, prec *MnMachinePrecision) *MinimumError {
	var err *MnAlgebraicSymMatrix = e.invHessian().Clone()
	if err.size() == 1 && err.get(0, 0) < prec.eps() {
		err.set(0, 0, 1.)
		return NewMinimumErrorFromMnMadePosDef(err, &MnMadePosDef{})
	}
	if err.size() == 1 && err.get(0, 0) > prec.eps() {
		return e
	}
	//   std::cout<<"MnPosDef init matrix= "<<err<<std::endl;
	var epspdf float64 = math.Max(1.e-6, prec.eps2())
	var dgmin float64 = err.get(0, 0)

	for i := 0; i < err.nrow(); i++ {
		if err.get(i, i) < prec.eps2() {
			log.Printf("negative or zero diagonal element %d in covariance matrix\n", i)
		}
		if err.get(i, i) < dgmin {
			dgmin = err.get(i, i)
		}
	}

	var dg float64 = 0.0
	if dgmin < prec.eps2() {
		dg = 1. + epspdf - dgmin
		//     dg = 0.5*(1. + epspdf - dgmin);
		log.Printf("added %f to diagonal of error matrix \n", dg)
	}

	var s *MnAlgebraicVector = NewMnAlgebraicVector(err.nrow())
	var p *MnAlgebraicSymMatrix = NewMnAlgebraicSymMatrix(err.nrow())
	for i := 0; i < err.nrow(); i++ {
		err.set(i, i, err.get(i, i)+dg)
		if err.get(i, i) < 0. {
			err.set(i, i, 1.)
		}
		s.set(i, 1./math.Sqrt(err.get(i, i)))
		for j := 0; j <= i; j++ {
			p.set(i, j, err.get(i, j)*s.get(i)*s.get(j))
		}
	}

	//   std::cout<<"MnPosDef p: "<<p<<std::endl
	var eval *MnAlgebraicVector = p.igenvalues()
	var pmin float64 = eval.get(0)
	var pmax float64 = eval.get(eval.size() - 1)
	//   std::cout<<"pmin= "<<pmin<<" pmax= "<<pmax<<std::endl;
	pmax = math.Max(math.Abs(pmax), 1.)
	if pmin > epspdf*pmax {
		return e
	}

	var padd float64 = 0.001*pmax - pmin
	log.Println("eigenvalues: ")
	for i := 0; i < err.nrow(); i++ {
		err.set(i, i, err.get(i, i)*(1.+padd))
		log.Printf("%5g\n", eval.get(i))
	}
	//   std::cout<<"MnPosDef final matrix: "<<err<<std::endl;
	log.Printf("matrix forced pos-def by adding %f to diagonal\n", padd)
	//   std::cout<<"eigenvalues: "<<eval<<std::endl;
	return NewMinimumErrorFromMnMadePosDef(err, &MnMadePosDef{})
}
