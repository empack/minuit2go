package minuit

import (
	"fmt"
	"math"
	"slices"
)

//  public class MnUserParameterState
//  {
// 	private boolean theValid;
// 	private boolean theCovarianceValid;
// 	private boolean theGCCValid;

// 	private double theFVal;
// 	private double theEDM;
// 	private int theNFcn;

// 	private MnUserParameters theParameters;
// 	private MnUserCovariance theCovariance;
// 	private MnGlobalCorrelationCoeff theGlobalCC;

// private List<Double> theIntParameters;
// private MnUserCovariance theIntCovariance;
type MnUserParameterState struct {
	theValid           bool
	theCovarianceValid bool
	theGCCValid        bool

	theFVal float64
	theEDM  float64
	theNFcn int

	theParameters *MnUserParameters
	theCovariance *MnUserCovariance
	theGlobalCC   *MnGlobalCorrelationCoeff

	theIntParameters []float64
	theIntCovariance *MnUserCovariance
}

// MnUserParameterState()
//
//	{
//	   theValid = false;
//	   theCovarianceValid = false;
//	   theParameters = new MnUserParameters();
//	   theCovariance = new MnUserCovariance();
//	   theIntParameters = new ArrayList<Double>();
//	   theIntCovariance =  new MnUserCovariance();
//	}
func NewMnUserParameterState() *MnUserParameterState {
	return &MnUserParameterState{
		theValid:           false,
		theCovarianceValid: false,

		theParameters: NewMnUserParameters(),
		theCovariance: NewMnUserCovariance(),

		theIntParameters: make([]float64, 0),
		theIntCovariance: NewMnUserCovariance(),
	}
}

// protected MnUserParameterState clone()
//
//	{
//	   return new MnUserParameterState(this);
//	}
func (this *MnUserParameterState) clone() *MnUserParameterState {
	return copyMnUserParameterState(this)
}

// 	private MnUserParameterState(MnUserParameterState other)
// 	{
// 	   theValid = other.theValid;
// 	   theCovarianceValid = other.theCovarianceValid;
// 	   theGCCValid = other.theGCCValid;

// 	   theFVal = other.theFVal;
// 	   theEDM = other.theEDM;
// 	   theNFcn = other.theNFcn;

// 	   theParameters = other.theParameters.clone();
// 	   theCovariance = other.theCovariance;
// 	   theGlobalCC = other.theGlobalCC;

//	   theIntParameters = new ArrayList<Double>(other.theIntParameters);
//	   theIntCovariance = other.theIntCovariance.clone();
//	}
func copyMnUserParameterState(other *MnUserParameterState) *MnUserParameterState {
	ups := &MnUserParameterState{
		theValid:           other.theValid,
		theCovarianceValid: other.theCovarianceValid,
		theGCCValid:        other.theGCCValid,

		theFVal: other.theFVal,
		theEDM:  other.theEDM,
		theNFcn: other.theNFcn,

		theParameters: other.theParameters.Clone(),
		//todo: handle copy
		theCovariance: other.theCovariance.clone(),
		theGlobalCC:   other.theGlobalCC.clone(),

		theIntCovariance: other.theIntCovariance.clone(),
	}
	copy(ups.theIntParameters, other.theIntParameters)
	return ups
}

// 	MnUserParameterState(double[] par, double[] err)
// 	{
// 	   theValid = true;
// 	   theParameters = new MnUserParameters(par, err);
// 	   theCovariance = new MnUserCovariance();
// 	   theGlobalCC = new MnGlobalCorrelationCoeff();
// 	   theIntParameters = new ArrayList<Double>(par.length);
// 	   for (int i=0; i<par.length; i++) theIntParameters.add(par[i]);
// 	   theIntCovariance = new MnUserCovariance();
// 	}

/** construct from user parameters (before minimization) */
func NewUserParamStateFromParamAndErrValues(par []float64, err []float64) *MnUserParameterState {
	ups := &MnUserParameterState{
		theValid:         true,
		theParameters:    NewMnUserParametersWithParErr(par, err),
		theCovariance:    NewMnUserCovariance(),
		theGlobalCC:      NewMnGlobalCorrelationCoeff(),
		theIntCovariance: NewMnUserCovariance(),
		theIntParameters: make([]float64, len(par)),
	}
	copy(ups.theIntParameters, par)
	return ups
}

// 	MnUserParameterState(MnUserParameters par)
// 	{
// 	   theValid = true;
// 	   theParameters = par;
// 	   theCovariance = new MnUserCovariance();
// 	   theGlobalCC = new MnGlobalCorrelationCoeff();
// 	   theIntParameters = new ArrayList<Double>(par.variableParameters());
// 	   theIntCovariance = new MnUserCovariance();

//	   int i = 0;
//	   for (MinuitParameter ipar : par.parameters())
//	   {
//		  if (ipar.isConst() || ipar.isFixed()) continue;
//		  if (ipar.hasLimits())
//			 theIntParameters.add(ext2int(ipar.number(),ipar.value()));
//		  else
//			 theIntParameters.add(ipar.value());
//	   }
//	}
func NewUserParameterStateFromUserParameter(par *MnUserParameters) *MnUserParameterState {

	ups := &MnUserParameterState{
		theValid:         true,
		theParameters:    par,
		theCovariance:    NewMnUserCovariance(),
		theGlobalCC:      NewMnGlobalCorrelationCoeff(),
		theIntParameters: make([]float64, par.VariableParameters()),
		theIntCovariance: NewMnUserCovariance(),
	}

	pos := 0
	for _, ipar := range par.Parameters() {
		if ipar.IsConst() || ipar.IsFixed() {
			continue
		}
		if ipar.HasLimits() {
			ups.theIntParameters[pos] = ups.ext2int(ipar.Number(), ipar.Value())
		} else {
			ups.theIntParameters[pos] = ipar.Value()
		}
		pos += 1
	}
	return ups
}

// 	MnUserParameterState(double[] par, double[] cov, int nrow)
// 	{
// 	   theValid = true;
// 	   theCovarianceValid = true;
// 	   theCovariance = new MnUserCovariance(cov, nrow);
// 	   theGlobalCC = new MnGlobalCorrelationCoeff();
// 	   theIntParameters = new ArrayList<Double>(par.length);
// 	   theIntCovariance = new MnUserCovariance(cov, nrow);

// 	   double[] err = new double[par.length];
// 	   for (int i = 0; i < par.length; i++)
// 	   {
// 		  assert(theCovariance.get(i,i) > 0.);
// 		  err[i] = Math.sqrt(theCovariance.get(i,i));
// 		  theIntParameters.add(par[i]);
// 	   }
// 	   theParameters = new MnUserParameters(par, err);
// 	   assert(theCovariance.nrow() == variableParameters());
// 	}
/** construct from user parameters + covariance (before minimization) */
func NewMnUserParameterStateFlFlIn(par []float64, cov []float64, nrow int) (*MnUserParameterState, error) {
	ups := &MnUserParameterState{
		theValid:         true,
		theParameters:    nil,
		theCovariance:    NewMnUserCovarianceWithDataNrow(cov, nrow),
		theGlobalCC:      NewMnGlobalCorrelationCoeff(),
		theIntCovariance: NewMnUserCovarianceWithDataNrow(cov, nrow),
	}

	ups.theIntParameters = make([]float64, len(par))
	var retErr error = nil
	err := make([]float64, len(par))
	for i := 0; i < len(par); i++ {
		if ups.theCovariance.Get(i, i) <= 0 {
			retErr = fmt.Errorf("covariance must be positive")
		}
		err[i] = math.Sqrt(ups.theCovariance.Get(i, i))
		ups.theIntParameters[i] = par[i]
	}
	ups.theParameters = NewMnUserParametersWithParErr(par, err)

	if ups.theCovariance.Nrow() != ups.VariableParameters() {
		retErr = fmt.Errorf("missmatch in covariance<->parameters")
	}
	return ups, retErr
}

// 	MnUserParameterState(double[] par, MnUserCovariance cov)
// 	{
// 	   theValid = true;
// 	   theCovarianceValid = true;
// 	   theCovariance = cov;
// 	   theGlobalCC = new MnGlobalCorrelationCoeff();
// 	   theIntParameters = new ArrayList<Double>(par.length);
// 	   theIntCovariance = cov.clone();

//	   if (theCovariance.nrow() != variableParameters()) throw new IllegalArgumentException("Bad covariance size");
//	   double[] err = new double[par.length];
//	   for (int i = 0; i < par.length; i++)
//	   {
//		  if (theCovariance.get(i,i) <= 0.) throw new IllegalArgumentException("Bad covariance");
//		  err[i] = Math.sqrt(theCovariance.get(i,i));
//		  theIntParameters.add(par[i]);
//	   }
//	   theParameters = new MnUserParameters(par, err);
//	}
func NewMnUserParameterStateFlUc(par []float64, cov *MnUserCovariance) (*MnUserParameterState, error) {

	ups := &MnUserParameterState{
		theValid:           true,
		theCovarianceValid: true,
		theCovariance:      cov,
		theGlobalCC:        NewMnGlobalCorrelationCoeff(),
		theIntCovariance:   cov.clone(),
	}

	ups.theIntParameters = make([]float64, len(par))
	if ups.theCovariance.Nrow() != ups.VariableParameters() {
		return nil, fmt.Errorf("Bad covariance size")
	}
	err := make([]float64, len(par))
	for i := 0; i < len(par); i++ {
		if ups.theCovariance.Get(i, i) <= 0 {
			return nil, fmt.Errorf("covariance must be positive")
		}
		err[i] = math.Sqrt(ups.theCovariance.Get(i, i))
		ups.theIntParameters = append(ups.theIntParameters, par[i])
	}
	ups.theParameters = NewMnUserParametersWithParErr(par, err)
	return ups, nil
}

// 	MnUserParameterState(MnUserParameters par, MnUserCovariance cov)
// 	{
// 	   theValid = true;
// 	   theCovarianceValid = true;
// 	   theParameters = par;
// 	   theCovariance = cov;
// 	   theGlobalCC = new MnGlobalCorrelationCoeff();
// 	   theIntParameters = new ArrayList<Double>();
// 	   theIntCovariance = cov.clone();

//	   theIntCovariance.scale(0.5);
//	   int i = 0;
//	   for (MinuitParameter ipar : par.parameters())
//	   {
//		  if(ipar.isConst() || ipar.isFixed()) continue;
//		  if(ipar.hasLimits())
//			 theIntParameters.add(ext2int(ipar.number(), ipar.value()));
//		  else
//			 theIntParameters.add(ipar.value());
//	   }
//	   assert(theCovariance.nrow() == variableParameters());
//	}
func NewUserParamStateFromUserParamCovariance(par *MnUserParameters, cov *MnUserCovariance) (*MnUserParameterState, error) {
	ups := &MnUserParameterState{
		theValid:           true,
		theCovarianceValid: true,
		theParameters:      par,
		theCovariance:      cov,
		theGlobalCC:        NewMnGlobalCorrelationCoeff(),
		theIntCovariance:   cov.clone(),
	}

	ups.theIntParameters = make([]float64, par.VariableParameters())
	ups.theIntCovariance.scale(0.5)
	pos := 0
	for _, ipar := range par.Parameters() {
		if ipar.IsConst() || ipar.IsFixed() {
			continue
		}
		if ipar.HasLimits() {
			ups.theIntParameters[pos] = ups.ext2int(ipar.Number(), ipar.Value())
		} else {
			ups.theIntParameters[pos] = ipar.Value()
		}
	}
	if ups.theCovariance.Nrow() != ups.VariableParameters() {
		return ups, fmt.Errorf("Bad covariance size")
	}
	return ups, nil
}

// 	MnUserParameterState(MinimumState st, double up, MnUserTransformation trafo)
// 	{
// 	   theValid = st.isValid();
// 	   theCovarianceValid = false;
// 	   theGCCValid = false;
// 	   theFVal = st.fval();
// 	   theEDM = st.edm();
// 	   theNFcn = st.nfcn();
// 	   theParameters = new MnUserParameters();
// 	   theCovariance = new MnUserCovariance();
// 	   theGlobalCC = new MnGlobalCorrelationCoeff();
// 	   theIntParameters = new ArrayList<Double>();
// 	   theIntCovariance = new MnUserCovariance();

// 	   for (MinuitParameter ipar : trafo.parameters())
// 	   {
// 		  if(ipar.isConst())
// 		  {
// 			 add(ipar.name(), ipar.value());
// 		  }
// 		  else if(ipar.isFixed())
// 		  {
// 			 add(ipar.name(), ipar.value(), ipar.error());
// 			 if(ipar.hasLimits())
// 			 {
// 				if(ipar.hasLowerLimit() && ipar.hasUpperLimit())
// 				   setLimits(ipar.name(), ipar.lowerLimit(),ipar.upperLimit());
// 				else if(ipar.hasLowerLimit() && !ipar.hasUpperLimit())
// 				   setLowerLimit(ipar.name(), ipar.lowerLimit());
// 				else
// 				   setUpperLimit(ipar.name(), ipar.upperLimit());
// 			 }
// 			 fix(ipar.name());
// 		  }
// 		  else if(ipar.hasLimits())
// 		  {
// 			 int i = trafo.intOfExt(ipar.number());
// 			 double err = st.hasCovariance() ? Math.sqrt(2.*up*st.error().invHessian().get(i,i)) : st.parameters().dirin().get(i);
// 			 add(ipar.name(), trafo.int2ext(i, st.vec().get(i)), trafo.int2extError(i, st.vec().get(i), err));
// 			 if(ipar.hasLowerLimit() && ipar.hasUpperLimit())
// 				setLimits(ipar.name(), ipar.lowerLimit(), ipar.upperLimit());
// 			 else if(ipar.hasLowerLimit() && !ipar.hasUpperLimit())
// 				setLowerLimit(ipar.name(), ipar.lowerLimit());
// 			 else
// 				setUpperLimit(ipar.name(), ipar.upperLimit());
// 		  }
// 		  else
// 		  {
// 			 int i = trafo.intOfExt(ipar.number());
// 			 double err = st.hasCovariance() ? Math.sqrt(2.*up*st.error().invHessian().get(i,i)) : st.parameters().dirin().get(i);
// 			 add(ipar.name(), st.vec().get(i), err);
// 		  }
// 	   }

// 	   theCovarianceValid = st.error().isValid();

// 	   if(theCovarianceValid)
// 	   {
// 		  theCovariance = trafo.int2extCovariance(st.vec(), st.error().invHessian());
// 		  theIntCovariance = new MnUserCovariance(st.error().invHessian().data().clone(), st.error().invHessian().nrow());
// 		  theCovariance.scale(2.*up);
// 		  theGlobalCC = new MnGlobalCorrelationCoeff(st.error().invHessian());
// 		  theGCCValid = true;

// 		  assert(theCovariance.nrow() == variableParameters());
// 	   }

// 	}
/** construct from internal parameters (after minimization) */
func NewMnUserParameterStateMsFlUt(st *MinimumState, up float64, trafo *MnUserTransformation) (*MnUserParameterState, error) {
	ups := &MnUserParameterState{
		theValid:           st.isValid(),
		theCovarianceValid: false,
		theGCCValid:        false,
		theFVal:            st.fval(),
		theEDM:             st.edm(),
		theNFcn:            st.nfcn(),
		theParameters:      NewMnUserParameters(),
		theCovariance:      NewMnUserCovariance(),
		theGlobalCC:        NewMnGlobalCorrelationCoeff(),
		theIntCovariance:   NewMnUserCovariance(),
	}

	var retErr error
	for _, ipar := range trafo.parameters() {
		if ipar.IsConst() {
			ups.AddStFl(ipar.Name(), ipar.Value())
		} else if ipar.IsFixed() {
			ups.AddStFlFl(ipar.Name(), ipar.Value(), ipar.Error())
			if ipar.HasLimits() {
				if ipar.HasLowerLimit() && ipar.HasUpperLimit() {
					ups.SetLimitsStFlFl(ipar.Name(), ipar.LowerLimit(), ipar.UpperLimit())
				} else if ipar.HasLowerLimit() && !ipar.HasUpperLimit() {
					ups.SetLowerLimitStFl(ipar.Name(), ipar.LowerLimit())
				} else {
					ups.SetUpperLimitStFl(ipar.Name(), ipar.UpperLimit())
				}
			}
			ups.FixSt(ipar.Name())
		} else if ipar.HasLimits() {
			i, _ := trafo.intOfExt(ipar.Number())
			var err float64
			if st.hasCovariance() {
				hessII, _ := st.error().invHessian().get(i, i)
				err = math.Sqrt(2 * up * hessII)
			} else {
				st.parameters().dirin().get(i)
			}
			ups.AddStFlFl(ipar.Name(), trafo.int2ext(i, st.vec().get(i)), trafo.int2extError(i, st.vec().get(i), err))
			if ipar.HasLowerLimit() && ipar.HasUpperLimit() {
				ups.SetLimitsStFlFl(ipar.Name(), ipar.LowerLimit(), ipar.UpperLimit())
			} else if ipar.HasLowerLimit() && !ipar.HasUpperLimit() {
				ups.SetLowerLimitStFl(ipar.Name(), ipar.LowerLimit())
			} else {
				ups.SetUpperLimitStFl(ipar.Name(), ipar.UpperLimit())
			}
		} else {
			i, _ := trafo.intOfExt(ipar.Number())
			var err float64
			if st.hasCovariance() {
				hessII, _ := st.error().invHessian().get(i, i)
				math.Sqrt(2 * up * hessII)
			} else {
				st.parameters().dirin().get(i)
			}
			ups.AddStFlFl(ipar.Name(), st.vec().get(i), err)
		}
	}

	ups.theCovarianceValid = st.error().isValid()
	if ups.theCovarianceValid {
		c_, fnErr := trafo.int2extCovariance(st.vec(), st.error().invHessian())
		if fnErr != nil {
			return nil, fnErr
		}
		ups.theCovariance = c_
		dataClone := slices.Clone(st.error().invHessian().data())
		ups.theIntCovariance = NewMnUserCovarianceWithDataNrow(dataClone, st.error().invHessian().nrow())
		ups.theCovariance.scale(2 * up)

		ups.theGlobalCC, _ = MnGlobalCorrelationCoeffFromAlgebraicSymMatrix(st.error().invHessian())
		ups.theGCCValid = true
		if ups.theCovariance.Nrow() != ups.VariableParameters() {
			retErr = fmt.Errorf("dimension missmatch covariance/variable params")
		}
	}

	return ups, retErr
}

// user external representation
func (this *MnUserParameterState) parameters() *MnUserParameters {
	return this.theParameters
}

func (this *MnUserParameterState) covariance() *MnUserCovariance {
	return this.theCovariance
}

func (this *MnUserParameterState) globalCC() *MnGlobalCorrelationCoeff {
	return this.theGlobalCC
}

/** Minuit internal representation */
func (this *MnUserParameterState) intParameters() []float64 {
	return this.theIntParameters
}

func (this *MnUserParameterState) intCovariance() *MnUserCovariance {
	return this.theIntCovariance
}

/** transformation internal <-> external */
func (this *MnUserParameterState) trafo() *MnUserTransformation {
	return this.theParameters.Trafo()
}

/**
 * Returns <CODE>true</CODE> if the the state is valid, <CODE>false</CODE> if not
 */
func (this *MnUserParameterState) IsValid() bool {
	return this.theValid
}

/**
 * Returns <CODE>true</CODE>
 * if the the state has a valid covariance, <CODE>false</CODE> otherwise.
 */
func (this *MnUserParameterState) HasCovariance() bool {
	return this.theCovarianceValid
}

func (this *MnUserParameterState) HasGlobalCC() bool {
	return this.theGCCValid
}

/**
 * returns the function value at the minimum
 */
func (this *MnUserParameterState) Fval() float64 {
	return this.theFVal
}

/**
 * Returns the expected vertival distance to the minimum (EDM)
 */
func (this *MnUserParameterState) Edm() float64 {
	return this.theEDM
}

/**
 * Returns the number of function calls during the minimization.
 */
func (this *MnUserParameterState) Nfcn() int {
	return this.theNFcn
}

// facade: forward interface of MnUserParameters and MnUserTransformation

/** access to parameters (row-wise) */
func (this *MnUserParameterState) MinuitParameters() []*MinuitParameter {
	return this.theParameters.Parameters()
}

/** access to parameters and errors in column-wise representation */
func (this *MnUserParameterState) Params() []float64 {
	return this.theParameters.Params()
}

func (this *MnUserParameterState) errors() []float64 {
	return this.theParameters.Errors()
}

func (this *MnUserParameterState) parameter(i int) *MinuitParameter {
	return this.theParameters.Parameter(i)
}

/** add free parameter name, value, error */
func (this *MnUserParameterState) AddStFlFl(name string, val float64, err float64) {
	this.theParameters.AddFree(name, val, err)
	this.theIntParameters = append(this.theIntParameters, val)
	this.theCovarianceValid = false
	this.theGCCValid = false
	this.theValid = true
}

/** add limited parameter name, value, lower bound, upper bound */
func (this *MnUserParameterState) AddStFlFlFlFl(name string, val float64, err float64, low float64, up float64) {
	this.theParameters.AddLimited(name, val, err, low, up)
	this.theCovarianceValid = false
	this.theIntParameters = append(this.theIntParameters, this.ext2int(this.Index(name), val))
	this.theGCCValid = false
	this.theValid = true
}

/** add const parameter name, value */
func (this *MnUserParameterState) AddStFl(name string, val float64) {
	this.theParameters.Add(name, val)
	this.theValid = true
}

//	public void fix(int e)
//	{
//	   int i = intOfExt(e);
//	   if(theCovarianceValid)
//	   {
//		  theCovariance = MnCovarianceSqueeze.squeeze(theCovariance, i);
//		  theIntCovariance = MnCovarianceSqueeze.squeeze(theIntCovariance, i);
//	   }
//	   theIntParameters.remove(i);
//	   theParameters.fix(e);
//	   theGCCValid = false;
//	}
//
// / interaction via external number of parameter
func (this *MnUserParameterState) Fix(e int) error {
	i, err := this.intOfExt(e)
	if err != nil {
		return err
	}
	if this.theCovarianceValid {
		this.theCovariance, _ = MnCovarianceSqueeze.Squeeze(this.theCovariance, i)
		this.theIntCovariance, _ = MnCovarianceSqueeze.Squeeze(this.theIntCovariance, i)
	}
	remove(this.theIntParameters, i)
	this.theParameters.Fix(e)
	this.theGCCValid = false
	return nil
}

// public void release(int e)
//
//	{
//	   theParameters.release(e);
//	   theCovarianceValid = false;
//	   theGCCValid = false;
//	   int i = intOfExt(e);
//	   if(parameter(e).hasLimits())
//		  theIntParameters.add(i, ext2int(e, parameter(e).value()));
//	   else
//		  theIntParameters.add(i, parameter(e).value());
//	}
func (this *MnUserParameterState) Release(e int) error {
	this.theParameters.Release(e)
	this.theCovarianceValid = false
	this.theGCCValid = false
	i, err := this.intOfExt(e)
	if err != nil {
		return err
	}
	if this.parameter(e).HasLimits() {
		slices.Insert(this.theIntParameters, i, this.ext2int(e, this.parameter(e).Value()))
	} else {
		slices.Insert(this.theIntParameters, i, this.parameter(e).Value())
	}
	return nil
}

// public void setValue(int e, double val)
//
//	{
//	   theParameters.setValue(e, val);
//	   if(!parameter(e).isFixed() && !parameter(e).isConst())
//	   {
//		  int i = intOfExt(e);
//		  if(parameter(e).hasLimits())
//			 theIntParameters.set(i,ext2int(e, val));
//		  else
//			 theIntParameters.set(i,val);
//	   }
//	}
func (this *MnUserParameterState) SetValue(e int, val float64) error {
	this.theParameters.SetValue(e, val)
	if !this.parameter(e).IsFixed() && !this.parameter(e).IsConst() {
		i, err := this.intOfExt(e)
		if err != nil {
			return err
		}
		if this.parameter(e).HasLimits() {
			this.theIntParameters[i] = this.ext2int(e, val)
		} else {
			this.theIntParameters[i] = val
		}
	}
	return nil
}

func (this *MnUserParameterState) SetError(e int, err float64) {
	this.theParameters.SetError(e, err)
}

// public void setLimits(int e, double low, double up)
//
//	{
//	   theParameters.setLimits(e, low, up);
//	   theCovarianceValid = false;
//	   theGCCValid = false;
//	   if(!parameter(e).isFixed() && !parameter(e).isConst())
//	   {
//		  int i = intOfExt(e);
//		  if(low < theIntParameters.get(i) && theIntParameters.get(i) < up)
//			 theIntParameters.set(i,ext2int(e, theIntParameters.get(i)));
//		  else
//			 theIntParameters.set(i,ext2int(e, 0.5*(low+up)));
//	   }
//	}
func (this *MnUserParameterState) SetLimits(e int, low float64, up float64) error {
	this.theParameters.SetLimits(e, low, up)
	this.theCovarianceValid = false
	this.theGCCValid = false
	if !this.parameter(e).IsFixed() && !this.parameter(e).IsConst() {
		i, err := this.intOfExt(e)
		if err != nil {
			return err
		}
		if low < this.theIntParameters[i] && this.theIntParameters[i] < up {
			this.theIntParameters[i] = this.ext2int(e, this.theIntParameters[i])
		} else {
			this.theIntParameters[i] = this.ext2int(e, 0.5*(low+up))
		}
	}
	return nil
}

// public void setUpperLimit(int e, double up)
//
//	{
//	   theParameters.setUpperLimit(e, up);
//	   theCovarianceValid = false;
//	   theGCCValid = false;
//	   if(!parameter(e).isFixed() && !parameter(e).isConst())
//	   {
//		  int i = intOfExt(e);
//		  if(theIntParameters.get(i) < up)
//			 theIntParameters.set(i,ext2int(e, theIntParameters.get(i)));
//		  else
//			 theIntParameters.set(i,ext2int(e, up - 0.5*Math.abs(up + 1.)));
//	   }
//	}
func (this *MnUserParameterState) SetUpperLimit(e int, up float64) error {
	this.theParameters.SetUpperLimit(e, up)
	this.theCovarianceValid = false
	this.theGCCValid = false
	if !this.parameter(e).IsFixed() && !this.parameter(e).IsConst() {
		i, err := this.intOfExt(e)
		if err != nil {
			return err
		}
		if this.theIntParameters[i] < up {
			this.theIntParameters[i] = this.ext2int(e, this.theIntParameters[i])
		} else {
			this.theIntParameters[i] = this.ext2int(e, up-0.5*math.Abs(up+1))
		}
	}
	return nil
}

// public void setLowerLimit(int e, double low)
//
//	{
//	   theParameters.setLowerLimit(e, low);
//	   theCovarianceValid = false;
//	   theGCCValid = false;
//	   if(!parameter(e).isFixed() && !parameter(e).isConst())
//	   {
//		  int i = intOfExt(e);
//		  if(low < theIntParameters.get(i))
//			 theIntParameters.set(i,ext2int(e, theIntParameters.get(i)));
//		  else
//			 theIntParameters.set(i,ext2int(e, low + 0.5*Math.abs(low + 1.)));
//	   }
//	}
func (this *MnUserParameterState) SetLowerLimit(e int, low float64) error {
	this.theParameters.SetLowerLimit(e, low)
	this.theCovarianceValid = false
	this.theGCCValid = false
	if !this.parameter(e).IsFixed() && !this.parameter(e).IsConst() {
		i, err := this.intOfExt(e)
		if err != nil {
			return err
		}
		if low < this.theIntParameters[i] {
			this.theIntParameters[i] = this.ext2int(e, this.theIntParameters[i])
		} else {
			this.theIntParameters[i] = this.ext2int(e, low+0.5*math.Abs(low+1))
		}
	}
	return nil
}

// public void removeLimits(int e)
//
//	{
//	   theParameters.removeLimits(e);
//	   theCovarianceValid = false;
//	   theGCCValid = false;
//	   if(!parameter(e).isFixed() && !parameter(e).isConst())
//		  theIntParameters.set(intOfExt(e),value(e));
//	}
func (this *MnUserParameterState) RemoveLimits(e int) error {
	this.theParameters.RemoveLimits(e)
	this.theCovarianceValid = false
	this.theGCCValid = false
	if !this.parameter(e).IsFixed() && !this.parameter(e).IsConst() {
		i, err := this.intOfExt(e)
		if err != nil {
			return err
		}
		this.theIntParameters[i] = this.Value(e)
	}
	return nil
}

func (this *MnUserParameterState) Value(index int) float64 {
	return this.theParameters.Value(index)
}

func (this *MnUserParameterState) Error(index int) float64 {
	return this.theParameters.Error(index)
}

// / interaction via name of parameter
func (this *MnUserParameterState) FixSt(name string) {
	this.Fix(this.Index(name))
}

func (this *MnUserParameterState) ReleaseSt(name string) {
	this.Release(this.Index(name))
}

func (this *MnUserParameterState) SetValueStFl(name string, val float64) {
	this.SetValue(this.Index(name), val)
}

func (this *MnUserParameterState) SetErrorStFl(name string, err float64) {
	this.SetError(this.Index(name), err)
}

func (this *MnUserParameterState) SetLimitsStFlFl(name string, low float64, up float64) {
	this.SetLimits(this.Index(name), low, up)
}

func (this *MnUserParameterState) SetUpperLimitStFl(name string, up float64) {
	this.SetUpperLimit(this.Index(name), up)
}

func (this *MnUserParameterState) SetLowerLimitStFl(name string, low float64) {
	this.SetLowerLimit(this.Index(name), low)
}

func (this *MnUserParameterState) RemoveLimitsSt(name string) {
	this.RemoveLimits(this.Index(name))
}

func (this *MnUserParameterState) ValueSt(name string) float64 {
	return this.Value(this.Index(name))
}

func (this *MnUserParameterState) ErrorSt(name string) float64 {
	return this.Error(this.Index(name))
}

/** convert name into external number of parameter */
func (this *MnUserParameterState) Index(name string) int {
	return this.theParameters.Index(name)
}

/** convert external number into name of parameter */
func (this *MnUserParameterState) Name(index int) string {
	return this.theParameters.Name(index)
}

// transformation internal <-> external
func (this *MnUserParameterState) int2ext(i int, val float64) float64 {
	return this.theParameters.Trafo().int2ext(i, val)
}
func (this *MnUserParameterState) ext2int(i int, val float64) float64 {
	return this.theParameters.Trafo().ext2int(i, val)
}
func (this *MnUserParameterState) intOfExt(ext int) (int, error) {
	return this.theParameters.Trafo().intOfExt(ext)
}
func (this *MnUserParameterState) ExtOfInt(internal int) int {
	return this.theParameters.Trafo().extOfInt(internal)
}
func (this *MnUserParameterState) VariableParameters() int {
	return this.theParameters.VariableParameters()
}
func (this *MnUserParameterState) Precision() *MnMachinePrecision {
	return this.theParameters.Precision()
}
func (this *MnUserParameterState) SetPrecision(eps float64) {
	this.theParameters.SetPrecision(eps)
}
func (this *MnUserParameterState) ToString() string {
	return MnPrint.toStringMnUserParameterState(this)
}

func remove(slice []float64, s int) []float64 {
	return append(slice[:s], slice[s+1:]...)
}
