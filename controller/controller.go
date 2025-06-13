// Package controller provides utilities for managing route handlers in an Echo framework application.
// It adds a richer context wrapper and helpers for validation, HTML rendering and HTMX-aware responses.
package controller

import (
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// defaultPageMaxSize defines an upper bound for page sizes when no custom configuration is supplied.
const defaultPageMaxSize = 50

// Context wraps echo.Context and exposes additional helper methods.
type Context[T any] struct {
	echo.Context
	data map[string]T
}

// Initialize returns a middleware that swaps echo.Context with our extended Context.
func Initialize() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &Context[any]{Context: c}
			return next(ctx)
		}
	}
}

// Use converts a handler that expects *controller.Context into a standard echo.HandlerFunc.
func Use[T any](handler func(*Context[any]) error) echo.HandlerFunc {
	return func(c echo.Context) error { return handler(c.(*Context[any])) }
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

		args := []reflect.Value{reflect.ValueOf(c.(*Context[any]))}
		if reflect.TypeOf(dto).String() != "*interface {}" {
			argVal := reflect.New(reflect.TypeOf(dto)).Elem()
			argVal.Set(reflect.ValueOf(dto))
			args = append(args, argVal)
		}

		result := reflect.ValueOf(handler).Call(args)
		if result[0].IsNil() {
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
func (ctx *Context[T]) Html(c templ.Component) error {
	return ctx.HtmlWithStatus(http.StatusOK, c)
}

// HtmlWithStatus renders a component and sends it with the provided HTTP status code.
// If the request is an HTMX request we render only the fragment.
func (ctx *Context[T]) HtmlWithStatus(code int, c templ.Component) error {
	if ctx.IsHtmx() {
		return c.Render(ctx.Request().Context(), ctx.Response())
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Writer.WriteHeader(code)
	return c.Render(ctx.Request().Context(), ctx.Response().Writer)
}

// Renders behaves like HtmlWithStatus but without HTMX logic – always renders raw component.
func (ctx *Context[T]) Renders(code int, c templ.Component) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Writer.WriteHeader(code)
	return c.Render(ctx.Request().Context(), ctx.Response().Writer)
}

// IsHtmx reports whether the incoming request is an HTMX request.
func (ctx *Context[T]) IsHtmx() bool {
	return ctx.Request().Header.Get("Hx-Request") == "true" && ctx.Request().Header.Get("hx-fullPage") != "true"
}

// Cookie utilities -----------------------------------------------------------

type Cookie struct {
	Key     string
	Value   string
	Expires time.Time
}

func (ctx *Context[T]) RemoveCookie(key string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(0, 0)
	ctx.SetCookie(cookie)
}

func (ctx *Context[T]) WriteCookie(data Cookie) {
	cookie := new(http.Cookie)
	cookie.Name = data.Key
	cookie.Value = data.Value
	cookie.Expires = data.Expires
	ctx.SetCookie(cookie)
}

func (ctx *Context[T]) ReadCookie(key string) Cookie {
	cookie, err := ctx.Cookie(key)
	if err != nil {
		cookie = &http.Cookie{Name: "", Value: "", Expires: time.Now()}
	}
	return Cookie{Key: cookie.Name, Value: cookie.Value, Expires: cookie.Expires}
}

// Pagination helpers ---------------------------------------------------------

type QueryPageParameter struct {
	Page     string `query:"page"`
	PageSize string `query:"pageSize"`
}

func (ctx *Context[T]) Page() int {
	var q QueryPageParameter
	if ctx.QueryParam("page") == "" {
		q.Page = "1"
	} else {
		ctx.Bind(&q)
	}
	p, _ := strconv.Atoi(q.Page)
	if p <= 0 {
		p = 1
	}
	return p
}

func (ctx *Context[T]) PageSize() int {
	var q QueryPageParameter
	if ctx.QueryParam("pageSize") == "" {
		q.PageSize = "-1"
	}
	ctx.Bind(&q)
	size, _ := strconv.Atoi(q.PageSize)
	if size <= 0 || size > defaultPageMaxSize {
		size = defaultPageMaxSize
	}
	return size
}

func (ctx *Context[T]) Set(key string, value T) {
	ctx.data[key] = value
}
func (ctx *Context[T]) User() T {
	return ctx.data["USER"]
}
