package minuit

type MnParabolaPoint struct {
	theX float64
	theY float64
}

func NewMnParabolaPoint(x, y float64) *MnParabolaPoint {
	return &MnParabolaPoint{
		theX: x,
		theY: y,
	}
}

func (this *MnParabolaPoint) x() float64 {
	return this.theX
}

func (this *MnParabolaPoint) y() float64 {
	return this.theY
}
