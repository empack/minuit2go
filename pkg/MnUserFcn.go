package minuit

type MnUserFcn struct {
	theTransform *MnUserTransformation
	ParentClass  *MnFcn
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
