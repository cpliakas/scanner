// Package scanner provides a recursive file scanner that is useful for
// efficiently processing relatively static datasets.
package scanner

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Scanner recursively scans a directory for files.
type Scanner struct {

	// Buffer is length of the files channel buffer.
	Buffer int

	// errs is the channel that errors are sent to.
	errs chan error

	// files is the challen that discovered files are sent to.
	files chan string

	// path is the directory being scanned.
	path string

	// separator is the OS's path separator represented as a string.
	separator string
}

// New returns a new Scanner instance.
func New(path string) *Scanner {
	return &Scanner{
		Buffer:    1,
		errs:      make(chan error),
		path:      path,
		separator: string(os.PathSeparator),
	}
}

// Scan recursively scans the directory and sends the files and errors to
// the passed Handler in goroutines.
func (s *Scanner) Scan(h Handler) {
	var wg sync.WaitGroup
	s.files = make(chan string, s.Buffer)

	go func() {
		s.scan(s.path)
		close(s.files)
		close(s.errs)
	}()

	if h != nil {
		wg.Add(2)

		go func() {
			defer wg.Done()
			for err := range s.errs {
				h.HandleError(err)
			}
		}()

		go func() {
			defer wg.Done()
			for f := range s.files {
				h.Handle(f)
			}
		}()
	}

	wg.Wait()
}

// scan recursively scans path and sends the discovered files and errors to
// the built-in channels.
func (s *Scanner) scan(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		s.errs <- err
		return
	}

	basedir := strings.TrimRight(path, s.separator)
	for _, f := range files {
		file := basedir + s.separator + f.Name()
		if f.IsDir() {
			s.scan(file)
		} else {
			s.files <- file
		}
	}
}

// Handler is an interface for handlers that process scanned files.
type Handler interface {
	Handle(string)
	HandleError(error)
}

// MemoryHandler returns a new memoryHandler instance.
func MemoryHandler() *memoryHandler {
	return &memoryHandler{}
}

// memoryHandler stores the scanned files and errors in slices.
type memoryHandler struct {
	Files  []string
	Errors []error
}

// Handle stores file in a struct.
func (h *memoryHandler) Handle(file string) {
	h.Files = append(h.Files, file)
}

// HandleError stores err in a struct.
func (h *memoryHandler) HandleError(err error) {
	h.Errors = append(h.Errors, err)
}
