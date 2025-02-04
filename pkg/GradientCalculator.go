package minuit

type GradientCalculator interface {
	Gradient(par *MinimumParameters) *FunctionGradient
	GradientWithGrad(par *MinimumParameters, grad *FunctionGradient) *FunctionGradient
}
