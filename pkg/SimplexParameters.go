package minuit

type SimplexParameters struct {
	theSimplexParameters []*Pair[float64, *MnAlgebraicVector]
	theJHigh             int
	theJLow              int
}

func NewSimplexParameters(simpl []*Pair[float64, *MnAlgebraicVector], jh, jl int) *SimplexParameters {
	return &SimplexParameters{
		theSimplexParameters: simpl,
		theJHigh:             jh,
		theJLow:              jl,
	}
}

func (this *SimplexParameters) update(y float64, p *MnAlgebraicVector) {
	this.theSimplexParameters[this.theJHigh] = NewPair[float64, *MnAlgebraicVector](y, p)
	if y < this.theSimplexParameters[this.theJLow].First {
		this.theJLow = this.theJHigh
	}

	var jh int = 0
	for i := 1; i < len(this.theSimplexParameters); i++ {
		if this.theSimplexParameters[i].First > this.theSimplexParameters[jh].First {
			jh = i
		}
	}
	this.theJHigh = jh
}

func (this *SimplexParameters) simplex() []*Pair[float64, *MnAlgebraicVector] {
	return this.theSimplexParameters
}

func (this *SimplexParameters) get(i int) *Pair[float64, *MnAlgebraicVector] {
	return this.theSimplexParameters[i]
}

func (this *SimplexParameters) jh() int {
	return this.theJHigh
}

func (this *SimplexParameters) jl() int {
	return this.theJLow
}

func (this *SimplexParameters) edm() float64 {
	return this.theSimplexParameters[this.theJHigh].First - this.theSimplexParameters[this.theJLow].First
}

func (this *SimplexParameters) dirin() *MnAlgebraicVector {
	var dirin *MnAlgebraicVector = NewMnAlgebraicVector(len(this.theSimplexParameters) - 1)
	for i := 0; i < len(this.theSimplexParameters)-1; i++ {
		var pbig float64 = this.theSimplexParameters[0].Second.get(i)
		var plit float64 = pbig
		for j := 0; j < len(this.theSimplexParameters); j++ {
			if this.theSimplexParameters[j].Second.get(i) < plit {
				plit = this.theSimplexParameters[j].Second.get(i)
			}
			if this.theSimplexParameters[j].Second.get(i) > pbig {
				pbig = this.theSimplexParameters[j].Second.get(i)
			}
		}
		dirin.set(i, pbig-plit)
	}
	return dirin
}
