package minuit

import (
	"errors"
	"log"
)

var MnCovarianceSqueeze = &mnCovarianceSqueezeStruct{}

type mnCovarianceSqueezeStruct struct{}

func (this *mnCovarianceSqueezeStruct) Squeeze(cov *MnUserCovariance, n int) (*MnUserCovariance, error) {
	if cov.Nrow() <= 0 {
		return nil, errors.New("assertion violation: n <= 0")
	}
	if n >= cov.Nrow() {
		return nil, errors.New("assertion violation: n has to be smaller than cov.Nrow()")
	}

	hess, err := NewMnAlgebraicSymMatrix(cov.Nrow())
	if err != nil {
		return nil, err
	}
	for i := 0; i < cov.Nrow(); i++ {
		for j := i; j < cov.Nrow(); j++ {
			err = hess.set(i, j, cov.Get(i, j))
			if err != nil {
				return nil, err
			}
		}
	}

	err = hess.invert()
	if err != nil {
		log.Println("MnUserCovariance inversion failed; return diagonal matrix;")
		var result *MnUserCovariance = NewMnUserCovarianceWithNrow(cov.Nrow() - 1)
		for i, j := 0, 0; i < cov.Nrow(); i++ {
			if i == n {
				continue
			}
			result.Set(j, j, cov.Get(i, i))
			j++
		}
		return result, nil
	}
	squeezed, err := this.SqueezeSymMatrix(hess, n)
	if err != nil {
		return nil, err
	}
	err = squeezed.invert()
	if err != nil {
		log.Println("MnUserCovariance back-inversion failed; return diagonal matrix;")
		var result *MnUserCovariance = NewMnUserCovarianceWithNrow(squeezed.nrow())
		for i := 0; i < squeezed.nrow(); i++ {
			v, err := squeezed.get(i, i)
			if err != nil {
				return nil, err
			}
			result.Set(i, i, 1./v)
		}
		return result, nil
	}
	return NewMnUserCovarianceWithDataNrow(squeezed.data(), squeezed.nrow()), nil
}

func (this *mnCovarianceSqueezeStruct) squeeze(err *MinimumError, n int) (*MinimumError, error) {
	hess, fnErr := err.hessian()
	if fnErr != nil {
		return nil, fnErr
	}
	squeezed, fnErr := this.SqueezeSymMatrix(hess, n)
	if fnErr != nil {
		return nil, fnErr
	}

	fnErr = squeezed.invert()
	if fnErr != nil {
		log.Println("MnCovarianceSqueeze: MinimumError inversion fails; return diagonal matrix.")
		tmp, fnErr := NewMnAlgebraicSymMatrix(squeezed.nrow())
		if fnErr != nil {
			return nil, fnErr
		}
		for i := 0; i < squeezed.nrow(); i++ {
			v, fnErr := squeezed.get(i, i)
			if fnErr != nil {
				return nil, fnErr
			}
			fnErr = tmp.set(i, i, 1./v)
			if fnErr != nil {
				return nil, fnErr
			}
		}
		return NewMinimumErrorFromMnInvertFailed(tmp, &MnInvertFailed{}), nil
	}

	return NewMinimumError(squeezed, err.dcovar()), nil
}

func (this *mnCovarianceSqueezeStruct) SqueezeSymMatrix(hess *MnAlgebraicSymMatrix, n int) (*MnAlgebraicSymMatrix, error) {
	if hess.nrow() <= 0 {
		return nil, errors.New("assertion violation: n <= 0")
	}
	if n >= hess.nrow() {
		return nil, errors.New("assertion violation: n has to be smaler than hess.Nrow()")
	}

	hs, err := NewMnAlgebraicSymMatrix(hess.nrow() - 1)
	if err != nil {
		return nil, err
	}
	for i, j := 0, 0; i < hess.nrow(); i++ {
		if i == n {
			continue
		}
		for k, l := i, j; k < hess.nrow(); k++ {
			if k == n {
				continue
			}
			v, err := hess.get(i, k)
			if err != nil {
				return nil, err
			}
			err = hs.set(j, l, v)
			if err != nil {
				return nil, err
			}
			l++
		}
		j++
	}
	return hs, nil
}
