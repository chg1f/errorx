errorx
---

`errorx` provides typed Go errors with mandatory codes and optional stack capture.

GoDoc: <https://pkg.go.dev/github.com/chg1f/errorx/v2>

## Install

```bash
go get github.com/chg1f/errorx/v2
```

## Overview

- Every `*errorx.Error[T]` carries a comparable code.
- `errors.Is` and `errors.As` continue to work.
- Stack capture is opt-in through `errorx/stacktrace`.
- Attributes are passed with `slog.Attr`.

## Basic Example

```go
package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/chg1f/errorx/v2"
)

func main() {
	base := errors.New("disk failure")

	err := errorx.WithCode("config.invalid").
		Wrap(base, "load config", slog.String("file", "app.yaml"))

	ex := errorx.Be[string](err)
	fmt.Println(ex.Code())
	fmt.Println(ex.String())
	fmt.Println(errorx.In(ex, "config.invalid"))
}
```

## Stack Example

```go
package main

import (
	"fmt"

	"github.com/chg1f/errorx/v2"
	_ "github.com/chg1f/errorx/v2/stacktrace"
)

func main() {
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	fmt.Println(ex.Stack().LogValue())
}
```

## Package Layout

- `errorx`: core error type, builder, and helper options
- `errorx/stacktrace`: runtime stack capture integration
