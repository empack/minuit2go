package example

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	minuit "github.com/empack/minuit2go/pkg"
)

// ExperimentalFcn implementiert eine Fit-Funktion mit experimentellen Daten
type ExperimentalFcn struct {
	xData []float64 // x-Werte der Messpunkte
	yData []float64 // y-Werte der Messpunkte
	yerr  []float64 // Fehler der y-Werte
}

// Erzeugt simulierte Testdaten
func generateTestData(nPoints int) ([]float64, []float64, []float64) {
	// Echte Parameter, die wir später wiederfinden wollen
	realA := 10.0    // Amplitude
	realMu := 5.0    // Mittelwert
	realSigma := 2.0 // Breite
	realBg := 2.0    // Hintergrund
	realSlope := 0.5 // Linearer Term

	xData := make([]float64, nPoints)
	yData := make([]float64, nPoints)
	yerr := make([]float64, nPoints)

	// Generiere Datenpunkte mit Rauschen
	for i := 0; i < nPoints; i++ {
		x := float64(i) * 10.0 / float64(nPoints)
		xData[i] = x

		// Modell: Gauß + linear + Hintergrund
		gauss := realA * math.Exp(-math.Pow(x-realMu, 2)/(2*realSigma*realSigma))
		linear := realSlope * x
		bg := realBg

		// Füge zufälliges Rauschen hinzu
		noise := rand.NormFloat64() * 0.5
		yData[i] = gauss + linear + bg + noise
		yerr[i] = 0.5 // konstanter Fehler für jeden Punkt
	}

	return xData, yData, yerr
}

func NewExperimentalFcn(xData, yData, yerr []float64) *ExperimentalFcn {
	return &ExperimentalFcn{
		xData: xData,
		yData: yData,
		yerr:  yerr,
	}
}

// ValueOf berechnet Chi-Quadrat zwischen Modell und Daten
func (f *ExperimentalFcn) ValueOf(par []float64) float64 {
	a := par[0]     // Amplitude
	mu := par[1]    // Mittelwert
	sigma := par[2] // Breite
	bg := par[3]    // Hintergrund
	slope := par[4] // Linearer Term

	chiSquare := 0.0

	for i := range f.xData {
		x := f.xData[i]
		y := f.yData[i]
		err := f.yerr[i]

		// Modell: Gauß + linear + Hintergrund
		gauss := a * math.Exp(-math.Pow(x-mu, 2)/(2*sigma*sigma))
		model := gauss + slope*x + bg

		// Chi-Quadrat Beitrag
		residual := (y - model) / err
		chiSquare += residual * residual
	}

	return chiSquare
}

func TestExperimentalFit(t *testing.T) {
	// Generiere Testdaten
	xData, yData, yerr := generateTestData(50)

	// Erstelle Fit-Funktion
	theFCN := NewExperimentalFcn(xData, yData, yerr)

	// Parameter Setup mit Startwerten
	upar := minuit.NewEmptyMnUserParameters()
	upar.AddFree("amplitude", 5.0, 1.0)  // Startpunkt weit vom echten Wert (10.0)
	upar.AddFree("mean", 4.0, 0.5)       // Startpunkt nahe am echten Wert (5.0)
	upar.AddFree("sigma", 1.0, 0.2)      // Startpunkt weit vom echten Wert (2.0)
	upar.AddFree("background", 0.0, 0.5) // Startpunkt weit vom echten Wert (2.0)
	upar.AddFree("slope", 0.0, 0.1)      // Startpunkt weit vom echten Wert (0.5)

	fmt.Printf("Initial parameters: %s\n", upar.String())

	// Minimierung mit MIGRAD
	println("start migrad")
	migrad := minuit.NewMnMigradWithParametersStra(theFCN, upar, minuit.PreciseStrategy)
	//migrad := minuit.MinimizeWithMaxfcnToler(theFCN, upar,)
	min, err := migrad.Minimize()
	if err != nil {
		t.Fatalf("minimize failed with:\n %s\n", err.Error())
	}

	// Falls erste Minimierung nicht erfolgreich, versuche höhere Strategie
	if !min.IsValid() {
		println("FM is invalid, try with strategy = 2.")
		migrad2 := minuit.NewMnMigradWithParameterStateStrategy(theFCN, min.UserState(), minuit.NewMnStrategyWithStra(2))
		min, err = migrad2.Minimize()
		if err != nil {
			t.Fatalf("minimize failed with:\n %s\n", err.Error())
		}
	}

	// Drucke Ergebnis
	var sbuff strings.Builder
	minuit.MnPrint.PrintFunctionMinimum(&sbuff, min)
	fmt.Printf("minimum: %s\n", sbuff.String())

	// Zeige die gefundenen Parameter
	params := min.UserState().Params()
	fmt.Printf("\nGefundene Parameter:\n")
	fmt.Printf("Amplitude:   %.3f (Erwartung: ~10.0)\n", params[0])
	fmt.Printf("Mittelwert:  %.3f (Erwartung: ~5.0)\n", params[1])
	fmt.Printf("Breite:      %.3f (Erwartung: ~2.0)\n", params[2])
	fmt.Printf("Hintergrund: %.3f (Erwartung: ~2.0)\n", params[3])
	fmt.Printf("Steigung:    %.3f (Erwartung: ~0.5)\n", params[4])
}
