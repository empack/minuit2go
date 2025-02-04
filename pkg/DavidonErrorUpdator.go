package minuit

type DavidonErrorUpdator struct {
}

func NewDavidonErrorUpdator() *DavidonErrorUpdator {
	return &DavidonErrorUpdator{}
}

func (this *DavidonErrorUpdator) Update(s0 *MinimumState, p1 *MinimumParameters,
	g1 *FunctionGradient) (*MinimumError, error) {
	V0 := s0.error().invHessian()
	dx, err := MnUtils.SubV(p1.vec(), s0.vec())
	if err != nil {
		return nil, err
	}
	dg, err := MnUtils.SubV(g1.vec(), s0.gradient().vec())
	if err != nil {
		return nil, err
	}
	delgam, err := MnUtils.InnerProduct(dx, dg)
	if err != nil {
		return nil, err
	}
	gvg, err := MnUtils.Similarity(dg, V0)
	if err != nil {
		return nil, err
	}
	vg, err := MnUtils.MulVSM(V0, dg)
	if err != nil {
		return nil, err
	}
	// chained implicit function calls for "Vupd" need error handling thus longer code from here
	outerDx, err := MnUtils.OuterProduct(dx)
	if err != nil {
		return nil, err
	}
	outerVg, err := MnUtils.OuterProduct(vg)
	if err != nil {
		return nil, err
	}

	Vupd, err := MnUtils.SubSM(MnUtils.DivSM(outerDx, delgam), MnUtils.DivSM(outerVg, gvg))
	if err != nil {
		return nil, err
	}

	if delgam > gvg {
		subDx, err := MnUtils.SubV(MnUtils.DivV(dx, delgam), MnUtils.DivV(vg, gvg))
		if err != nil {
			return nil, err
		}
		outer, err := MnUtils.OuterProduct(subDx)
		if err != nil {
			return nil, err
		}
		Vupd, err = MnUtils.AddSM(Vupd, MnUtils.MulSM(outer, gvg))
		if err != nil {
			return nil, err
		}
	}

	sum_upd := MnUtils.AbsoluteSumOfElements(Vupd)
	Vupd, err = MnUtils.AddSM(Vupd, V0)
	if err != nil {
		return nil, err
	}
	dcov := 0.5 * (s0.error().dcovar() + sum_upd/MnUtils.AbsoluteSumOfElements(Vupd))
	return NewMinimumError(Vupd, dcov), nil
}
