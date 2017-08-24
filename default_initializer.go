package cproxy

type DefaultInitializer struct {
}

func NewInitializer() *DefaultInitializer {
	return &DefaultInitializer{}
}

func (this *DefaultInitializer) Initialize(_, _ Socket) bool {
	return true
}
