# Scanner

[![Build Status](https://travis-ci.org/cpliakas/scanner.svg?branch=master)](https://travis-ci.org/cpliakas/scanner)
[![codecov](https://codecov.io/gh/cpliakas/scanner/branch/master/graph/badge.svg)](https://codecov.io/gh/cpliakas/scanner)
[![GoDoc](https://godoc.org/github.com/cpliakas/scanner?status.svg)](https://godoc.org/github.com/cpliakas/scanner)
[![Go Report Card](https://goreportcard.com/badge/github.com/cpliakas/scanner)](https://goreportcard.com/report/github.com/cpliakas/scanner)

Package scanner provides a recursive file scanner that is useful for
efficiently processing large and relatively static datasets.

## Why?

Although recursively scanning a directory might seem trivial, this package
does the following things so you don't have to:

* Provides a framework that decouples file discovery from processing
* Uses Go's concurrency primitives for efficient file and error handling
* Handles all the fun that comes with recursive functions and symlinks
* Implements testing with high code coverage, because let's be honest, would
  you really want to write the tests for this code?

## Installation

Assuming a [correctly configured](https://golang.org/doc/install#testing) Go
toolchain:

```shell
go get github.com/cpliakas/scanner
```

## Usage

Refer to [GoDoc.org](https://godoc.org/github.com/cpliakas/scanner) for
usage examples and code snippets.
