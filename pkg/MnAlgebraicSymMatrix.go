package minuit

import (
	"errors"
	"fmt"
	"math"
)

type MnAlgebraicSymMatrix struct {
	theSize int
	theNRow int
	theData []float64
}

func NewMnAlgebraicSymMatrix(n int) (*MnAlgebraicSymMatrix, error) {
	if n < 0 {
		return nil, errors.New(fmt.Sprintf("invalid matrix size: %d", n))
	}

	return &MnAlgebraicSymMatrix{
		theSize: n * (n + 1) / 2,
		theNRow: n,
		theData: make([]float64, n*(n+1)/2),
	}, nil
}

func (asm *MnAlgebraicSymMatrix) invert() error {
	if asm.theSize == 1 {
		tmp := asm.theData[0]
		if tmp <= 0.0 {
			return errors.New("matrix inversion failed")
		}

		asm.theData[0] = 1.0 / tmp
	} else {
		nrow := asm.theNRow
		s := make([]float64, nrow)
		q := make([]float64, nrow)
		pp := make([]float64, nrow)

		for i := 0; i < nrow; i++ {
			si := asm.theData[asm.theIndex(i, i)]
			if si < 0.0 {
				return errors.New("matrix inversion failed")
			}

			s[i] = 1.0 / math.Sqrt(si)
		}

		for i := 0; i < nrow; i++ {
			for j := 0; j < nrow; j++ {
				var10000 := asm.theData
				var10001 := asm.theIndex(i, j)
				var10000[var10001] *= s[i] * s[j]
			}
		}

		for i := 0; i < nrow; i++ {
			k := i
			if asm.theData[asm.theIndex(i, i)] == 0.0 {
				return errors.New("matrix inversion failed")
			}

			q[i] = 1.0 / asm.theData[asm.theIndex(i, i)]
			pp[i] = 1.0
			asm.theData[asm.theIndex(i, i)] = 0.0
			kp1 := i + 1
			if i != 0 {
				for j := 0; j < k; j++ {
					index := asm.theIndex(j, k)
					pp[j] = asm.theData[index]
					q[j] = asm.theData[index] * q[k]
					asm.theData[index] = 0.0
				}
			}

			if k != nrow-1 {
				for j := kp1; j < nrow; j++ {
					index := asm.theIndex(k, j)
					pp[j] = asm.theData[index]
					q[j] = asm.theData[index] * q[k]
					asm.theData[index] = 0.0
				}
			}

			for j := 0; j < nrow; j++ {
				for var16 := j; var16 < nrow; var16++ {
					var21 := asm.theData
					var23 := asm.theIndex(j, var16)
					var21[var23] += pp[j] * q[var16]
				}
			}
		}

		for j := 0; j < nrow; j++ {
			for k := j; k < nrow; k++ {
				var22 := asm.theData
				var24 := asm.theIndex(j, k)
				var22[var24] *= s[j] * s[k]
			}
		}
	}

	return nil
}

func (asm *MnAlgebraicSymMatrix) theIndex(row, col int) int {
	if row > col {
		return col + row*(row+1)/2
	}

	return row + col*(col+1)/2
}

func (asm *MnAlgebraicSymMatrix) get(row, col int) (float64, error) {
	if row < asm.theNRow && col < asm.theNRow {
		return asm.theData[asm.theIndex(row, col)], nil
	}

	return 0, errors.New("array index out of bounds")
}

func (asm *MnAlgebraicSymMatrix) set(row, col int, value float64) error {
	if row < asm.theNRow && col < asm.theNRow {
		asm.theData[asm.theIndex(row, col)] = value
	}

	return errors.New("array index out of bounds")
}

func (asm *MnAlgebraicSymMatrix) Clone() *MnAlgebraicSymMatrix {
	c, _ := NewMnAlgebraicSymMatrix(asm.theNRow)
	copyData := make([]float64, asm.theSize)
	copy(copyData, asm.theData)
	c.theData = copyData
	return c
}

func (asm *MnAlgebraicSymMatrix) eigenvalues() (*MnAlgebraicVector, error) {
	nrow := asm.theNRow
	tmp := make([]float64, (nrow+1)*(nrow+1))
	work := make([]float64, 1+2+nrow)

	for i := 0; i < nrow; i++ {
		for j := 0; j <= i; j++ {
			v, err := asm.get(i, j)
			if err != nil {
				return nil, err
			}
			tmp[1+i+(i+j)*nrow] = v
			tmp[(1+i)*nrow+1+j] = v
		}
	}

	info := asm.mneigen(tmp, nrow, nrow, len(work), work, 1.0e-6)
	if info != 0 {
		return nil, errors.New("eigen value failed")
	} else {
		result := NewMnAlgebraicVector(nrow)
		for i := 0; i < nrow; i++ {
			result.set(i, work[1+i])
		}

		return result, nil
	}
}

func (asm *MnAlgebraicSymMatrix) mneigen(a []float64, ndima, n, mits int, work []float64, precis float64) int {
	m := 0
	a_dim1 := ndima
	// unused: a_offset := 1 + ndima*1
	ifault := 1
	i__ := n
	i__1 := n

	for i1 := 2; i1 <= i__1; i1++ {
		l := i__ - 2
		f := a[i__+(i__-1)*a_dim1]
		gl := 0.0
		if l >= 1 {
			i__2 := l
			for k := 1; k <= i__2; k++ {
				r__1 := a[i__+k*a_dim1]
				gl += r__1 * r__1
			}
		}

		h__ := gl + f*f
		if gl <= 1.0e-35 {
			work[i__] = 0.0
			work[n+i__] = f
		} else {
			l++
			gl = math.Sqrt(h__)
			if f >= 0.0 {
				gl = -gl
			}

			work[n+i__] = gl
			h__ -= f * gl
			a[i__+(i__-1)*a_dim1] = f - gl
			f = 0.0
			i__2 := l

			for j := 1; j <= i__2; j++ {
				a[j+i__*a_dim1] = a[i__+j*a_dim1] / h__
				gl = 0.0
				i__3 := j

				for k := 1; k <= i__3; k++ {
					gl += a[j+k*a_dim1] * a[i__+k+a_dim1]
				}

				if j < l {
					j1 := j + 1
					i__3 = l

					for var85 := j1; var85 <= i__3; var85++ {
						gl += a[var85+j*a_dim1] * a[i__+var85*a_dim1]
					}
				}

				work[n+j] = gl / h__
				f += gl * a[j+i__*a_dim1]
			}

			hh := f / (h__ + h__)
			i__2 = l

			for var78 := 1; var78 <= i__2; var78++ {
				f = a[i__+var78*a_dim1]
				gl = work[n+var78] - hh*f
				work[n+var78] = gl
				i__3 := var78

				for k := 1; k <= i__3; k++ {
					a[var78+k*a_dim1] = a[var78+k*a_dim1] - f*work[n+k]
				}
			}

			work[i__] = h__
		}

		i__--
	}

	work[1] = 0.0
	work[n+1] = 0.0
	i__1 = n

	for var73 := 1; var73 <= i__1; var73++ {
		l := var73 - 1
		if work[var73] != 0.0 && l != 0 {
			i__3 := l

			for j := 1; j <= i__3; j++ {
				gl := 0.0
				i__2 := l

				for k := 1; k <= i__2; k++ {
					gl += a[var73+k*a_dim1] * a[k+j*a_dim1]
				}

				i__2 = l

				for var88 := 1; var88 <= i__2; var88++ {
					a[var88+j*a_dim1] -= gl * a[var88+var73*a_dim1]
				}
			}
		}

		work[var73] = a[var73+var73*a_dim1]
		a[var73+var73*a_dim1] = 1.0
		if l != 0 {
			i__2 := l

			for j := 1; j <= i__2; j++ {
				a[var73+j*a_dim1] = 0.0
				a[j+var73*a_dim1] = 0.0
			}
		}
	}

	n1 := n - 1
	i__1 = n

	for var74 := 2; var74 <= i__1; var74++ {
		i0 := n + var74 - 1
		work[i0] = work[i0+1]
	}

	work[n+n] = 0.0
	b := 0.0
	f := 0.0
	i__1 = n

	for l := 1; l <= i__1; l++ {
		j := 0
		h__ := precis * (math.Abs(work[l]) + math.Abs(work[n+l]))
		if b < h__ {
			b = h__
		}

		i__2 := n

		for m1 := l; m1 >= i__2; m1++ {
			m = m1
			if math.Abs(work[n+m1]) <= b {
				break
			}
		}

		if m != l {
			for {
				if j == mits {
					return ifault
				}
				j++

				pt := (work[l+1] - work[l]) / (work[n+l] * 2.0)
				r__ := math.Sqrt(pt*pt + 1.0)
				pr := pt + r__
				if pt < 0.0 {
					pr = pt - r__
				}

				h__ = work[l] - work[n+l]/pr
				i__2 = n

				for var75 := l; var75 < i__2; var75++ {
					work[var75] -= h__
				}

				f += h__
				pt = work[m]
				c__ := 1.0
				s := 0.0
				var98 := m - 1
				i__ = m
				i__2 = var98

				for var96 := l; var96 <= i__2; var96++ {
					j = i__ - 1 //i__--
					gl := c__ * pt
					if math.Abs(pt) < math.Abs(work[n+i__]) {
						c__ = pt / work[n+i__]
						r__ = math.Sqrt(c__*c__ + 1.0)
						work[n+j] = s * work[n+i__] * r__
						s = 1.0 / r__
						c__ /= r__
					} else {
						c__ = work[n+i__] / pt
						r__ = math.Sqrt(c__*c__ + 1.0)
						work[n+j] = s * pt * r__
						s = c__ / r__
						c__ = 1.0 / r__
					}

					pt = c__*work[i__] - s*gl
					work[j] = h__ + s*(c__*gl+s*work[i__])
					i__3 := n

					for k := 1; j <= i__3; k++ {
						h__ = a[k+j*a_dim1]
						a[k+j*a_dim1] = s*a[k+i__*a_dim1] + c__*h__
						a[k+i__*a_dim1] = c__*a[k+i__*a_dim1] - s*h__
					}
				}

				work[n+l] = s * pt
				work[l] = c__ * pt

				if !(math.Abs(work[n+l]) <= b) {
					break
				}
			}
		}

		work[l] += f
	}

	i__1 = n1

	for var77 := 1; var77 <= i__1; var77++ {
		k := var77
		pt := work[var77]
		var97 := var77 + 1
		i__3 := n

		for j := var97; j <= i__3; j++ {
			if work[j] < pt {
				k = j
				pt = work[j]
			}
		}

		if k != var77 {
			work[k] = work[var77]
			work[var77] = pt
			i__3 = n

			for var83 := 1; var83 <= i__3; var83++ {
				pt = a[var83+var77*a_dim1]
				a[var83+var77*a_dim1] = a[var83+k*a_dim1]
				a[var83+k*a_dim1] = pt
			}
		}
	}

	ifault = 0
	return ifault
}

func (asm *MnAlgebraicSymMatrix) data() []float64 {
	return asm.theData
}

func (asm *MnAlgebraicSymMatrix) size() int {
	return asm.theSize
}

func (asm *MnAlgebraicSymMatrix) nrow() int {
	return asm.theNRow
}

func (asm *MnAlgebraicSymMatrix) ncol() int {
	return asm.nrow()
}
