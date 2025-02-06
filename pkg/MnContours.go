package minuit

import (
	"errors"
	"log"
	"math"
)

type MnContours struct {
	theFCN      FCNBase
	theMinimum  *FunctionMinimum
	theStrategy *MnStrategy
}

func NewMnContours(fcn FCNBase, min *FunctionMinimum) *MnContours {
	return NewMnContoursWithStra(fcn, min, DEFAULT_STRATEGY)
}

/** construct from FCN + minimum + strategy */
func NewMnContoursWithStra(fcn FCNBase, min *FunctionMinimum, stra int) *MnContours {
	return NewMnContoursWithStrategy(fcn, min, NewMnStrategyWithStra(stra))
}

/** construct from FCN + minimum + strategy */
func NewMnContoursWithStrategy(fcn FCNBase, min *FunctionMinimum, stra *MnStrategy) *MnContours {
	return &MnContours{
		theFCN:      fcn,
		theMinimum:  min,
		theStrategy: stra,
	}
}

func (this *MnContours) Points(px, py int) ([]*Point, error) {
	return this.PointsWithError(px, py, 1)
}
func (this *MnContours) PointsWithError(px, py int, errDef float64) ([]*Point, error) {
	return this.PointsWithErrorN(px, py, errDef, 20)
}

// PointsWithErrorN
/*
 * Calculates one function contour of FCN with respect to parameters
 * parx and pary. The return value is a list of (x,y)
 * points. FCN minimized always with respect to all other n - 2 variable parameters
 * (if any). MINUIT will try to find n points on the contour (default 20). To
 * calculate more than one contour, the user needs to set the error definition in
 * its FCN to the appropriate value for the desired confidence level and call this method for each contour.
 */
func (this *MnContours) PointsWithErrorN(px, py int, errDef float64, npoints int) ([]*Point, error) {
	cont, err := this.ContourWithErrorN(px, py, errDef, npoints)
	if err != nil {
		return nil, err
	}
	return cont.Points(), nil
}

func (this *MnContours) Contour(px, py int) (*ContoursError, error) {
	return this.ContourWithError(px, py, 1)
}
func (this *MnContours) ContourWithError(px, py int, errDef float64) (*ContoursError, error) {
	return this.ContourWithErrorN(px, py, errDef, 20)
}

/**
 * Causes a CONTOURS error analysis and returns the result in form of ContoursError. As
 * a by-product ContoursError keeps the MinosError information of parameters parx and
 * pary. The result ContoursError can be easily printed using MnPrint or toString().
 */
func (this *MnContours) ContourWithErrorN(px, py int, errDef float64, npoints int) (*ContoursError, error) {
	errDef *= this.theMinimum.ErrorDef()
	if npoints <= 3 {
		return nil, errors.New("assertion violation: number of points must be greater than 3")
	}
	var maxcalls int = 100 * (npoints + 5) * (this.theMinimum.UserState().VariableParameters() + 1)
	var nfcn int = 0

	var result []*Point = make([]*Point, 0, npoints)
	var states []*MnUserParameterState = make([]*MnUserParameterState, 0) //TODO remove when not used
	var toler float64 = 0.05

	//get first four points
	var minos *MnMinos = NewMnMinosWithFunctionMinimumStrategy(this.theFCN, this.theMinimum, this.theStrategy)

	var valx float64 = this.theMinimum.UserState().Value(px)
	var valy float64 = this.theMinimum.UserState().Value(py)

	mex, _ := minos.minosWithErrDef(px, errDef)
	nfcn += mex.Nfcn()
	if !mex.IsValid() {
		log.Println("MnContours is unable to find first two points.")
		return NewContoursError(px, py, result, mex, mex, nfcn), nil
	}
	var ex *Point = mex.Range()

	mey, _ := minos.minosWithErrDef(py, errDef)
	nfcn += mey.Nfcn()
	if !mey.IsValid() {
		log.Println("MnContours is unable to find second two points.")
		return NewContoursError(px, py, result, mex, mey, nfcn), nil
	}
	var ey *Point = mey.Range()

	var migrad *MnMigrad = NewMnMigrad(this.theFCN, this.theMinimum.UserState().clone(), NewMnStrategyWithStra(max(0, this.theStrategy.Strategy()-1)))

	migrad.Fix(px)
	migrad.SetValue(px, valx+ex.second)
	exy_up, _ := migrad.Minimize()
	nfcn += exy_up.Nfcn()
	if !exy_up.IsValid() {
		log.Printf("MnContours is unable to find upper y value for x parameter %d.\n", px)
		return NewContoursError(px, py, result, mex, mey, nfcn), nil
	}

	migrad.SetValue(px, valx+ex.first)
	exy_lo, _ := migrad.Minimize()
	nfcn += exy_lo.Nfcn()
	if !exy_lo.IsValid() {
		log.Printf("MnContours is unable to find lower y value for x parameter %d.\n", px)
		return NewContoursError(px, py, result, mex, mey, nfcn), nil
	}

	var migrad1 *MnMigrad = NewMnMigrad(this.theFCN, this.theMinimum.UserState().clone(), NewMnStrategyWithStra(max(0, this.theStrategy.Strategy()-1)))
	migrad1.Fix(py)
	migrad1.SetValue(py, valy+ey.second)
	eyx_up, _ := migrad1.Minimize()
	nfcn += eyx_up.Nfcn()
	if !eyx_up.IsValid() {
		log.Printf("MnContours is unable to find upper x value for y parameter %d.\n", py)
		return NewContoursError(px, py, result, mex, mey, nfcn), nil
	}

	migrad1.SetValue(py, valy+ey.first)
	eyx_lo, _ := migrad1.Minimize()
	nfcn += eyx_lo.Nfcn()
	if !eyx_lo.IsValid() {
		log.Printf("MnContours is unable to find lower x value for y parameter %d.\n", py)
		return NewContoursError(px, py, result, mex, mey, nfcn), nil
	}

	var scalx float64 = 1. / (ex.second - ex.first)
	var scaly float64 = 1. / (ey.second - ey.first)

	result = append(result, NewPoint(valx+ex.first, exy_lo.UserState().Value(py)))
	result = append(result, NewPoint(eyx_lo.UserState().Value(px), valy+ey.first))
	result = append(result, NewPoint(valx+ex.second, exy_up.UserState().Value(py)))
	result = append(result, NewPoint(eyx_up.UserState().Value(px), valy+ey.second))

	var upar *MnUserParameterState = this.theMinimum.UserState().clone()
	upar.Fix(px)
	upar.Fix(py)

	var par []int = []int{px, py}
	var cross *MnFunctionCross = NewMnFunctionCross(this.theFCN, upar, this.theMinimum.Fval(), this.theStrategy, errDef)

	for i := 4; i < npoints; i++ {
		var idist1 *Point = result[len(result)-1]
		var idist2 *Point = result[0]
		var pos2 int = 0
		var distx float64 = idist1.first - idist2.first
		var disty float64 = idist1.second - idist2.second
		var bigdis float64 = scalx*scalx*distx*distx + scaly*scaly*disty*disty

		for j := 0; j < len(result)-1; j++ {
			var ipair *Point = result[j]
			var distx2 float64 = ipair.first - result[j+1].first
			var disty2 float64 = ipair.second - result[j+1].second
			var dist float64 = scalx*scalx*distx2*distx2 + scaly*scaly*disty2*disty2
			if dist > bigdis {
				bigdis = dist
				idist1 = ipair
				idist2 = result[j+1]
				pos2 = j + 1
			}
		}

		var a1 float64 = 0.5
		var a2 float64 = 0.5
		var sca float64 = 1.0

		for {
			if nfcn > maxcalls {
				log.Println("MnContours: maximum number of function calls exhausted.")
				return NewContoursError(px, py, result, mex, mey, nfcn), nil
			}

			var xmidcr float64 = a1*idist1.first + a2*idist2.first
			var ymidcr float64 = a1*idist1.second + a2*idist2.second
			var xdir float64 = idist2.second - idist1.second
			var ydir float64 = idist1.first - idist2.first
			var scalfac float64 = sca * math.Max(math.Abs(xdir*scalx), math.Abs(ydir*scaly))
			var xdircr float64 = xdir / scalfac
			var ydircr float64 = ydir / scalfac
			var pmid []float64 = []float64{xmidcr, ymidcr}
			var pdir []float64 = []float64{xdircr, ydircr}

			var opt *MnCross = cross.cross(par, pmid, pdir, toler, maxcalls)
			nfcn += opt.nfcn()
			if opt.isValid() {
				var aopt float64 = opt.value()
				if pos2 == 0 {
					result = append(result, NewPoint(xmidcr+(aopt)*xdircr, ymidcr+(aopt)*ydircr))
				} else {
					result = append(result, NewPoint(xmidcr+(aopt)*xdircr, ymidcr+(aopt)*ydircr))
				}
				break
			}
			if sca < 0.0 {
				log.Printf("MnContours is unable to find point %d on contour.\n", i+1)
				log.Printf("MnContours finds only %d points.\n", i)
				return NewContoursError(px, py, result, mex, mey, nfcn), nil
			}
			sca = -1.0
		}
	}
	return NewContoursError(px, py, result, mex, mey, nfcn), nil
}

func (this *MnContours) strategy() *MnStrategy {
	return this.theStrategy
}
