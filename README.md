# Go FlatZinc

[![Go Reference](https://pkg.go.dev/badge/github.com/rhartert/gofzn.svg)](https://pkg.go.dev/github.com/rhartert/gofzn)
[![Go Report Card](https://goreportcard.com/badge/github.com/rhartert/gofzn)](https://goreportcard.com/report/github.com/rhartert/gofzn)
[![Tests](https://github.com/rhartert/gofzn/actions/workflows/test.yml/badge.svg)](https://github.com/rhartert/gofzn/actions/workflows/test.yml)

GoFZN is a parser for FlatZinc models written in Go. The goal of this project is 
to foster the development of new constraint solvers in Go by providing 
researchers and constraint programming practitioners with a convenient way to 
interface their solvers with MiniZinc.

Under the hood, GoFZN is a handwritten [recursive descent parser] with a 
structure that closely mirrors the FlatZinc grammar. Its tokenizer (or lexer) 
is inspired by Rob Pike's talk [Lexical Scanning in Go]. 

## What's FlatZinc?

FlatZinc is a subset of [MiniZinc], a high-level constraint modeling language
designed to easily express and solve discrete optimization problems. Think of 
MiniZinc as the language humans use to write models, while FlatZinc serves as
the language constraint solvers use to read and process these models.

## Usage

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

Note that GoFZN only takes care of verifying that components in the `fzn.Model` 
are *syntactically* correct. For example, the following variable declaration 
will be parsed succesfully despite having an inconsistent domain.

```
var 10..0: X; // domain is inconsistent
```

### Interfacing Directly with a Solver

You can interface your solver directly with GoFZN by providing the `fzn.Parse` 
function with your own implementation of the `fzn.Handler` interface. The main 
advantage of this solution is that it completely removes the steps to build 
the `fzn.Model` struct, thus enabling a slightly more efficient use of the 
library.

## Contributions

Contributions are welcome! Please feel free to submit a pull request or open an 
issue. Also, don't hesitate to reach out at [ren.hartert@gmail.com] if you plan 
to use this parser for your own project.

[ren.hartert@gmail.com]: mailto:ren.hartert@gmail.com
[MiniZinc]: https://www.minizinc.org/
[Lexical Scanning in Go]: https://www.youtube.com/watch?v=HxaD_trXwRE
[recursive descent parser]: https://en.wikipedia.org/wiki/Recursive_descent_parser