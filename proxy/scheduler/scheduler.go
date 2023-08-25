package scheduler

import "time"

type DelayedTask struct {
	Delay  int64
	action func()
}

type DelayedRepeatingTask struct {
	Delay   int64
	Seconds int64
	action  func()
}

type RepeatingTask struct {
	Seconds int64
	action  func()
}

func NewDelayedTask(delay int64, action func()) {
	dt := &DelayedTask{
		Delay:  delay,
		action: action,
	}

	go func() {
		time.Sleep(time.Duration(dt.Delay) * time.Second)
		dt.onRun()
	}()
}

func (dt *DelayedTask) onRun() {
	dt.action()
}

func NewDelayedRepeatingTask(delay, seconds int64, action func()) {
	drt := &DelayedRepeatingTask{
		Delay:   delay,
		Seconds: seconds,
		action:  action,
	}

	go func() {
		time.Sleep(time.Duration(drt.Delay) * time.Second)
		drt.onRun()

		for {
			time.Sleep(time.Duration(drt.Seconds) * time.Second)
			drt.onRun()
		}
	}()
}

func (drt *DelayedRepeatingTask) onRun() {
	drt.action()
}

func NewRepeatingTask(seconds int64, action func()) {
	drt := &RepeatingTask{
		Seconds: seconds,
		action:  action,
	}

	go func() {
		for {
			drt.onRun()
			time.Sleep(time.Duration(drt.Seconds) * time.Second)
		}
	}()
}

func (drt *RepeatingTask) onRun() {
	drt.action()
}
