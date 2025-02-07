package minuit

type MnFcnInterface interface {
	valueOf(v *MnAlgebraicVector) float64

	numOfCalls() int

	errorDef() float64

	fcn() FCNBase
}
type MnFcn struct {
	theFCN      FCNBase
	theNumCall  int
	theErrorDef float64
}

func NewMnFcn(fcn FCNBase, errorDef float64) *MnFcn {
	return &MnFcn{
		theFCN:      fcn,
		theNumCall:  0,
		theErrorDef: errorDef,
	}
}

func (this *MnFcn) valueOf(v *MnAlgebraicVector) float64 {
	this.theNumCall++
	return this.theFCN.ValueOf(v.asArray())
}

func (this *MnFcn) numOfCalls() int {
	return this.theNumCall
}

func (this *MnFcn) errorDef() float64 {
	return this.theErrorDef
}

func (this *MnFcn) fcn() FCNBase {
	return this.theFCN
}
