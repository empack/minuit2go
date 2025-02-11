package minuit

import (
	"fmt"
	"math"
	"strings"
)

var MnPrint = &MnPrintStruct{}

type MnPrintStruct struct{}

func (this *MnPrintStruct) toStringMnAlgebraicVector(x *MnAlgebraicVector) string {
	var builder strings.Builder
	this.PrintMnAlgebraicVector(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMnAlgebraicVector(builder *strings.Builder, vec *MnAlgebraicVector) {
	builder.WriteString("LAVector parameters:\n\n")
	nrow := vec.size()

	for i := 0; i < nrow; i++ {
		builder.WriteString(fmt.Sprintf("%g ", vec.get(i)))
	}

	builder.WriteString("\n")
}

func (this *MnPrintStruct) toStringMnAlgebraicSymMatrix(x *MnAlgebraicSymMatrix) string {
	var builder strings.Builder
	this.PrintMnAlgebraicSymMatrix(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMnAlgebraicSymMatrix(builder *strings.Builder, matrix *MnAlgebraicSymMatrix) {
	builder.WriteString("LASymMatrix parameters:\n\n")
	n := matrix.nrow()

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			m, _ := matrix.get(i, j)
			builder.WriteString(fmt.Sprintf("%10g ", m))
		}

		builder.WriteString("\n")
	}
}

func (this *MnPrintStruct) ToStringFunctionMinimum(min *FunctionMinimum) string {
	var builder strings.Builder
	this.PrintFunctionMinimum(&builder, min)
	return builder.String()
}

func (this *MnPrintStruct) PrintFunctionMinimum(builder *strings.Builder, min *FunctionMinimum) {
	builder.WriteString("\n")
	if !min.IsValid() {
		builder.WriteString("\n")
		builder.WriteString("WARNING: Minuit did not converge.\n")
		builder.WriteString("\n")
	} else {
		builder.WriteString("\n")
		builder.WriteString("Minuit did successfully converge.\n")
		builder.WriteString("\n")
	}

	builder.WriteString(fmt.Sprintf("# of function calls: %d\n", min.Nfcn()))
	builder.WriteString(fmt.Sprintf("minimum function value: %g\n", min.Fval()))
	builder.WriteString(fmt.Sprintf("minimum edm: %g\n", min.Edm()))
	builder.WriteString("minimum internal state vector: ")
	this.PrintMnAlgebraicVector(builder, min.parameters().vec())
	builder.WriteString("\n")
	if min.hasValidCovariance() {
		builder.WriteString(fmt.Sprintf("minimum internal covariance matrix:  %s\n", min.error().matrix()))
	}

	//builder.WriteString(fmt.Sprintf("%g\n", min.UserParameters()))
	this.PrintMnUserParameters(builder, min.UserParameters())
	//builder.WriteString(fmt.Sprintf("%g\n", min.UserCovariance()))
	this.PrintMnUserCovariance(builder, min.UserCovariance())
	//builder.WriteString(fmt.Sprintf("%g\n", min.UserState().globalCC()))
	this.PrintMnGlobalCorrelationCoeff(builder, min.UserState().globalCC())
	if !min.IsValid() {
		builder.WriteString("WARNING: FunctionMinimum is invalid.\n")
	}

	builder.WriteString("\n")
}

func (this *MnPrintStruct) toStringMinimumState(x *MinimumState) string {
	var builder strings.Builder
	this.PrintMinimumState(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMinimumState(builder *strings.Builder, min *MinimumState) {
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("minimum function value: %g\n", min.fval()))
	builder.WriteString(fmt.Sprintf("minimum edm: %g\n", min.edm()))
	builder.WriteString(fmt.Sprintf("minimum internal state vector:  %g\n", min.vec()))
	builder.WriteString(fmt.Sprintf("minimum internal gradient vector: %g\n", min.gradient().vec()))
	if min.hasCovariance() {
		builder.WriteString(fmt.Sprintf("minimum internal covariance matrix: %s\n", min.error().matrix()))
	}

	builder.WriteString("\n")
}

func (this *MnPrintStruct) toStringMnUserParameters(x *MnUserParameters) string {
	var builder strings.Builder
	this.PrintMnUserParameters(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMnUserParameters(builder *strings.Builder, par *MnUserParameters) {
	builder.WriteString("\n")
	builder.WriteString("# ext. ||   name    ||   type  ||   value   ||  error +/- \n")
	builder.WriteString("\n")

	atLoLim := false
	atHiLim := false

	for _, ipar := range par.Parameters() {
		builder.WriteString(fmt.Sprintf(" %5d || %9s || ", ipar.Number(), ipar.Name()))
		if ipar.IsConst() {
			builder.WriteString(fmt.Sprintf("         || %10g   ||", ipar.Value()))
		} else if ipar.IsFixed() {
			builder.WriteString(fmt.Sprintf("  fixed  || %10g   ||\n", ipar.Value()))
		} else if ipar.HasLimits() {
			if ipar.Error() > 0.0 {
				builder.WriteString(fmt.Sprintf(" limited || %10g", ipar.Value()))
				if math.Abs(ipar.Value()-ipar.LowerLimit()) < par.Precision().eps2() {
					builder.WriteString("* ")
					atLoLim = true
				}

				if math.Abs(ipar.Value()-ipar.UpperLimit()) < par.Precision().eps2() {
					builder.WriteString("**")
					atHiLim = true
				}

				builder.WriteString(fmt.Sprintf(" || %10g\n", ipar.Error()))
			} else {
				builder.WriteString(fmt.Sprintf("  free   || %10g || no\n", ipar.Value()))
			}
		} else if ipar.Error() > 0.0 {
			builder.WriteString(fmt.Sprintf("  free   || %10g || %10g\n", ipar.Value(), ipar.Error()))
		} else {
			builder.WriteString(fmt.Sprintf("  free   || %10g || no\n", ipar.Value()))
		}
	}

	builder.WriteString("\n")
	if atLoLim {
		builder.WriteString("* parameter is at lower limit")
	}

	if atHiLim {
		builder.WriteString("** parameter is at upper limit")
	}

	builder.WriteString("\n")
}

func (this *MnPrintStruct) toStringMnUserCovariance(x *MnUserCovariance) string {
	var builder strings.Builder
	this.PrintMnUserCovariance(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMnUserCovariance(builder *strings.Builder, matrix *MnUserCovariance) {
	builder.WriteString("\n")
	builder.WriteString("MnUserCovariance: \n")
	builder.WriteString("\n")
	n := matrix.Nrow()

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			builder.WriteString(fmt.Sprintf("%10g ", matrix.Get(i, j)))
		}

		builder.WriteString("\n")
	}

	builder.WriteString("\n")
	builder.WriteString("MnUserCovariance parameter correlations: \n")
	builder.WriteString("\n")
	n = matrix.Nrow()

	for i := 0; i < n; i++ {
		di := matrix.Get(i, i)

		for j := 0; j < n; j++ {
			dj := matrix.Get(j, j)
			builder.WriteString(fmt.Sprintf("%g ", matrix.Get(i, j)/math.Sqrt(math.Abs(di*dj))))
		}

		builder.WriteString("\n")
	}
}

func (this *MnPrintStruct) toStringMnGlobalCorrelationCoeff(x *MnGlobalCorrelationCoeff) string {
	var builder strings.Builder
	this.PrintMnGlobalCorrelationCoeff(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMnGlobalCorrelationCoeff(builder *strings.Builder, coeff *MnGlobalCorrelationCoeff) {
	builder.WriteString("\n")
	builder.WriteString("MnGlobalCorrelationCoeff: \n")
	builder.WriteString("\n")

	for i := 0; i < len(coeff.GlobalCC()); i++ {
		builder.WriteString(fmt.Sprintf("%g\n", coeff.GlobalCC()[i]))
	}
}

func (this *MnPrintStruct) toStringMnUserParameterState(x *MnUserParameterState) string {
	var builder strings.Builder
	this.PrintMnUserParameterState(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMnUserParameterState(builder *strings.Builder, state *MnUserParameterState) {
	builder.WriteString("\n")
	if !state.IsValid() {
		builder.WriteString("\n")
		builder.WriteString("WARNING: MnUserParameterState is not valid.\n")
		builder.WriteString("\n")
	}

	builder.WriteString(fmt.Sprintf("# of function calls: %g\n", state.Nfcn()))
	builder.WriteString(fmt.Sprintf("function value: %g\n", state.Fval()))
	builder.WriteString(fmt.Sprintf("expected distance to the minimum (edm): %g\n", state.Edm()))
	builder.WriteString(fmt.Sprintf("external parameters: %g\n", state.parameters()))
	if state.HasCovariance() {
		builder.WriteString(fmt.Sprintf("covariance matrix: %g\n", state.covariance()))
	}

	if state.HasGlobalCC() {
		builder.WriteString(fmt.Sprintf("global correlation coefficients : %g\n", state.globalCC()))
	}

	if state.IsValid() {
		builder.WriteString("WARNING: MnUserParameterState is not valid.\n")
	}

	builder.WriteString("\n")
}

func (this *MnPrintStruct) toStringMinosError(x *MinosError) string {
	var builder strings.Builder
	this.PrintMinosError(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintMinosError(builder *strings.Builder, me *MinosError) {
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Minos # of function calls: %d\n", me.Nfcn()))
	if !me.IsValid() {
		builder.WriteString("Minos error is not valid.\n")
	}

	if !me.LowerValid() {
		builder.WriteString("lower Minos error is not valid.\n")
	}

	if !me.UpperValid() {
		builder.WriteString("upper Minos error is not valid.\n")
	}

	if me.AtLowerLimit() {
		builder.WriteString(fmt.Sprintf("Minos error is lower limit of parameter %d\n", me.Parameter()))
	}

	if me.AtUpperLimit() {
		builder.WriteString(fmt.Sprintf("Minos error is upper limit of parameter %d\n", me.Parameter()))
	}

	if me.AtLowerMaxFcn() {
		builder.WriteString("Minos number of function calls for lower error exhausted.\n")
	}

	if me.AtUpperMaxFcn() {
		builder.WriteString("Minos number of function calls for upper error exhausted.\n")
	}

	if me.LowerNewMin() {
		builder.WriteString("Minos found a new minimum in negative direction.\n")
		builder.WriteString(fmt.Sprintf("%g\n", me.LowerState()))
	}

	if me.UpperNewMin() {
		builder.WriteString("Minos found a new minimum in positive direction.\n")
		builder.WriteString(fmt.Sprintf("%g\n", me.UpperState()))
	}

	builder.WriteString("# ext. ||   name    || value@min ||  negative || positive  \n")
	builder.WriteString(fmt.Sprintf("%4d||%10s||%10g||%10g||%10g\n", me.Parameter(),
		me.LowerState().Name(me.Parameter()), me.Min(), me.Lower(), me.Upper()))
	builder.WriteString("\n")
}

func (this *MnPrintStruct) toStringContoursError(x *ContoursError) string {
	var builder strings.Builder
	this.PrintContoursError(&builder, x)
	return builder.String()
}

func (this *MnPrintStruct) PrintContoursError(builder *strings.Builder, ce *ContoursError) {
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Contours # of function calls: %d\n", ce.Nfcn()))
	builder.WriteString("MinosError in x: \n")
	builder.WriteString(fmt.Sprintf("%s\n", ce.XMinosError()))
	builder.WriteString("MinosError in y: \n")
	builder.WriteString(fmt.Sprintf("%s\n", ce.YMinosError()))
	plot := NewMnPlot()
	plot.PlotWithMin(ce.Xmin(), ce.Ymin(), ce.Points())
	i := 0

	for _, ipoint := range ce.Points() {
		builder.WriteString(fmt.Sprintf("%d %10g %10g\n", i+1, ipoint.first, ipoint.second))
	}

	builder.WriteString("\n")
}
