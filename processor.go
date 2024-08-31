package tpack

import (
	"sync"
)

// Processor represents a generic data stream processor.
type Processor interface {
	// InChan returns the input communication channel.
	InChan() chan []byte

	// OutChan returns the output communication channel.
	OutChan() chan []byte

	// ErrChan returns the error output communication channel.
	ErrChan() chan error
}

// data is an interface that represents a type that can be either a string
// or a slice of bytes.
type data interface {
	~string | ~[]byte
}

// procOpts provides additional configuration properties to processor.
type procOpts struct {
	parallelism int
}

// newProcOptsDefault returns a new procOpts with default values.
func newProcOptsDefault() *procOpts {
	return &procOpts{
		parallelism: 1,
	}
}

// ProcOpt represents a customization option to configure a processor.
type ProcOpt func(*procOpts)

// Parallel configures the processor parallelism. The specified value is
// required to be greater than zero. Otherwise, it is ignored.
// The default parallelism is 1.
//
// Parallel workflows can be useful when the processing order is not
// important - for further counting or any other type of aggregation command.
func Parallel(p int) ProcOpt {
	return func(opts *procOpts) {
		if p > 0 {
			opts.parallelism = p
		}
	}
}

// processor implements the [Processor] interface. It consumes data from the
// in channel and executes the transformation function on each chunk of data.
// The transformed object is then forwarded to the out channel or, if the mapper
// returns an error, it is sent to the err channel.
type processor[T data] struct {
	in        chan []byte
	out       chan []byte
	err       chan error
	transform func(T) ([]T, error)
	options   *procOpts
}

var _ Processor = (*processor[[]byte])(nil)

// NewProcessor returns a new Processor with the specified transformation function.
// The transformation function type can be either string or slice of bytes.
// Additional options can be provided to customize the processor, e.g.
//
//	tpack.Parallel(2)
func NewProcessor[T data](transform func(T) ([]T, error), opts ...ProcOpt) Processor {
	options := newProcOptsDefault()

	// apply specified options
	for _, opt := range opts {
		opt(options)
	}

	processor := &processor[T]{
		in:        make(chan []byte, options.parallelism),
		out:       make(chan []byte, options.parallelism),
		err:       make(chan error, options.parallelism),
		transform: transform,
		options:   options,
	}

	// start data processing goroutine
	go processor.init()

	return processor
}

// init starts processing incoming data.
func (p *processor[T]) init() {
	var wg sync.WaitGroup
	for i := 0; i < p.options.parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for input := range p.in {
				switch f := any(p.transform).(type) {
				case func([]byte) ([][]byte, error):
					output, err := f(input)
					if err != nil {
						p.err <- err
					}
					for _, data := range output {
						p.out <- data
					}
				case func(string) ([]string, error):
					output, err := f(string(input))
					if err != nil {
						p.err <- err
					}
					for _, data := range output {
						p.out <- []byte(data)
					}
				default: // cannot happen
					panic("processor: unsupported type")
				}
			}
		}()
	}
	// wait for all processors to exit
	wg.Wait()
	// close the out and the err channels
	close(p.out)
	close(p.err)
}

// InChan returns the input communication channel.
func (p *processor[T]) InChan() chan []byte {
	return p.in
}

// OutChan returns the output communication channel.
func (p *processor[T]) OutChan() chan []byte {
	return p.out
}

// ErrChan returns the error output communication channel.
func (p *processor[T]) ErrChan() chan error {
	return p.err
}
