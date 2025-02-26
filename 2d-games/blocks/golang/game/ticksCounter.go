package game

type TicksCounter struct {
	ticks int
	value int
}

func NewTicksCounter(ticks int) *TicksCounter {
	return &TicksCounter{
		ticks: ticks,
		value: 0,
	}
}

func (t *TicksCounter) Update() bool {
	t.value = (t.value + 1) % t.ticks

	return t.value == 0
}

func (t *TicksCounter) SetTicks(ticks int) {
	if ticks <= 0 {
		panic("ticks must be bigger than 0")
	}

	t.ticks = ticks
}
