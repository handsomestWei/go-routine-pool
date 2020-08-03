package pool

import (
	"sync"
)

type basePoolExecutor struct {
	corePoolSize chan int
	fsm          *simpleFSM
	oc           sync.Once
}

type poolExecutor struct {
	basePoolExecutor
}

type waitPoolExecutor struct {
	basePoolExecutor
	gp *sync.WaitGroup
}

func newBasePoolExecutor(cap int) basePoolExecutor {
	return basePoolExecutor{
		corePoolSize: make(chan int, cap),
		fsm:          newSimpleFSM(),
		oc:           sync.Once{},
	}
}

func NewPool(cap int) *poolExecutor {
	checkCap(cap)
	return &poolExecutor{
		basePoolExecutor: newBasePoolExecutor(cap),
	}
}

func NewWaitPool(cap int) *waitPoolExecutor {
	checkCap(cap)
	return &waitPoolExecutor{
		basePoolExecutor: newBasePoolExecutor(cap),
		gp:               new(sync.WaitGroup),
	}
}

func checkCap(cap int) {
	if cap < 0 {
		panic("The pool cap cannot lower zero")
	}
}

func (b *basePoolExecutor) checkSubmit(f func()) {
	if f == nil {
		panic("The submit func is nil")
	}
	b.oc.Do(func() {
		b.fsm.actEvent(running)
	})
	if !b.fsm.isRunning() {
		panic("The pool is not running")
	}
}

func (b *basePoolExecutor) ShutDown() {
	b.fsm.actEvent(shutdown)
}

func (b *basePoolExecutor) IsShutDown() bool {
	return b.fsm.Current() >= shutdown
}

func (b *basePoolExecutor) IsTerminated() bool {
	return b.fsm.Current() >= stop
}

func (t *poolExecutor) Submit(f func()) {
	t.checkSubmit(f)
	t.corePoolSize <- 1
	go func() {
		defer func() {
			if err := recover(); err != nil {
				panic(err)
			}
			<-t.corePoolSize
		}()
		if t.IsTerminated() {
			return
		}
		f()
	}()
}

func (t *waitPoolExecutor) Submit(f func()) {
	t.checkSubmit(f)
	t.gp.Add(1)
	t.corePoolSize <- 1
	go func() {
		defer func() {
			if err := recover(); err != nil {
				panic(err)
			}
			<-t.corePoolSize
			t.gp.Done()
		}()
		if t.IsTerminated() {
			return
		}
		f()
	}()
}

func (t *waitPoolExecutor) Wait() {
	t.gp.Wait()
	t.fsm.actEvent(stop)
	close(t.corePoolSize)
}
