package scanner_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cpliakas/scanner"
)

func TestScanner(t *testing.T) {
	s := scanner.New("test/data")

	h := scanner.NewMemoryHandler()
	s.Scan(h)

	if len(h.Files) != 3 {
		t.Fatalf("expected 3 files, got %v", len(h.Files))
	}
	if h.Files[0] != "test/data/file1.txt" {
		t.Fatalf("expected 'test/data/file1.txt', got '%s'", h.Files[0])
	}
	if h.Files[1] != "test/data/subdir/file2.txt" {
		t.Fatalf("expected 'test/data/subdir/file2.txt', got '%s'", h.Files[1])
	}
	if h.Files[2] != "test/data/subdir/file3.txt" {
		t.Fatalf("expected 'test/data/subdir/file3.txt', got '%s'", h.Files[2])
	}
}

func TestScannerError(t *testing.T) {
	s := scanner.New("test/data/baddir")

	h := scanner.NewMemoryHandler()
	s.Scan(h)

	if len(h.Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", len(h.Errors))
	}
	if !os.IsNotExist(h.Errors[0]) {
		t.Fatal("expected 'file does not exist' error")
	}
}

func TestScannerNilHandler(t *testing.T) {
	s := scanner.New("test/data")
	s.Scan(nil)
}

func ExampleScanner_Scan() {
	s := scanner.New("test/data")
	h := scanner.NewMemoryHandler()

	s.Scan(h)

	if len(h.Files) > 0 {
		fmt.Println(h.Files[0])
	}

	// Output: test/data/file1.txt
}

func BenchmarkScanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := scanner.New("test/data")
		s.Scan(nil)
	}
}
