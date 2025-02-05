package minuit

import "math"

type MnParabola struct {
	theA float64
	theB float64
	theC float64
}

func NewMnParabola(a, b, c float64) *MnParabola {
	return &MnParabola{
		theA: a,
		theB: b,
		theC: c,
	}
}

func (this *MnParabola) y(x float64) float64 {
	return this.theA*x*x + this.theB*x + this.theC
}

func (this *MnParabola) x_pos(y float64) float64 {
	return math.Sqrt(y/this.theA+this.min()*this.min()-this.theC/this.theA) + this.min()
}

func (this *MnParabola) x_neg(y float64) float64 {
	return -math.Sqrt(y/this.theA+this.min()*this.min()-this.theC/this.theA) + this.min()
}

func (this *MnParabola) min() float64 {
	return -this.theB / (2.0 * this.theA)
}

func (this *MnParabola) ymin() float64 {
	return -this.theB*this.theB/(4.0*this.theA) + this.theC
}

func (this *MnParabola) a() float64 {
	return this.theA
}

func (this *MnParabola) b() float64 {
	return this.theB
}

func (this *MnParabola) c() float64 {
	return this.theC
}
