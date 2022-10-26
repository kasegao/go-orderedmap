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
	om := omap.New[string, string]()

	om.Set("a", "foo")
	om.Set("b", "bar")
	om.Set("c", "baz")

	if v, ok := om.Get("a"); ok {
		fmt.Println(v) // print "foo"
	}

	if e, ok := om.GetAt(1); ok {
		fmt.Println(e.Key, e.Value) // print "bar"
	}

	om.DeleteAt(-1) // delete "c:baz"
	if _, ok := om.Get("c"); !ok {
		fmt.Println("not found") // print "not found"
	}

	fmt.Println(om) // print "OrderedMap[a:foo b:bar]"
}
```
