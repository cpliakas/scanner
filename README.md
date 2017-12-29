# Scanner

[![Build Status](https://travis-ci.org/cpliakas/scanner.svg?branch=master)](https://travis-ci.org/cpliakas/scanner)
[![codecov](https://codecov.io/gh/cpliakas/scanner/branch/master/graph/badge.svg)](https://codecov.io/gh/cpliakas/scanner)
[![GoDoc](https://godoc.org/github.com/cpliakas/scanner?status.svg)](https://godoc.org/github.com/cpliakas/scanner)
[![Go Report Card](https://goreportcard.com/badge/github.com/cpliakas/scanner)](https://goreportcard.com/report/github.com/cpliakas/scanner)

Package scanner provides a recursive file scanner that is useful for
efficiently processing relatively static datasets.

By using Go's concurrency primitives, this package provides a framework to
decouple file discovery from processing of the discovered files through
implementing the `scanner.Handler` interface.

## Installation

Assuming a [correctly configured](https://golang.org/doc/install#testing) Go
toolchain:

```shell
go get github.com/cpliakas/scanner
```

## Usage

The code below recursively discovers files in the `/path/to/dir` directory and
writes their paths to STDOUT.

```go
package main

import (
	"fmt"

	"github.com/cpliakas/scanner"
)

func main() {

	// Set the directory that will be recursively scanned for files.
	s := scanner.New("/path/to/dir")

	// Scan the directory, capture the file paths in memory.
	h := scanner.NewMemoryHandler()
	s.Scan(h)

	// Print the paths of the discovered files to STDOUT.
	for _, f := range h.Files {
		fmt.Println(f)
	}
}
```
