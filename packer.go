package tpack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

// Packer facilitates writing programs that manipulate text streams, a
// fundamental concept in the Unix philosophy. It enables seamless data
// transfer between processes, allowing the output of one process to be
// used as input for another.
type Packer struct {
	wg        sync.WaitGroup
	in        io.Reader
	out       io.Writer
	err       io.Writer
	processor Processor
	options   *packerOpts
}

// packerOpts provides additional configuration properties to Packer.
type packerOpts struct {
	errWriteFunc func(error)
}

// newPackerOptsDefault returns a new packerOpts with default values.
func newPackerOptsDefault() *packerOpts {
	return &packerOpts{
		errWriteFunc: func(err error) { panic(err) },
	}
}

// PackerOpt represents a customization option to configure a packer.
type PackerOpt func(*packerOpts)

// WithErrWriteHandler configures a custom error handler for writing to err.
// The default error handler will panic on an error.
func WithErrWriteHandler(f func(error)) PackerOpt {
	return func(opts *packerOpts) {
		opts.errWriteFunc = f
	}
}

// NewPacker returns a new [Packer] using custom communication channels.
func NewPacker(in io.Reader, out, err io.Writer, processor Processor,
	opts ...PackerOpt) *Packer {
	options := newPackerOptsDefault()

	// apply specified options
	for _, opt := range opts {
		opt(options)
	}

	return &Packer{
		in:        in,
		out:       out,
		err:       err,
		processor: processor,
		options:   options,
	}
}

// NewPackerStdOut returns a new [Packer] that uses the standard output and
// error streams as output channels and the provided [io.Reader] as the
// input communication channel.
func NewPackerStdOut(in io.Reader, processor Processor,
	opts ...PackerOpt) *Packer {
	return NewPacker(in, os.Stdout, os.Stderr, processor, opts...)
}

// NewPackerStd returns a new [Packer] that uses the standard input, output,
// and error streams as communication channels.
func NewPackerStd(processor Processor, opts ...PackerOpt) *Packer {
	return NewPackerStdOut(os.Stdin, processor, opts...)
}

// Execute starts processing data stream.
func (p *Packer) Execute() {
	p.wg.Add(2)
	go p.writeOut()
	go p.writeErr()

	// Read newline-delimited lines of text from the input reader.
	// In a Unix pipeline, the input text stream is delimited by a newline
	// character. Each line of input is passed as a separate argument
	// to the subsequent commands in the pipeline.
	scanner := bufio.NewScanner(p.in)

	for scanner.Scan() {
		p.processor.InChan() <- scanner.Bytes()
	}

	// check for scanner error
	if scanner.Err() != nil {
		p.processor.ErrChan() <- scanner.Err()
	}

	close(p.processor.InChan())
	p.wg.Wait()
}

// writeOut processes the output channel.
func (p *Packer) writeOut() {
	for message := range p.processor.OutChan() {
		_, err := p.out.Write(lf(message))
		if err != nil {
			p.handleOutWriteError(fmt.Errorf("%w: %s", err, string(message)))
		}
	}
	p.wg.Done()
}

// writeErr processes the error channel.
func (p *Packer) writeErr() {
	for err := range p.processor.ErrChan() {
		_, werr := p.err.Write(lf([]byte(err.Error())))
		if werr != nil {
			p.options.errWriteFunc(fmt.Errorf("%w: %w", werr, err))
		}
	}
	p.wg.Done()
}

// handleOutWriteError attempts to write the error using the error writer.
// If it fails, it will utilize the function from options to manage the error.
func (p *Packer) handleOutWriteError(err error) {
	_, werr := p.err.Write(lf([]byte(err.Error())))
	if werr != nil {
		p.options.errWriteFunc(fmt.Errorf("%w: %w", werr, err))
	}
}

// lf appends a newline that has been trimmed by the scanner.
func lf(bytes []byte) []byte {
	return append(bytes, '\n')
}
