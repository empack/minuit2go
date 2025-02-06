package minuit

import "log"

type MnNotPosDef struct{}
type MnMadePosDef struct{}
type MnHesseFailed struct{}
type MnInvertFailed struct{}

type MinimumError struct {
	theMatrix       *MnAlgebraicSymMatrix
	theDCovar       float64
	theValid        bool
	thePosDef       bool
	theMadePosDef   bool
	theHesseFailed  bool
	theInvertFailed bool
	theAvailable    bool
}

func NewMinimumErrorFromNumber(n int) (*MinimumError, error) {
	m, err := NewMnAlgebraicSymMatrix(n)
	if err != nil {
		return nil, err
	}
	return &MinimumError{
		theMatrix: m,
		theDCovar: 1.0,
	}, nil
}

func NewMinimumError(mat *MnAlgebraicSymMatrix, dcov float64) *MinimumError {
	return &MinimumError{
		theMatrix:    mat,
		theDCovar:    dcov,
		theValid:     true,
		thePosDef:    true,
		theAvailable: true,
	}
}

func NewMinimumErrorFromHesse(mat *MnAlgebraicSymMatrix, _ *MnHesseFailed) *MinimumError {
	return &MinimumError{
		theMatrix:       mat,
		theDCovar:       1.0,
		theValid:        false,
		thePosDef:       false,
		theMadePosDef:   false,
		theHesseFailed:  true,
		theInvertFailed: false,
		theAvailable:    true,
	}
}

func NewMinimumErrorFromMnMadePosDef(mat *MnAlgebraicSymMatrix, _ *MnMadePosDef) *MinimumError {
	return &MinimumError{
		theMatrix:       mat,
		theDCovar:       1.0,
		theValid:        false,
		thePosDef:       false,
		theMadePosDef:   true,
		theHesseFailed:  false,
		theInvertFailed: false,
		theAvailable:    true,
	}
}

func NewMinimumErrorFromMnInvertFailed(mat *MnAlgebraicSymMatrix, _ *MnInvertFailed) *MinimumError {
	return &MinimumError{
		theMatrix:       mat,
		theDCovar:       1.0,
		theValid:        false,
		thePosDef:       true,
		theMadePosDef:   false,
		theHesseFailed:  false,
		theInvertFailed: true,
		theAvailable:    true,
	}
}

func NewMinimumErrorFromMnNotPosDef(mat *MnAlgebraicSymMatrix, _ *MnNotPosDef) *MinimumError {
	return &MinimumError{
		theMatrix:       mat,
		theDCovar:       1.,
		theValid:        false,
		thePosDef:       false,
		theMadePosDef:   false,
		theHesseFailed:  false,
		theInvertFailed: false,
		theAvailable:    true,
	}
}

func (this *MinimumError) matrix() *MnAlgebraicSymMatrix {
	return MnUtils.MulSM(this.theMatrix, 2)
}

func (this *MinimumError) invHessian() *MnAlgebraicSymMatrix {
	return this.theMatrix
}

func (this *MinimumError) hessian() (*MnAlgebraicSymMatrix, error) {
	var tmp *MnAlgebraicSymMatrix = this.theMatrix.Clone()
	err := tmp.invert()
	if err == nil {
		return tmp, nil
	} else {
		log.Println("BasicMinimumError inversion fails; return diagonal matrix.")
		tmp, err := NewMnAlgebraicSymMatrix(this.theMatrix.nrow())
		if err != nil {
			return nil, err
		}
		for i := 0; i < this.theMatrix.nrow(); i++ {
			v, err := this.theMatrix.get(i, i)
			if err != nil {
				return nil, err
			}
			err = tmp.set(i, i, 1.0/v)
			if err != nil {
				return nil, err
			}
		}
		return tmp, nil
	}
}

func (this *MinimumError) dcovar() float64 {
	return this.theDCovar
}

func (this *MinimumError) isAccurate() bool {
	return this.theDCovar < 0.1
}

func (this *MinimumError) isValid() bool {
	return this.theValid
}

func (this *MinimumError) isPosDef() bool {
	return this.thePosDef
}

func (this *MinimumError) isMadePosDef() bool {
	return this.theMadePosDef
}

func (this *MinimumError) hesseFailed() bool {
	return this.theHesseFailed
}

func (this *MinimumError) invertFailed() bool {
	return this.theInvertFailed
}

func (this *MinimumError) isAvailable() bool {
	return this.theAvailable
}
