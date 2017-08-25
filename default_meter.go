package cproxy

type DefaultMeter struct {
}

func NewMeter() *DefaultMeter {
	return &DefaultMeter{}
}

func (this *DefaultMeter) Measure(int) {
}
