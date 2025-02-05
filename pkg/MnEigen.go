package minuit

// MnEigen
/*
 *  Calculates and the eigenvalues of the user covariance matrix MnUserCovariance.
 */
var MnEigen = &mnEigenStruct{}

type mnEigenStruct struct{}

/* Calculate eigenvalues of the covariance matrix.
 * Will perform the calculation of the eigenvalues of the covariance matrix
 * and return the result in the form of a double array.
 * The eigenvalues are ordered from the smallest to the largest eigenvalue.
 */
func (this *mnEigenStruct) eigenvalues(covar *MnUserCovariance) ([]float64, error) {
	cov, err := NewMnAlgebraicSymMatrix(covar.Nrow())
	if err != nil {
		return nil, err
	}
	for i := 0; i < covar.Nrow(); i++ {
		for j := i; j < covar.Nrow(); j++ {
			err = cov.set(i, j, covar.Get(i, j))
			if err != nil {
				return nil, err
			}
		}
	}

	eigen, err := cov.eigenvalues()
	if err != nil {
		return nil, err
	}
	return eigen.data, nil
}
