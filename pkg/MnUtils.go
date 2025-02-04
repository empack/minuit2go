package minuit

import (
	"errors"
	"math"
)

var MnUtils = &MnUtilsStruct{}

type MnUtilsStruct struct {
}

func (this *MnUtilsStruct) Similarity(avec *MnAlgebraicVector, mat *MnAlgebraicSymMatrix) (float64, error) {
	n := avec.size()
	tmp, err := this.MulVSM(mat, avec)
	if err != nil {
		return 0, err
	}
	result := 0.0

	for i := 0; i < n; i++ {
		result += tmp.get(i) * avec.get(i)
	}

	return result, nil
}

func (this *MnUtilsStruct) AddV(v1, v2 *MnAlgebraicVector) (*MnAlgebraicVector, error) {
	if v1.size() != v2.size() {
		return nil, errors.New("incompatible vectors")
	} else {
		result := v1.Clone()
		a := result.data
		b := v2.data

		for i := 0; i < len(a); i++ {
			a[i] += b[i]
		}

		return result, nil
	}
}

func (this *MnUtilsStruct) AddSM(m1, m2 *MnAlgebraicSymMatrix) (*MnAlgebraicSymMatrix, error) {
	if m1.size() != m2.size() {
		return nil, errors.New("incompatible matrices")
	} else {
		result := m1.Clone()
		a := result.data()
		b := m2.data()

		for i := 0; i < len(a); i++ {
			a[i] += b[i]
		}

		return result, nil
	}
}

func (this *MnUtilsStruct) SubV(v1, v2 *MnAlgebraicVector) (*MnAlgebraicVector, error) {
	if v1.size() != v2.size() {
		return nil, errors.New("incompatible vectors")
	} else {
		result := v1.Clone()
		a := result.data
		b := v2.data

		for i := 0; i < len(a); i++ {
			a[i] -= b[i]
		}

		return result, nil
	}
}

func (this *MnUtilsStruct) SubSM(m1, m2 *MnAlgebraicSymMatrix) (*MnAlgebraicSymMatrix, error) {
	if m1.size() != m2.size() {
		return nil, errors.New("incompatible matrices")
	} else {
		result := m1.Clone()
		a := result.data()
		b := m2.data()

		for i := 0; i < len(a); i++ {
			a[i] -= b[i]
		}

		return result, nil
	}
}

func (this *MnUtilsStruct) MulV(v1 *MnAlgebraicVector, scale float64) *MnAlgebraicVector {
	result := v1.Clone()
	a := result.data

	for i := 0; i < len(a); i++ {
		a[i] *= scale
	}

	return result
}

func (this *MnUtilsStruct) MulSM(m1 *MnAlgebraicSymMatrix, scale float64) *MnAlgebraicSymMatrix {
	result := m1.Clone()
	a := result.data()

	for i := 0; i < len(a); i++ {
		a[i] *= scale
	}

	return result
}

func (this *MnUtilsStruct) MulVSM(m1 *MnAlgebraicSymMatrix, v1 *MnAlgebraicVector) (*MnAlgebraicVector, error) {
	if m1.nrow() != v1.size() {
		return nil, errors.New("incompatible arguments")
	} else {
		result := NewMnAlgebraicVector(m1.nrow())
		a := result.data

		for i := 0; i < len(a); i++ {
			total := 0.0

			for k := 0; k < len(a); k++ {
				v, err := m1.get(i, k)
				if err != nil {
					return nil, err
				}
				total += v * v1.get(k)
			}

			a[i] = total
		}

		return result, nil
	}
}

func (this *MnUtilsStruct) MulSMs(m1, m2 *MnAlgebraicSymMatrix) (*MnAlgebraicSymMatrix, error) {
	if m1.size() != m2.size() {
		return nil, errors.New("incompatible matrices")
	} else {
		n := m1.nrow()
		result, err := NewMnAlgebraicSymMatrix(n)
		if err != nil {
			return nil, err
		}

		for i := 0; i < n; i++ {
			for j := 0; j < i; j++ {
				total := 0.0

				for k := 0; k < n; k++ {
					v1, m1err := m1.get(i, k)
					if m1err != nil {
						return nil, m1err
					}
					v2, m2err := m2.get(k, j)
					if m2err != nil {
						return nil, m2err
					}
					total += v1 * v2
				}

				setErr := result.set(i, j, total)
				if setErr != nil {
					return nil, setErr
				}
			}
		}

		return result, nil
	}
}

func (this *MnUtilsStruct) InnerProduct(v1, v2 *MnAlgebraicVector) (float64, error) {
	if v1.size() != v2.size() {
		return 0, errors.New("incompatible vectors")
	} else {
		a := v1.data
		b := v2.data
		total := 0.0

		for i := 0; i < len(a); i++ {
			total += a[i] * b[i]
		}

		return total, nil
	}
}

func (this *MnUtilsStruct) DivSM(m *MnAlgebraicSymMatrix, scale float64) *MnAlgebraicSymMatrix {
	return this.MulSM(m, 1.0/scale)
}

func (this *MnUtilsStruct) DivV(m *MnAlgebraicVector, scale float64) *MnAlgebraicVector {
	return this.MulV(m, 1.0/scale)
}

func (this *MnUtilsStruct) OuterProduct(v2 *MnAlgebraicVector) (*MnAlgebraicSymMatrix, error) {
	n := v2.size()
	result, err := NewMnAlgebraicSymMatrix(n)
	if err != nil {
		return nil, err
	}
	data := v2.data

	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			setErr := result.set(i, j, data[i]*data[j])
			if setErr != nil {
				return nil, setErr
			}
		}
	}

	return result, nil
}

func (this *MnUtilsStruct) AbsoluteSumOfElements(m *MnAlgebraicSymMatrix) float64 {
	data := m.data()
	result := 0.0

	for i := 0; i < len(data); i++ {
		result += math.Abs(data[i])
	}

	return result
}
