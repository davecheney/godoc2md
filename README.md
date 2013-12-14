
	
		
		
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
    <span id="Dev">Dev</span>  <a href="/pkg/builtin/#string">string</a> = &#34;development&#34;
    <span id="Prod">Prod</span> <a href="/pkg/builtin/#string">string</a> = &#34;production&#34;
    <span id="Test">Test</span> <a href="/pkg/builtin/#string">string</a> = &#34;test&#34;
)</pre>

Envs





## Variables

<pre>var <span id="Env">Env</span> = <a href="#Dev">Dev</a></pre>

Env is the environment that Martini is executing in. The MARTINI_ENV is read on initialization to set this variable.








## type BeforeFunc
<pre>type BeforeFunc func(<a href="#ResponseWriter">ResponseWriter</a>)</pre>

BeforeFunc is a function that is called before the ResponseWriter has been written to.



			

			

			

			

			
		


## type ClassicMartini
<pre>type ClassicMartini struct {
    *<a href="#Martini">Martini</a>
    <a href="#Router">Router</a>
}</pre>

ClassicMartini represents a Martini with some reasonable defaults. Embeds the router functions for convenience.



			

			

			

			
				
				<h3 id="Classic">func <a href="/target/martini.go?s=2898:2928#L87">Classic</a></h3>
				<pre>func Classic() *<a href="#ClassicMartini">ClassicMartini</a></pre>
				<p>
Classic creates a classic Martini with some basic default middleware - martini.Logger, martini.Recovery, and martini.Static.
</p>

				
			

			
		


## type Context
<pre>type Context interface {
    <a href="/pkg/github.com/codegangsta/inject/">inject</a>.<a href="/pkg/github.com/codegangsta/inject/#Injector">Injector</a>
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



			

			

			

			
				
				<h3 id="Logger">func <a href="/target/logger.go?s=164:185#L1">Logger</a></h3>
				<pre>func Logger() <a href="#Handler">Handler</a></pre>
				<p>
Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
</p>

				
			
				
				<h3 id="Recovery">func <a href="/target/recovery.go?s=163:186#L1">Recovery</a></h3>
				<pre>func Recovery() <a href="#Handler">Handler</a></pre>
				<p>
Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
</p>

				
			
				
				<h3 id="Static">func <a href="/target/static.go?s=155:192#L1">Static</a></h3>
				<pre>func Static(directory <a href="/pkg/builtin/#string">string</a>) <a href="#Handler">Handler</a></pre>
				<p>
Static returns a middleware handler that serves static files in the given directory.
</p>

				
			

			
		


## type Martini
<pre>type Martini struct {
    <a href="/pkg/github.com/codegangsta/inject/">inject</a>.<a href="/pkg/github.com/codegangsta/inject/#Injector">Injector</a>
    <span class="comment">// contains filtered or unexported fields</span>
}</pre>

Martini represents the top level web application. inject.Injector methods can be invoked to map services on a global level.



			

			

			

			
				
				<h3 id="New">func <a href="/target/martini.go?s=844:863#L27">New</a></h3>
				<pre>func New() *<a href="#Martini">Martini</a></pre>
				<p>
New creates a bare bones Martini instance. Use this method if you want to have full control over the middleware that is used.
</p>

				
			

			
				
				<h3 id="Martini.Action">func (*Martini) <a href="/target/martini.go?s=1629:1670#L46">Action</a></h3>
				<pre>func (m *<a href="#Martini">Martini</a>) Action(handler <a href="#Handler">Handler</a>)</pre>
				<p>
Action sets the handler that will be called after all the middleware has been invoked. This is set to martini.Router in a martini.Classic().
</p>

				
				
			
				
				<h3 id="Martini.Handlers">func (*Martini) <a href="/target/martini.go?s=2176:2223#L64">Handlers</a></h3>
				<pre>func (m *<a href="#Martini">Martini</a>) Handlers(handlers ...<a href="#Handler">Handler</a>)</pre>
				<p>
Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
Will panic if any of the handlers is not a callable function
</p>

				
				
			
				
				<h3 id="Martini.Run">func (*Martini) <a href="/target/martini.go?s=1797:1820#L52">Run</a></h3>
				<pre>func (m *<a href="#Martini">Martini</a>) Run()</pre>
				<p>
Run the http server. Listening on os.GetEnv(&#34;PORT&#34;) or 3000 by default.
</p>

				
				
			
				
				<h3 id="Martini.ServeHTTP">func (*Martini) <a href="/target/martini.go?s=1375:1446#L41">ServeHTTP</a></h3>
				<pre>func (m *<a href="#Martini">Martini</a>) ServeHTTP(res <a href="/pkg/net/http/">http</a>.<a href="/pkg/net/http/#ResponseWriter">ResponseWriter</a>, req *<a href="/pkg/net/http/">http</a>.<a href="/pkg/net/http/#Request">Request</a>)</pre>
				<p>
ServeHTTP is the HTTP Entry point for a Martini instance. Useful if you want to control your own HTTP server.
</p>

				
				
			
				
				<h3 id="Martini.Use">func (*Martini) <a href="/target/martini.go?s=1149:1187#L34">Use</a></h3>
				<pre>func (m *<a href="#Martini">Martini</a>) Use(handler <a href="#Handler">Handler</a>)</pre>
				<p>
Use adds a middleware Handler to the stack. Will panic if the handler is not a callable func. Middleware Handlers are invoked in the order that they are added.
</p>

				
				
			
		


## type Params
<pre>type Params map[<a href="/pkg/builtin/#string">string</a>]<a href="/pkg/builtin/#string">string</a></pre>

Params is a map of name/value pairs for named routes. An instance of martini.Params is available to be injected into any route handler.



			

			

			

			

			
		


## type ResponseWriter
<pre>type ResponseWriter interface {
    <a href="/pkg/net/http/">http</a>.<a href="/pkg/net/http/#ResponseWriter">ResponseWriter</a>
    <span class="comment">// Status returns the status code of the response or 0 if the response has not been written.</span>
    Status() <a href="/pkg/builtin/#int">int</a>
    <span class="comment">// Written returns whether or not the ResponseWriter has been written.</span>
    Written() <a href="/pkg/builtin/#bool">bool</a>
    <span class="comment">// Size returns the size of the response body.</span>
    Size() <a href="/pkg/builtin/#int">int</a>
    <span class="comment">// Before allows for a function to be called before the ResponseWriter has been written to. This is</span>
    <span class="comment">// useful for setting headers or any other operations that must happen before a response has been written.</span>
    Before(<a href="#BeforeFunc">BeforeFunc</a>)
}</pre>

ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
if the functionality calls for it.



			

			

			

			
				
				<h3 id="NewResponseWriter">func <a href="/target/response_writer.go?s=1051:1112#L20">NewResponseWriter</a></h3>
				<pre>func NewResponseWriter(rw <a href="/pkg/net/http/">http</a>.<a href="/pkg/net/http/#ResponseWriter">ResponseWriter</a>) <a href="#ResponseWriter">ResponseWriter</a></pre>
				<p>
NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
</p>

				
			

			
		


## type Route
<pre>type Route interface {
    <span class="comment">// URLWith returns a rendering of the Route&#39;s url with the given string params.</span>
    URLWith([]<a href="/pkg/builtin/#string">string</a>) <a href="/pkg/builtin/#string">string</a>
}</pre>

Route is an interface representing a Route in Martini&#39;s routing layer.



			

			

			

			

			
		


## type Router
<pre>type Router interface {
    <span class="comment">// Get adds a route for a HTTP GET request to the specified matching pattern.</span>
    Get(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Patch adds a route for a HTTP PATCH request to the specified matching pattern.</span>
    Patch(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Post adds a route for a HTTP POST request to the specified matching pattern.</span>
    Post(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Put adds a route for a HTTP PUT request to the specified matching pattern.</span>
    Put(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Delete adds a route for a HTTP DELETE request to the specified matching pattern.</span>
    Delete(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Options adds a route for a HTTP OPTIONS request to the specified matching pattern.</span>
    Options(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Head adds a route for a HTTP HEAD request to the specified matching pattern.</span>
    Head(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>
    <span class="comment">// Any adds a route for any HTTP method request to the specified matching pattern.</span>
    Any(<a href="/pkg/builtin/#string">string</a>, ...<a href="#Handler">Handler</a>) <a href="#Route">Route</a>

    <span class="comment">// NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.</span>
    NotFound(...<a href="#Handler">Handler</a>)

    <span class="comment">// Handle is the entry point for routing. This is used as a martini.Handler</span>
    Handle(<a href="/pkg/net/http/">http</a>.<a href="/pkg/net/http/#ResponseWriter">ResponseWriter</a>, *<a href="/pkg/net/http/">http</a>.<a href="/pkg/net/http/#Request">Request</a>, <a href="#Context">Context</a>)
}</pre>

Router is Martini&#39;s de-facto routing interface. Supports HTTP verbs, stacked handlers, and dependency injection.



			

			

			

			
				
				<h3 id="NewRouter">func <a href="/target/router.go?s=1720:1743#L37">NewRouter</a></h3>
				<pre>func NewRouter() <a href="#Router">Router</a></pre>
				<p>
NewRouter creates a new Router instance.
</p>

				
			

			
		


## type Routes
<pre>type Routes interface {
    <span class="comment">// URLFor returns a rendered URL for the given route. Optional params can be passed to fulfill named parameters in the route.</span>
    URLFor(route <a href="#Route">Route</a>, params ...interface{}) <a href="/pkg/builtin/#string">string</a>
}</pre>

Routes is a helper service for Martini&#39;s routing layer.



			

			

			

			

			
		
	

	


