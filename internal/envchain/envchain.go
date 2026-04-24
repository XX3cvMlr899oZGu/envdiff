// Package envchain provides a pipeline builder for chaining multiple
// env transformations (filter, resolve, rename, cast, trim, etc.) in sequence.
package envchain

import "fmt"

// StepFunc is a transformation applied to an env map.
// It returns a new map and any error encountered.
type StepFunc func(env map[string]string) (map[string]string, error)

// Chain holds an ordered list of transformation steps.
type Chain struct {
	steps []namedStep
}

type namedStep struct {
	name string
	fn   StepFunc
}

// New returns an empty Chain.
func New() *Chain {
	return &Chain{}
}

// Add appends a named step to the chain.
func (c *Chain) Add(name string, fn StepFunc) *Chain {
	c.steps = append(c.steps, namedStep{name: name, fn: fn})
	return c
}

// Run executes all steps in order, passing the output of each step as the
// input to the next. It returns the final map or the first error encountered,
// annotated with the name of the failing step.
func (c *Chain) Run(env map[string]string) (map[string]string, error) {
	current := copyMap(env)
	for _, s := range c.steps {
		result, err := s.fn(current)
		if err != nil {
			return nil, fmt.Errorf("envchain: step %q failed: %w", s.name, err)
		}
		current = result
	}
	return current, nil
}

// StepCount returns the number of steps registered in the chain.
func (c *Chain) StepCount() int {
	return len(c.steps)
}

// StepNames returns the ordered list of step names.
func (c *Chain) StepNames() []string {
	names := make([]string, len(c.steps))
	for i, s := range c.steps {
		names[i] = s.name
	}
	return names
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
