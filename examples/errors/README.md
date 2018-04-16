

# errors
`import "github.com/pkg/errors"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package errors provides simple error handling primitives.

The traditional error handling idiom in Go is roughly akin to


	if err != nil {
	        return err
	}

which applied recursively up the call stack results in error reports
without context or debugging information. The errors package allows
programmers to add context to the failure path in their code in a way
that does not destroy the original value of the error.

### Adding context to an error
The errors.Wrap function returns a new error that adds context to the
original error by recording a stack trace at the point Wrap is called,
and the supplied message. For example


	_, err := ioutil.ReadAll(r)
	if err != nil {
	        return errors.Wrap(err, "read failed")
	}

If additional control is required the errors.WithStack and errors.WithMessage
functions destructure errors.Wrap into its component operations of annotating
an error with a stack trace and an a message, respectively.

### Retrieving the cause of an error
Using errors.Wrap constructs a stack of errors, adding context to the
preceding error. Depending on the nature of the error it may be necessary
to reverse the operation of errors.Wrap to retrieve the original error
for inspection. Any error value which implements this interface


	type causer interface {
	        Cause() error
	}

can be inspected by errors.Cause. errors.Cause will recursively retrieve
the topmost error which does not implement causer, which is assumed to be
the original cause. For example:


	switch err := errors.Cause(err).(type) {
	case *MyError:
	        // handle specifically
	default:
	        // unknown error
	}

causer interface is not exported by this package, but is considered a part
of stable public API.

### Formatted printing of errors
All error values returned from this package implement fmt.Formatter and can
be formatted by the fmt package. The following verbs are supported


	%s    print the error. If the error has a Cause it will be
	      printed recursively
	%v    see %s
	%+v   extended format. Each Frame of the error's StackTrace will
	      be printed in detail.

### Retrieving the stack trace of an error or wrapper
New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
invoked. This information can be retrieved with the following interface.


	type stackTracer interface {
	        StackTrace() errors.StackTrace
	}

Where errors.StackTrace is defined as


	type StackTrace []Frame

The Frame type represents a call site in the stack trace. Frame supports
the fmt.Formatter interface that can be used for printing information about
the stack trace of this error. For example:


	if err, ok := err.(stackTracer); ok {
	        for _, f := range err.StackTrace() {
	                fmt.Printf("%+s:%d", f)
	        }
	}

stackTracer interface is not exported by this package, but is considered a part
of stable public API.

See the documentation for Frame.Format for more details.


##### Example  (StackTrace):
``` go
type stackTracer interface {
    StackTrace() errors.StackTrace
}

err, ok := errors.Cause(fn()).(stackTracer)
if !ok {
    panic("oops, err does not implement stackTracer")
}

st := err.StackTrace()
fmt.Printf("%+v", st[0:2]) // top two frames

// Example output:
// github.com/pkg/errors_test.fn
//	/home/dfc/src/github.com/pkg/errors/example_test.go:47
// github.com/pkg/errors_test.Example_stackTrace
//	/home/dfc/src/github.com/pkg/errors/example_test.go:127
```



## <a name="pkg-index">Index</a>
* [func Cause(err error) error](#Cause)
* [func Errorf(format string, args ...interface{}) error](#Errorf)
* [func New(message string) error](#New)
* [func WithMessage(err error, message string) error](#WithMessage)
* [func WithStack(err error) error](#WithStack)
* [func Wrap(err error, message string) error](#Wrap)
* [func Wrapf(err error, format string, args ...interface{}) error](#Wrapf)
* [type Frame](#Frame)
  * [func (f Frame) Format(s fmt.State, verb rune)](#Frame.Format)
* [type StackTrace](#StackTrace)
  * [func (st StackTrace) Format(s fmt.State, verb rune)](#StackTrace.Format)

#### <a name="pkg-examples">Examples</a>
* [Cause](#example-cause)
* [Cause (Printf)](#example-cause-printf)
* [Errorf (Extended)](#example-errorf-extended)
* [New](#example-new)
* [New (Printf)](#example-new-printf)
* [WithMessage](#example-withmessage)
* [WithStack](#example-withstack)
* [WithStack (Printf)](#example-withstack-printf)
* [Wrap](#example-wrap)
* [Wrap (Extended)](#example-wrap-extended)
* [Wrapf](#example-wrapf)
* [Package (StackTrace)](#example--stacktrace)

#### <a name="pkg-files">Package files</a>
[errors.go](https://github.com/pkg/errors/tree/master/errors.go) [stack.go](https://github.com/pkg/errors/tree/master/stack.go)





## <a name="Cause">func</a> [Cause](https://github.com/pkg/errors/tree/master/errors.go?s=6654:6681#L256)
``` go
func Cause(err error) error
```
Cause returns the underlying cause of the error, if possible.
An error value has a cause if it implements the following
interface:


	type causer interface {
	       Cause() error
	}

If the error does not implement Cause, the original error will
be returned. If the error is nil, nil will be returned without further
investigation.


##### Example Cause:
``` go
err := fn()
fmt.Println(err)
fmt.Println(errors.Cause(err))

// Output: outer: middle: inner: error
// error
```

Output:

```
outer: middle: inner: error
error
```


##### Example Cause (Printf):
``` go
err := errors.Wrap(func() error {
    return func() error {
        return errors.Errorf("hello %s", fmt.Sprintf("world"))
    }()
}(), "failed")

fmt.Printf("%v", err)

// Output: failed: hello world
```

Output:

```
failed: hello world
```


## <a name="Errorf">func</a> [Errorf](https://github.com/pkg/errors/tree/master/errors.go?s=3695:3748#L111)
``` go
func Errorf(format string, args ...interface{}) error
```
Errorf formats according to a format specifier and returns the string
as a value that satisfies error.
Errorf also records the stack trace at the point it was called.


##### Example Errorf (Extended):
``` go
err := errors.Errorf("whoops: %s", "foo")
fmt.Printf("%+v", err)

// Example output:
// whoops: foo
// github.com/pkg/errors_test.ExampleErrorf
//         /home/dfc/src/github.com/pkg/errors/example_test.go:101
// testing.runExample
//         /home/dfc/go/src/testing/example.go:114
// testing.RunExamples
//         /home/dfc/go/src/testing/example.go:38
// testing.(*M).Run
//         /home/dfc/go/src/testing/testing.go:744
// main.main
//         /github.com/pkg/errors/_test/_testmain.go:102
// runtime.main
//         /home/dfc/go/src/runtime/proc.go:183
// runtime.goexit
//         /home/dfc/go/src/runtime/asm_amd64.s:2059
```


## <a name="New">func</a> [New](https://github.com/pkg/errors/tree/master/errors.go?s=3420:3450#L101)
``` go
func New(message string) error
```
New returns an error with the supplied message.
New also records the stack trace at the point it was called.


##### Example New:
``` go
err := errors.New("whoops")
fmt.Println(err)

// Output: whoops
```

Output:

```
whoops
```


##### Example New (Printf):
``` go
err := errors.New("whoops")
fmt.Printf("%+v", err)

// Example output:
// whoops
// github.com/pkg/errors_test.ExampleNew_printf
//         /home/dfc/src/github.com/pkg/errors/example_test.go:17
// testing.runExample
//         /home/dfc/go/src/testing/example.go:114
// testing.RunExamples
//         /home/dfc/go/src/testing/example.go:38
// testing.(*M).Run
//         /home/dfc/go/src/testing/testing.go:744
// main.main
//         /github.com/pkg/errors/_test/_testmain.go:106
// runtime.main
//         /home/dfc/go/src/runtime/proc.go:183
// runtime.goexit
//         /home/dfc/go/src/runtime/asm_amd64.s:2059
```


## <a name="WithMessage">func</a> [WithMessage](https://github.com/pkg/errors/tree/master/errors.go?s=5698:5747#L213)
``` go
func WithMessage(err error, message string) error
```
WithMessage annotates err with a new message.
If err is nil, WithMessage returns nil.


##### Example WithMessage:
``` go
cause := errors.New("whoops")
err := errors.WithMessage(cause, "oh noes")
fmt.Println(err)

// Output: oh noes: whoops
```

Output:

```
oh noes: whoops
```


## <a name="WithStack">func</a> [WithStack](https://github.com/pkg/errors/tree/master/errors.go?s=4406:4437#L144)
``` go
func WithStack(err error) error
```
WithStack annotates err with a stack trace at the point WithStack was called.
If err is nil, WithStack returns nil.


##### Example WithStack:
``` go
cause := errors.New("whoops")
err := errors.WithStack(cause)
fmt.Println(err)

// Output: whoops
```

Output:

```
whoops
```


##### Example WithStack (Printf):
``` go
cause := errors.New("whoops")
err := errors.WithStack(cause)
fmt.Printf("%+v", err)

// Example Output:
// whoops
// github.com/pkg/errors_test.ExampleWithStack_printf
//         /home/fabstu/go/src/github.com/pkg/errors/example_test.go:55
// testing.runExample
//         /usr/lib/go/src/testing/example.go:114
// testing.RunExamples
//         /usr/lib/go/src/testing/example.go:38
// testing.(*M).Run
//         /usr/lib/go/src/testing/testing.go:744
// main.main
//         github.com/pkg/errors/_test/_testmain.go:106
// runtime.main
//         /usr/lib/go/src/runtime/proc.go:183
// runtime.goexit
//         /usr/lib/go/src/runtime/asm_amd64.s:2086
// github.com/pkg/errors_test.ExampleWithStack_printf
//         /home/fabstu/go/src/github.com/pkg/errors/example_test.go:56
// testing.runExample
//         /usr/lib/go/src/testing/example.go:114
// testing.RunExamples
//         /usr/lib/go/src/testing/example.go:38
// testing.(*M).Run
//         /usr/lib/go/src/testing/testing.go:744
// main.main
//         github.com/pkg/errors/_test/_testmain.go:106
// runtime.main
//         /usr/lib/go/src/runtime/proc.go:183
// runtime.goexit
//         /usr/lib/go/src/runtime/asm_amd64.s:2086
```


## <a name="Wrap">func</a> [Wrap](https://github.com/pkg/errors/tree/master/errors.go?s=5050:5092#L180)
``` go
func Wrap(err error, message string) error
```
Wrap returns an error annotating err with a stack trace
at the point Wrap is called, and the supplied message.
If err is nil, Wrap returns nil.


##### Example Wrap:
``` go
cause := errors.New("whoops")
err := errors.Wrap(cause, "oh noes")
fmt.Println(err)

// Output: oh noes: whoops
```

Output:

```
oh noes: whoops
```


##### Example Wrap (Extended):
``` go
err := fn()
fmt.Printf("%+v\n", err)

// Example output:
// error
// github.com/pkg/errors_test.fn
//         /home/dfc/src/github.com/pkg/errors/example_test.go:47
// github.com/pkg/errors_test.ExampleCause_printf
//         /home/dfc/src/github.com/pkg/errors/example_test.go:63
// testing.runExample
//         /home/dfc/go/src/testing/example.go:114
// testing.RunExamples
//         /home/dfc/go/src/testing/example.go:38
// testing.(*M).Run
//         /home/dfc/go/src/testing/testing.go:744
// main.main
//         /github.com/pkg/errors/_test/_testmain.go:104
// runtime.main
//         /home/dfc/go/src/runtime/proc.go:183
// runtime.goexit
//         /home/dfc/go/src/runtime/asm_amd64.s:2059
// github.com/pkg/errors_test.fn
// 	  /home/dfc/src/github.com/pkg/errors/example_test.go:48: inner
// github.com/pkg/errors_test.fn
//        /home/dfc/src/github.com/pkg/errors/example_test.go:49: middle
// github.com/pkg/errors_test.fn
//      /home/dfc/src/github.com/pkg/errors/example_test.go:50: outer
```


## <a name="Wrapf">func</a> [Wrapf](https://github.com/pkg/errors/tree/master/errors.go?s=5384:5447#L197)
``` go
func Wrapf(err error, format string, args ...interface{}) error
```
Wrapf returns an error annotating err with a stack trace
at the point Wrapf is call, and the format specifier.
If err is nil, Wrapf returns nil.


##### Example Wrapf:
``` go
cause := errors.New("whoops")
err := errors.Wrapf(cause, "oh noes #%d", 2)
fmt.Println(err)

// Output: oh noes #2: whoops
```

Output:

```
oh noes #2: whoops
```



## <a name="Frame">type</a> [Frame](https://github.com/pkg/errors/tree/master/stack.go?s=131:149#L12)
``` go
type Frame uintptr
```
Frame represents a program counter inside a stack frame.










### <a name="Frame.Format">func</a> (Frame) [Format](https://github.com/pkg/errors/tree/master/stack.go?s=1204:1249#L52)
``` go
func (f Frame) Format(s fmt.State, verb rune)
```
Format formats the frame according to the fmt.Formatter interface.


	%s    source file
	%d    source line
	%n    function name
	%v    equivalent to %s:%d

Format accepts flags that alter the printing of some verbs, as follows:


	%+s   function name and path of source file relative to the compile time
	      GOPATH separated by \n\t (<funcname>\n\t<path>)
	%+v   equivalent to %+s:%d




## <a name="StackTrace">type</a> [StackTrace](https://github.com/pkg/errors/tree/master/stack.go?s=1854:1877#L81)
``` go
type StackTrace []Frame
```
StackTrace is stack of Frames from innermost (newest) to outermost (oldest).










### <a name="StackTrace.Format">func</a> (StackTrace) [Format](https://github.com/pkg/errors/tree/master/stack.go?s=2258:2309#L91)
``` go
func (st StackTrace) Format(s fmt.State, verb rune)
```
Format formats the stack of Frames according to the fmt.Formatter interface.


	%s	lists source files for each Frame in the stack
	%v	lists the source file and line number for each Frame in the stack

Format accepts flags that alter the printing of some verbs, as follows:


	%+v   Prints filename, function, and line number for each Frame in the stack.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
