# Tamarin Language Features

Here is an overview of Tamarin Language Features. This is not comprehensive.

## Print

Print to stdout:

```
print("Hello gophers!")
```

Print any number of variables:

```
print("x:", x, "y:", y)
```

Equivalent to `fmt.Println`.

## Assignment Statements

Both `let` and `const` statements are supported:

```
let x = 42
const y = "this is a constant"
```

Using the `:=` operator instead of `let` is encouraged:

```
x := 42
```

## Dynamic Typing

Variables may change type, similar to Python.

```
let x = 42
x = "now a string"
print(x)
```

## Semicolons

Semicolons are optional, so statements are ended by newlines if semicolons are not present.

```
let foo = "bar"; let baz = "qux"
```

## Comments

Lines are commented using `//`.

```
// This line is commented out
```

## Strings

Strings come in three varieties, with some different features available depending
on each type. The most basic form is the double-quoted string, which behaves very
similarly to Go's string type. Single quoted strings have the additional behavior
of string templating using variables like Python f-strings. Finally, strings defined
using backticks opt out of all character escaping and support multi-line strings.

```
salutation := "hello there"         // double-quoted string
count := 12
count_str := "the count is {count}" // single-quoted string with variable
backtick_str := `\t\r\n`            // backtick string
```

## String Templates

Arbitrary Tamarin code can be specified in string templates. The only limitation
is that `{` and `}` are disallowed, since those are used to delineate the template
variables. In practice we think this limitation doesn't get in the way of most
common uses of string templating.

A few quick examples:

```
len_msg := 'list length: {len([1, 2, 3])}' // "list length: 3"
count := 42
count_msg := 'the count is {count}'        // "the count is 42"
type_msg := 'type: {type(1.2)}'            // "type: float"
```

## Functions

Functions are defined using the `func` keyword. They may be passed around as values.
The `return` keyword is optional. If not present, the value of the last statement or
expression in the block is understood to be the return value. Expressions that do not
evaluate to a value will result in an `*object.Null` being returned.

```

func addOne(x) {
x + 1
}

// This way of defining a function is equivalent to the above
const subOne = func(x) {
return x - 1
}

addOne(100)
subOne(100)

```

Default parameter values are supported:

```

func increment(value, amount=1) {
return value + amount
}

print(increment(100)) // 101

```

## Conditionals

Go style if-else statements are supported.

```

name := "ben"

if name == "noa" {
print("the name is noa")
} else {
print("the name is something else")
}

```

## Switch Statements

Go style switch statements are supported.

```

name := "ben"

switch name {
case "ben":
print("matched ben")
case "noa":
print("matched noa")
default:
print("default")
}

```

## Loops

Two forms of for loops are accepted. The `break` keyword may be used to
stop looping in either form.

```

for i := 0; i < 10; i++ {
print(i)
}

for {
if condition {
break
}
}

```

## Operations that may fail

`Result` objects wrap `Ok` and `Err` values for operations that may fail.

```

obj := json.unmarshal("true")
obj.unwrap() // returns true

failed := json.unmarshal("/not-valid/")
failed.is_err() // returns true
failed.unwrap() // raises error that stops execution

```

## Pipe Expressions

These execute a series of function calls, passing the result from one stage
in as the first argument to the next.

This pipe expression evalutes to the string `"HELLO"`.

```

"hello" | strings.to_upper

```

## Array Methods

Arrays offer `map` and `filter` methods:

```

arr := [1, 2, 3, 4].filter(func(x) { x < 3 })
arr = arr.map(func(x) { x \* x })
// arr is now [1, 4]

```

## Builtins

```

type(x) // returns the string type name of x
len(s) // returns the size of the string, array, hash, or set
any(arr) // true if any item in arr is truthy
all(arr) // true if all items in arr are truthy
match(regex, str) // check if the string matches the regex
sprintf(msg, ...) // equivalent to fmt.Sprintf
keys(hash) // returns an array of keys in the given hash
delete(hash, key) // delete an item from the hash
string(obj) // convert an object to its string representation
bool(obj) // evaluates an object's truthiness
ok(result) // create a Result object containing the given object
err(message) // create a Result error object
unwrap(result) // unwraps the ok value from the Result if allowed
unwrap_or(obj) // unwraps but returns the provided obj if the Result is an Error
sorted(obj) // works with arrays, hashes, and sets
reversed(arr) // returns a reversed version of the given array
assert(obj, msg) // raises an error if obj is falsy
print(...) // equivalent to fmt.Println
printf(...) // equivalent to fmt.Printf

```

## Types

A variety of built-in types are available.

```

101 // integer
1.1 // float
"1" // string
[1,2,3] // array
{1:2} // hash
{1,2} // set
false // boolean
null // null
func() {} // function
time.now() // time

```

There are also `HttpResponse` and `DatabaseConnection` types in progress.

## Standard Library

Documentation for this is a work in progress. For now, browse the modules [here](../internal/modules).

## Proxying Calls to Go Objects

You can expose arbitrary Go objects to Tamarin code in order to enable method
calls on those objects. This allows you to expose existing structs in your
application as Tamarin objects that scripts can be written against. Tamarin
automatically discovers public methods on your Go types and converts inputs and
outputs for primitive types and for structs that you register.

In order to do this, create a ProxyManager and indicate which Go types you
want to support working with:

```go
	proxyMgr, err := object.NewProxyManager(object.ProxyManagerOpts{
		Types: []any{
			&MyService{},    // a service type
			MyServiceOpts{}, // a struct used as a parameter to service methods
		},
	})
```

Then create a global scope for your script executions that includes a proxy
to your Go service:

```go
	svc := &MyService{}
	s := scope.New(scope.Opts{})
	s.Declare("svc", object.NewProxy(proxyMgr, svc), true)

	result, err := exec.Execute(ctx, exec.Opts{
		Input: string(scriptSourceCode),
		Scope: s,
	})
```

Now in your Tamarin script you can make calls against the service. See
[example-proxy](../cmd/example-proxy/main.go) for a complete example.
