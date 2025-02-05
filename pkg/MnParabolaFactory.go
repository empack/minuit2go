package minuit

var MnParabolaFactory = MnParabolaFactoryStruct{}

type MnParabolaFactoryStruct struct {
}

func (this *MnParabolaFactoryStruct) createWith3Points(p1, p2, p3 *MnParabolaPoint) *MnParabola {
	x1 := p1.x()
	x2 := p2.x()
	x3 := p3.x()
	dx12 := x1 - x2
	dx13 := x1 - x3
	dx23 := x2 - x3
	xm := (x1 + x2 + x3) / 3.0
	x1 -= xm
	x2 -= xm
	x3 -= xm
	y1 := p1.y()
	y2 := p2.y()
	y3 := p3.y()
	a := y1/(dx12*dx13) - y2/(dx12*dx23) + y3/(dx13*dx23)
	b := -y1*(x2+x3)/(dx12*dx13) + y2*(x1+x3)/(dx12*dx23) - y3*(x1+x2)/(dx13*dx23)
	c := y1 - a*x1*x1 - b*x1
	c += xm * (xm*a - b)
	b -= 2.0 * xm * a

	return NewMnParabola(a, b, c)
}

func (this *MnParabolaFactoryStruct) createWith2Points(p1 *MnParabolaPoint, dxdy1 float64,
	p2 *MnParabolaPoint) *MnParabola {
	x1 := p1.x()
	xx1 := x1 * x1
	x2 := p2.x()
	xx2 := x2 * x2
	y1 := p1.y()
	y12 := p1.y() - p2.y()
	det := xx1 - xx2 - 2.0*x1*(x1-x2)
	a := -(y12 + (x2-x1)*dxdy1) / det
	b := -(-2.0*x1*y12 + (xx1-xx2)*dxdy1) / det
	c := y1 - a*xx1 - b*x1

	return NewMnParabola(a, b, c)
}
