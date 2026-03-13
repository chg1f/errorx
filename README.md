errorx
---

`errorx` provides typed Go errors with mandatory codes, optional localization, and optional stack capture.

GoDoc: <https://pkg.go.dev/github.com/chg1f/errorx/v2>

## Install

```bash
go get github.com/chg1f/errorx/v2
```

## Overview

- Every `*errorx.Error[T]` carries a comparable code.
- `errors.Is` and `errors.As` continue to work.
- Localization is opt-in through `errorx/i18n`.
- Stack capture is opt-in through `errorx/stack`.
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

## Localization Example

```go
package main

import (
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/chg1f/errorx/v2"
	"github.com/chg1f/errorx/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	var locales fs.FS

	if err := i18n.LoadFiles(locales); err != nil {
		panic(err)
	}

	err := errorx.WithCode("invalid").
		New("invalid", slog.String("field", "email"))

	ex := errorx.Be[string](err)
	fmt.Println(ex.Localize(language.MustParse("zh-CN")))
	fmt.Println(ex.String())
}
```

## Stack Example

```go
package main

import (
	"fmt"

	"github.com/chg1f/errorx/v2"
	_ "github.com/chg1f/errorx/v2/stack"
)

func main() {
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	fmt.Println(len(ex.Stack().Frames()) > 0)
}
```

## Package Layout

- `errorx`: core error type, builder, and helper options
- `errorx/i18n`: `go-i18n` integration and message loading from `fs.FS`
- `errorx/stack`: runtime stack capture integration
