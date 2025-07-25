# Goyard Web Framework

**Version: 1.0.0**

Goyard is a modern, lightweight Go web framework built on top of Echo v4, designed for rapid development of server-side rendered applications with HTMX support, automatic validation, and seamless templ template integration.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Routing & Controllers](#routing--controllers)
- [Request Handling & DTOs](#request-handling--dtos)
- [Validation](#validation)
- [Template Rendering](#template-rendering)
- [HTMX Integration](#htmx-integration)
- [Middleware](#middleware)
- [File Uploads](#file-uploads)
- [Utilities](#utilities)
- [Examples](#examples)
- [API Reference](#api-reference)

## Features

- 🚀 **Built on Echo v4** - Leverages Echo's performance and ecosystem
- 🎯 **Type-Safe Controllers** - Automatic DTO binding and validation
- 📝 **Templ Integration** - First-class support for Go's templ templating
- ⚡ **HTMX Ready** - Built-in HTMX helpers and smart rendering
- 🔒 **Automatic Validation** - Built-in validation with go-playground/validator
- 📁 **File Upload Support** - Secure file handling with SHA-256 hashing
- 🍪 **Cookie Helpers** - Simplified cookie management
- 📊 **Pagination Support** - Built-in pagination utilities
- 🔧 **Utility Pipes** - Data transformation and conversion helpers
- 🎨 **Frontend Integration** - Seamless integration with Tailwind CSS, Alpine.js

## Quick Start

### Installation

```bash
go mod init your-app
go get github.com/UchaBokeria/goyard
```

### Basic Application

```go
package main

import (
    "log"
    "github.com/UchaBokeria/goyard/goyard"
    "github.com/UchaBokeria/goyard/controller"
)

func main() {
    app := goyard.New()
    
    // Basic route
    app.GET("/", controller.Use(func(ctx *controller.Context[any]) error {
        return ctx.String(200, "Hello, Goyard!")
    }))
    
    log.Fatal(app.Run(":3000"))
}
```

### With Templates

```go
package main

import (
    "log"
    "github.com/UchaBokeria/goyard/goyard"
    "github.com/UchaBokeria/goyard/controller"
)

// Define your templ component
templ HomePage() {
    <html>
        <body>
            <h1>Welcome to Goyard!</h1>
        </body>
    </html>
}

func home(ctx *controller.Context[any]) error {
    return ctx.Html(HomePage())
}

func main() {
    app := goyard.New()
    app.GET("/", controller.Use(home))
    log.Fatal(app.Run(":3000"))
}
```

## Core Concepts

### Framework Structure

Goyard extends Echo with a custom context that provides additional functionality:

```go
type Context[T any] struct {
    echo.Context
    data map[string]T
}
```

### Application Initialization

```go
// Create new Goyard application
app := goyard.New()

// The app wraps Echo and provides:
// - Extended context with helper methods
// - Automatic validation middleware
// - HTMX-aware rendering
// - Custom run function with port adjustment
```

## Routing & Controllers

### Basic Routing

```go
app := goyard.New()

// Simple handlers
app.GET("/", controller.Use(indexHandler))
app.POST("/users", controller.Use(createUserHandler))

// Route groups
api := app.Group("/api")
api.GET("/users", controller.Use(getUsersHandler))
```

### Type-Safe Controllers with DTOs

```go
type CreateUserDTO struct {
    Email    string `json:"email" form:"email" validate:"required,email"`
    Password string `json:"password" form:"password" validate:"required,strongpwd"`
    Age      int    `json:"age" form:"age" validate:"min=18"`
}

// Handler automatically receives validated DTO
func createUser(ctx *controller.Context[any], dto *CreateUserDTO) error {
    // DTO is automatically bound and validated
    // Process user creation
    return ctx.JSON(200, map[string]string{"status": "created"})
}

// Register with automatic binding
app.POST("/users", controller.Set[CreateUserDTO](createUser))
```

### Context Methods

```go
func handler(ctx *controller.Context[any]) error {
    // Basic responses
    ctx.String(200, "Hello")
    ctx.JSON(200, data)
    ctx.HTML(200, "<h1>HTML</h1>")
    
    // Template rendering (HTMX-aware)
    ctx.Html(component)
    ctx.HtmlWithStatus(201, component)
    ctx.Renders(200, component) // Always full render
    
    // HTMX detection
    if ctx.IsHtmx() {
        // Handle HTMX request
    }
    
    // Pagination helpers
    page := ctx.Page()         // Gets page number from query
    pageSize := ctx.PageSize() // Gets page size from query
    
    // Data storage
    ctx.Set("key", value)
    user := ctx.User() // Gets stored user data
    
    return nil
}
```

## Request Handling & DTOs

### DTO Structure

```go
type UserDTO struct {
    // JSON and form binding
    Name     string `json:"name" form:"name" validate:"required,min=2"`
    Email    string `json:"email" form:"email" validate:"required,email"`
    Password string `json:"password" form:"password" validate:"required,strongpwd"`
    Age      *int   `json:"age" form:"age" validate:"omitempty,min=18"`
    
    // Query parameters
    Page     string `query:"page"`
    PageSize string `query:"pageSize"`
}
```

### Automatic Binding Process

1. **Binding**: Request data (JSON/form) automatically bound to DTO
2. **Validation**: Built-in validation using struct tags
3. **Error Handling**: Automatic error responses for binding/validation failures
4. **Handler Invocation**: Clean DTO passed to handler function

```go
// This handler automatically gets validated DTO
func updateUser(ctx *controller.Context[any], dto *UserDTO) error {
    // dto is guaranteed to be valid
    // Handle business logic
    return ctx.JSON(200, dto)
}

app.PUT("/users/:id", controller.Set[UserDTO](updateUser))
```

## Validation

### Built-in Validators

Goyard uses go-playground/validator with additional custom validators:

```go
type UserDTO struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,strongpwd"` // Custom strong password validator
    Age      int    `validate:"min=18,max=120"`
    Name     string `validate:"required,min=2,max=50"`
}
```

### Strong Password Validator

The `strongpwd` validator ensures passwords meet security requirements:
- Minimum 8 characters
- Contains lowercase letter
- Contains uppercase letter  
- Contains number
- Contains special character

```go
Password string `validate:"required,strongpwd"`
```

### Manual Validation

```go
import "github.com/UchaBokeria/goyard/controller"

func someHandler(ctx *controller.Context[any]) error {
    var dto UserDTO
    // ... populate dto ...
    
    if err := controller.Validate(dto); err != nil {
        return ctx.String(400, err.Error())
    }
    
    return nil
}
```

## Template Rendering

### Templ Integration

Goyard has first-class support for [templ](https://github.com/a-h/templ):

```go
// Define templ component
templ UserCard(user User) {
    <div class="user-card">
        <h3>{ user.Name }</h3>
        <p>{ user.Email }</p>
    </div>
}

// Render in handler
func showUser(ctx *controller.Context[any]) error {
    user := getUserFromDB()
    return ctx.Html(UserCard(user))
}
```

### HTMX-Aware Rendering

```go
func userProfile(ctx *controller.Context[any]) error {
    user := getCurrentUser()
    
    // Automatically detects HTMX requests
    // HTMX requests: renders component only
    // Regular requests: renders full page with layout
    return ctx.Html(UserProfile(user))
}
```

### Manual Render Control

```go
func handler(ctx *controller.Context[any]) error {
    // Always render full component (ignores HTMX)
    return ctx.Renders(200, component)
    
    // Manual HTMX check
    if ctx.IsHtmx() {
        return ctx.Html(PartialComponent())
    }
    return ctx.Html(FullPageComponent())
}
```

## HTMX Integration

### HTMX Helper Functions

```go
import "github.com/UchaBokeria/goyard/htmx"

// Generate HTMX link attributes
attrs := htmx.Link("/users/123")
// Results in:
// hx-get="/users/123"
// hx-swap="innerHTML show:window:top"  
// hx-push-url="true"
// hx-target="#Content"
// hx-encoding="text/html"
// hx-indicator=".Loading"
// hx-trigger="click"

// POST link with parameters
postAttrs := htmx.PostLink("/users", "name,email")
```

### In Templ Templates

```go
templ UserLink(userID string) {
    <a { htmx.Link("/users/" + userID)... }>
        View User
    </a>
}

templ UserForm() {
    <form { htmx.PostLink("/users", "name,email")... }>
        <input name="name" type="text"/>
        <input name="email" type="email"/>
        <button type="submit">Create</button>
    </form>
}
```

### HTMX Request Detection

```go
func handler(ctx *controller.Context[any]) error {
    if ctx.IsHtmx() {
        // Return partial content for HTMX
        return ctx.Html(UserListPartial())
    }
    
    // Return full page for regular requests
    return ctx.Html(UserListPage())
}
```

## Middleware

### Built-in Middleware

```go
import "github.com/UchaBokeria/goyard/controller"

app := goyard.New()

// Initialize middleware (automatically included)
app.Use(controller.Initialize())

// Custom middleware
app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Middleware logic
        return next(c)
    }
})
```

### Example Middleware

```go
// HTMX middleware (from examples)
func HtmxMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if c.Request().Header.Get("HX-Request") == "true" {
                // HTMX-specific logic
            }
            return next(c)
        }
    }
}

// Interceptor middleware
func InterceptorMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return controller.Set[any](func(ctx *controller.Context[any]) error {
            // Pre-processing logic
            return next(ctx)
        })
    }
}
```

## File Uploads

### Upload Configuration

```go
import "github.com/UchaBokeria/goyard/upload"

// Configure upload directory (default: "./public/uploads/")
upload.Dir = "./custom/upload/path/"
```

### Handling File Uploads

```go
func uploadHandler(ctx *controller.Context[any]) error {
    // Get file from form
    file, err := ctx.FormFile("file")
    if err != nil {
        return ctx.String(400, "No file uploaded")
    }
    
    // Save file (automatically hashed with SHA-256)
    response := upload.Save(file)
    
    if !response.Success {
        return ctx.String(500, response.Message)
    }
    
    return ctx.JSON(200, response)
}

app.POST("/upload", controller.Use(uploadHandler))
```

### Upload Response

```go
type Response struct {
    ID      int    // File ID (0 for success, -1 for error)
    Message string // Status message
    Success bool   // Upload success status
}
```

### File Naming

Files are automatically renamed using SHA-256 hash of content plus original extension:
- `image.jpg` → `a1b2c3d4e5f6...xyz.jpg`
- Prevents naming conflicts
- Enables deduplication
- Maintains file integrity

## Utilities

### Common Utilities

```go
import "github.com/UchaBokeria/goyard/common"

// Pointer utilities
ptr := common.ToPointer("value")

// Random token generation
token := common.GenerateToken(32) // 32 character token

// Pretty print JSON
common.Print(data) // Prints formatted JSON to stdout

// Convert slice to select options
users := []User{{ID: 1, Name: "John"}, {ID: 2, Name: "Jane"}}
options, err := common.ToSelect(users, "ID", "Name")
// Returns: []SelectDataType{{Key: "1", Val: "John"}, {Key: "2", Val: "Jane"}}
```

### Data Conversion Pipes

```go
import "github.com/UchaBokeria/goyard/pipes"

// String/number conversions
num := pipes.Int("123")        // string to int
str := pipes.String(123)       // int to string
ustr := pipes.Ustring(uint(42)) // uint to string
bstr := pipes.Bstring(true)    // bool to string ("true"/"false")

// JSON conversion
jsonStr := pipes.Json(data)    // struct to JSON string

// Alpine.js data binding
alpineData := pipes.Alpine(data) // struct to Alpine.js format

// Pointer utilities
intPtr := pipes.Pint(42)       // int to *int
uintPtr := pipes.Puint(42)     // int to *uint
str := pipes.PtIntToString(intPtr) // *int to string (safe)

// Type conversions
u := pipes.Uint("42")          // string to uint
uintPtr := pipes.UintToPtInt(u) // uint to *int
```

### Cookie Management

```go
func handler(ctx *controller.Context[any]) error {
    // Write cookie
    ctx.WriteCookie(controller.Cookie{
        Key:     "session",
        Value:   "abc123",
        Expires: time.Now().Add(24 * time.Hour),
    })
    
    // Read cookie
    cookie := ctx.ReadCookie("session")
    if cookie.Value != "" {
        // Cookie exists
    }
    
    // Remove cookie
    ctx.RemoveCookie("session")
    
    return nil
}
```

### Pagination

```go
func listUsers(ctx *controller.Context[any]) error {
    page := ctx.Page()         // Gets 'page' query param (default: 1)
    pageSize := ctx.PageSize() // Gets 'pageSize' query param (default: 50, max: 50)
    
    offset := (page - 1) * pageSize
    users := getUsersFromDB(offset, pageSize)
    
    return ctx.JSON(200, users)
}

// Usage: GET /users?page=2&pageSize=10
```

## Examples

### Complete CRUD Application

```go
package main

import (
    "log"
    "github.com/UchaBokeria/goyard/goyard"
    "github.com/UchaBokeria/goyard/controller"
)

type User struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

type CreateUserDTO struct {
    Name  string `json:"name" form:"name" validate:"required,min=2"`
    Email string `json:"email" form:"email" validate:"required,email"`
}

type UpdateUserDTO struct {
    Name  string `json:"name" form:"name" validate:"omitempty,min=2"`
    Email string `json:"email" form:"email" validate:"omitempty,email"`
}

// Handlers
func listUsers(ctx *controller.Context[any]) error {
    page := ctx.Page()
    pageSize := ctx.PageSize()
    
    users := getUsersFromDB(page, pageSize)
    return ctx.Json(200, users)
}

func createUser(ctx *controller.Context[any], dto *CreateUserDTO) error {
    user := User{
        Name:  dto.Name,
        Email: dto.Email,
    }
    
    savedUser := saveUserToDB(user)
    return ctx.JSON(201, savedUser)
}

func updateUser(ctx *controller.Context[any], dto *UpdateUserDTO) error {
    userID := ctx.Param("id")
    // Update logic here
    return ctx.JSON(200, updatedUser)
}

func deleteUser(ctx *controller.Context[any]) error {
    userID := ctx.Param("id")
    // Delete logic here
    return ctx.NoContent(204)
}

func main() {
    app := goyard.New()
    
    // CRUD routes
    users := app.Group("/users")
    users.GET("", controller.Use(listUsers))
    users.POST("", controller.Set[CreateUserDTO](createUser))
    users.PUT("/:id", controller.Set[UpdateUserDTO](updateUser))
    users.DELETE("/:id", controller.Use(deleteUser))
    
    log.Fatal(app.Run(":3000"))
}
```

### HTMX-Powered Frontend

```go
// main.go
func main() {
    app := goyard.New()
    
    // Serve static files
    app.Static("/static", "public")
    
    // Routes
    app.GET("/", controller.Use(homePage))
    app.GET("/users", controller.Use(usersList))
    app.POST("/users", controller.Set[CreateUserDTO](createUser))
    
    log.Fatal(app.Run(":3000"))
}

func homePage(ctx *controller.Context[any]) error {
    return ctx.Html(HomePage())
}

func usersList(ctx *controller.Context[any]) error {
    users := getAllUsers()
    
    if ctx.IsHtmx() {
        return ctx.Html(UsersListPartial(users))
    }
    
    return ctx.Html(UsersPage(users))
}

func createUser(ctx *controller.Context[any], dto *CreateUserDTO) error {
    user := createNewUser(dto)
    
    // Return the new user component for HTMX to inject
    return ctx.Html(UserCard(user))
}
```

```go
// templates.templ
package main

import "github.com/UchaBokeria/goyard/htmx"

templ HomePage() {
    <html>
        <head>
            <title>HTMX App</title>
            <script src="https://unpkg.com/htmx.org@1.9.6"></script>
        </head>
        <body>
            <h1>Users Management</h1>
            <div id="users-container" { htmx.Link("/users")... }>
                Loading users...
            </div>
            
            <form { htmx.PostLink("/users", "name,email")... } 
                  hx-target="#users-container" 
                  hx-swap="beforeend">
                <input name="name" placeholder="Name" required/>
                <input name="email" placeholder="Email" type="email" required/>
                <button type="submit">Add User</button>
            </form>
        </body>
    </html>
}

templ UserCard(user User) {
    <div class="user-card">
        <h3>{ user.Name }</h3>
        <p>{ user.Email }</p>
    </div>
}

templ UsersListPartial(users []User) {
    for _, user := range users {
        @UserCard(user)
    }
}
```

## API Reference

### Core Types

```go
// Goyard application
type Goyard struct {
    *echo.Echo
    Run func(address string) error
}

// Extended context
type Context[T any] struct {
    echo.Context
    data map[string]T
}

// Upload response
type Response struct {
    ID      int
    Message string
    Success bool
}

// Cookie structure
type Cookie struct {
    Key     string
    Value   string
    Expires time.Time
}

// Select data for dropdowns
type SelectDataType struct {
    Key string `json:"key"`
    Val string `json:"val"`
}
```

### Framework Functions

```go
// Create new application
goyard.New() *types.Goyard

// Controller utilities
controller.Use[T any](handler func(*Context[any]) error) echo.HandlerFunc
controller.Set[T any](handler interface{}) echo.HandlerFunc
controller.Initialize() echo.MiddlewareFunc
controller.Validate(dto interface{}) error

// HTMX helpers
htmx.Link(path string) templ.Attributes
htmx.PostLink(path string, params string) templ.Attributes

// Upload utilities
upload.Save(file *multipart.FileHeader) *Response

// Common utilities
common.ToPointer[T any](data T) *T
common.GenerateToken(n int) string
common.Print(data interface{}) error
common.ToSelect[T any](data []T, keyField, valField string) ([]SelectDataType, error)

// Pipes (data conversion)
pipes.Int(s string) int
pipes.String(i int) string
pipes.Json(data any) string
pipes.Alpine(data any) string
pipes.Pint(num int) *int
pipes.Puint(num int) *uint
```

### Context Methods

```go
// Template rendering
ctx.Html(c templ.Component) error
ctx.HtmlWithStatus(code int, c templ.Component) error
ctx.Renders(code int, c templ.Component) error

// HTMX detection
ctx.IsHtmx() bool

// Cookie management
ctx.WriteCookie(data Cookie)
ctx.ReadCookie(key string) Cookie
ctx.RemoveCookie(key string)

// Pagination
ctx.Page() int
ctx.PageSize() int

// Data storage
ctx.Set(key string, value T)
ctx.User() T
```

---

**Goyard** - Build modern web applications with Go, HTMX, and templ. Simple, fast, and productive.
