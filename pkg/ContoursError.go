package minuit

type ContoursError struct {
	theParX   int
	theParY   int
	thePoints []*Point
	theXMinos *MinosError
	theYMinos *MinosError
	theNFcn   int
}

func NewContoursError(parx, pary int, points []*Point, xmnos, ymnos *MinosError, nfcn int) *ContoursError {
	return &ContoursError{
		theParX:   parx,
		theParY:   pary,
		thePoints: points,
		theXMinos: xmnos,
		theYMinos: ymnos,
		theNFcn:   nfcn,
	}
}

func (this *ContoursError) Points() []*Point {
	return this.thePoints
}

func (this *ContoursError) XRange() *Point {
	return this.theXMinos.Range()
}

func (this *ContoursError) YRange() *Point {
	return this.theYMinos.Range()
}

func (this *ContoursError) Xpar() int {
	return this.theParX
}

func (this *ContoursError) YPar() int {
	return this.theParY
}

func (this *ContoursError) XMinosError() *MinosError {
	return this.theXMinos
}

func (this *ContoursError) YMinosError() *MinosError {
	return this.theYMinos
}

func (this *ContoursError) Nfcn() int {
	return this.theNFcn
}

func (this *ContoursError) Xmin() float64 {
	return this.theXMinos.Min()
}

func (this *ContoursError) Ymin() float64 {
	return this.theYMinos.Min()
}
