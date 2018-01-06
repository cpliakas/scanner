// Package scanner provides a recursive file scanner that is useful for
// efficiently processing relatively static datasets.
package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Scanner recursively scans a directory for files.
type Scanner struct {

	// Concurrency is maximum number of handling operations that are run
	// at the same time. The default value is 1, resulting in single-threaded
	// file handling.
	Concurrency int

	// Path is the directory being scanned.
	Path string

	// errs is the channel through which Scanner passes errors to Handler.
	errs chan error

	// errs is the channel through which Scanner passes files to Handler.
	files chan string

	// wait is the channel that limits concurrent handling of files and errors.
	wait chan bool

	// wg is the wait group that the scanner uses to wait for all goroutines to
	// complete.
	wg *sync.WaitGroup
}

// New returns a new Scanner instance.
func New(path string) *Scanner {
	return &Scanner{
		Concurrency: 1,
		Path:        path,
	}
}

// Scan recursively scans the directory and sends the files and errors to
// the passed Handler in goroutines.
func (s *Scanner) Scan(h Handler) {

	// Concurrency cannot be less than 1, because a negative buffer argument
	// when making a channel is not allowed (the argument is Concurrency - 1).
	if s.Concurrency < 1 {
		panic("Scanner.Concurrency must be >= 1")
	}

	// Default to NullHandler.
	if h == nil {
		h = &NullHandler{}
	}

	s.errs = make(chan error)
	s.files = make(chan string)
	s.wait = make(chan bool, s.Concurrency-1)
	s.wg = &sync.WaitGroup{}

	// Start the scanning pipeline.
	s.scan()
	go s.handleFiles(h)
	go s.handleErrors(h)

	s.wg.Wait()
}

// scan starts the scanning process and returns the channels that discovered
// files and errors are sent to.
func (s *Scanner) scan() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(s.files)
		defer close(s.errs)
		s.readDir(s.Path)
	}()
	return
}

// handleFiles listens for errors and passes them to the handler. This method
// is expected to be run in a goroutine.
func (s *Scanner) handleFiles(h Handler) {
	for fname := range s.files {
		s.wg.Add(1)
		go func(fname string) {
			defer s.wg.Done()
			h.Handle(fname)
			<-s.wait
		}(fname)
	}
}

// handleErrors listens for errors and passes them to the handler. This method
// is expected to be run in a goroutine.
func (s *Scanner) handleErrors(h Handler) {
	for err := range s.errs {
		s.wg.Add(1)
		go func(err error) {
			defer s.wg.Done()
			h.HandleError(err)
			<-s.wait
		}(err)
	}
}

// readDir recursively scans path and sends the discovered files and errors to
// the files and errs channels respectively.
func (s *Scanner) readDir(dirname string) {
	dir, err := ioutil.ReadDir(dirname)
	if err != nil {
		s.errs <- err
		s.wait <- true
		return
	}

	for _, file := range dir {
		path := filepath.Join(dirname, file.Name())
		if file.IsDir() {
			s.readDir(path)
		} else if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			// TODO Figure out how to handle symlinks.
		} else {
			s.files <- path
			s.wait <- true
		}
	}
}

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

// NullHandler implements Handler and performs no-ops.
type NullHandler struct{}

// Handle implements Handler.Handle and performs a no-op.
func (h *NullHandler) Handle(_ string) {}

// HandleError implements Handler.HandleError and performs a no-op.
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
	Files  []string
	Errors []error
}

// Handle implements Handler.Handle and stores file in a struct.
func (h *MemoryHandler) Handle(file string) {
	h.Files = append(h.Files, file)
}

// HandleError implements Handler.HandleError stores err in a struct.
func (h *MemoryHandler) HandleError(err error) {
	h.Errors = append(h.Errors, err)
}
