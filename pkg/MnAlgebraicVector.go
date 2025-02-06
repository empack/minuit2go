package minuit

import "slices"

type MnAlgebraicVector struct {
	data []float64
}

func NewMnAlgebraicVector(size int) *MnAlgebraicVector {
	return &MnAlgebraicVector{
		data: make([]float64, size),
	}
}

func (this *MnAlgebraicVector) Clone() *MnAlgebraicVector {
	var result *MnAlgebraicVector = NewMnAlgebraicVector(len(this.data))
	result.data = slices.Clone(this.data)
	return result
}

func (this *MnAlgebraicVector) size() int {
	return len(this.data)
}

func (this *MnAlgebraicVector) get(i int) float64 {
	return this.data[i]
}

func (this *MnAlgebraicVector) set(i int, value float64) {
	this.data[i] = value
}

func (this *MnAlgebraicVector) asArray() []float64 {
	return this.data
}
func (this *MnAlgebraicVector) String() string { return MnPrint.toStringMnAlgebraicVector(this) }
