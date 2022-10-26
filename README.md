# Ordered Map for golang

Python's [OrderedDict](https://docs.python.org/3/library/collections.html?highlight=ordereddict) for Golang.

## Install

```bash
go install github.com/kasegao/go-orderedmap@latest
```

## Examples

```go
package main

import (
	"fmt"

	omap "github.com/kasegao/go-orderedmap"
)

func main() {
	om := omap.New[int, string]()

	om.Set(0, "foo")
	om.Set(1, "bar")

	if v, ok := om.Get(0); ok {
		fmt.Println(v) // print "foo"
	}

	om.Delete(0)
	if _, ok := om.Get(0); !ok {
		fmt.Println("not found")
	}

	fmt.Println(om) // print "OrderedMap[1:bar]"
}
```
