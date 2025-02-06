package minuit

import (
	"bytes"
	"fmt"
	"math"
)

type MnPlot struct {
	thePageWidth  int
	thePageLength int
	bl            float64
	bh            float64
	nb            int
	bwid          float64
}

func NewMnPlot() *MnPlot {
	return NewMnPlotWithWidthLength(80, 30)
}

func NewMnPlotWithWidthLength(width, length int) *MnPlot {
	if width > 120 {
		width = 120
	}
	if length > 56 {
		length = 56
	}

	return &MnPlot{
		thePageWidth:  width,
		thePageLength: length,
		bl:            0,
		bh:            0,
		nb:            0,
		bwid:          0,
	}
}

func (this *MnPlot) Plot(points []*Point) {
	x := make([]float64, len(points))
	y := make([]float64, len(points))
	var chpt bytes.Buffer

	for i, point := range points {
		x[i] = point.first
		y[i] = point.second
		chpt.WriteByte('*')
	}

	this.mnplot(x, y, chpt, len(points), this.width(), this.length())
}

func (this *MnPlot) PlotWithMin(xmin, ymin float64, points []*Point) {
	x := make([]float64, len(points)+2)
	x[0] = xmin
	x[1] = xmin
	y := make([]float64, len(points)+2)
	y[0] = ymin
	y[1] = ymin
	var chpt bytes.Buffer
	chpt.WriteByte(' ')
	chpt.WriteByte('X')

	i := 2
	for _, point := range points {
		x[i] = point.first
		y[i] = point.second
		chpt.WriteByte('*')
		i++
	}

	this.mnplot(x, y, chpt, len(points)+2, this.width(), this.length())
}

func (this *MnPlot) width() int {
	return this.thePageWidth
}

func (this *MnPlot) length() int {
	return this.thePageLength
}

func (this *MnPlot) mnplot(xpt, ypt []float64, chpt bytes.Buffer, nxypt, npagwd, npagln int) {
	xvalus := make([]float64, 12)
	var cline bytes.Buffer

	for ii := 0; ii < npagwd; ii++ {
		cline.WriteByte(' ')
	}

	var maxnx int
	if npagwd-20 < 100 {
		maxnx = npagwd - 20
	} else {
		maxnx = 100
	}
	if maxnx < 10 {
		maxnx = 10
	}

	maxny := npagln
	if npagln < 10 {
		maxny = 10
	}

	if nxypt > 1 {
		xbest := xpt[0]
		ybest := ypt[0]
		chbest := chpt.Bytes()[0]
		km1 := nxypt - 1

		for i := 1; i <= km1; i++ {
			iquit := 0
			ni := nxypt - i

			for j := 1; j <= ni; j++ {
				if !(ypt[j-1] > ypt[j]) {
					savx := xpt[j-1]
					xpt[j-1] = xpt[j]
					xpt[j] = savx
					savy := ypt[j-1]
					ypt[j-1] = ypt[j]
					ypt[j] = savy
					chsav := chpt.Bytes()[j-1]
					chpt.Bytes()[j-1] = chpt.Bytes()[j]
					chpt.Bytes()[j] = chsav
					iquit = 1
				}
			}

			if iquit == 0 {
				break
			}
		}

		xmax := xpt[0]
		xmin := xmax

		for var75 := 1; var75 <= nxypt; var75++ {
			if xpt[var75-1] > xmax {
				xmax = xpt[var75-1]
			}

			if xpt[var75-1] < xmin {
				xmin = xpt[var75-1]
			}
		}

		dxx := (xmax - xmin) * 0.001
		xmax += dxx
		xmin -= dxx
		this.mnbins(xmin, xmax, maxnx)
		xmin = this.bl
		xmax = this.bh
		nx := this.nb
		bwidx := this.bwid
		ymax := ypt[0]
		ymin := ypt[nxypt-1]
		if ymax == ymin {
			ymax = ymin + 1.0
		}

		dyy := (ymax - ymin) * 0.001
		ymax += dyy
		ymin -= dyy
		this.mnbins(ymin, ymax, maxny)
		ymin = this.bl
		ymax = this.bh
		ny := this.nb
		bwidy := this.bwid
		anyValue := float64(ny)
		if chbest != ' ' {
			xbest = (xmax + xmin) * 0.5
			ybest = (ymax + ymin) * 0.5
		}

		ax := 1.0 / bwidx
		ay := 1.0 / bwidy
		bx := -ax*xmin + 2.0
		by := -ay*ymin - 2.0

		for var76 := 1; var76 <= nxypt; var76++ {
			xpt[var76-1] = ax*xpt[var76-1] + bx
			ypt[var76-1] = anyValue - ay*ypt[var76-1] - by
		}

		nxbest := int(ax*xbest + bx)
		nybest := int(anyValue - ay*ybest - by)
		ny += 2
		nx += 2
		isp1 := 1
		linodd := 1
		overpr := false

		for var77 := 1; var77 <= ny; var77++ {
			for ibk := 1; ibk <= nx; ibk++ {
				cline.Bytes()[ibk-1] = ' '
			}

			cline.Bytes()[0] = '.'
			cline.Bytes()[nx-1] = '.'
			cline.Bytes()[nxbest-1] = '.'
			if var77 == 1 || var77 == nybest || var77 == ny {
				for j := 1; j <= nx; j++ {
					cline.Bytes()[j-1] = '.'
				}
			}

			yprt := ymax - (float64(var77)-1.0)*bwidy
			isplset := false
			if isp1 <= nxypt {
				for k := isp1; k <= nxypt; k++ {
					ks := int(ypt[k-1])
					if ks > var77 {
						isp1 = k
						isplset = true
						break
					}

					ix := int(xpt[k-1])
					if cline.Bytes()[ix-1] != '.' && cline.Bytes()[ix-1] != ' ' {
						if cline.Bytes()[ix-1] != chpt.Bytes()[k-1] {
							overpr = true
							cline.Bytes()[ix-1] = '&'
						}
					} else {
						cline.Bytes()[ix-1] = chpt.Bytes()[k-1]
					}
				}

				if !isplset {
					isp1 = nxypt + 1
				}
			}

			if linodd != 1 && var77 != ny {
				fmt.Printf("                  %s", string(cline.Bytes()[:60]))
			} else {
				fmt.Printf(" %14.7g ..%s", yprt, string(cline.Bytes()[:60]))
				linodd = 0
			}

			fmt.Println()
		}

		for ibk := 1; ibk <= nx; ibk++ {
			cline.Bytes()[ibk-1] = ' '
			if ibk%10 == 1 {
				cline.Bytes()[ibk-1] = '/'
			}
		}

		fmt.Printf("                  %s", string(cline.Bytes()))
		fmt.Printf("\n")

		for var80 := 1; var80 <= 12; var80++ {
			xvalus[var80-1] = xmin + (float64(var80)-1.0)*10.0*bwidx
		}

		fmt.Print("           ")
		iten := (nx + 9) / 10

		for var81 := 1; var81 <= iten; var81++ {
			fmt.Printf(" %9.4g", xvalus[var81-1])
		}

		fmt.Printf("\n")
		if overpr {
			chmess := "   Overprint character is &"
			fmt.Printf("                         ONE COLUMN=%13.7g%s", bwidx, chmess)
		} else {
			chmess := " "
			fmt.Printf("                         ONE COLUMN=%13.7g%s", bwidx, chmess)
		}

		fmt.Println()
	}
}

func (this *MnPlot) mnbins(a1, a2 float64, naa int) {
	na := 0
	var al, ah float64
	if a1 < a2 {
		al = a1
		ah = a2
	} else {
		al = a2
		ah = a1
	}
	if al == ah {
		ah = al + 1.0
	}

	skip := naa == -1 && this.bwid > 0.0
	if !skip {
		na = naa - 1
		if na < 1 {
			na = 1
		}
	}

	for {
		if !skip {
			awid := (ah - al) / float64(na)
			log_ := int(math.Log10(awid))
			if awid <= 1.0 {
				log_--
			}

			sigfig := awid * math.Pow(10.0, float64(-log_))
			var sigrnd float64
			if sigfig <= 2.0 {
				sigrnd = 2.0
			} else if sigfig <= 2.5 {
				sigrnd = 2.0
			} else if sigfig <= 5.0 {
				sigrnd = 5.0
			} else {
				sigrnd = 1.0
				log_++
			}

			this.bwid = sigrnd * math.Pow(10.0, float64(log_))
		}

		alb := al / this.bwid
		lwid := int(alb)
		if alb < 0.0 {
			lwid--
		}

		this.bl = this.bwid * float64(lwid)
		alb = ah/this.bwid + 1.0
		kwid := int(alb)
		if alb < 0.0 {
			kwid--
		}

		this.bh = this.bwid * float64(kwid)
		this.nb = kwid - lwid
		if naa <= 5 {
			if naa == -1 {
				return
			}

			if naa <= 1 && this.nb != 1 {
				this.bwid += 2.0
				this.nb = 1
				return
			}

			return
		}

		if this.nb<<1 != naa {
			return
		}

		na++
		skip = false
	}

}
