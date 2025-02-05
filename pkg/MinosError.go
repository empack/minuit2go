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
	return -1.0 * this.LowerState().Error(this.Parameter()) * (1.0 + this.theLower.Value())
}

func (this *MinosError) Upper() float64 {
	return this.UpperState().Error(this.Paramter()) * (1.0 + this.theUpper.Value())
}

func (this *MinosError) Paramter() int {
	return this.theParameter
}

func (this *MinosError) LowerState() *MnUserParameterState {
	return this.theLower.State()
}

func (this *MinosError) UpperState() *MnUserParameterState {
	return this.theUpper.State()
}

func (this *MinosError) IsValid() bool {
	return this.theLower.IsValid() && this.theUpper.IsValid()
}

func (this *MinosError) LowerValid() bool {
	return this.theLower.IsValid()
}

func (this *MinosError) UpperValid() bool {
	return this.theUpper.IsValid()
}

func (this *MinosError) AtLowerLimit() bool {
	return this.theLower.AtLimit()
}

func (this *MinosError) AtUpperLimit() bool {
	return this.theUpper.AtLimit()
}

func (this *MinosError) AtLowerMaxFcn() bool {
	return this.theLower.AtMaxFcn()
}

func (this *MinosError) AtUpperMaxFcn() bool {
	return this.theUpper.AtMaxFcn()
}

func (this *MinosError) LowerNewMin() bool {
	return this.theLower.NewMinimum()
}

func (this *MinosError) UpperNewMin() bool {
	return this.theUpper.NewMinimum()
}

func (this *MinosError) Nfcn() int {
	return this.theUpper.Nfcn() + this.theLower.Nfcn()
}

func (this *MinosError) Min() float64 {
	return this.theMinValue
}
