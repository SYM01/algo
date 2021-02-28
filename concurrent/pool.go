// Package concurrent contains some easy-to-use concurrent control utilities.
package concurrent

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// NewActionPool creates a new concurrent action pool. if maxConcurrent <= 0,
// then the real maxConcurrent will be runtime.NumCPU()
func NewActionPool(maxConcurrent int) ActionPool {
	if maxConcurrent <= 0 {
		maxConcurrent = runtime.NumCPU()
	}

	return &actionPool{
		slot:   make(chan struct{}, maxConcurrent),
		logger: os.Stderr,
	}
}

// Action to be run concurrently
type Action struct {
	// Runner is a func need to be executed. It may return a list of new
	// actions. These actions will be executed after the current action had
	// finished
	Runner func() []*Action

	// Dependencies for the current action. The current action will only be
	// executed after all dependencies had been finished
	Dependencies []*Action
}

// ActionPool is an action pool based on goroutines
type ActionPool interface {
	// Do an action. It's panic-free if using it correctly.
	// It will panic if you do a nil action, or do an action in a closed pool.
	Do(*Action)

	// Wait until all the actions had been finished, and close the pool.
	WaitAndClose()

	// Set a logger to output some err msg when panic. os.Stderr will be used
	// by default
	SetLogger(io.Writer)
}

type actionPool struct {
	closed uint32
	wg     sync.WaitGroup // wg for all goroutines
	slot   chan struct{}

	logger io.Writer
}

// Do implements ActionPool's method.
func (p *actionPool) Do(a *Action) {
	if a == nil {
		panic("ActionPool: the action is nil")
	}

	if atomic.LoadUint32(&p.closed) > 0 {
		panic("ActionPool: the current pool had already been closed.")
	}

	p.wg.Add(1)
	// blocked until there is an empty slot
	p.slot <- struct{}{}

	go p.do(a, nil)
}

// do an action and free a slot when finished executing.
func (p *actionPool) do(a *Action, parentWg *sync.WaitGroup) {
	// clean up
	defer func() {
		// parent waitgroup
		if parentWg != nil {
			parentWg.Done()
		}
		<-p.slot
		p.wg.Done()
	}()
	defer p.panicCatcher()

	if len(a.Dependencies) > 0 {
		depWg := new(sync.WaitGroup)
		for _, dep := range a.Dependencies {
			p.wg.Add(1)
			depWg.Add(1)
			go p.do(dep, depWg)

			// p.do will free one slot, so we need to require a new slot again
			p.slot <- struct{}{}
		}
		depWg.Wait()
	}

	if a.Runner == nil {
		return
	}

	// run action
	newActions := a.Runner()

	for _, a := range newActions {
		p.wg.Add(1)
		go p.do(a, nil)
		// p.do will free one slot, so we need to require a new slot again
		p.slot <- struct{}{}
	}
}

func (p *actionPool) panicCatcher() {
	if r := recover(); r != nil {
		fmt.Fprintf(p.logger, "%s:\n%s\n", r, debug.Stack())
	}
}

// WaitAndClose implements ActionPool's method.
func (p *actionPool) WaitAndClose() {
	atomic.StoreUint32(&p.closed, 1)
	p.wg.Wait()
}

// SetLogger implements ActionPool's method.
func (p *actionPool) SetLogger(l io.Writer) {
	p.logger = l
}
