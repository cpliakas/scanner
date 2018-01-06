package scanner_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cpliakas/scanner"
)

func TestScanner(t *testing.T) {
	s := scanner.New("fixtures/data")

	h := scanner.NewMemoryHandler()
	s.Scan(h)

	if len(h.Files) != 3 {
		t.Fatalf("expected 3 files, got %v", len(h.Files))
	}
	if h.Files[0] != "fixtures/data/file1.txt" {
		t.Fatalf("expected 'fixtures/data/file1.txt', got '%s'", h.Files[0])
	}
	if h.Files[1] != "fixtures/data/subdir/file2.txt" {
		t.Fatalf("expected 'fixtures/data/subdir/file2.txt', got '%s'", h.Files[1])
	}
	if h.Files[2] != "fixtures/data/subdir/file3.txt" {
		t.Fatalf("expected 'fixtures/data/subdir/file3.txt', got '%s'", h.Files[2])
	}
}

func TestScannerNoConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	done := make(chan bool)
	go func() {
		s := scanner.New("fixtures/data")
		h := &delayHandler{}
		s.Scan(h)
		done <- true
	}()

	select {
	case <-done:
		t.Error("expected timeout when handling with no concurrency")
	case <-time.After(time.Second * 2):
		return
	}
}

func TestScannerWithConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	done := make(chan bool)
	go func() {
		s := scanner.New("fixtures/data")
		s.Concurrency = 3
		h := &delayHandler{}
		s.Scan(h)
		done <- true
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second * 2):
		t.Error("timeout, expected all files to be handled concurrently")
	}
}

func TestScannerError(t *testing.T) {
	s := scanner.New("fixtures/data/baddir")

	h := scanner.NewMemoryHandler()
	s.Scan(h)

	if len(h.Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", len(h.Errors))
	}
	if !os.IsNotExist(h.Errors[0]) {
		t.Fatal("expected 'file does not exist' error")
	}
}

func TestScannerNegativeConcurrency(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	s := scanner.New("fixtures/data")
	s.Concurrency = 0
	s.Scan(nil)
}

func TestScannerNilHandler(t *testing.T) {
	s := scanner.New("fixtures/data")
	s.Scan(nil)
}

func ExampleScanner_Scan() {
	s := scanner.New("fixtures/data")

	h := scanner.NewMemoryHandler()
	s.Scan(h)

	if len(h.Files) > 0 {
		fmt.Println(h.Files[0])
	}
	// Output: fixtures/data/file1.txt
}

func BenchmarkScanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := scanner.New("fixtures/data")
		s.Scan(nil)
	}
}

// delayHandler implements Handler and simulates handling that takes some time
// in order to test concurrency.
type delayHandler struct{}

// Handle implements Handler.Handle and sleeps for one second.
func (h *delayHandler) Handle(_ string) {
	time.Sleep(time.Second * 1)
}

// HandleError implements Handler.HandleError and performs a no-op.
func (h *delayHandler) HandleError(_ error) {}
