# GLS: Goroutine Local Storage

## simple usage

```go
package main

import (
	"fmt"
	
	"github.com/xiezhenye/gls"
)


func main() {
	gls.Set("key1", "value1")
	test()
}

func test() {
	fmt.Println(gls.Get("key1"))
}
```

## work with Context
```go
package main

import (
	"context"
	"fmt"
	
	"github.com/xiezhenye/gls"
)


func main() {
	ctx := context.Background()
	ctx = gls.AttachGls(ctx)
	gls.Set("key1", "value1")
	test(ctx)
}

func test(ctx context.Context) {
	ctx = context.WithValue(ctx, "k", "v")
	fmt.Println(gls.Get("key1"))
}
```

## goroutine util
```go
package main

import (
	"context"
	"fmt"
	
	"github.com/xiezhenye/gls"
)


func main() {
	gls.Set("key1", "value1")
    gls.Go(func() {
		fmt.Println(gls.Get("key1"))
    })
	
	// with context
	ctx := gls.AttachGls(context.Background())
	gls.Set("key2", "value2")
	gls.GoWithContext(ctx, func(newCtx context.Context) {
		fmt.Println(gls.Get("key2"))
	})
}

```

