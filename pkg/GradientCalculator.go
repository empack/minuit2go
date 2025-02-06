package minuit

type GradientCalculator interface {
	Gradient(par *MinimumParameters) (*FunctionGradient, error)
	GradientWithGrad(par *MinimumParameters, grad *FunctionGradient) (*FunctionGradient, error)
}
