## simple file share and upload server

## host

```go
 http://localhost:9090/
```

## Example

```go
package main

import (
	. "github.com/soekchl/myUtils"
	sfs "github.com/soekchl/myUtils/simpleFileServer"
)

func main() {
	err := sfs.Start(":9090", ".", 10)
	if err != nil {
		Error(err)
	}
}
```