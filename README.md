# Go FlatZinc

[![Go Reference](https://pkg.go.dev/badge/github.com/rhartert/gofzn.svg)](https://pkg.go.dev/github.com/rhartert/gofzn)
[![Go Report Card](https://goreportcard.com/badge/github.com/rhartert/gofzn)](https://goreportcard.com/report/github.com/rhartert/gofzn)
[![Tests](https://github.com/rhartert/gofzn/actions/workflows/test.yml/badge.svg)](https://github.com/rhartert/gofzn/actions/workflows/test.yml)

GoFZN is a parser for FlatZinc models written in Golang. The project has been
motivated by the absence of FlatZinc parser in Go.

GoFZN is a [recursive descent parser] with a structure that closely mirrors the 
FlatZinc grammar. GoFZN has been entirely handwritten and does not rely on 
parser generators like YACC or Bison. Why you'd ask? Because writting tokenizers 
and top-down parsers is actually [quite fun]! FlatZinc is a rather simple 
language (in terms of grammar that is) and is well-suited for simple top-down 
handwritten parsers. Having a handwritten tokenizer and parser also makes it 
easy to adjust GoFZN (e.g. to grammar changes) while enabling high-quality
 error messages in case of syntax errors. 

> ⚠️ While the repository already provides the functionality needed to interface 
> with FlatZinc, it is still in its alpha stage and likely to undergo changes. 
> However, we do not anticipate these changes to be fundamental.

## What's FlatZinc?

FlatZinc is a subset of [MiniZinc], a high-level constraint modeling language
designed to easily express and solve discrete optimization problems. Think of 
MiniZinc as the language humans use to write models, while FlatZinc serves as
the language constraint solvers use to read and process these models.

## Usage

### Installation

To install the GoFZN library, simply use the following go get command:

```bash
go get github.com/rhartert/gofzn
```

### Reading a FlatZinc Model

The easiest way to use GoFZN is to use the `fzm.ParseModel` function to read a 
`fzn.Model` struct from a reader (e.g. a file descriptor).

```go
package main 

import (
    "fmt"
    "log"
    "os"
    
    "github.com/rhartert/gofzn/fzn"
)

func main() {
    file, err := os.Open("model.fzn")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    model, err := fzn.ParseModel(file)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Parsed model: %+v\n", model)
}
```

Alternatively, you can interface your solver directly with GoFZN by providing 
the `fzn.Parse` function with your own implementation of the `fzn.Handler`
interface. This has two main advantages: (i) it allows you to use custom 
validation rules for the declared entities, and (ii) it enables a slightly more 
memory-efficient use of the library by avoiding the creation of the `fzn.Model` 
struct.

## Contributions

Contributions are welcome! Please feel free to submit a pull request or open an 
issue. Also, don't hesitate to reach out at [ren@ptrh.io] if you plan to use
this parser for your own project.

[ren@ptrh.io]: mailto:ren@ptrh.io
[MiniZinc]: https://www.minizinc.org/
[quite fun]: https://www.youtube.com/watch?v=HxaD_trXwRE
[recursive descent parser]: https://en.wikipedia.org/wiki/Recursive_descent_parser