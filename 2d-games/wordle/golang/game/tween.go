package wordle

import "fmt"

type TweenFunc func(elapsed, begin, end, duration float64) float64
type TweenUpdateFunc func(v float64)
type TweenCompletedFunc func()

type Tween struct {
	begin         float64
	end           float64
	duration      float64
	elapsed       float64
	tweenFunc     TweenFunc
	updateFunc    TweenUpdateFunc
	completedFunc TweenCompletedFunc
	isCompleted   bool
}

func LinearTweenFunc(elapsed, begin, end, duration float64) float64 {
	return begin + ((end - begin) * (elapsed / duration))
}

func NewTween(begin, end, duration float64, tweenFunc TweenFunc, updateFunc TweenUpdateFunc, completedFunc TweenCompletedFunc) *Tween {
	if begin >= end {
		panic(fmt.Sprintf("begin cannot be bigger than end. begin: %v, end: %v", begin, end))
	}

	return &Tween{
		begin:         begin,
		end:           end,
		duration:      duration,
		tweenFunc:     tweenFunc,
		updateFunc:    updateFunc,
		completedFunc: completedFunc,
		elapsed:       0,
		isCompleted:   false,
	}
}

func (t *Tween) Update(dt float64) {
	if t.isCompleted {
		return
	}

	t.elapsed += dt
	if t.elapsed >= t.duration {
		t.elapsed = t.duration
	}

	v := t.tweenFunc(t.elapsed, t.begin, t.end, t.duration)

	if t.updateFunc != nil {
		t.updateFunc(v)
	}

	if v >= t.end {
		t.isCompleted = true

		if t.completedFunc != nil {
			t.completedFunc()
		}
	}
}

func (t *Tween) Complete() {
	t.elapsed = t.duration
	t.isCompleted = true

	if t.completedFunc != nil {
		t.completedFunc()
	}
}

func UpdateTween(current *Tween, new *Tween) *Tween {
	if current != nil {
		current.Complete()
	}

	return new
}
