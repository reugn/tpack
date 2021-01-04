package tpack

import (
	"bufio"
	"io"
	"os"
	"sync"
)

// Packer is the representation of a packed processing unit.
type Packer struct {
	sync.WaitGroup
	in        io.Reader
	out       io.Writer
	err       io.Writer
	processor Processor
	validate  func()
}

// NewPacker returns a new Packer.
func NewPacker(in io.Reader, out io.Writer, err io.Writer, processor Processor) *Packer {
	return &Packer{
		in:        in,
		out:       out,
		err:       err,
		processor: processor,
		validate:  nil,
	}
}

// NewPackerStdOut returns a new Packer with the standard output and standard error
// as output communication channels and the specified input io.Reader.
func NewPackerStdOut(in io.Reader, processor Processor) *Packer {
	return &Packer{
		in:        in,
		out:       os.Stdout,
		err:       os.Stderr,
		processor: processor,
		validate:  nil,
	}
}

// NewPackerStd returns a new Packer with the standard streams as communication channels.
func NewPackerStd(processor Processor) *Packer {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	validate := func() {
		info, err := stdin.Stat()
		if err != nil {
			stderr.Write([]byte(err.Error()))
			os.Exit(-1)
		}
		if info.Mode()&os.ModeNamedPipe == 0 {
			stderr.Write([]byte("named pipe (FIFO)"))
			os.Exit(-1)
		}
	}

	return &Packer{
		in:        stdin,
		out:       stdout,
		err:       stderr,
		processor: processor,
		validate:  validate,
	}
}

// Execute starts processing the incoming messages.
func (p *Packer) Execute() {
	if p.validate != nil {
		p.validate()
	}

	go p.writeOut()
	go p.writeErr()
	p.Add(2)

	scanner := bufio.NewScanner(p.in)
	for scanner.Scan() {
		p.processor.InChannel() <- scanner.Bytes()
	}
	close(p.processor.InChannel())
	p.Wait()
}

func (p *Packer) writeOut() {
	for message := range p.processor.OutChannel() {
		_, err := p.out.Write(lf(message))
		if err != nil {
			p.handleErrorOnWrite(err)
		}
	}
	p.Done()
}

func (p *Packer) writeErr() {
	for message := range p.processor.ErrChannel() {
		_, err := p.err.Write(lf(message))
		if err != nil {
			p.handleErrorOnWrite(err)
		}
	}
	p.Done()
}

func (p *Packer) handleErrorOnWrite(err error) {
	_, err = p.err.Write(lf([]byte(err.Error())))
	if err != nil {
		os.Exit(-1)
	}
}

func lf(bytes []byte) []byte {
	return append(bytes, '\n')
}
