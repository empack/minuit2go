package minuit

type MnUserCovariance struct {
	theData []float64
	theNRow int
}

func NewMnUserCovariance() *MnUserCovariance {
	return &MnUserCovariance{
		theData: make([]float64, 0),
		theNRow: 0,
	}
}

func NewMnUserCovarianceWithDataNrow(data []float64, nrow int) *MnUserCovariance {
	if len(data) != nrow*(nrow+1)/2 {
		panic("MnUserCovariance: Inconsistent arguments")
	} else {
		return &MnUserCovariance{
			theData: data,
			theNRow: nrow,
		}
	}
}

func NewMnUserCovarianceWithNrow(nrow int) *MnUserCovariance {
	return &MnUserCovariance{
		theData: make([]float64, nrow*(nrow+1)/2),
		theNRow: nrow,
	}
}

func (this *MnUserCovariance) clone() *MnUserCovariance {
	return &MnUserCovariance{
		theData: this.theData,
		theNRow: this.theNRow,
	}
}

func (this *MnUserCovariance) Get(row, col int) float64 {
	if row < this.theNRow && col < this.theNRow {
		if row > col {
			return this.theData[col+row*(row+1)/2]
		} else {
			return this.theData[row+col*(col+1)/2]
		}
	} else {
		panic("MnUserCovariance: illegal arguments")
	}
}

func (this *MnUserCovariance) Set(row, col int, value float64) {
	if row < this.theNRow && col < this.theNRow {
		if row > col {
			this.theData[col+row*(row+1)/2] = value
		} else {
			this.theData[row+col*(col+1)/2] = value
		}
	} else {
		panic("MnUserCovariance: illegal arguments")
	}
}

func (this *MnUserCovariance) scale(f float64) {
	for i := 0; i < len(this.theData); i++ {
		var10000 := this.theData
		var10000[i] *= f
	}
}

func (this *MnUserCovariance) data() []float64 {
	return this.theData
}

func (this *MnUserCovariance) Nrow() int {
	return this.theNRow
}

func (this *MnUserCovariance) NCol() int {
	return this.theNRow
}

func (this *MnUserCovariance) size() int {
	return len(this.theData)
}

func (this *MnUserCovariance) String() string {
	return MnPrint.toStringUserCovariance(this)
}
