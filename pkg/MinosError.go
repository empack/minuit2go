package minuit

type MinosError struct {
	theParameter int
	theMinValue  float64
	theUpper     *MnCross
	theLower     *MnCross
}

func NewMinosError() *MinosError {
	return &MinosError{
		theParameter: 0,
		theMinValue:  0,
		theUpper:     NewMnCross(),
		theLower:     NewMnCross(),
	}
}

func NewMinosErrorWithValues(par int, min float64, low, up *MnCross) *MinosError {
	return &MinosError{
		theParameter: par,
		theMinValue:  min,
		theUpper:     up,
		theLower:     low,
	}
}

func (this *MinosError) Range() *Point {
	return NewPoint(this.Lower(), this.Upper())
}

func (this *MinosError) Lower() float64 {
	return -1.0 * this.LowerState().Error(this.Parameter()) * (1.0 + this.theLower.value())
}

func (this *MinosError) Upper() float64 {
	return this.UpperState().Error(this.Parameter()) * (1.0 + this.theUpper.value())
}

func (this *MinosError) Parameter() int {
	return this.theParameter
}

func (this *MinosError) LowerState() *MnUserParameterState {
	return this.theLower.state()
}

func (this *MinosError) UpperState() *MnUserParameterState {
	return this.theUpper.state()
}

func (this *MinosError) IsValid() bool {
	return this.theLower.isValid() && this.theUpper.isValid()
}

func (this *MinosError) LowerValid() bool {
	return this.theLower.isValid()
}

func (this *MinosError) UpperValid() bool {
	return this.theUpper.isValid()
}

func (this *MinosError) AtLowerLimit() bool {
	return this.theLower.atLimit()
}

func (this *MinosError) AtUpperLimit() bool {
	return this.theUpper.atLimit()
}

func (this *MinosError) AtLowerMaxFcn() bool {
	return this.theLower.atMaxFcn()
}

func (this *MinosError) AtUpperMaxFcn() bool {
	return this.theUpper.atMaxFcn()
}

func (this *MinosError) LowerNewMin() bool {
	return this.theLower.newMinimum()
}

func (this *MinosError) UpperNewMin() bool {
	return this.theUpper.newMinimum()
}

func (this *MinosError) Nfcn() int {
	return this.theUpper.nfcn() + this.theLower.nfcn()
}

func (this *MinosError) Min() float64 {
	return this.theMinValue
}

func (this *MinosError) String() string {
	return MnPrint.toStringMinosError(this)
}
