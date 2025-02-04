package minuit

import (
	"errors"
	"math"
)

type MinuitParameter struct {
	theNum        int
	theName       string
	theValue      float64
	theError      float64
	theConst      bool
	theFix        bool
	theLoLimit    float64
	theUpLimit    float64
	theLoLimValid bool
	theUpLimValid bool
}

// NewMinuitParameterConstant
/* constructor for constant parameter */
func NewMinuitParameterConstant(num int, name string, val float64) *MinuitParameter {
	return &MinuitParameter{
		theNum:   num,
		theName:  name,
		theValue: val,
		theConst: true,
	}
}

// NewMinuitParameterStandard
/* constructor for standard parameter */
func NewMinuitParameterStandard(num int, name string, val, err float64) *MinuitParameter {
	return &MinuitParameter{
		theNum:   num,
		theName:  name,
		theValue: val,
		theError: err,
	}
}

/** constructor for limited parameter */
func NewMinuitParameterLimited(num int, name string, val, err, min, max float64) (*MinuitParameter, error) {
	if min == max {
		return nil, errors.New("IllegalArgumentException: min == max")
	}
	return &MinuitParameter{
		theNum:        num,
		theName:       name,
		theValue:      val,
		theError:      err,
		theLoLimit:    math.Min(min, max),
		theUpLimit:    math.Max(min, max),
		theLoLimValid: true,
		theUpLimValid: true,
	}, nil
}

func (this *MinuitParameter) Clone() *MinuitParameter {
	return &MinuitParameter{
		theNum:        this.theNum,
		theName:       this.theName,
		theValue:      this.theValue,
		theError:      this.theError,
		theConst:      this.theConst,
		theFix:        this.theFix,
		theLoLimit:    this.theLoLimit,
		theUpLimit:    this.theUpLimit,
		theLoLimValid: this.theLoLimValid,
		theUpLimValid: this.theUpLimValid,
	}
}

// access methods
func (this *MinuitParameter) Number() int {
	return this.theNum
}
func (this *MinuitParameter) Name() string {
	return this.theName
}
func (this *MinuitParameter) Value() float64 {
	return this.theValue
}
func (this *MinuitParameter) Error() float64 {
	return this.theError
}

// interaction
func (this *MinuitParameter) SetValue(val float64) {
	this.theValue = val
}
func (this *MinuitParameter) SetError(err float64) {
	this.theError = err
}
func (this *MinuitParameter) SetLimits(low, up float64) error {
	if low == up {
		return errors.New("IllegalArgumentException: min == max")
	}
	this.theLoLimit = math.Min(low, up)
	this.theUpLimit = math.Max(low, up)
	this.theLoLimValid = true
	this.theUpLimValid = true
	return nil
}

func (this *MinuitParameter) SetUpperLimit(up float64) {
	this.theLoLimit = 0.0
	this.theUpLimit = up
	this.theLoLimValid = false
	this.theUpLimValid = true
}

func (this *MinuitParameter) SetLowerLimit(low float64) {
	this.theLoLimit = low
	this.theUpLimit = 0.0
	this.theLoLimValid = true
	this.theUpLimValid = false
}

func (this *MinuitParameter) RemoveLimits() {
	this.theLoLimit = 0.0
	this.theUpLimit = 0.0
	this.theLoLimValid = false
	this.theUpLimValid = false
}

func (this *MinuitParameter) Fix() {
	this.theFix = true
}

func (this *MinuitParameter) Release() {
	this.theFix = false
}

// state of parameter (fixed/const/limited)
func (this *MinuitParameter) IsConst() bool {
	return this.theConst
}
func (this *MinuitParameter) IsFixed() bool {
	return this.theFix
}

func (this *MinuitParameter) HasLimits() bool {
	return this.theLoLimValid || this.theUpLimValid
}

func (this *MinuitParameter) HasLowerLimit() bool {
	return this.theLoLimValid
}

func (this *MinuitParameter) HasUpperLimit() bool {
	return this.theUpLimValid
}
func (this *MinuitParameter) LowerLimit() float64 {
	return this.theLoLimit
}
func (this *MinuitParameter) UpperLimit() float64 {
	return this.theUpLimit
}
