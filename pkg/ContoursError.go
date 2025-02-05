package minuit

var ContoursError = &ContoursErrorStruct{}

type ContoursErrorStruct struct {
	theParX   int
	theParY   int
	thePoints []*Point
	theXMinos *MinosError
	theYMinos *MinosError
	theNFcn   int
}

func NewContoursError(parx, pary int, points []*Point, xmnos, ymnos *MinosError, nfcn int) *ContoursErrorStruct {
	return &ContoursErrorStruct{
		theParX:   parx,
		theParY:   pary,
		thePoints: points,
		theXMinos: xmnos,
		theYMinos: ymnos,
		theNFcn:   nfcn,
	}
}

func (this *ContoursErrorStruct) Points() []*Point {
	return this.thePoints
}

func (this *ContoursErrorStruct) XRange() *Point {
	return this.theXMinos.range()
}

func (this *ContoursErrorStruct) YRange() *Point {
	return this.theYMinos.range()
}

func (this *ContoursErrorStruct) Xpar() int {
	return this.theParX
}

func (this *ContoursErrorStruct) YPar() int {
	return this.theParY
}

func (this *ContoursErrorStruct) XMinosError() *MinosError {
	return this.theXMinos
}

func (this *ContoursErrorStruct) YMinosError() *MinosError {
	return this.theYMinos
}

func (this *ContoursErrorStruct) Nfcn() int {
	return this.theNFcn
}

func (this *ContoursErrorStruct) Xmin() float64 {
	return this.theXMinos.min()
}

func (this *ContoursErrorStruct) Ymin() float64 {
	return this.theYMinos.min()
}
