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

func (this *MnAlgebraicSymMatrix) invert() error {
	if this.theSize == 1 {
		tmp := this.theData[0]
		if tmp <= 0.0 {
			return errors.New("matrix inversion failed")
		}

		this.theData[0] = 1.0 / tmp
	} else {
		nrow := this.theNRow
		s := make([]float64, nrow)
		q := make([]float64, nrow)
		pp := make([]float64, nrow)

		for i := 0; i < nrow; i++ {
			si := this.theData[this.theIndex(i, i)]
			if si < 0.0 {
				return errors.New("matrix inversion failed")
			}

			s[i] = 1.0 / math.Sqrt(si)
		}

		for i := 0; i < nrow; i++ {
			for j := i; j < nrow; j++ {
				this.theData[this.theIndex(i, j)] *= s[i] * s[j]
			}
		}

		for i := 0; i < nrow; i++ {
			var k int = i
			if this.theData[this.theIndex(k, k)] == 0. {
				return errors.New("matrix inversion failed")
			}
			q[k] = 1. / this.theData[this.theIndex(k, k)]
			pp[k] = 1.
			this.theData[this.theIndex(k, k)] = 0.
			var kp1 int = k + 1
			if k != 0 {
				for j := 0; j < k; j++ {
					var index int = this.theIndex(j, k)
					pp[j] = this.theData[index]
					q[j] = this.theData[index] * q[k]
					this.theData[index] = 0.
				}
			}
			if k != nrow-1 {
				for j := kp1; j < nrow; j++ {
					var index int = this.theIndex(k, j)
					pp[j] = this.theData[index]
					q[j] = -this.theData[index] * q[k]
					this.theData[index] = 0.
				}
			}
			for j := 0; j < nrow; j++ {
				for k = j; k < nrow; k++ {
					this.theData[this.theIndex(j, k)] += pp[j] * q[k]
				}
			}

		}

		for j := 0; j < nrow; j++ {
			for k := j; k < nrow; k++ {
				this.theData[this.theIndex(j, k)] *= s[j] * s[k]
			}
		}
	}

	return nil
}

func (this *MnAlgebraicSymMatrix) theIndex(row, col int) int {
	if row > col {
		return col + row*(row+1)/2
	}

	return row + col*(col+1)/2
}

func (this *MnAlgebraicSymMatrix) get(row, col int) (float64, error) {
	if row < this.theNRow && col < this.theNRow {
		return this.theData[this.theIndex(row, col)], nil
	}

	return 0, errors.New("array index out of bounds")
}

func (this *MnAlgebraicSymMatrix) set(row, col int, value float64) error {
	if row < this.theNRow && col < this.theNRow {
		this.theData[this.theIndex(row, col)] = value
		return nil
	}

	return errors.New("array index out of bounds")
}

func (this *MnAlgebraicSymMatrix) Clone() *MnAlgebraicSymMatrix {
	c, _ := NewMnAlgebraicSymMatrix(this.theNRow)
	copyData := make([]float64, this.theSize)
	copy(copyData, this.theData)
	c.theData = copyData
	return c
}

func (this *MnAlgebraicSymMatrix) eigenvalues() (*MnAlgebraicVector, error) {
	nrow := this.theNRow
	tmp := make([]float64, (nrow+1)*(nrow+1))
	work := make([]float64, 1+2*nrow)

	for i := 0; i < nrow; i++ {
		for j := 0; j <= i; j++ {
			v, err := this.get(i, j)
			if err != nil {
				return nil, err
			}
			tmp[1+i+((1+j)*nrow)] = v
			tmp[(1+i)*nrow+1+j] = v
		}
	}

	info := this.mneigen(tmp, nrow, nrow, len(work), work, 1.0e-6)
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

func (this *MnAlgebraicSymMatrix) mneigen(a []float64, ndima, n, mits int, work []float64, precis float64) int { /* System generated locals */
	var a_dim1, i__1, i__2, i__3 int

	/* Local variables */
	var b, c__, f, h__ float64
	var i__, j, k, l, m int
	var r__, s float64
	var i0, i1, j1, m1, n1 int
	var hh, gl, pr, pt float64

	/* PRECIS is the machine precision EPSMAC */
	/* Parameter adjustments */
	a_dim1 = ndima

	/* Function Body */
	var ifault int = 1

	i__ = n
	i__1 = n
	for i1 = 2; i1 <= i__1; i1++ {
		l = i__ - 2
		f = a[i__+(i__-1)*a_dim1]
		gl = 0.

		if l >= 1 {
			i__2 = l
			for k = 1; k <= i__2; k++ {
				/* Computing 2nd power */
				var r__1 float64 = a[i__+k*a_dim1]
				gl += r__1 * r__1
			}
		}
		/* Computing 2nd power */
		h__ = gl + f*f
		if gl <= 1e-35 {
			work[i__] = 0.
			work[n+i__] = f
		} else {
			l++

			gl = math.Sqrt(h__)

			if f >= 0. {
				gl = -gl
			}

			work[n+i__] = gl
			h__ -= f * gl
			a[i__+(i__-1)*a_dim1] = f - gl
			f = 0.
			i__2 = l
			for j = 1; j <= i__2; j++ {
				a[j+i__*a_dim1] = a[i__+j*a_dim1] / h__
				gl = 0.
				i__3 = j
				for k = 1; k <= i__3; k++ {
					gl += a[j+k*a_dim1] * a[i__+k*a_dim1]
				}

				if j < l {
					j1 = j + 1
					i__3 = l
					for k = j1; k <= i__3; k++ {
						gl += a[k+j*a_dim1] * a[i__+k*a_dim1]
					}
				}
				work[n+j] = gl / h__
				f += gl * a[j+i__*a_dim1]
			}
			hh = f / (h__ + h__)
			i__2 = l
			for j = 1; j <= i__2; j++ {
				f = a[i__+j*a_dim1]
				gl = work[n+j] - hh*f
				work[n+j] = gl
				i__3 = j
				for k = 1; k <= i__3; k++ {
					a[j+k*a_dim1] = a[j+k*a_dim1] - f*work[n+k] - gl*a[i__+k*a_dim1]
				}
			}
			work[i__] = h__
		}
		i__--
	}
	work[1] = 0.
	work[n+1] = 0.
	i__1 = n
	for i__ = 1; i__ <= i__1; i__++ {
		l = i__ - 1

		if work[i__] != 0. && l != 0 {
			i__3 = l
			for j = 1; j <= i__3; j++ {
				gl = 0.
				i__2 = l
				for k = 1; k <= i__2; k++ {
					gl += a[i__+k*a_dim1] * a[k+j*a_dim1]
				}
				i__2 = l
				for k = 1; k <= i__2; k++ {
					a[k+j*a_dim1] -= gl * a[k+i__*a_dim1]
				}
			}
		}
		work[i__] = a[i__+i__*a_dim1]
		a[i__+i__*a_dim1] = 1.

		if l != 0 {
			i__2 = l
			for j = 1; j <= i__2; j++ {
				a[i__+j*a_dim1] = 0.
				a[j+i__*a_dim1] = 0.
			}
		}

	}

	n1 = n - 1
	i__1 = n
	for i__ = 2; i__ <= i__1; i__++ {
		i0 = n + i__ - 1
		work[i0] = work[i0+1]
	}
	work[n+n] = 0.
	b = 0.
	f = 0.
	i__1 = n
	for l = 1; l <= i__1; l++ {
		j = 0
		h__ = precis * (math.Abs(work[l]) + math.Abs(work[n+l]))

		if b < h__ {
			b = h__
		}

		i__2 = n
		for m1 = l; m1 <= i__2; m1++ {
			m = m1

			if math.Abs(work[n+m]) <= b {
				break
			}
		}

		if m != l {
			for {
				if j == mits {
					return ifault
				}

				j++
				pt = (work[l+1] - work[l]) / (work[n+l] * 2.)
				r__ = math.Sqrt(pt*pt + 1.)
				pr = pt + r__

				if pt < 0. {
					pr = pt - r__
				}

				h__ = work[l] - work[n+l]/pr
				i__2 = n
				for i__ = l; i__ <= i__2; i__++ {
					work[i__] -= h__
				}
				f += h__
				pt = work[m]
				c__ = 1.
				s = 0.
				m1 = m - 1
				i__ = m
				i__2 = m1
				for i1 = l; i1 <= i__2; i1++ {
					j = i__
					i__--
					gl = c__ * work[n+i__]
					h__ = c__ * pt

					if math.Abs(pt) < math.Abs(work[n+i__]) {
						c__ = pt / work[n+i__]
						r__ = math.Sqrt(c__*c__ + 1.)
						work[n+j] = s * work[n+i__] * r__
						s = 1. / r__
						c__ /= r__
					} else {
						c__ = work[n+i__] / pt
						r__ = math.Sqrt(c__*c__ + 1.)
						work[n+j] = s * pt * r__
						s = c__ / r__
						c__ = 1. / r__
					}
					pt = c__*work[i__] - s*gl
					work[j] = h__ + s*(c__*gl+s*work[i__])
					i__3 = n
					for k = 1; k <= i__3; k++ {
						h__ = a[k+j*a_dim1]
						a[k+j*a_dim1] = s*a[k+i__*a_dim1] + c__*h__
						a[k+i__*a_dim1] = c__*a[k+i__*a_dim1] - s*h__
					}
				}
				work[n+l] = s * pt
				work[l] = c__ * pt

				if math.Abs(work[n+l]) <= b {
					break
				}
			}
		}
		work[l] += f
	}
	i__1 = n1
	for i__ = 1; i__ <= i__1; i__++ {
		k = i__
		pt = work[i__]
		i1 = i__ + 1
		i__3 = n
		for j = i1; j <= i__3; j++ {
			if work[j] < pt {
				k = j
				pt = work[j]
			}
		}

		if k != i__ {
			work[k] = work[i__]
			work[i__] = pt
			i__3 = n
			for j = 1; j <= i__3; j++ {
				pt = a[j+i__*a_dim1]
				a[j+i__*a_dim1] = a[j+k*a_dim1]
				a[j+k*a_dim1] = pt
			}
		}
	}
	ifault = 0

	return ifault
} /* mneig_ */

func (this *MnAlgebraicSymMatrix) data() []float64 {
	return this.theData
}

func (this *MnAlgebraicSymMatrix) size() int {
	return this.theSize
}

func (this *MnAlgebraicSymMatrix) nrow() int {
	return this.theNRow
}

func (this *MnAlgebraicSymMatrix) ncol() int {
	return this.nrow()
}

func (this *MnAlgebraicSymMatrix) String() string { return MnPrint.toStringMnAlgebraicSymMatrix(this) }
