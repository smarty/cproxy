package cproxy

type defaultInitializer struct{}

func newInitializer() *defaultInitializer { return &defaultInitializer{} }

func (this *defaultInitializer) Initialize(_, _ Socket) bool { return true }
