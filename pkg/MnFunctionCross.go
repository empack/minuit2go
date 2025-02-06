package minuit

import (
	"log"
	"math"
)

type MnFunctionCross struct {
	theFCN      FCNBase
	theState    *MnUserParameterState
	theFval     float64
	theStrategy *MnStrategy
	theErrorDef float64
}

func NewMnFunctionCross(fcn FCNBase, state *MnUserParameterState, fval float64, stra *MnStrategy, errorDef float64) *MnFunctionCross {
	return &MnFunctionCross{
		theFCN:      fcn,
		theState:    state,
		theFval:     fval,
		theStrategy: stra,
		theErrorDef: errorDef,
	}
}

func (this *MnFunctionCross) cross(par []int, pmid []float64, pdir []float64, tlr float64, maxcalls int) *MnCross {
	var npar int = len(par)
	var nfcn int = 0
	var prec *MnMachinePrecision = this.theState.Precision()

	var tlf float64 = tlr * this.theErrorDef
	var tla float64 = tlr
	var maxitr int = 15
	var ipt int = 0
	var aminsv float64 = this.theFval
	var aim float64 = aminsv + this.theErrorDef

	var aopt float64 = 0.0
	var limset bool = false
	var alsb []float64 = make([]float64, 3)
	var flsb []float64 = make([]float64, 3)
	var up float64 = this.theErrorDef

	var aulim float64 = 100.0
	for i := 0; i < len(par); i++ {
		var kex int = par[i]
		if this.theState.parameter(kex).HasLimits() {
			var zmid float64 = pmid[i]
			var zdir float64 = pdir[i]
			if math.Abs(zdir) < this.theState.Precision().eps() {
				continue
			}

			if zdir > 0. && this.theState.parameter(kex).HasUpperLimit() {
				var zlim float64 = this.theState.parameter(kex).UpperLimit()
				aulim = math.Min(aulim, (zlim-zmid)/zdir)
			} else if zdir < 0. && this.theState.parameter(kex).HasLowerLimit() {
				var zlim float64 = this.theState.parameter(kex).LowerLimit()
				aulim = math.Min(aulim, (zlim-zmid)/zdir)
			}
		}
	}

	if aulim < aopt+tla {
		limset = true
	}

	var migrad *MnMigrad = NewMnMigrad(this.theFCN, this.theState, NewMnStrategyWithStra(max(0, this.theStrategy.Strategy()-1)))

	for i := 0; i < npar; i++ {
		migrad.SetValue(par[i], pmid[i])
	}

	min0, _ := migrad.MinimizeWithMaxfcnToler(maxcalls, tlr)
	nfcn += min0.Nfcn()

	if min0.hasReachedCallLimit() {
		return NewMnCrossFcnLimit(min0.UserState(), nfcn)
	}

	if !min0.IsValid() {
		return NewMnCrossWithNfcn(nfcn)
	}
	if limset == true && min0.Fval() < aim {
		return NewMnCrossParLimit(min0.UserState(), nfcn)
	}
	ipt++
	alsb[0] = 0.
	flsb[0] = min0.Fval()
	flsb[0] = math.Max(flsb[0], aminsv+0.1*up)
	aopt = math.Sqrt(up/(flsb[0]-aminsv)) - 1.
	if math.Abs(flsb[0]-aim) < tlf {
		return NewMnCrossWithValueStateNfcn(aopt, min0.UserState(), nfcn)
	}

	if aopt > 1.0 {
		aopt = 1.0
	}
	if aopt < -0.5 {
		aopt = -0.5
	}
	limset = false
	if aopt > aulim {
		aopt = aulim
		limset = true
	}

	for i := 0; i < npar; i++ {
		migrad.SetValue(par[i], pmid[i]+(aopt)*pdir[i])
	}

	min1, _ := migrad.MinimizeWithMaxfcnToler(maxcalls, tlr)
	nfcn += min1.Nfcn()

	if min1.hasReachedCallLimit() {
		return NewMnCrossFcnLimit(min1.UserState(), nfcn)
	}
	if !min1.IsValid() {
		return NewMnCrossWithNfcn(nfcn)
	}
	if limset == true && min1.Fval() < aim {
		return NewMnCrossParLimit(min1.UserState(), nfcn)
	}

	ipt++
	alsb[1] = aopt
	flsb[1] = min1.Fval()
	var dfda float64 = (flsb[1] - flsb[0]) / (alsb[1] - alsb[0])

	var ecarmn float64
	var ecarmx float64
	var ibest int = 0
	var iworst int
	var noless int
	var min2 *FunctionMinimum

L300:
	for {
		if dfda < 0. {
			var maxlk int = maxitr - ipt
			for it := 0; it < maxlk; it++ {
				alsb[0] = alsb[1]
				flsb[0] = flsb[1]
				aopt = alsb[0] + 0.2*float64(it)
				limset = false
				if aopt > aulim {
					aopt = aulim
					limset = true
				}
				for i := 0; i < npar; i++ {
					migrad.SetValue(par[i], pmid[i]+(aopt)*pdir[i])
				}
				min1, _ = migrad.MinimizeWithMaxfcnToler(maxcalls, tlr)
				nfcn += min1.Nfcn()

				if min1.hasReachedCallLimit() {
					return NewMnCrossFcnLimit(min1.UserState(), nfcn)
				}
				if !min1.IsValid() {
					return NewMnCrossWithNfcn(nfcn)
				}
				if limset == true && min1.Fval() < aim {
					return NewMnCrossParLimit(min1.UserState(), nfcn)
				}

				ipt++
				alsb[1] = aopt
				flsb[1] = min1.Fval()
				dfda = (flsb[1] - flsb[0]) / (alsb[1] - alsb[0])
				if dfda > 0.0 {
					break
				}
			}
			if ipt > maxitr {
				return NewMnCrossWithNfcn(nfcn)
			}

		}

	L460:
		for {
			aopt = alsb[1] + (aim-flsb[1])/dfda
			var fdist float64 = math.Min(math.Abs(aim-flsb[0]), math.Abs(aim-flsb[1]))
			var adist float64 = math.Min(math.Abs(aopt-alsb[0]), math.Abs(aopt-alsb[1]))
			tla = tlr
			if math.Abs(aopt) > 1. {
				tla = tlr * math.Abs(aopt)
			}
			if adist < tla && fdist < tlf {
				return NewMnCrossWithValueStateNfcn(aopt, min1.UserState(), nfcn)
			}
			if ipt > maxitr {
				return NewMnCrossWithNfcn(nfcn)
			}
			var bmin float64 = math.Min(alsb[0], alsb[1]) - 1.0
			if aopt < bmin {
				aopt = bmin
			}
			var bmax float64 = math.Max(alsb[0], alsb[1]) + 1.0
			if aopt > bmax {
				aopt = bmax
			}

			limset = false
			if aopt > aulim {
				aopt = aulim
				limset = true
			}

			for i := 0; i < npar; i++ {
				migrad.SetValue(par[i], pmid[i]+(aopt)*pdir[i])
			}
			min2, _ = migrad.MinimizeWithMaxfcnToler(maxcalls, tlr)
			nfcn += min2.Nfcn()

			if min2.hasReachedCallLimit() {
				return NewMnCrossFcnLimit(min2.UserState(), nfcn)
			}

			if !min2.IsValid() {
				return NewMnCrossWithNfcn(nfcn)
			}
			if limset == true && min2.Fval() < aim {
				return NewMnCrossParLimit(min2.UserState(), nfcn)
			}

			ipt++
			alsb[2] = aopt
			flsb[2] = min2.Fval()

			ecarmn = math.Abs(flsb[2] - aim)
			ecarmx = 0.0
			ibest = 2
			iworst = 0
			noless = 0

			for i := 0; i < 3; i++ {
				var ecart float64 = math.Abs(flsb[i] - aim)
				if ecart > ecarmx {
					ecarmx = ecart
					iworst = i
				}
				if ecart < ecarmn {
					ecarmn = ecart
					ibest = i
				}
				if flsb[i] < aim {
					noless++
				}
			}

			if noless == 1 || noless == 2 {
				break L300
			}
			if noless == 0 && ibest != 2 {
				return NewMnCrossWithNfcn(nfcn)
			}
			if noless == 3 && ibest != 2 {
				alsb[1] = alsb[2]
				flsb[1] = flsb[2]
				continue L300
			}

			flsb[iworst] = flsb[2]
			alsb[iworst] = alsb[2]
			dfda = (flsb[1] - flsb[0]) / (alsb[1] - alsb[0])
		}
	}

	for ok := true; ok; {
		var parbol *MnParabola = MnParabolaFactory.createWith3Points(NewMnParabolaPoint(alsb[0], flsb[0]), NewMnParabolaPoint(alsb[1], flsb[1]), NewMnParabolaPoint(alsb[2], flsb[2]))

		var coeff1 float64 = parbol.c()
		var coeff2 float64 = parbol.b()
		var coeff3 float64 = parbol.a()
		var determ float64 = coeff2*coeff2 - 4.*coeff3*(coeff1-aim)
		if determ < prec.eps() {
			return NewMnCrossWithNfcn(nfcn)
		}
		var rt float64 = math.Sqrt(determ)
		var x1 float64 = (-coeff2 + rt) / (2. * coeff3)
		var x2 float64 = (-coeff2 - rt) / (2. * coeff3)
		var s1 float64 = coeff2 + 2.*x1*coeff3
		var s2 float64 = coeff2 + 2.*x2*coeff3

		if s1*s2 > 0.0 {
			log.Println("MnFunctionCross problem 1")
		}
		aopt = x1
		var slope float64 = s1
		if s2 > 0.0 {
			aopt = x2
			slope = s2
		}

		tla = tlr
		if math.Abs(aopt) > 1.0 {
			tla = tlr * math.Abs(aopt)
		}
		if math.Abs(aopt-alsb[ibest]) < tla && math.Abs(flsb[ibest]-aim) < tlf {
			return NewMnCrossWithValueStateNfcn(aopt, min2.UserState(), nfcn)
		}

		var ileft int = 3
		var iright int = 3
		var iout int = 3
		ibest = 0
		ecarmx = 0.0
		ecarmn = math.Abs(aim - flsb[0])
		for i := 0; i < 3; i++ {
			var ecart float64 = math.Abs(flsb[i] - aim)
			if ecart < ecarmn {
				ecarmn = ecart
				ibest = i
			}
			if ecart > ecarmx {
				ecarmx = ecart
			}
			if flsb[i] > aim {
				if iright == 3 {
					iright = i
				} else if flsb[i] > flsb[iright] {
					iout = i
				} else {
					iout = iright
					iright = i
				}
			} else if ileft == 3 {
				ileft = i
			} else if flsb[i] < flsb[ileft] {
				iout = i
			} else {
				iout = ileft
				ileft = i
			}
		}

		if ecarmx > 10.*math.Abs(flsb[iout]-aim) {
			aopt = 0.5 * (aopt + 0.5*(alsb[iright]+alsb[ileft]))
		}
		var smalla float64 = 0.1 * tla
		if slope*smalla > tlf {
			smalla = tlf / slope
		}
		var aleft float64 = alsb[ileft] + smalla
		var aright float64 = alsb[iright] - smalla
		if aopt < aleft {
			aopt = aleft
		}
		if aopt > aright {
			aopt = aright
		}
		if aleft > aright {
			aopt = 0.5 * (aleft + aright)
		}

		limset = false
		if aopt > aulim {
			aopt = aulim
			limset = true
		}

		for i := 0; i < npar; i++ {
			migrad.SetValue(par[i], pmid[i]+(aopt)*pdir[i])
		}
		min2 = migrad.minimize(maxcalls, tlr)
		nfcn += min2.Nfcn()

		if min2.hasReachedCallLimit() {
			return NewMnCrossFcnLimit(min2.UserState(), nfcn)
		}
		if !min2.IsValid() {
			return NewMnCrossWithNfcn(nfcn)
		}
		if limset == true && min2.Fval() < aim {
			return NewMnCrossParLimit(min2.UserState(), nfcn)
		}

		ipt++
		alsb[iout] = aopt
		flsb[iout] = min2.Fval()
		ibest = iout
		ok = ipt < maxitr
	}

	return NewMnCrossWithNfcn(nfcn)
}
