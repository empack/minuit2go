package minuit

import "math"

var MnLineSearch = &mnLineSearchStruct{}

type mnLineSearchStruct struct {
}

func (this *mnLineSearchStruct) search(fcn *MnFcn, st *MinimumParameters, step *MnAlgebraicVector, gdel float64, prec *MnMachinePrecision) (*MnParabolaPoint, error) {
	var overal float64 = 1000.0
	var undral float64 = -100.0
	var toler float64 = 0.050
	var slamin float64 = 0.0
	var slambg float64 = 5.0
	var alpha float64 = 2.0
	var maxiter int = 12
	var niter int = 0

	for i := 0; i < step.size(); i++ {
		if math.Abs(step.get(i)) < prec.eps() {
			continue
		}
		var ratio float64 = math.Abs(st.vec().get(i) / step.get(i))
		if math.Abs(slamin) < prec.eps() {
			slamin = ratio
		}
		if ratio < slamin {
			slamin = ratio
		}
	}
	if math.Abs(slamin) < prec.eps() {
		slamin = prec.eps()
	}
	slamin *= prec.eps2()

	var F0 float64 = st.fval()
	v_, err := MnUtils.AddV(st.vec(), step)
	if err != nil {
		return nil, err
	}
	var F1 float64 = fcn.valueOf(v_)
	var fvmin float64 = st.fval()
	var xvmin float64 = 0.0

	if F1 < F0 {
		fvmin = F1
		xvmin = 1.
	}
	var toler8 float64 = toler
	var slamax float64 = slambg
	var flast float64 = F1
	var slam float64 = 1.0

	var iterate bool = false
	var p0 *MnParabolaPoint = NewMnParabolaPoint(0.0, F0)
	var p1 *MnParabolaPoint = NewMnParabolaPoint(slam, flast)
	var F2 float64 = 0.0

	for ok := true; ok; {
		// cut toler8 as function goes up
		iterate = false
		var denom float64 = 2. * (flast - F0 - gdel*slam) / (slam * slam)
		if math.Abs(denom) < prec.eps() {
			denom = -0.1 * gdel
			slam = 1.0
		}
		if math.Abs(denom) > prec.eps() {
			slam = -gdel / denom
		}
		if slam < 0.0 {
			slam = slamax
		}

		if slam > slamax {
			slam = slamax
		}
		if slam < toler8 {
			slam = toler8
		}
		if slam < slamin {
			return NewMnParabolaPoint(xvmin, fvmin), nil
		}
		if math.Abs(slam-1.) < toler8 && p1.y() < p0.y() {
			return NewMnParabolaPoint(xvmin, fvmin), nil
		}
		if math.Abs(slam-1.) < toler8 {
			slam = 1. + toler8
		}

		v_, err = MnUtils.AddV(st.vec(), MnUtils.MulV(step, slam))
		if err != nil {
			return nil, err
		}
		F2 = fcn.valueOf(v_)
		if F2 < fvmin {
			fvmin = F2
			xvmin = slam
		}
		if p0.y()-prec.eps() < fvmin && fvmin < p0.y()+prec.eps() {
			iterate = true
			flast = F2
			toler8 = toler * slam
			overal = slam - toler8
			slamax = overal
			p1 = NewMnParabolaPoint(slam, flast)
			niter++
		}
		ok = iterate && niter < maxiter
	}

	if niter >= maxiter {
		// exhausted max number of iterations
		return NewMnParabolaPoint(xvmin, fvmin), nil
	}

	var p2 *MnParabolaPoint = NewMnParabolaPoint(slam, F2)

	for ok := true; ok; ok = ok {
		slamax = math.Max(slamax, alpha*math.Abs(xvmin))
		var pb *MnParabola = MnParabolaFactory.createWith3Points(p0, p1, p2)
		if pb.a() < prec.eps2() {
			var slopem float64 = 2.*pb.a()*xvmin + pb.b()
			if slopem < 0. {
				slam = xvmin + slamax
			} else {
				slam = xvmin - slamax
			}
		} else {
			slam = pb.min()
			if slam > xvmin+slamax {
				slam = xvmin + slamax
			}
			if slam < xvmin-slamax {
				slam = xvmin - slamax
			}
		}
		if slam > 0.0 {
			if slam > overal {
				slam = overal
			}
		} else {
			if slam < undral {
				slam = undral
			}
		}

		var F3 float64 = 0.0

		for ok2 := true; ok2; {
			iterate = false
			var toler9 float64 = math.Max(toler8, math.Abs(toler8*slam))
			// min. of parabola at one point
			if math.Abs(p0.x()-slam) < toler9 ||
				math.Abs(p1.x()-slam) < toler9 ||
				math.Abs(p2.x()-slam) < toler9 {
				return NewMnParabolaPoint(xvmin, fvmin), nil
			}
			v_, err = MnUtils.AddV(st.vec(), MnUtils.MulV(step, slam))
			if err != nil {
				return nil, err
			}
			F3 = fcn.valueOf(v_)
			// if latest point worse than all three previous, cut step
			if F3 > p0.y() && F3 > p1.y() && F3 > p2.y() {
				if slam > xvmin {
					overal = math.Min(overal, slam-toler8)
				}
				if slam < xvmin {
					undral = math.Max(undral, slam+toler8)
				}
				slam = 0.5 * (slam + xvmin)
				iterate = true
				niter++
			}
			ok2 = iterate && niter < maxiter
		}

		if niter >= maxiter {
			// exhausted max number of iterations
			return NewMnParabolaPoint(xvmin, fvmin), nil
		}

		// find worst previous point out of three and replace
		var p3 *MnParabolaPoint = NewMnParabolaPoint(slam, F3)
		if p0.y() > p1.y() && p0.y() > p2.y() {
			p0 = p3
		} else if p1.y() > p0.y() && p1.y() > p2.y() {
			p1 = p3
		} else {
			p2 = p3
		}
		if F3 < fvmin {
			fvmin = F3
			xvmin = slam
		} else {
			if slam > xvmin {
				overal = math.Min(overal, slam-toler8)
			}
			if slam < xvmin {
				undral = math.Max(undral, slam+toler8)
			}
		}

		niter++
		ok = niter < maxiter
	}
	return NewMnParabolaPoint(xvmin, fvmin), nil
}
