package apm

type DebuggableInterval interface {
	Start(name string) Interval
}

type Interval func()

func (i Interval) End() {
	if i == nil {
		return
	}

	i()
}
