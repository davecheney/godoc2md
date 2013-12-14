
	
# martini
		
		
Package martini is a powerful package for quickly writing modular web applications/services in Golang.


For a full guide visit a href="http://github.com/codegangsta/martini">http://github.com/codegangsta/martini</a>

<pre>package main

import &#34;github.com/codegangsta/martini&#34;

func main() {
  m := martini.Classic()

  m.Get(&#34;/&#34;, func() string {
    return &#34;Hello world!&#34;
  })

  m.Run()
}
</pre>

		

		


## Constants

<pre>const (
    Dev  string = &#34;development&#34;
    Prod string = &#34;production&#34;
    Test string = &#34;test&#34;
)</pre>

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
<pre>func Classic() *ClassicMartini</pre>

Classic creates a classic Martini with some basic default middleware - martini.Logger, martini.Recovery, and martini.Static.









## type Context
<pre>type Context interface {
    inject.Injector
    <span class="comment">// Next is an optional function that Middleware Handlers can call to yield the until after</span>
    <span class="comment">// the other Handlers have been executed. This works really well for any operations that must</span>
    <span class="comment">// happen after an http request</span>
    Next()
    <span class="comment">// contains filtered or unexported methods</span>
}</pre>

Context represents a request context. Services can be mapped on the request level from this interface.















## type Handler
<pre>type Handler interface{}</pre>

Handler can be any callable function. Martini attempts to inject services into the handler&#39;s argument list.
Martini will panic if an argument could not be fullfilled via dependency injection.











### func Logger
<pre>func Logger() Handler</pre>

Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.





### func Recovery
<pre>func Recovery() Handler</pre>

Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.





### func Static
<pre>func Static(directory string) Handler</pre>

Static returns a middleware handler that serves static files in the given directory.









## type Martini
<pre>type Martini struct {
    inject.Injector
    <span class="comment">// contains filtered or unexported fields</span>
}</pre>

Martini represents the top level web application. inject.Injector methods can be invoked to map services on a global level.











### func New
<pre>func New() *Martini</pre>

New creates a bare bones Martini instance. Use this method if you want to have full control over the middleware that is used.







### func (*Martini) Action
<pre>func (m *Martini) Action(handler Handler)</pre>
<p>
Action sets the handler that will be called after all the middleware has been invoked. This is set to martini.Router in a martini.Classic().
</p>





### func (*Martini) Handlers
<pre>func (m *Martini) Handlers(handlers ...Handler)</pre>
<p>
Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
Will panic if any of the handlers is not a callable function
</p>





### func (*Martini) Run
<pre>func (m *Martini) Run()</pre>
<p>
Run the http server. Listening on os.GetEnv(&#34;PORT&#34;) or 3000 by default.
</p>





### func (*Martini) ServeHTTP
<pre>func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request)</pre>
<p>
ServeHTTP is the HTTP Entry point for a Martini instance. Useful if you want to control your own HTTP server.
</p>





### func (*Martini) Use
<pre>func (m *Martini) Use(handler Handler)</pre>
<p>
Use adds a middleware Handler to the stack. Will panic if the handler is not a callable func. Middleware Handlers are invoked in the order that they are added.
</p>







## type Params
<pre>type Params map[string]string</pre>

Params is a map of name/value pairs for named routes. An instance of martini.Params is available to be injected into any route handler.















## type ResponseWriter
<pre>type ResponseWriter interface {
    http.ResponseWriter
    <span class="comment">// Status returns the status code of the response or 0 if the response has not been written.</span>
    Status() int
    <span class="comment">// Written returns whether or not the ResponseWriter has been written.</span>
    Written() bool
    <span class="comment">// Size returns the size of the response body.</span>
    Size() int
    <span class="comment">// Before allows for a function to be called before the ResponseWriter has been written to. This is</span>
    <span class="comment">// useful for setting headers or any other operations that must happen before a response has been written.</span>
    Before(BeforeFunc)
}</pre>

ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
if the functionality calls for it.











### func NewResponseWriter
<pre>func NewResponseWriter(rw http.ResponseWriter) ResponseWriter</pre>

NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter









## type Route
<pre>type Route interface {
    <span class="comment">// URLWith returns a rendering of the Route&#39;s url with the given string params.</span>
    URLWith([]string) string
}</pre>

Route is an interface representing a Route in Martini&#39;s routing layer.















## type Router
<pre>type Router interface {
    <span class="comment">// Get adds a route for a HTTP GET request to the specified matching pattern.</span>
    Get(string, ...Handler) Route
    <span class="comment">// Patch adds a route for a HTTP PATCH request to the specified matching pattern.</span>
    Patch(string, ...Handler) Route
    <span class="comment">// Post adds a route for a HTTP POST request to the specified matching pattern.</span>
    Post(string, ...Handler) Route
    <span class="comment">// Put adds a route for a HTTP PUT request to the specified matching pattern.</span>
    Put(string, ...Handler) Route
    <span class="comment">// Delete adds a route for a HTTP DELETE request to the specified matching pattern.</span>
    Delete(string, ...Handler) Route
    <span class="comment">// Options adds a route for a HTTP OPTIONS request to the specified matching pattern.</span>
    Options(string, ...Handler) Route
    <span class="comment">// Head adds a route for a HTTP HEAD request to the specified matching pattern.</span>
    Head(string, ...Handler) Route
    <span class="comment">// Any adds a route for any HTTP method request to the specified matching pattern.</span>
    Any(string, ...Handler) Route

    <span class="comment">// NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.</span>
    NotFound(...Handler)

    <span class="comment">// Handle is the entry point for routing. This is used as a martini.Handler</span>
    Handle(http.ResponseWriter, *http.Request, Context)
}</pre>

Router is Martini&#39;s de-facto routing interface. Supports HTTP verbs, stacked handlers, and dependency injection.











### func NewRouter
<pre>func NewRouter() Router</pre>

NewRouter creates a new Router instance.









## type Routes
<pre>type Routes interface {
    <span class="comment">// URLFor returns a rendered URL for the given route. Optional params can be passed to fulfill named parameters in the route.</span>
    URLFor(route Route, params ...interface{}) string
}</pre>

Routes is a helper service for Martini&#39;s routing layer.



















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)