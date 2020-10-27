package cproxy

type defaultInitializer struct{}

func newInitializer() *defaultInitializer { return &defaultInitializer{} }

func (this *defaultInitializer) Initialize(_, _ socket) bool { return true }
