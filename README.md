# Jack

**WARNING**: this README is highly asipriational. Almost none of this is
implemented.

Jack is a system for writing plugins in Go. A sample follows.

Say you have a "math" package, to complete your mathematical operations from the
command line. Each operation is implemented as a plugin. So, there's an "add"
plugin, a "subtract" plugin, etc. Here's how you would write one of those:

```
package main

import "github.com/BrianHicks/jack"

type Mul struct {}

func (m *Mul) Multiply(a, b int) (int, error) {
    return a * b, nil
}

func main() {
    jack.Run(Mul{})
}
```

Assuming the output binary from that example is `mul`, the calling code would
look like this:

```
package main

import "github.com/BrianHicks/jack"

func main() {
    mul, _ := jack.Load("mul") // assuming `err` is nil for this example

    out.(int), _ := mul.Call("Multiply", 2, 2) // `out` should be 4
}
```

if you want a nicer calling interface, you can wrap `jack.Jack`:

```
package main

type Mul struct {
    mul jack.Jack
}

func NewMul() (*Mul, err) {
    inner, err := jack.Load("mul")
    return &Mul{inner}, err
}

func (m *Mul) Multiply(a, b int) int {
    out.(int), err := mul.Call("multiply", a, b)
    if err != nil {
        panic(err)
    }

    return out
}
```
