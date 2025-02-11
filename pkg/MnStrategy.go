package minuit

type MnStrategy struct {
	theStrategy     int
	theGradNCyc     int
	theGradTlrStp   float64
	theGradTlr      float64
	theHessNCyc     int
	theHessTlrStp   float64
	theHessTlrG2    float64
	theHessGradNCyc int
}

type StrategyType = int

const (
	FastStrategy StrategyType = iota
	StandardStrategy
	PreciseStrategy
)

func NewMnStrategy() *MnStrategy {
	return &MnStrategy{
		theStrategy:     1,
		theGradNCyc:     3,
		theGradTlrStp:   0.3,
		theGradTlr:      0.05,
		theHessNCyc:     5,
		theHessTlrStp:   0.3,
		theHessTlrG2:    0.05,
		theHessGradNCyc: 2,
	}
}

func NewMnStrategyWithStra(stra StrategyType) *MnStrategy {
	if stra == FastStrategy {
		return &MnStrategy{
			theStrategy:     0,
			theGradNCyc:     2,
			theGradTlrStp:   0.5,
			theGradTlr:      0.1,
			theHessNCyc:     3,
			theHessTlrStp:   0.5,
			theHessTlrG2:    0.1,
			theHessGradNCyc: 1,
		}
	}

	if stra == PreciseStrategy {
		return &MnStrategy{
			theStrategy:     2,
			theGradNCyc:     5,
			theGradTlrStp:   0.1,
			theGradTlr:      0.02,
			theHessNCyc:     7,
			theHessTlrStp:   0.1,
			theHessTlrG2:    0.02,
			theHessGradNCyc: 6,
		}
	}

	return NewMnStrategy()
}

func (this *MnStrategy) Strategy() int {
	return this.theStrategy
}

func (this *MnStrategy) GradientNCycles() int {
	return this.theGradNCyc
}

func (this *MnStrategy) GradientStepTolerance() float64 {
	return this.theGradTlrStp
}

func (this *MnStrategy) GradientTolerance() float64 {
	return this.theGradTlr
}

func (this *MnStrategy) HessianNCycles() int {
	return this.theHessNCyc
}

func (this *MnStrategy) HessianStepTolerance() float64 {
	return this.theHessTlrStp
}

func (this *MnStrategy) HessianG2Tolerance() float64 {
	return this.theHessTlrG2
}

func (this *MnStrategy) HessianGradientNCycles() int {
	return this.theHessGradNCyc
}

func (this *MnStrategy) IsLow() bool {
	return this.theStrategy <= 0
}

func (this *MnStrategy) IsMedium() bool {
	return this.theStrategy == 1
}

func (this *MnStrategy) IsHigh() bool {
	return this.theStrategy >= 2
}

func (this *MnStrategy) SetLowStrategy() {
	this.theStrategy = 0
	this.SetGradientNCycles(2)
	this.SetGradientStepTolerance(0.5)
	this.SetGradientTolerance(0.1)
	this.SetHessianNCycles(3)
	this.SetHessianStepTolerance(0.5)
	this.SetHessianG2Tolerance(0.1)
	this.SetHessianGradientNCycles(1)
}

func (this *MnStrategy) SetMediumStrategy() {
	this.theStrategy = 1
	this.SetGradientNCycles(3)
	this.SetGradientStepTolerance(0.3)
	this.SetGradientTolerance(0.05)
	this.SetHessianNCycles(5)
	this.SetHessianStepTolerance(0.3)
	this.SetHessianG2Tolerance(0.05)
	this.SetHessianGradientNCycles(2)
}

func (this *MnStrategy) SetHighStrategy() {
	this.theStrategy = 2
	this.SetGradientNCycles(5)
	this.SetGradientStepTolerance(0.1)
	this.SetGradientTolerance(0.02)
	this.SetHessianNCycles(7)
	this.SetHessianStepTolerance(0.1)
	this.SetHessianG2Tolerance(0.02)
	this.SetHessianGradientNCycles(6)
}

func (this *MnStrategy) SetGradientNCycles(n int) {
	this.theGradNCyc = n
}

func (this *MnStrategy) SetGradientStepTolerance(stp float64) {
	this.theGradTlrStp = stp
}

func (this *MnStrategy) SetGradientTolerance(toler float64) {
	this.theGradTlr = toler
}

func (this *MnStrategy) SetHessianNCycles(n int) {
	this.theHessNCyc = n
}

func (this *MnStrategy) SetHessianStepTolerance(stp float64) {
	this.theHessTlrStp = stp
}

func (this *MnStrategy) SetHessianG2Tolerance(toler float64) {
	this.theHessTlrG2 = toler
}

func (this *MnStrategy) SetHessianGradientNCycles(n int) {
	this.theHessNCyc = n
}
