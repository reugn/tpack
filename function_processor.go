package tpack

// FunctionProcessor implements Processor using a function to process messages.
type FunctionProcessor struct {
	in       chan []byte
	out      chan []byte
	err      chan []byte
	function func(in []byte) ([][]byte, error)
}

// NewFunctionProcessor returns a new FunctionProcessor with the specified function.
func NewFunctionProcessor(function func(in []byte) ([][]byte, error)) *FunctionProcessor {
	processor := &FunctionProcessor{
		in:       make(chan []byte),
		out:      make(chan []byte),
		err:      make(chan []byte),
		function: function,
	}

	go processor.init()
	return processor
}

func (fp *FunctionProcessor) init() {
	for message := range fp.in {
		out, err := fp.function(message)
		if err != nil {
			fp.err <- []byte(err.Error())
		}
		for _, outMsg := range out {
			fp.out <- outMsg
		}
	}
	close(fp.out)
	close(fp.err)
}

// InChannel returns the input communication channel.
func (fp *FunctionProcessor) InChannel() chan []byte {
	return fp.in
}

// OutChannel returns the output communication channel.
func (fp *FunctionProcessor) OutChannel() chan []byte {
	return fp.out
}

// ErrChannel returns the error output communication channel.
func (fp *FunctionProcessor) ErrChannel() chan []byte {
	return fp.err
}
