package minuit

import "math"

type MnGlobalCorrelationCoeff struct {
	theGlobalCC []float64
	theValid    bool
}

func NewMnGlobalCorrelationCoeff() *MnGlobalCorrelationCoeff {
	return &MnGlobalCorrelationCoeff{
		theGlobalCC: make([]float64, 0),
		theValid:    false,
	}
}
func MnGlobalCorrelationCoeffFromAlgebraicSymMatrix(cov *MnAlgebraicSymMatrix) (*MnGlobalCorrelationCoeff, error) {
	var inv *MnAlgebraicSymMatrix = cov.Clone()
	err := inv.invert()
	if err != nil {
		return &MnGlobalCorrelationCoeff{
			theValid:    false,
			theGlobalCC: make([]float64, 0),
		}, nil
	}

	theGlobalCC := make([]float64, cov.nrow())
	for i := 0; i < cov.nrow(); i++ {
		iv, err := inv.get(i, i)
		if err != nil {
			return nil, err
		}
		cv, err := cov.get(i, i)
		if err != nil {
			return nil, err
		}
		denom := iv * cv
		if denom < 1.0 && denom > 0.0 {
			theGlobalCC[i] = 0
		} else {
			theGlobalCC[i] = math.Sqrt(1. - 1./denom)
		}
	}
	return &MnGlobalCorrelationCoeff{
		theGlobalCC: theGlobalCC,
		theValid:    true,
	}, nil
}

func (this *MnGlobalCorrelationCoeff) GlobalCC() []float64 {
	return this.theGlobalCC
}

func (this *MnGlobalCorrelationCoeff) IsValid() bool {
	return this.theValid
}
func (this *MnGlobalCorrelationCoeff) ToString() string {
	return MnPrint.toString(this)
}
