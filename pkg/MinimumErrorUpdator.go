package minuit

type MinimumErrorUpdator interface {
	update(state *MinimumState, par *MinimumParameters, grad *FunctionGradient) *MinimumError
}
