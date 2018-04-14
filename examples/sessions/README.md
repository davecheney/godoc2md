

# sessions
`import "github.com/gorilla/sessions"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package sessions provides cookie and filesystem sessions and
infrastructure for custom session backends.

The key features are:


	* Simple API: use it as an easy way to set signed (and optionally
	  encrypted) cookies.
	* Built-in backends to store sessions in cookies or the filesystem.
	* Flash messages: session values that last until read.
	* Convenient way to switch session persistency (aka "remember me") and set
	  other attributes.
	* Mechanism to rotate authentication and encryption keys.
	* Multiple sessions per request, even using different backends.
	* Interfaces and infrastructure for custom session backends: sessions from
	  different stores can be retrieved and batch-saved using a common API.

Let's start with an example that shows the sessions API in a nutshell:


	import (
		"net/http"
		"github.com/gorilla/sessions"
	)
	
	var store = sessions.NewCookieStore([]byte("something-very-secret"))
	
	func MyHandler(w http.ResponseWriter, r *http.Request) {
		// Get a session. Get() always returns a session, even if empty.
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		// Save it before we write to the response/return from the handler.
		session.Save(r, w)
	}

First we initialize a session store calling NewCookieStore() and passing a
secret key used to authenticate the session. Inside the handler, we call
store.Get() to retrieve an existing session or a new one. Then we set some
session values in session.Values, which is a map[interface{}]interface{}.
And finally we call session.Save() to save the session in the response.

Note that in production code, we should check for errors when calling
session.Save(r, w), and either display an error message or otherwise handle it.

Save must be called before writing to the response, otherwise the session
cookie will not be sent to the client.

Important Note: If you aren't using gorilla/mux, you need to wrap your handlers
with context.ClearHandler as or else you will leak memory! An easy way to do this
is to wrap the top-level mux when calling http.ListenAndServe:


	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))

The ClearHandler function is provided by the gorilla/context package.

That's all you need to know for the basic usage. Let's take a look at other
options, starting with flash messages.

Flash messages are session values that last until read. The term appeared with
Ruby On Rails a few years back. When we request a flash message, it is removed
from the session. To add a flash, call session.AddFlash(), and to get all
flashes, call session.Flashes(). Here is an example:


	func MyHandler(w http.ResponseWriter, r *http.Request) {
		// Get a session.
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Get the previous flashes, if any.
		if flashes := session.Flashes(); len(flashes) > 0 {
			// Use the flash values.
		} else {
			// Set a new flash.
			session.AddFlash("Hello, flash messages world!")
		}
		session.Save(r, w)
	}

Flash messages are useful to set information to be read after a redirection,
like after form submissions.

There may also be cases where you want to store a complex datatype within a
session, such as a struct. Sessions are serialised using the encoding/gob package,
so it is easy to register new datatypes for storage in sessions:


	import(
		"encoding/gob"
		"github.com/gorilla/sessions"
	)
	
	type Person struct {
		FirstName	string
		LastName 	string
		Email		string
		Age			int
	}
	
	type M map[string]interface{}
	
	func init() {
	
		gob.Register(&Person{})
		gob.Register(&M{})
	}

As it's not possible to pass a raw type as a parameter to a function, gob.Register()
relies on us passing it a value of the desired type. In the example above we've passed
it a pointer to a struct and a pointer to a custom type representing a
map[string]interface. (We could have passed non-pointer values if we wished.) This will
then allow us to serialise/deserialise values of those types to and from our sessions.

Note that because session values are stored in a map[string]interface{}, there's
a need to type-assert data when retrieving it. We'll use the Person struct we registered above:


	func MyHandler(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Retrieve our struct and type-assert it
		val := session.Values["person"]
		var person = &Person{}
		if person, ok := val.(*Person); !ok {
			// Handle the case that it's not an expected type
		}
	
		// Now we can use our person object
	}

By default, session cookies last for a month. This is probably too long for
some cases, but it is easy to change this and other attributes during
runtime. Sessions can be configured individually or the store can be
configured and then all sessions saved using it will use that configuration.
We access session.Options or store.Options to set a new configuration. The
fields are basically a subset of http.Cookie fields. Let's change the
maximum age of a session to one week:


	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

Sometimes we may want to change authentication and/or encryption keys without
breaking existing sessions. The CookieStore supports key rotation, and to use
it you just need to set multiple authentication and encryption keys, in pairs,
to be tested in order:


	var store = sessions.NewCookieStore(
		[]byte("new-authentication-key"),
		[]byte("new-encryption-key"),
		[]byte("old-authentication-key"),
		[]byte("old-encryption-key"),
	)

New sessions will be saved using the first pair. Old sessions can still be
read because the first pair will fail, and the second will be tested. This
makes it easy to "rotate" secret keys and still be able to validate existing
sessions. Note: for all pairs the encryption key is optional; set it to nil
or omit it and and encryption won't be used.

Multiple sessions can be used in the same request, even with different
session backends. When this happens, calling Save() on each session
individually would be cumbersome, so we have a way to save all sessions
at once: it's sessions.Save(). Here's an example:


	var store = sessions.NewCookieStore([]byte("something-very-secret"))
	
	func MyHandler(w http.ResponseWriter, r *http.Request) {
		// Get a session and set a value.
		session1, _ := store.Get(r, "session-one")
		session1.Values["foo"] = "bar"
		// Get another session and set another value.
		session2, _ := store.Get(r, "session-two")
		session2.Values[42] = 43
		// Save all sessions.
		sessions.Save(r, w)
	}

This is possible because when we call Get() from a session store, it adds the
session to a common registry. Save() uses it to save all registered sessions.




## <a name="pkg-index">Index</a>
* [func NewCookie(name, value string, options *Options) *http.Cookie](#NewCookie)
* [func Save(r *http.Request, w http.ResponseWriter) error](#Save)
* [type CookieStore](#CookieStore)
  * [func NewCookieStore(keyPairs ...[]byte) *CookieStore](#NewCookieStore)
  * [func (s *CookieStore) Get(r *http.Request, name string) (*Session, error)](#CookieStore.Get)
  * [func (s *CookieStore) MaxAge(age int)](#CookieStore.MaxAge)
  * [func (s *CookieStore) New(r *http.Request, name string) (*Session, error)](#CookieStore.New)
  * [func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter, session *Session) error](#CookieStore.Save)
* [type FilesystemStore](#FilesystemStore)
  * [func NewFilesystemStore(path string, keyPairs ...[]byte) *FilesystemStore](#NewFilesystemStore)
  * [func (s *FilesystemStore) Get(r *http.Request, name string) (*Session, error)](#FilesystemStore.Get)
  * [func (s *FilesystemStore) MaxAge(age int)](#FilesystemStore.MaxAge)
  * [func (s *FilesystemStore) MaxLength(l int)](#FilesystemStore.MaxLength)
  * [func (s *FilesystemStore) New(r *http.Request, name string) (*Session, error)](#FilesystemStore.New)
  * [func (s *FilesystemStore) Save(r *http.Request, w http.ResponseWriter, session *Session) error](#FilesystemStore.Save)
* [type MultiError](#MultiError)
  * [func (m MultiError) Error() string](#MultiError.Error)
* [type Options](#Options)
* [type Registry](#Registry)
  * [func GetRegistry(r *http.Request) *Registry](#GetRegistry)
  * [func (s *Registry) Get(store Store, name string) (session *Session, err error)](#Registry.Get)
  * [func (s *Registry) Save(w http.ResponseWriter) error](#Registry.Save)
* [type Session](#Session)
  * [func NewSession(store Store, name string) *Session](#NewSession)
  * [func (s *Session) AddFlash(value interface{}, vars ...string)](#Session.AddFlash)
  * [func (s *Session) Flashes(vars ...string) []interface{}](#Session.Flashes)
  * [func (s *Session) Name() string](#Session.Name)
  * [func (s *Session) Save(r *http.Request, w http.ResponseWriter) error](#Session.Save)
  * [func (s *Session) Store() Store](#Session.Store)
* [type Store](#Store)


#### <a name="pkg-files">Package files</a>
[doc.go](https://github.com/gorilla/sessions/tree/master/doc.go) [lex.go](https://github.com/gorilla/sessions/tree/master/lex.go) [sessions.go](https://github.com/gorilla/sessions/tree/master/sessions.go) [store.go](https://github.com/gorilla/sessions/tree/master/store.go)





## <a name="NewCookie">func</a> [NewCookie](https://github.com/gorilla/sessions/tree/master/sessions.go?s=5493:5558#L197)
``` go
func NewCookie(name, value string, options *Options) *http.Cookie
```
NewCookie returns an http.Cookie with the options set. It also sets
the Expires field calculated based on the MaxAge value, for Internet
Explorer compatibility.



## <a name="Save">func</a> [Save](https://github.com/gorilla/sessions/tree/master/sessions.go?s=5231:5286#L190)
``` go
func Save(r *http.Request, w http.ResponseWriter) error
```
Save saves all sessions used during the current request.




## <a name="CookieStore">type</a> [CookieStore](https://github.com/gorilla/sessions/tree/master/store.go?s=1999:2099#L67)
``` go
type CookieStore struct {
    Codecs  []securecookie.Codec
    Options *Options // default configuration
}
```
CookieStore stores sessions using secure cookies.







### <a name="NewCookieStore">func</a> [NewCookieStore](https://github.com/gorilla/sessions/tree/master/store.go?s=1704:1756#L53)
``` go
func NewCookieStore(keyPairs ...[]byte) *CookieStore
```
NewCookieStore returns a new CookieStore.

Keys are defined in pairs to allow key rotation, but the common case is
to set a single authentication key and optionally an encryption key.

The first key in a pair is used for authentication and the second for
encryption. The encryption key can be set to nil or omitted in the last
pair, but the authentication key is required in all pairs.

It is recommended to use an authentication key with 32 or 64 bytes.
The encryption key, if set, must be either 16, 24, or 32 bytes to select
AES-128, AES-192, or AES-256 modes.

Use the convenience function securecookie.GenerateRandomKey() to create
strong keys.





### <a name="CookieStore.Get">func</a> (\*CookieStore) [Get](https://github.com/gorilla/sessions/tree/master/store.go?s=2418:2491#L79)
``` go
func (s *CookieStore) Get(r *http.Request, name string) (*Session, error)
```
Get returns a session for the given name after adding it to the registry.

It returns a new session if the sessions doesn't exist. Access IsNew on
the session to check if it is an existing session or a new one.

It returns a new session and an error if the session exists but could
not be decoded.




### <a name="CookieStore.MaxAge">func</a> (\*CookieStore) [MaxAge](https://github.com/gorilla/sessions/tree/master/store.go?s=3734:3771#L119)
``` go
func (s *CookieStore) MaxAge(age int)
```
MaxAge sets the maximum age for the store and the underlying cookie
implementation. Individual sessions can be deleted by setting Options.MaxAge
= -1 for that session.




### <a name="CookieStore.New">func</a> (\*CookieStore) [New](https://github.com/gorilla/sessions/tree/master/store.go?s=2807:2880#L88)
``` go
func (s *CookieStore) New(r *http.Request, name string) (*Session, error)
```
New returns a session for the given name without adding it to the registry.

The difference between New() and Get() is that calling New() twice will
decode the session data twice, while Get() registers and reuses the same
decoded session after the first call.




### <a name="CookieStore.Save">func</a> (\*CookieStore) [Save](https://github.com/gorilla/sessions/tree/master/store.go?s=3254:3345#L105)
``` go
func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter,
    session *Session) error
```
Save adds a single session to the response.




## <a name="FilesystemStore">type</a> [FilesystemStore](https://github.com/gorilla/sessions/tree/master/store.go?s=4822:4942#L162)
``` go
type FilesystemStore struct {
    Codecs  []securecookie.Codec
    Options *Options // default configuration
    // contains filtered or unexported fields
}
```
FilesystemStore stores sessions in the filesystem.

It also serves as a reference for custom stores.

This store is still experimental and not well tested. Feedback is welcome.







### <a name="NewFilesystemStore">func</a> [NewFilesystemStore](https://github.com/gorilla/sessions/tree/master/store.go?s=4309:4382#L140)
``` go
func NewFilesystemStore(path string, keyPairs ...[]byte) *FilesystemStore
```
NewFilesystemStore returns a new FilesystemStore.

The path argument is the directory where sessions will be saved. If empty
it will use os.TempDir().

See NewCookieStore() for a description of the other parameters.





### <a name="FilesystemStore.Get">func</a> (\*FilesystemStore) [Get](https://github.com/gorilla/sessions/tree/master/store.go?s=5401:5478#L182)
``` go
func (s *FilesystemStore) Get(r *http.Request, name string) (*Session, error)
```
Get returns a session for the given name after adding it to the registry.

See CookieStore.Get().




### <a name="FilesystemStore.MaxAge">func</a> (\*FilesystemStore) [MaxAge](https://github.com/gorilla/sessions/tree/master/store.go?s=7361:7402#L246)
``` go
func (s *FilesystemStore) MaxAge(age int)
```
MaxAge sets the maximum age for the store and the underlying cookie
implementation. Individual sessions can be deleted by setting Options.MaxAge
= -1 for that session.




### <a name="FilesystemStore.MaxLength">func</a> (\*FilesystemStore) [MaxLength](https://github.com/gorilla/sessions/tree/master/store.go?s=5133:5175#L171)
``` go
func (s *FilesystemStore) MaxLength(l int)
```
MaxLength restricts the maximum length of new sessions to l.
If l is 0 there is no limit to the size of a session, use with caution.
The default for a new FilesystemStore is 4096.




### <a name="FilesystemStore.New">func</a> (\*FilesystemStore) [New](https://github.com/gorilla/sessions/tree/master/store.go?s=5628:5705#L189)
``` go
func (s *FilesystemStore) New(r *http.Request, name string) (*Session, error)
```
New returns a session for the given name without adding it to the registry.

See CookieStore.New().




### <a name="FilesystemStore.Save">func</a> (\*FilesystemStore) [Save](https://github.com/gorilla/sessions/tree/master/store.go?s=6373:6468#L213)
``` go
func (s *FilesystemStore) Save(r *http.Request, w http.ResponseWriter,
    session *Session) error
```
Save adds a single session to the response.

If the Options.MaxAge of the session is <= 0 then the session file will be
deleted from the store path. With this process it enforces the properly
session cookie handling so no need to trust in the cookie management in the
web browser.




## <a name="MultiError">type</a> [MultiError](https://github.com/gorilla/sessions/tree/master/sessions.go?s=6165:6188#L222)
``` go
type MultiError []error
```
MultiError stores multiple errors.

Borrowed from the App Engine SDK.










### <a name="MultiError.Error">func</a> (MultiError) [Error](https://github.com/gorilla/sessions/tree/master/sessions.go?s=6190:6224#L224)
``` go
func (m MultiError) Error() string
```



## <a name="Options">type</a> [Options](https://github.com/gorilla/sessions/tree/master/sessions.go?s=516:843#L24)
``` go
type Options struct {
    Path   string
    Domain string
    // MaxAge=0 means no Max-Age attribute specified and the cookie will be
    // deleted after the browser session ends.
    // MaxAge<0 means delete cookie immediately.
    // MaxAge>0 means Max-Age attribute present and given in seconds.
    MaxAge   int
    Secure   bool
    HttpOnly bool
}
```
Options stores configuration for a session or session store.

Fields are a subset of http.Cookie fields.










## <a name="Registry">type</a> [Registry](https://github.com/gorilla/sessions/tree/master/sessions.go?s=3801:3882#L141)
``` go
type Registry struct {
    // contains filtered or unexported fields
}
```
Registry stores sessions used during a request.







### <a name="GetRegistry">func</a> [GetRegistry](https://github.com/gorilla/sessions/tree/master/sessions.go?s=3456:3499#L127)
``` go
func GetRegistry(r *http.Request) *Registry
```
GetRegistry returns a registry instance for the current request.





### <a name="Registry.Get">func</a> (\*Registry) [Get](https://github.com/gorilla/sessions/tree/master/sessions.go?s=4042:4120#L149)
``` go
func (s *Registry) Get(store Store, name string) (session *Session, err error)
```
Get registers and returns a session for the given name and session store.

It returns a new session if there are no sessions registered for the name.




### <a name="Registry.Save">func</a> (\*Registry) [Save](https://github.com/gorilla/sessions/tree/master/sessions.go?s=4538:4590#L165)
``` go
func (s *Registry) Save(w http.ResponseWriter) error
```
Save saves all sessions registered for the current request.




## <a name="Session">type</a> [Session](https://github.com/gorilla/sessions/tree/master/sessions.go?s=1256:1530#L49)
``` go
type Session struct {
    // The ID of the session, generated by stores. It should not be used for
    // user data.
    ID string
    // Values contains the user-data for the session.
    Values  map[interface{}]interface{}
    Options *Options
    IsNew   bool
    // contains filtered or unexported fields
}
```
Session stores the values and optional configuration for a session.







### <a name="NewSession">func</a> [NewSession](https://github.com/gorilla/sessions/tree/master/sessions.go?s=1002:1052#L39)
``` go
func NewSession(store Store, name string) *Session
```
NewSession is called by session stores to create a new session instance.





### <a name="Session.AddFlash">func</a> (\*Session) [AddFlash](https://github.com/gorilla/sessions/tree/master/sessions.go?s=2211:2272#L83)
``` go
func (s *Session) AddFlash(value interface{}, vars ...string)
```
AddFlash adds a flash message to the session.

A single variadic argument is accepted, and it is optional: it defines
the flash key. If not defined "_flash" is used by default.




### <a name="Session.Flashes">func</a> (\*Session) [Flashes](https://github.com/gorilla/sessions/tree/master/sessions.go?s=1734:1789#L65)
``` go
func (s *Session) Flashes(vars ...string) []interface{}
```
Flashes returns a slice of flash messages from the session.

A single variadic argument is accepted, and it is optional: it defines
the flash key. If not defined "_flash" is used by default.




### <a name="Session.Name">func</a> (\*Session) [Name](https://github.com/gorilla/sessions/tree/master/sessions.go?s=2837:2868#L103)
``` go
func (s *Session) Name() string
```
Name returns the name used to register the session.




### <a name="Session.Save">func</a> (\*Session) [Save](https://github.com/gorilla/sessions/tree/master/sessions.go?s=2678:2746#L98)
``` go
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error
```
Save is a convenience method to save this session. It is the same as calling
store.Save(request, response, session). You should call Save before writing to
the response or returning from the handler.




### <a name="Session.Store">func</a> (\*Session) [Store](https://github.com/gorilla/sessions/tree/master/sessions.go?s=2954:2985#L108)
``` go
func (s *Session) Store() Store
```
Store returns the session store used to register the session.




## <a name="Store">type</a> [Store](https://github.com/gorilla/sessions/tree/master/store.go?s=425:930#L22)
``` go
type Store interface {
    // Get should return a cached session.
    Get(r *http.Request, name string) (*Session, error)

    // New should create and return a new session.
    //
    // Note that New should never return a nil session, even in the case of
    // an error if using the Registry infrastructure to cache the session.
    New(r *http.Request, name string) (*Session, error)

    // Save should persist session to the underlying store implementation.
    Save(r *http.Request, w http.ResponseWriter, s *Session) error
}
```
Store is an interface for custom session stores.

See CookieStore and FilesystemStore for examples.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
