package example

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"testing"

	minuit "github.com/empack/minuit2go/pkg"
)

type ReneFcn struct {
	theMeasurements []float64
}

func NewReneFcn(theMeasurements []float64) *ReneFcn {
	return &ReneFcn{
		theMeasurements,
	}
}

func (this *ReneFcn) errorDef() float64 {
	return 1.0
}

func (this *ReneFcn) ValueOf(par []float64) float64 {
	var a float64 = par[2]
	var b float64 = par[1]
	var c float64 = par[0]
	var p0 float64 = par[3]
	var p1 float64 = par[4]
	var p2 float64 = par[5]
	var fval float64 = 0.0
	for i := 0; i < len(this.theMeasurements); i++ {
		var ni float64 = this.theMeasurements[i]
		if ni < 1.e-10 {
			continue
		}
		var xi float64 = (float64(i)+1.0)/40.0 - 1.0/80.0 //xi=0-3
		var ei float64 = ni
		var nexp float64 = a*xi*xi + b*xi + c + (0.5*p0*p1/math.Pi)/math.Max(1.e-10, (xi-p2)*(xi-p2)+0.25*p1*p1)
		fval += (ni - nexp) * (ni - nexp) / ei
	}
	return fval
}

func TestRene(t *testing.T) {
	var tmp []float64 = []float64{38., 36., 46., 52., 54., 52., 61., 52., 64., 77.,
		60., 56., 78., 71., 81., 83., 89., 96., 118., 96.,
		109., 111., 107., 107., 135., 156., 196., 137.,
		160., 153., 185., 222., 251., 270., 329., 422.,
		543., 832., 1390., 2835., 3462., 2030., 1130.,
		657., 469., 411., 375., 295., 281., 281., 289.,
		273., 297., 256., 274., 287., 280., 274., 286.,
		279., 293., 314., 285., 322., 307., 313., 324.,
		351., 314., 314., 301., 361., 332., 342., 338.,
		396., 356., 344., 395., 416., 406., 411., 422.,
		393., 393., 409., 455., 427., 448., 459., 403.,
		441., 510., 501., 502., 482., 487., 506., 506.,
		526., 517., 534., 509., 482., 591., 569., 518.,
		609., 569., 598., 627., 617., 610., 662., 666.,
		652., 671., 647., 650., 701.}

	var measurements []float64 = slices.Clone(tmp)

	var theFCN *ReneFcn = NewReneFcn(measurements)

	var upar *minuit.MnUserParameters = minuit.NewEmptyMnUserParameters()
	upar.AddFree("p0", 100., 10.)
	upar.AddFree("p1", 100., 10.)
	upar.AddFree("p2", 100., 10.)
	upar.AddFree("p3", 100., 10.)
	upar.AddFree("p4", 1., 0.3)
	upar.AddFree("p5", 1., 0.3)

	fmt.Printf("Initial parameters: %s\n", upar.String())

	println("start migrad")
	var migrad *minuit.MnMigrad = minuit.NewMnMigradWithParameters(theFCN, upar)
	min, err := migrad.Minimize()
	if err != nil {
		t.Fatalf("minimize failed with:\n %s\n", err.Error())
	}
	if !min.IsValid() {
		//try with higher strategy
		println("FM is invalid, try with strategy = 2.")
		var migrad2 *minuit.MnMigrad = minuit.NewMnMigradWithParameterStateStrategy(theFCN, min.UserState(), minuit.NewMnStrategyWithStra(2))
		min, err = migrad2.Minimize()
		if err != nil {
			t.Fatalf("minimize failed with:\n %s\n", err.Error())
		}
	}
	var sbuff strings.Builder
	minuit.MnPrint.PrintFunctionMinimum(&sbuff, min)
	fmt.Printf("minimum: %s\n", sbuff.String())

	{
		var params []float64 = []float64{1, 1, 1, 1, 1, 1}
		var error []float64 = []float64{1, 1, 1, 1, 1, 1}
		var scan *minuit.MnScan = minuit.NewMnScan(theFCN, params, error)
		fmt.Printf("scan parameters: %s\n", scan.Parameters())
		var plot *minuit.MnPlot = minuit.NewMnPlot()
		for i := 0; i < upar.VariableParameters(); i++ {
			xy, fnErr := scan.Scan(i)
			if fnErr != nil {
				t.Fatalf("scan failed with:\n %s\n", fnErr.Error())
			}
			plot.Plot(xy)
		}
		fmt.Printf("scan parameters: %s\n", scan.Parameters())
	}

	{
		var params []float64 = []float64{1, 1, 1, 1, 1, 1}
		var error []float64 = []float64{1, 1, 1, 1, 1, 1}
		var scan *minuit.MnScan = minuit.NewMnScan(theFCN, params, error)
		fmt.Printf("scan parameters: %s\n", scan.Parameters())
		min2, err := scan.Minimize()
		if err != nil {
			t.Fatalf("minimize failed with:\n %s\n", err.Error())
		}
		//     std::cout<<min<<std::endl;
		var sbuff2 strings.Builder
		minuit.MnPrint.PrintFunctionMinimum(&sbuff2, min2)
		fmt.Printf("%s\n", sbuff2.String())
		fmt.Printf("scan parameters: %s\n", scan.Parameters())
	}
}
