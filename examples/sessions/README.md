
# sessions
    import "github.com/gorilla/sessions"

Package gorilla/sessions provides cookie and filesystem sessions and
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
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		// Save it.
		session.Save(r, w)
	}

First we initialize a session store calling NewCookieStore() and passing a
secret key used to authenticate the session. Inside the handler, we call
store.Get() to retrieve an existing session or a new one. Then we set some
session values in session.Values, which is a map[interface{}]interface{}.
And finally we call session.Save() to save the session in the response.

Note that in production code, we should check for errors when calling
session.Save(r, w), and either display an error message or otherwise handle it.

That's all you need to know for the basic usage. Let's take a look at other
options, starting with flash messages.

Flash messages are session values that last until read. The term appeared with
Ruby On Rails a few years back. When we request a flash message, it is removed
from the session. To add a flash, call session.AddFlash(), and to get all
flashes, call session.Flashes(). Here is an example:


	func MyHandler(w http.ResponseWriter, r *http.Request) {
		// Get a session.
		session, _ := store.Get(r, "session-name")
		// Get the previously flashes, if any.
		if flashes := session.Flashes(); len(flashes) > 0 {
			// Just print the flash values.
			fmt.Fprint(w, "%v", flashes)
		} else {
			// Set a new flash.
			session.AddFlash("Hello, flash messages world!")
			fmt.Fprint(w, "No flashes found.")
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
relies on us passing it an empty pointer to the type as a parameter. In the example
above we've passed it a pointer to a struct and a pointer to a custom type
representing a map[string]interface. This will then allow us to serialise/deserialise
values of those types to and from our sessions.

By default, session cookies last for a month. This is probably too long for
some cases, but it is easy to change this and other attributes during
runtime. Sessions can be configured individually or the store can be
configured and then all sessions saved using it will use that configuration.
We access session.Options or store.Options to set a new configuration. The
fields are basically a subset of http.Cookie fields. Let's change the
maximum age of a session to one week:


	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
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









## func NewCookie
<pre>func NewCookie(name, value string, options *Options) *http.Cookie</pre>
NewCookie returns an http.Cookie with the options set. It also sets
the Expires field calculated based on the MaxAge value, for Internet
Explorer compatibility.






## func Save
<pre>func Save(r *http.Request, w http.ResponseWriter) error</pre>
Save saves all sessions used during the current request.







## type CookieStore
<pre>type CookieStore struct {
    Codecs  []securecookie.Codec
    Options *Options // default configuration
}</pre>
CookieStore stores sessions using secure cookies.











### func NewCookieStore

    func NewCookieStore(keyPairs ...[]byte) *CookieStore

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







### func (\*CookieStore) Get

    func (s *CookieStore) Get(r *http.Request, name string) (*Session, error)

Get returns a session for the given name after adding it to the registry.

It returns a new session if the sessions doesn't exist. Access IsNew on
the session to check if it is an existing session or a new one.

It returns a new session and an error if the session exists but could
not be decoded.






### func (\*CookieStore) New

    func (s *CookieStore) New(r *http.Request, name string) (*Session, error)

New returns a session for the given name without adding it to the registry.

The difference between New() and Get() is that calling New() twice will
decode the session data twice, while Get() registers and reuses the same
decoded session after the first call.






### func (\*CookieStore) Save

    func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter,
    session *Session) error

Save adds a single session to the response.








## type FilesystemStore
<pre>type FilesystemStore struct {
    Codecs  []securecookie.Codec
    Options *Options // default configuration
    // contains filtered or unexported fields
}</pre>
FilesystemStore stores sessions in the filesystem.

It also serves as a referece for custom stores.

This store is still experimental and not well tested. Feedback is welcome.











### func NewFilesystemStore

    func NewFilesystemStore(path string, keyPairs ...[]byte) *FilesystemStore

NewFilesystemStore returns a new FilesystemStore.

The path argument is the directory where sessions will be saved. If empty
it will use os.TempDir().

See NewCookieStore() for a description of the other parameters.







### func (\*FilesystemStore) Get

    func (s *FilesystemStore) Get(r *http.Request, name string) (*Session, error)

Get returns a session for the given name after adding it to the registry.

See CookieStore.Get().






### func (\*FilesystemStore) MaxLength

    func (s *FilesystemStore) MaxLength(l int)

MaxLength restricts the maximum length of new sessions to l.
If l is 0 there is no limit to the size of a session, use with caution.
The default for a new FilesystemStore is 4096.






### func (\*FilesystemStore) New

    func (s *FilesystemStore) New(r *http.Request, name string) (*Session, error)

New returns a session for the given name without adding it to the registry.

See CookieStore.New().






### func (\*FilesystemStore) Save

    func (s *FilesystemStore) Save(r *http.Request, w http.ResponseWriter,
    session *Session) error

Save adds a single session to the response.








## type MultiError
<pre>type MultiError []error</pre>
MultiError stores multiple errors.

Borrowed from the App Engine SDK.













### func (MultiError) Error

    func (m MultiError) Error() string








## type Options
<pre>type Options struct {
    Path   string
    Domain string
    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
    // MaxAge>0 means Max-Age attribute present and given in seconds.
    MaxAge   int
    Secure   bool
    HttpOnly bool
}</pre>
Options stores configuration for a session or session store.

Fields are a subset of http.Cookie fields.















## type Registry
<pre>type Registry struct {
    // contains filtered or unexported fields
}</pre>
Registry stores sessions used during a request.











### func GetRegistry

    func GetRegistry(r *http.Request) *Registry

GetRegistry returns a registry instance for the current request.







### func (\*Registry) Get

    func (s *Registry) Get(store Store, name string) (session *Session, err error)

Get registers and returns a session for the given name and session store.

It returns a new session if there are no sessions registered for the name.






### func (\*Registry) Save

    func (s *Registry) Save(w http.ResponseWriter) error

Save saves all sessions registered for the current request.








## type Session
<pre>type Session struct {
    ID      string
    Values  map[interface{}]interface{}
    Options *Options
    IsNew   bool
    // contains filtered or unexported fields
}</pre>
Session stores the values and optional configuration for a session.











### func NewSession

    func NewSession(store Store, name string) *Session

NewSession is called by session stores to create a new session instance.







### func (\*Session) AddFlash

    func (s *Session) AddFlash(value interface{}, vars ...string)

AddFlash adds a flash message to the session.

A single variadic argument is accepted, and it is optional: it defines
the flash key. If not defined "_flash" is used by default.






### func (\*Session) Flashes

    func (s *Session) Flashes(vars ...string) []interface{}

Flashes returns a slice of flash messages from the session.

A single variadic argument is accepted, and it is optional: it defines
the flash key. If not defined "_flash" is used by default.






### func (\*Session) Name

    func (s *Session) Name() string

Name returns the name used to register the session.






### func (\*Session) Save

    func (s *Session) Save(r *http.Request, w http.ResponseWriter) error

Save is a convenience method to save this session. It is the same as calling
store.Save(request, response, session)






### func (\*Session) Store

    func (s *Session) Store() Store

Store returns the session store used to register the session.








## type Store
<pre>type Store interface {
    // Get should return a cached session.
    Get(r *http.Request, name string) (*Session, error)

    // New should create and return a new session.
    //
    // Note that New should never return a nil session, even in the case of
    // an error if using the Registry infrastructure to cache the session.
    New(r *http.Request, name string) (*Session, error)

    // Save should persist session to the underlying store implementation.
    Save(r *http.Request, w http.ResponseWriter, s *Session) error
}</pre>
Store is an interface for custom session stores.

See CookieStore and FilesystemStore for examples.



















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)