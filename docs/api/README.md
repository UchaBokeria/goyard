# GoYard API Reference

## Core Package

### `goyard.New()`

Creates a new GoYard application instance.

```go
app := goyard.New()
```

### `goyard.Version`

A constant containing the current version of the framework.

```go
fmt.Printf("Running GoYard %s", goyard.Version)
```

## Package: pkg

The `pkg` package contains the core utilities and components of the framework.

### Controller

The Controller provides enhanced request handling with an extended context.

```go
import "github.com/yourorg/goyard/pkg"

// Use controller in your handler
handler := pkg.Use(func(ctx *pkg.Context) error {
    return ctx.String(200, "Hello World")
})
```

### HTML Rendering

Render HTML using Templ components.

```go
func handler(ctx *pkg.Context) error {
    return ctx.Html(view.HomePage())
}
```

## Examples

See the [examples directory](../../examples) for complete working examples. 