package workgroup

import (
	"sync"

	"golang.org/x/net/context"
)

// A group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero group is valid and does not cancel on error.
type group struct {
	cancel  func()

	wg      sync.WaitGroup

	errOnce sync.Once
	err     error

	results map[string]interface{}
	errors  map[string]interface{}
}

// WithContext returns a new group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &group{
		cancel: cancel,
		results:make(map[string]interface{}),
		errors: make(map[string]interface{}),
	}, ctx
}

// New returns a new group
func New() *group {
	return &group{
		results:make(map[string]interface{}),
		errors: make(map[string]interface{}),
	}
}


// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *group) Wait() (results map[string]interface{}, errors map[string]interface{}) {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.results, g.errors
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *group) Go(label string, f func() (interface{}, error)) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		result, err := f()
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.errors[label] = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}

		g.results[label] = result
	}()
}
