package minuit

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

// MnUserTransformation
/*
 * knows how to transform between user specified parameters (external) and
 * internal parameters used for minimization
 */
type MnUserTransformation struct {
	thePrecision  *MnMachinePrecision
	theParameters []*MinuitParameter
	theExtOfInt   []int
	nameMap       map[string]int

	theDoubleLimTrafo *SinParameterTransformation
	theUpperLimTrafo  *SqrtUpParameterTransformation
	theLowerLimTrafo  *SqrtLowParameterTransformation

	theCache []float64
}

func NewMnUserTransformation() *MnUserTransformation {
	return &MnUserTransformation{
		thePrecision:      NewMnMachinePrecision(),
		theParameters:     make([]*MinuitParameter, 0),
		theExtOfInt:       make([]int, 0),
		theDoubleLimTrafo: NewSinParameterTransformation(),
		theUpperLimTrafo:  NewSqrtUpParameterTransformation(),
		theLowerLimTrafo:  NewSqrtLowParameterTransformation(),
		theCache:          make([]float64, 0),
		nameMap:           make(map[string]int),
	}
}

func (this *MnUserTransformation) Clone() *MnUserTransformation {
	return &MnUserTransformation{
		thePrecision:      this.thePrecision,
		theParameters:     slices.Clone(this.theParameters),
		theExtOfInt:       slices.Clone(this.theExtOfInt),
		theDoubleLimTrafo: NewSinParameterTransformation(),
		theUpperLimTrafo:  NewSqrtUpParameterTransformation(),
		theLowerLimTrafo:  NewSqrtLowParameterTransformation(),
		theCache:          slices.Clone(this.theCache),
		nameMap:           make(map[string]int),
	}
}

func (this *MnUserTransformation) transform(pstates *MnAlgebraicVector) *MnAlgebraicVector {
	// FixMe: Worry about efficiency here
	var result *MnAlgebraicVector = NewMnAlgebraicVector(len(this.theCache))
	for i := 0; i < result.size(); i++ {
		result.set(i, this.theCache[i])
	}
	for i := 0; i < pstates.size(); i++ {
		if this.theParameters[this.theExtOfInt[i]].HasLimits() {
			result.set(this.theExtOfInt[i], this.int2ext(i, pstates.get(i)))
		} else {
			result.set(this.theExtOfInt[i], pstates.get(i))
		}
	}
	return result
}

func NewMnUserTransformationFrom(par, err []float64) *MnUserTransformation {
	res := &MnUserTransformation{
		thePrecision:      NewMnMachinePrecision(),
		theDoubleLimTrafo: NewSinParameterTransformation(),
		theUpperLimTrafo:  NewSqrtUpParameterTransformation(),
		theLowerLimTrafo:  NewSqrtLowParameterTransformation(),
		theParameters:     make([]*MinuitParameter, len(par)),
		theExtOfInt:       make([]int, len(par)),
		theCache:          make([]float64, 0),
		nameMap:           make(map[string]int),
	}
	for i := 0; i < len(par); i++ {
		res.addFree(fmt.Sprintf("i+%d", i), par[i], err[i])
	}
	return res
}

// forwarded interface
func (this *MnUserTransformation) precision() *MnMachinePrecision {
	return this.thePrecision
}

func (this *MnUserTransformation) setPrecision(eps float64) {
	this.thePrecision.SetPrecision(eps)
}

func (this *MnUserTransformation) parameters() []*MinuitParameter {
	return this.theParameters
}

func (this *MnUserTransformation) variableParameters() int {
	return len(this.theExtOfInt)
}

// access to parameters and errors in column-wise representation
func (this *MnUserTransformation) params() []float64 {
	var result []float64 = make([]float64, len(this.theParameters))
	for i := 0; i < len(this.theParameters); i++ {
		result[i] = this.theParameters[i].Value()
	}
	return result
}
func (this *MnUserTransformation) errors() []float64 {
	var result []float64 = make([]float64, len(this.theParameters))
	for i := 0; i < len(this.theParameters); i++ {
		result[i] = this.theParameters[i].Error()
	}
	return result
}

/** access to single parameter */
func (this *MnUserTransformation) parameter(index int) *MinuitParameter {
	return this.theParameters[index]
}

/** addFree add free parameter */
func (this *MnUserTransformation) addFree(name string, val, err float64) error {
	if _, ok := this.nameMap[name]; ok == true {
		return errors.New("IllegalArgumentException: duplicate name:" + name)
	}
	this.nameMap[name] = len(this.theParameters)
	this.theExtOfInt = append(this.theExtOfInt, len(this.theParameters))
	this.theCache = append(this.theCache, val)
	this.theParameters = append(this.theParameters, NewMinuitParameterStandard(len(this.theParameters), name, val, err))
	return nil
}

/** add limited parameter */
func (this *MnUserTransformation) addLimited(name string, val, err, low, up float64) error {
	if _, ok := this.nameMap[name]; ok == true {
		return errors.New("IllegalArgumentException: duplicate name:" + name)
	}
	this.nameMap[name] = len(this.theParameters)
	this.theExtOfInt = append(this.theExtOfInt, len(this.theParameters))
	this.theCache = append(this.theCache, val)
	if param, paramErr := NewMinuitParameterLimited(len(this.theParameters), name, val, err, low, up); paramErr != nil {
		return paramErr
	} else {
		this.theParameters = append(this.theParameters, param)
	}
	return nil
}

/** add  parameter */
func (this *MnUserTransformation) add(name string, val float64) error {
	if _, ok := this.nameMap[name]; ok == true {
		return errors.New("IllegalArgumentException: duplicate name:" + name)
	}
	this.nameMap[name] = len(this.theParameters)
	this.theCache = append(this.theCache, val)
	this.theParameters = append(this.theParameters, NewMinuitParameterConstant(len(this.theParameters), name, val))
	return nil
}

/** interaction via external number of parameter */
func (this *MnUserTransformation) fix(index int) error {
	iind, err := this.intOfExt(index)
	if err != nil {
		return err
	}
	this.theExtOfInt = append(this.theExtOfInt[:iind], this.theExtOfInt[iind+1:]...)
	this.theParameters[index].Fix()
	return nil
}

func (this *MnUserTransformation) release(index int) error {
	if slices.Contains(this.theExtOfInt, index) {
		return errors.New("IllegalArgumentException: index=" + fmt.Sprint(index))
	}
	this.theExtOfInt = append(this.theExtOfInt, index)
	slices.Sort(this.theExtOfInt)
	this.theParameters[index].Release()
	return nil
}

func (this *MnUserTransformation) setValue(index int, val float64) {
	this.theParameters[index].SetValue(val)
	this.theCache[index] = val
}

func (this *MnUserTransformation) setError(index int, err float64) {
	this.theParameters[index].SetError(err)
}

func (this *MnUserTransformation) setLimits(index int, low, up float64) {
	this.theParameters[index].SetLimits(low, up)
}

func (this *MnUserTransformation) setUpperLimit(index int, up float64) {
	this.theParameters[index].SetUpperLimit(up)
}

func (this *MnUserTransformation) setLowerLimit(index int, low float64) {
	this.theParameters[index].SetLowerLimit(low)
}

func (this *MnUserTransformation) removeLimits(index int) {
	this.theParameters[index].RemoveLimits()
}

func (this *MnUserTransformation) value(index int) float64 {
	return this.theParameters[index].Value()
}

func (this *MnUserTransformation) error(index int) float64 {
	return this.theParameters[index].Error()
}

/** interaction via name of parameter */
func (this *MnUserTransformation) fixByName(name string) {
	this.fix(this.index(name))
}

func (this *MnUserTransformation) releaseByName(name string) error {
	return this.release(this.index(name))
}

func (this *MnUserTransformation) setValueByName(name string, val float64) {
	this.setValue(this.index(name), val)
}

func (this *MnUserTransformation) setErrorByName(name string, err float64) {
	this.setError(this.index(name), err)
}

func (this *MnUserTransformation) setLimitsByName(name string, low, up float64) {
	this.setLimits(this.index(name), low, up)
}

func (this *MnUserTransformation) setLowerLimitByName(name string, low float64) {
	this.setLowerLimit(this.index(name), low)
}

func (this *MnUserTransformation) setUpperLimitByName(name string, up float64) {
	this.setUpperLimit(this.index(name), up)
}

func (this *MnUserTransformation) removeLimitsByName(name string) {
	this.removeLimits(this.index(name))
}

func (this *MnUserTransformation) valueByName(name string) float64 {
	return this.value(this.index(name))
}
func (this *MnUserTransformation) errorByName(name string) float64 {
	return this.error(this.index(name))
}

/** convert name into external number of parameter */
func (this *MnUserTransformation) index(name string) int {
	if val, ok := this.nameMap[name]; ok == true {
		return val
	} else {
		return -1
	}
}

/** convert external number into name of parameter */
func (this *MnUserTransformation) name(index int) string {
	return this.theParameters[index].Name()
}

func (this *MnUserTransformation) int2ext(i int, val float64) float64 {
	var parm *MinuitParameter = this.theParameters[this.theExtOfInt[i]]
	if parm.HasLimits() {
		if parm.HasUpperLimit() && parm.HasLowerLimit() {
			return this.theDoubleLimTrafo.int2ext(val, parm.UpperLimit(), parm.LowerLimit())
		} else if parm.HasUpperLimit() && !parm.HasLowerLimit() {
			return this.theUpperLimTrafo.int2ext(val, parm.UpperLimit())
		} else {
			return this.theLowerLimTrafo.int2ext(val, parm.LowerLimit())
		}
	}
	return val
}

func (this *MnUserTransformation) int2extError(i int, val, err float64) float64 {
	dx := err
	var parm *MinuitParameter = this.theParameters[this.theExtOfInt[i]]
	if parm.HasLimits() {
		var ui float64 = this.int2ext(i, val)
		var du1 float64 = this.int2ext(i, val+dx) - ui
		var du2 float64 = this.int2ext(i, val-dx) - ui
		if parm.HasUpperLimit() && parm.HasLowerLimit() {
			if dx > 1.0 {
				du1 = parm.UpperLimit() - parm.LowerLimit()
			}
			dx = 0.5 * (math.Abs(du1) + math.Abs(du2))
		} else {
			dx = 0.5 * (math.Abs(du1) + math.Abs(du2))
		}
	}
	return dx
}

func (this *MnUserTransformation) int2extCovariance(vec *MnAlgebraicVector, cov *MnAlgebraicSymMatrix) (*MnUserCovariance, error) {
	var result *MnUserCovariance = NewMnUserCovarianceWithNrow(cov.nrow())
	for i := 0; i < vec.size(); i++ {
		var dxdi float64 = 1.0
		if this.theParameters[this.theExtOfInt[i]].HasLimits() {
			dxdi = this.dInt2Ext(i, vec.get(i))
		}
		for j := i; j < vec.size(); j++ {
			var dxdj float64 = 1.0
			if this.theParameters[this.theExtOfInt[j]].HasLimits() {
				dxdj = this.dInt2Ext(j, vec.get(j))
			}
			v_, fnErr := cov.get(i, j)
			if fnErr != nil {
				return nil, fnErr
			}
			result.Set(i, j, dxdi*v_*dxdj)
		}
	}
	return result, nil
}

func (this *MnUserTransformation) ext2int(i int, val float64) float64 {
	var parm *MinuitParameter = this.theParameters[i]
	if parm.HasLimits() {
		if parm.HasUpperLimit() && parm.HasLowerLimit() {
			return this.theDoubleLimTrafo.ext2int(val, parm.UpperLimit(), parm.LowerLimit(), this.precision())
		} else if parm.HasUpperLimit() && !parm.HasLowerLimit() {
			return this.theUpperLimTrafo.ext2int(val, parm.UpperLimit(), this.precision())
		} else {
			return this.theLowerLimTrafo.ext2int(val, parm.LowerLimit(), this.precision())
		}
	}
	return val
}

func (this *MnUserTransformation) dInt2Ext(i int, val float64) float64 {
	var dd float64 = 1.0
	var parm *MinuitParameter = this.theParameters[this.theExtOfInt[i]]
	if parm.HasLimits() {
		if parm.HasUpperLimit() && parm.HasLowerLimit() {
			dd = this.theDoubleLimTrafo.dInt2Ext(val, parm.UpperLimit(), parm.LowerLimit())
		} else if parm.HasUpperLimit() && !parm.HasLowerLimit() {
			dd = this.theUpperLimTrafo.dInt2Ext(val, parm.UpperLimit())
		} else {
			dd = this.theLowerLimTrafo.dInt2Ext(val, parm.LowerLimit())
		}
	}
	return dd
}

func (this *MnUserTransformation) intOfExt(ext int) (int, error) {
	for iind := 0; iind < len(this.theExtOfInt); iind++ {
		if ext == this.theExtOfInt[iind] {
			return iind, nil
		}
	}
	return 0, errors.New("IllegalArgumentException: ext=" + fmt.Sprintf("%d", ext))
}
func (this *MnUserTransformation) extOfInt(internal int) int {
	return this.theExtOfInt[internal]
}
