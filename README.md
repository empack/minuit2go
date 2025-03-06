# minuit2go

Welcome to **minuit2go**, a Golang port of the Minuit optimization library.

## Overview

Minuit is a popular numerical optimization library originally written in Fortran and later ported to C++ AND THEN ported to Java. This project aims to provide a Go implementation of the Minuit library, making it accessible to Go developers.

## Features

- **Optimization Algorithms**: Implements various optimization algorithms from the original Minuit library.
- **Easy Integration**: Seamlessly integrates with Go projects.
- **Performance**: Efficient and optimized for performance.

## Installation

To install minuit2go, use the following `go get` command:

```bash
go get github.com/empack/minuit2go
```

## Usage

Here is a basic example of how to use minuit2go in your project:

```go
package main

import (
    "fmt"
    "github.com/empack/minuit2go"
)

func main() {
    // Example usage of minuit2go
    optimizer := minuit2go.NewOptimizer()
    result := optimizer.Optimize(func(x []float64) float64 {
        // Define your objective function here
        return x[0]*x[0] + x[1]*x[1]
    }, []float64{1.0, 2.0})

    fmt.Printf("Optimization Result: %+v\n", result)
}
```

## Contact

For any questions or suggestions, feel free to open an issue or contact the project maintainers.
