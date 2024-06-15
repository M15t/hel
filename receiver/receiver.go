package receiver

// Service struct
type Service struct {
	Receiver
}

// Receiver interface
type Receiver interface {
	getLocation(string) error
}

// GetLocation calls the GetLocation method of the provided Receiver interface with the given input string and returns any error encountered during the operation.
func GetLocation(r Receiver, input string) error {
	return r.getLocation(input)
}
