
# martini
    import "github.com/codegangsta/martini"

Package martini is a powerful package for quickly writing modular web applications/services in Golang.

For a full guide visit a href="http://github.com/codegangsta/martini">http://github.com/codegangsta/martini</a>


	package main
	
	import "github.com/codegangsta/martini"
	
	func main() {
	  m := martini.Classic()
	
	  m.Get("/", func() string {
	    return "Hello world!"
	  })
	
	  m.Run()
	}




## Constants
``` go
const (
    Dev  string = "development"
    Prod string = "production"
    Test string = "test"
)
```

Envs


## Variables

<pre>var Env = Dev</pre>
Env is the environment that Martini is executing in. The MARTINI_ENV is read on initialization to set this variable.





## type BeforeFunc
<pre>type BeforeFunc func(ResponseWriter)</pre>
BeforeFunc is a function that is called before the ResponseWriter has been written to.












## type ClassicMartini
<pre>type ClassicMartini struct {
    *Martini
    Router
}</pre>
ClassicMartini represents a Martini with some reasonable defaults. Embeds the router functions for convenience.









### func Classic

    func Classic() *ClassicMartini

Classic creates a classic Martini with some basic default middleware - martini.Logger, martini.Recovery, and martini.Static.





## type Context
<pre>type Context interface {
    inject.Injector
    // Next is an optional function that Middleware Handlers can call to yield the until after
    // the other Handlers have been executed. This works really well for any operations that must
    // happen after an http request
    Next()
    // contains filtered or unexported methods
}</pre>
Context represents a request context. Services can be mapped on the request level from this interface.












## type Handler
<pre>type Handler interface{}</pre>
Handler can be any callable function. Martini attempts to inject services into the handler's argument list.
Martini will panic if an argument could not be fullfilled via dependency injection.









### func Logger

    func Logger() Handler

Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.


### func Recovery

    func Recovery() Handler

Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.


### func Static

    func Static(directory string) Handler

Static returns a middleware handler that serves static files in the given directory.





## type Martini
<pre>type Martini struct {
    inject.Injector
    // contains filtered or unexported fields
}</pre>
Martini represents the top level web application. inject.Injector methods can be invoked to map services on a global level.









### func New

    func New() *Martini

New creates a bare bones Martini instance. Use this method if you want to have full control over the middleware that is used.




### func (\*Martini) Action

    func (m *Martini) Action(handler Handler)

Action sets the handler that will be called after all the middleware has been invoked. This is set to martini.Router in a martini.Classic().



### func (\*Martini) Handlers

    func (m *Martini) Handlers(handlers ...Handler)

Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
Will panic if any of the handlers is not a callable function



### func (\*Martini) Run

    func (m *Martini) Run()

Run the http server. Listening on os.GetEnv("PORT") or 3000 by default.



### func (\*Martini) ServeHTTP

    func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request)

ServeHTTP is the HTTP Entry point for a Martini instance. Useful if you want to control your own HTTP server.



### func (\*Martini) Use

    func (m *Martini) Use(handler Handler)

Use adds a middleware Handler to the stack. Will panic if the handler is not a callable func. Middleware Handlers are invoked in the order that they are added.




## type Params
<pre>type Params map[string]string</pre>
Params is a map of name/value pairs for named routes. An instance of martini.Params is available to be injected into any route handler.












## type ResponseWriter
<pre>type ResponseWriter interface {
    http.ResponseWriter
    // Status returns the status code of the response or 0 if the response has not been written.
    Status() int
    // Written returns whether or not the ResponseWriter has been written.
    Written() bool
    // Size returns the size of the response body.
    Size() int
    // Before allows for a function to be called before the ResponseWriter has been written to. This is
    // useful for setting headers or any other operations that must happen before a response has been written.
    Before(BeforeFunc)
}</pre>
ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
if the functionality calls for it.









### func NewResponseWriter

    func NewResponseWriter(rw http.ResponseWriter) ResponseWriter

NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter





## type ReturnHandler
<pre>type ReturnHandler func(http.ResponseWriter, []reflect.Value)</pre>
ReturnHandler is a service that Martini provides that is called
when a route handler returns something. The ReturnHandler is
responsible for writing to the ResponseWriter based on the values
that are passed into this function.












## type Route
<pre>type Route interface {
    // URLWith returns a rendering of the Route's url with the given string params.
    URLWith([]string) string
}</pre>
Route is an interface representing a Route in Martini's routing layer.












## type Router
<pre>type Router interface {
    // Get adds a route for a HTTP GET request to the specified matching pattern.
    Get(string, ...Handler) Route
    // Patch adds a route for a HTTP PATCH request to the specified matching pattern.
    Patch(string, ...Handler) Route
    // Post adds a route for a HTTP POST request to the specified matching pattern.
    Post(string, ...Handler) Route
    // Put adds a route for a HTTP PUT request to the specified matching pattern.
    Put(string, ...Handler) Route
    // Delete adds a route for a HTTP DELETE request to the specified matching pattern.
    Delete(string, ...Handler) Route
    // Options adds a route for a HTTP OPTIONS request to the specified matching pattern.
    Options(string, ...Handler) Route
    // Head adds a route for a HTTP HEAD request to the specified matching pattern.
    Head(string, ...Handler) Route
    // Any adds a route for any HTTP method request to the specified matching pattern.
    Any(string, ...Handler) Route

    // NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.
    NotFound(...Handler)

    // Handle is the entry point for routing. This is used as a martini.Handler
    Handle(http.ResponseWriter, *http.Request, Context)
}</pre>
Router is Martini's de-facto routing interface. Supports HTTP verbs, stacked handlers, and dependency injection.









### func NewRouter

    func NewRouter() Router

NewRouter creates a new Router instance.





## type Routes
<pre>type Routes interface {
    // URLFor returns a rendered URL for the given route. Optional params can be passed to fulfill named parameters in the route.
    URLFor(route Route, params ...interface{}) string
}</pre>
Routes is a helper service for Martini's routing layer.



















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)