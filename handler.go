package scanner

// Handler is an interface for handlers that process scanned files.
type Handler interface {

	// Handle is passed the path to the file that was discovered during the
	// scan operation. Files are passed in a single-threaded manner, so
	// concurrency is the responsibility of the Handler implementation.
	Handle(string)

	// HandleError is passed the errors that occur during scan operations.
	// Like the Handle method, errors are passed in a single-threaded
	// manner.
	HandleError(error)
}

// NullHandler implements Handler and performs no-ops when it is passed files
// and errors. This is most often used in testing and benchmarking.
type NullHandler struct{}

// Handle implements Handler.Handle and performs a no-op when a file is passed
// to it.
func (h *NullHandler) Handle(_ string) {}

// HandleError implements Handler.HandleError and performs a no-op when an
// error is passed to it.
func (h *NullHandler) HandleError(_ error) {}

// NewMemoryHandler returns a new MemoryHandler instance, which stores the
// scanned files and errors in the exported Files and Errors fields
// respectively.
func NewMemoryHandler() *MemoryHandler {
	return &MemoryHandler{}
}

// MemoryHandler implements Handler and stores the scanned files and errors
// in slices.
type MemoryHandler struct {

	// Files is the slice that stores discovered files in memory.
	Files []string

	// Files is the slice that stores errors in memory.
	Errors []error
}

// Handle implements Handler.Handle and stores file in a slice.
func (h *MemoryHandler) Handle(file string) {
	h.Files = append(h.Files, file)
}

// HandleError implements Handler.HandleError stores err in a slice.
func (h *MemoryHandler) HandleError(err error) {
	h.Errors = append(h.Errors, err)
}
