package simple_context

import (
	"errors"
	"mime"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-amwk/core"
)

type contextImpl core.Context

type Context struct {
	contextImpl

	state    sync.Map
	index    int
	isAbort  bool
	handlers []core.HandlerFunc
}

func InitContext(ctx *Context, impl core.Context) {
	ctx.contextImpl = impl
	ctx.index = -1
	ctx.isAbort = false
	ctx.handlers = make([]core.HandlerFunc, 0)
}

// Get returns the value associated with the key in the context.
func (ctx *Context) Get(key string) (any, bool) {
	return ctx.state.Load(key)
}

// Set sets the value for the key in the context and returns the previous value.
func (ctx *Context) Set(key string, value any) any {
	oldValue, _ := ctx.state.Swap(key, value)
	return oldValue
}

// Abort marks the context as aborted, and subsequent handlers will not be executed.
func (ctx *Context) Abort() {
	ctx.isAbort = true
}

// IsAbort checks if the context is marked as aborted.
func (ctx *Context) IsAbort() bool {
	return ctx.isAbort
}

// Next calls the next handler in the chain.
func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < len(ctx.handlers) && !ctx.isAbort {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

// Use adds handlers to the context, which will be executed in the order they are added.
func (ctx *Context) Use(handlers ...core.HandlerFunc) {
	ctx.handlers = append(ctx.handlers, handlers...)
}

// BasicAuth returns the username and password from the Basic Authentication header if present,
// or empty strings and false if not present.
func (ctx *Context) BasicAuth() (string, string, bool) {
	return ctx.Request().BasicAuth()
}

// Body returns the request body as a byte slice.
func (ctx *Context) Body() ([]byte, error) {
	return ctx.Request().Body()
}

// ClientIP returns the IP address of the client making the request.
func (ctx *Context) ClientIP() string {
	proxyIp := ctx.Header("X-Forwarded-For")
	if proxyIp != "" {
		return proxyIp
	}

	return ctx.Request().ClientIP()
}

// ContentLength returns the length of the request body in bytes.
func (ctx *Context) ContentLength() int64 {
	return ctx.Request().ContentLength()
}

// ContentType returns the Content-Type header of the request.
func (ctx *Context) ContentType() string {
	contentType, _, _ := mime.ParseMediaType(ctx.Header("Content-Type"))
	return contentType
}

// Cookie retrieves a cookie by name from the request.
func (ctx *Context) Cookie(name string) (*http.Cookie, error) {
	return ctx.Request().Cookie(name)
}

// Cookies returns all cookies from the request.
func (ctx *Context) Cookies() []*http.Cookie {
	return ctx.Request().Cookies()
}

// Header retrieves a header value by name from the request.
func (ctx *Context) Header(key string) string {
	return ctx.Request().Header(key)
}

// HeaderValues retrieves all values for a header by name from the request.
func (ctx *Context) HeaderValues(key string) []string {
	return ctx.Request().HeaderValues(key)
}

// Headers returns all headers from the request.
func (ctx *Context) Headers() http.Header {
	return ctx.Request().Headers()
}

// Method returns the HTTP method of the request (e.g., GET, POST).
func (ctx *Context) Method() string {
	return ctx.Request().Method()
}

// Protocol returns the HTTP protocol version of the request (e.g., HTTP/1.1).
func (ctx *Context) Protocol() string {
	return ctx.Request().Protocol()
}

// Path returns the request path.
func (ctx *Context) Path() string {
	return ctx.Request().Path()
}

// PathValue retrieves a path parameter value by name from the request.
func (ctx *Context) PathValue(name string) string {
	return ctx.Request().PathValue(name)
}

// Resource returns the resource pattern of the request.
func (ctx *Context) Resource() string {
	return ctx.Request().Resource()
}

// Query retrieves a query parameter value by name from the request.
func (ctx *Context) Query(key string) string {
	return ctx.Request().Queries().Get(key)
}

// QueryValues retrieves all values for a query parameter by name from the request.
func (ctx *Context) QueryValues(key string) []string {
	return ctx.Request().Queries()[key]
}

// Queries returns all query parameters from the request as a url.Values.
func (ctx *Context) Queries() url.Values {
	return ctx.Request().Queries()
}

// AddHeader adds a header to the response.
func (ctx *Context) AddHeader(key, value string) {
	ctx.Response().AddHeader(key, value)
}

// SetHeader sets a header in the response.
func (ctx *Context) SetHeader(key, value string) {
	ctx.Response().SetHeader(key, value)
}

// GetHeader retrieves a header value by name from the response.
func (ctx *Context) GetHeader(key string) string {
	return ctx.Response().GetHeader(key)
}

func (ctx *Context) DelHeader(key string) {
	ctx.Response().DelHeader(key)
}

// Status sets the HTTP status code for the response and returns an error if it fails.
func (ctx *Context) Status(code int) error {
	if code < 100 || code > 999 {
		return errors.New("invalid status code")
	}

	return ctx.Response().Status(code)
}

// Write writes data to the response body.
func (ctx *Context) Write(data []byte) (int, error) {
	return ctx.Response().Write(data)
}
