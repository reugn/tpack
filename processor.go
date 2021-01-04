package tpack

// Processor represents a generic stream processor.
type Processor interface {

	// InChannel returns the input communication channel.
	InChannel() chan []byte

	// OutChannel returns the output communication channel.
	OutChannel() chan []byte

	// ErrChannel returns the error output communication channel.
	ErrChannel() chan []byte
}
