package cproxy

type DefaultInitializer struct {
}

func NewInitializer() *DefaultInitializer {
	return &DefaultInitializer{}
}

func (it *DefaultInitializer) Initialize(_, _ Socket) bool {
	return true
}
