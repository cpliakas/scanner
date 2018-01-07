/*
Package scanner provides a recursive file scanner that is useful for
efficiently processing relatively static datasets.

The example below captures discovered files into memory using the built-in
MemoryHandler and writes them to STDOUT:

	s := scanner.New("/path/to/dir")

	h := scanner.NewMemoryHandler()
	s.Scan(h)

	for _, f := range h.Files {
		fmt.Println(f)
	}

Get the same output as the code above by implementing the Handler interface:

	type PrintlnHandler struct{}
	func (h *PrintlnHandler) Handle(file string) { fmt.Println(file) }
	func (h *PrintlnHandler) HandleError(_ error) {}

	func scan() {
		s := scanner.New("/path/to/dir")
		h := &PrintlnHandler{}
		s.Scan(h)
	}

Configure the concurrency of file and error handling (note that the order in
which files are written will likely change each time the code is run):

	s := scanner.New("/path/to/dir")
	s.Concurrency = 5

	h := &PrintlnHandler{}
	s.Scan(h)

*/
package scanner
