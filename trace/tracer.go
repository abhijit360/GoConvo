package trace

import (
	"fmt"
	"io"
)

// tracer is the interface that describes an object capable of
// tracing events through code.
type Tracer interface {
	// ...interface{} : is syntax for zero or more arguments of any type
	Trace(...interface{})
}

type nilTracer struct{}

func (t *nilTracer) Trace (a ...interface{}){}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}){
	fmt.Fprint(t.out,a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func Off() Tracer {
	return &nilTracer{}
}