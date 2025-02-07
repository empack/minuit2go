package minuit

type MnUserFcn struct {
	theTransform *MnUserTransformation
	ParentClass  *MnFcn
}

func (this *MnUserFcn) numOfCalls() int {
	return this.ParentClass.numOfCalls()
}

func (this *MnUserFcn) errorDef() float64 {
	return this.ParentClass.theErrorDef
}

func (this *MnUserFcn) fcn() FCNBase {
	return this.ParentClass.fcn()
}

func NewMnUserFcn(fcn FCNBase, errDef float64, trafo *MnUserTransformation) *MnUserFcn {
	return &MnUserFcn{
		theTransform: trafo,
		ParentClass:  NewMnFcn(fcn, errDef),
	}
}

func (this *MnUserFcn) valueOf(v *MnAlgebraicVector) float64 {
	return this.ParentClass.valueOf(this.theTransform.transform(v))
}
