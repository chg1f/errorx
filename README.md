errorx
---

`errorx` provides typed Go errors with a mandatory code, optional locale lookup, and optional stack capture.

## Features

- Every `*errorx.Error[T]` carries a comparable code.
- Builders stay small and composable.
- `i18n` is opt-in through `errorx/i18n`.
- `stack` is opt-in through `errorx/stack`.
- Standard `errors.Is` and `errors.As` continue to work.

## Install

```bash
go get github.com/chg1f/errorx/v2
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/chg1f/errorx/v2"
	"github.com/chg1f/errorx/v2/i18n"
	_ "github.com/chg1f/errorx/v2/stack"
	"golang.org/x/text/language"
)

func main() {
	if err := i18n.LoadFiles(os.DirFS("locales")); err != nil {
		panic(err)
	}

	err = errorx.WithCode("invalid").
		WithLocale("invalid").
		WithValues(map[string]any{"field": "email"}).
		New("email is invalid")

	ex := errorx.Be[string](err)
	fmt.Println(ex.Code())
	fmt.Println(ex.String())
	fmt.Println(ex.Localize(language.MustParse("zh-CN")))
	fmt.Println(len(ex.Stack().Frames()) > 0)
}
```

## Package Layout

- `errorx`: core error type, builder, and helper options.
- `errorx/i18n`: locale bundle loading from `fs.FS`, backed by `go-i18n`.
- `errorx/stack`: runtime stack capture integration.

## GoDoc Notes

- Import the module as `github.com/chg1f/errorx/v2`.
- See exported builders such as `WithCode`, `WithMessage`, and `WithLocale`.
- See helper options in `Be` such as `Empty`.
