# go-dict

`go-dict` is a golang package to manage and lookup a local dictionary.

## Documentation

[http://godoc.org/github.com/ddliu/go-dict](View documentaion on godoc.org)

## Installation

```bash
go get github.com/ddliu/go-dict
```

## Usage

### Simple dict(words list)

```go
package main
import (
   "github.com/ddliu/go-dict" 
   "fmt"
)

func main() {
    d := dict.NewSimpleDict()
    d.Load("/usr/share/dict/words")

    // Lookup dictionary with regexp, and return first 20 results
    list := d.Lookup("a[aeiou]{2}", 0, 20)

    // Loop through the dictionary
    d.Walk(func(w string) bool {
        fmt.Println(w)
        return true
    })
}
```

### Dict with properties

```go
package main

import (
    "github.com/ddliu/go-dict"
    "fmt"
)

func main() {
    d := NewDict()

    d.AddMap("Duck", map[string]interface{}{
        "Legs": 2,
        "Swim": true,
        "Fly": false,
    })

    // prop = NewDictWord()
    d.AddMap("Dog", map[string]interface{}{
        "Legs": 4,
        "Swim": true,
        "Fly": false,
    })

    d.AddMap("Snake", map[string]interface{}{
        "Legs": 0,
        "Swim": false,
        "Fly": false,
    })

    d.AddMap("Bird", map[string]interface{}{
        "Legs": 2,
        "Swim": false,
        "Fly": true,
    })

    d.AddMap("Lion", map[string]interface{}{
        "Legs": 4,
        "Swim": false,
        "Fly": false,
        "Color": "yellow",
    })

    color := d.MustGet("Lion").MustPropString("Color")

    d.Export("/tmp/dict")
}
```

## Changelog

### v0.1.0 (2013-10-09)

Initial release