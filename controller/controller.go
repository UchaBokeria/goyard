// Package controller provides utilities for managing route handlers in an Echo framework application.
// It adds a richer context wrapper and helpers for validation, HTML rendering and HTMX-aware responses.
package controller

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Context wraps echo.Context and exposes additional helper methods.
type Context struct {
	echo.Context
}

func Data[T any](ctx *Context, key string) T {
	return ctx.Get(key).(T)
}

// Initialize returns a middleware that swaps echo.Context with our extended Context.
func Initialize() echo.MiddlewareFunc {
	fmt.Println("Goyard Initialized")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &Context{Context: c}
			return next(ctx)
		}
	}
}

// Use converts a handler that expects *controller.Context into a standard echo.HandlerFunc.
func Use(handler func(*Context) error) echo.HandlerFunc {
	return func(c echo.Context) error { return handler(c.(*Context)) }
}

// Set automatically binds request data into a DTO, validates it, and then forwards
// the call to the provided handler via reflection (so the handler can accept the DTO directly).
func Set[T any](handler interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := new(T)
		if err := c.Bind(dto); err != nil {
			return c.String(http.StatusBadRequest, "Parameters Binding Problem: "+err.Error())
		}

		if reflect.TypeOf(dto).String() != "*interface {}" {
			if err := Validate(dto); err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
		}

		// Injection logic
		rv := reflect.ValueOf(dto)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem() // Dereference to the struct value
		}
		if rv.Kind() == reflect.Struct {
			t := rv.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				injectKey := field.Tag.Get("inject")
				if injectKey != "" {
					xval := c.Get(injectKey) // From context, already interface{}
					if xval == nil {
						// Handle missing value; e.g., skip or error
						return c.String(http.StatusBadRequest, "missing inject value for key: "+injectKey)
					}

					fv := rv.Field(i)
					if !fv.IsValid() || !fv.CanSet() {
						// Skip unexported or invalid fields
						continue
					}

					valRv := reflect.ValueOf(xval)
					fieldType := fv.Type()

					// Check assignability
					if !valRv.Type().AssignableTo(fieldType) {
						// Attempt conversion if possible
						if valRv.Type().ConvertibleTo(fieldType) {
							valRv = valRv.Convert(fieldType)
						} else {
							return c.String(http.StatusBadRequest, "inject value for '"+injectKey+"' (type "+reflect.TypeOf(xval).String()+") not assignable to field '"+field.Name+"' (type "+fieldType.String()+")")
						}
					}

					// Set the field
					fv.Set(valRv)
				}
			}
		}

		args := []reflect.Value{reflect.ValueOf(c.(*Context))}
		if reflect.TypeOf(dto).String() != "*interface {}" {
			args = append(args, reflect.ValueOf(dto))
		}

		result := reflect.ValueOf(handler).Call(args)
		if len(result) == 0 || result[0].IsNil() {
			return nil
		}
		return result[0].Interface().(error)
	}
}

// strongpwd is a custom validator ensuring a field matches a strong password policy.
func strongpwd(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	checks := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, pattern := range checks {
		matched, _ := regexp.MatchString(pattern, field)
		if !matched {
			return false
		}
	}
	return true
}

// Validate runs validation on any struct using go-playground/validator.
func Validate(dto interface{}) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("strongpwd", strongpwd)
	return v.Struct(dto)
}

// Html renders the given templ component with status 200.
func (ctx *Context) Html(c templ.Component) error {
	return ctx.HtmlWithStatus(http.StatusOK, c)
}

// HtmlWithStatus renders a component and sends it with the provided HTTP status code.
// If the request is an HTMX request we render only the fragment.
func (ctx *Context) HtmlWithStatus(code int, c templ.Component) error {
	if ctx.IsHtmx() {
		return c.Render(ctx.Request().Context(), ctx.Response())
	}

	var base templ.Component
	wrapper, ok := ctx.Get("LayoutRenderNoHtmx").(func(c templ.Component) templ.Component)
	if wrapper != nil && ok {
		base = wrapper(c)
	} else {
		base = c
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Writer.WriteHeader(code)
	return base.Render(ctx.Request().Context(), ctx.Response().Writer)
}

// Simple page renderer
func (ctx *Context) Page(page templ.Component) any {
	return Set[any](func(ctx *Context) error {
		return ctx.Html(page)
	})
}

// Renders behaves like HtmlWithStatus but without HTMX logic – always renders raw component.
func (ctx *Context) Renders(code int, c templ.Component) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Writer.WriteHeader(code)
	return c.Render(ctx.Request().Context(), ctx.Response().Writer)
}

// IsHtmx reports whether the incoming request is an HTMX request.
func (ctx *Context) IsHtmx() bool {
	return ctx.Request().Header.Get("Hx-Request") == "true" && ctx.Request().Header.Get("hx-fullPage") != "true"
}

// Cookie utilities -----------------------------------------------------------
type Cookie struct {
	Key      string     // required
	Value    string     // required
	Expires  *time.Time // optional
	Path     *string    // optional
	MaxAge   *int       // optional
	HttpOnly *bool      // optional
	SameSite *string    // optional (map to http.SameSite later)
	Domain   *string    // optional
	Quoted   *bool      // optional
}

func (ctx *Context) RemoveCookie(key string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(0, 0)

	ctx.SetCookie(cookie)
}

func Ptr[T any](v T) *T { return &v }

func quoteIf(shouldQuote bool, val string) string {
	if shouldQuote {
		return `"` + val + `"`
	}
	return val
}

func (ctx *Context) WriteCookie(data Cookie) {
	cookie := new(http.Cookie)
	cookie.Name = data.Key
	cookie.Value = data.Value

	if data.Expires != nil {
		cookie.Expires = *data.Expires
	}
	if data.Path != nil {
		cookie.Path = *data.Path
	}
	if data.MaxAge != nil {
		cookie.MaxAge = *data.MaxAge
	}
	if data.HttpOnly != nil {
		cookie.HttpOnly = *data.HttpOnly
	}
	if data.SameSite != nil {
		switch *data.SameSite {
		case "Strict":
			cookie.SameSite = http.SameSiteStrictMode
		case "Lax":
			cookie.SameSite = http.SameSiteLaxMode
		case "None":
			cookie.SameSite = http.SameSiteNoneMode
		default:
			// invalid values fallback to default (unset)
			cookie.SameSite = http.SameSiteStrictMode
		}
	}
	if data.Domain != nil {
		cookie.Domain = *data.Domain
	}
	if data.Quoted != nil {
		cookie.Raw = quoteIf(*data.Quoted, data.Value)
	}

	ctx.SetCookie(cookie)
}

func (ctx *Context) ReadCookie(key string) Cookie {
	cookie, err := ctx.Cookie(key)
	if err != nil {
		cookie = &http.Cookie{Name: "", Value: "", Expires: time.Now()}
	}
	return Cookie{Key: cookie.Name, Value: cookie.Value, Expires: Ptr(cookie.Expires)}
}
