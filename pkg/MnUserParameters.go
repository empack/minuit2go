package minuit

type MnUserParameters struct {
	TheTransformation *MnUserTransformation
}

func NewMnUserParameters() *MnUserParameters {
	return &MnUserParameters{
		TheTransformation: NewMnUserTransformation(),
	}
}

func NewMnUserParametersWithParErr(par, err []float64) *MnUserParameters {
	return &MnUserParameters{
		TheTransformation: NewMnUserTransformationFrom(par, err),
	}
}

func NewEmptyMnUserParameters() *MnUserParameters {
	return &MnUserParameters{
		TheTransformation: NewMnUserTransformation(),
	}
}

func (p *MnUserParameters) Clone() *MnUserParameters {
	return &MnUserParameters{
		TheTransformation: p.TheTransformation.Clone(),
	}
}

// ? Idk about this one
/* private MnUserParameters(MnUserParameters other){
	theTransformation = other.theTransformation.clone();
} */

func (p *MnUserParameters) Trafo() *MnUserTransformation {
	return p.TheTransformation
}

func (p *MnUserParameters) VariableParameters() int {
	return p.TheTransformation.variableParameters()
}

// access to parameters (row-wise)
func (p *MnUserParameters) Parameters() []*MinuitParameter {
	return p.TheTransformation.parameters()
}

/** access to parameters and errors in column-wise representation */
func (p *MnUserParameters) Params() []float64 {
	return p.TheTransformation.params()
}

func (p *MnUserParameters) Errors() []float64 {
	return p.TheTransformation.errors()
}

/** access to single parameter */
/* MinuitParameter parameter(int index)
   {
      return theTransformation.parameter(index);
   } */
func (p *MnUserParameters) Parameter(index int) *MinuitParameter {
	return p.TheTransformation.parameter(index)
}

/**
 * Add free parameter name, value, error
 * <p>
 * When adding parameters, MINUIT assigns indices to each parameter which will be
 * the same as in the double[] in the FCNBase.valueOf(). That means the
 * first parameter the user adds gets index 0, the second index 1, and so on. When
 * calculating the function value inside FCN, MINUIT will call FCNBase.valueOf() with
 * the elements at their respective positions.
 */

func (p *MnUserParameters) AddFree(name string, val, err float64) {
	p.TheTransformation.addFree(name, val, err)
}

/**
 * Add limited parameter name, value, lower bound, upper bound
 */

func (p *MnUserParameters) AddLimited(name string, val, err, low, up float64) {
	p.TheTransformation.addLimited(name, val, err, low, up)
}

/**
 * Add const parameter name, value
 */

func (p *MnUserParameters) Add(name string, val float64) {
	p.TheTransformation.add(name, val)
}

/// interaction via external number of parameter
/**
 * Fixes the specified parameter (so that the minimizer will no longer vary it)
 */

func (p *MnUserParameters) Fix(index int) {
	p.TheTransformation.fix(index)
}

/**
 * Releases the specified parameter (so that the minimizer can vary it)
 */

func (p *MnUserParameters) Release(index int) {
	p.TheTransformation.release(index)
}

/**
 * Set the value of parameter. The parameter in
 * question may be variable, fixed, or constant, but must be defined.
 */

func (p *MnUserParameters) SetValue(index int, val float64) {
	p.TheTransformation.setValue(index, val)
}

func (p *MnUserParameters) SetError(index int, err float64) {
	p.TheTransformation.setError(index, err)
}

/**
 * Set the lower and upper bound on the specified variable.
 */

func (p *MnUserParameters) SetLimits(index int, low, up float64) {
	p.TheTransformation.setLimits(index, low, up)
}

func (p *MnUserParameters) SetUpperLimit(index int, up float64) {
	p.TheTransformation.setUpperLimit(index, up)
}

func (p *MnUserParameters) SetLowerLimit(index int, low float64) {
	p.TheTransformation.setLowerLimit(index, low)
}

func (p *MnUserParameters) RemoveLimits(index int) {
	p.TheTransformation.removeLimits(index)
}

func (p *MnUserParameters) Value(index int) float64 {
	return p.TheTransformation.value(index)
}

func (p *MnUserParameters) Error(index int) float64 {
	return p.TheTransformation.error(index)
}

/// interaction via name of parameter
/**
 * Fixes the specified parameter (so that the minimizer will no longer vary it)
 */

func (p *MnUserParameters) FixByName(name string) {
	p.TheTransformation.fixByName(name)
}

/**
 * Releases the specified parameter (so that the minimizer can vary it)
 */

func (p *MnUserParameters) ReleaseByName(name string) {
	p.TheTransformation.releaseByName(name)
}

/**
 * Set the value of parameter. The parameter in
 * question may be variable, fixed, or constant, but must be defined.
 */

func (p *MnUserParameters) SetValueByName(name string, val float64) {
	p.TheTransformation.setValueByName(name, val)
}

func (p *MnUserParameters) SetErrorByName(name string, err float64) {
	p.TheTransformation.setErrorByName(name, err)
}

/**
 * Set the lower and upper bound on the specified variable.
 */

func (p *MnUserParameters) SetLimitsByName(name string, low, up float64) {
	p.TheTransformation.setLimitsByName(name, low, up)
}

func (p *MnUserParameters) SetUpperLimitByName(name string, up float64) {
	p.TheTransformation.setUpperLimitByName(name, up)
}

func (p *MnUserParameters) SetLowerLimitByName(name string, low float64) {
	p.TheTransformation.setLowerLimitByName(name, low)
}

func (p *MnUserParameters) RemoveLimitsByName(name string) {
	p.TheTransformation.removeLimitsByName(name)
}

func (p *MnUserParameters) ValueByName(name string) float64 {
	return p.TheTransformation.valueByName(name)
}

func (p *MnUserParameters) ErrorByName(name string) float64 {
	return p.TheTransformation.errorByName(name)
}

// convert name into external number of parameter
func (p *MnUserParameters) Index(name string) int {
	return p.TheTransformation.index(name)
}

// convert external number into name of parameter
func (p *MnUserParameters) Name(index int) string {
	return p.TheTransformation.name(index)
}

func (p *MnUserParameters) Precision() *MnMachinePrecision {
	return p.TheTransformation.precision()
}

func (p *MnUserParameters) SetPrecision(eps float64) {
	p.TheTransformation.setPrecision(eps)
}

func (p *MnUserParameters) String() string {
	return MnPrint.toStringMnUserParameters(p)
}
