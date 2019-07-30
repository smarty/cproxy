package cproxy

type DefaultMeter struct {
}

func NewMeter() *DefaultMeter {
	return &DefaultMeter{}
}

func (it *DefaultMeter) Measure(int) {
}
