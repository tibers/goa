# goa v2

goa v2 brings a host of fixes and has a cleaner more composable overall design:

* The top level design package focuses solely on the transport agnostic DSL
  expressions such as data structure definitions.
* The DSL engine has a clear interface and can run arbitrary DSLs.
* The generated code follows a strict separation of concern where the actual
  service implementation is isolated from the transport code.

The import path for goa v2 has changed from `github.com/goadesign/goa` to
`goa.design/goa.v2`.

## Separation of Concerns

The DSL in goa v2 makes it possible to describe the services in a transport
agnostic way. The service methods DSLs each describe the method input and output
types. Transport specific DSL then describe how the method input is built from
incoming data and how the output is serialized. For example a method may specify
that it accepts an object composed of two fields as input then the HTTP specific
DSL may specify that one of the attribute is read from the incoming request
headers while the other from the request body.

This clean decoupling means that the same service implementation can expose
endpoints accessible via multiple transport (e.g. HTTP and gRPC). goa takes care
of generating all the transport specific code including marshalling,
unmarshalling, validation etc. so that user code can focus on the actual
implementation.

The new `http` directory contains packages implementing the DSL, design objects,
code generation and runtime support for HTTP/HTTP APIs. The HTTP DSL is built on
top of the core DSL package and add transport specific keywords to describe
aspects specific to HTTP requests and responses.

## New Data Types

The primitive types now include `Int`, `Int32`, `Int64`, `UInt` `UInt32`,
`UInt64`, `Float32`, `Float64` and `Bytes`. This makes it possible to support
transports such as gRPC but also makes HTTP interface definitions crisper. The
v1 types `Integer` and `Float` have been removed in favor of these new types.

## Composable Code Generation

Code generation now follows a 2-phase process where the first phase produces a
set of writers each exposing templates that can be further modified prior to
running the last phase which generates the final artefacts. This makes it
possible for plugins to alter the code generated by the built-in code
generators.

## Getting Started

> Note: v2 is work in progress, it is not ready for production usage yet.

Install the code generation tool:

```bash
cd cmd/goa
go install
```

* The code generation for OpenAPI and HTTP servers is implemented. It can be invoked using:

```bash
goa gen [IMPORT]
```

where `goa` is the code generation tool for goa v2 installed by doing:

```bash
go install goa.design/goa.v2/cmd/goa
```

and `IMPORT` is the Go import path to the design package.

* An example server can be generated using:

```bash
goa example [IMPORT]
```

* The core and HTTP DSLs are stable, see:
  - The [core API DSL spec](https://github.com/goadesign/goa/blob/v2/dsl/_spec/dsl_spec_test.go)
  - The [type DSL spec](https://github.com/goadesign/goa/blob/v2/dsl/_spec/type_spec_test.go)
  - The [HTTP DSL spec](https://github.com/goadesign/goa/blob/v2/dsl/http/_spec/dsl_spec_test.go)

* The cellar example is functional, see:
  - The [design](https://github.com/goadesign/goa/blob/v2/examples/cellar/design)
  - The generated services: [storage](https://github.com/goadesign/goa/blob/v2/examples/cellar/gen/storage/service.go) and
    [sommelier](https://github.com/goadesign/goa/blob/v2/examples/cellar/gen/sommelier/service.go)
  - The generated endpoints: [storage](https://github.com/goadesign/goa/blob/v2/examples/cellar/gen/storage/endpoint.go) and
    [sommelier](https://github.com/goadesign/goa/blob/v2/examples/cellar/gen/sommelier/endpoint.go)
  - The generated HTTP servers: [storage](https://github.com/goadesign/goa/tree/v2/examples/cellar/gen/storage/http/server) and
    [sommelier](https://github.com/goadesign/goa/blob/v2/examples/cellar/gen/sommelier/http/server)
  - The generated HTTP clients: [storage](https://github.com/goadesign/goa/tree/v2/examples/cellar/gen/storage/http/client) and
    [sommelier](https://github.com/goadesign/goa/blob/v2/examples/cellar/gen/sommelier/http/client)
  - The generated [OpenAPIv2 spec](https://github.com/goadesign/goa.v2/tree/v2/examples/cellar/gen/openapi.json)

## Design Docs

### `API` Expression

Like in v1 the top level DSL function in v2 is `API`. The `API` DSL lists the
global properties of the API such as its hostname, its version number etc.
One change compared to v1 is the use of `Server` instead of `Host` and `Scheme`
to define API hosts. This provides a more flexible way to list multiple hosts
and is inline with the OpenAPI v3 spec.

```go
var _ = API("cellar", func() {
    Title("The virtual wine cellar")
    Version("1.0")
    Description("An example of an API implemented with goa")
	Server("https://service.goa.design:443", func() {
		Description("Production host")
	})
	Server("https://service.test.goa.design:443", func() {
		Description("Integration host")
	})
    Docs(func() {
        Description("goa guide")
        URL("http://goa.design/getting-started.html")
    })
    Contact(func() {
        Name("goa team")
        Email("admin@goa.design")
        URL("http://goa.design")
                })
    License(func() {
        Name("MIT")
    })
})
```

### `Service` Expression

The `Service` DSL defines a group of methods. This maps to a resource in REST or
a `service` declaration in gRPC. A service may define common error responses to
all the service methods, more on error responses in the next section.

```go
// The "account" service.
var _ = Service("account", func() {
    // Error which applies to all methods.
    Error(ErrUnauthorized, Unauthorized)
    
    // HTTP transport properties.
    HTTP(func() {
        Path("/accounts")
    })
}
```

The `HTTP` function makes it possible to define HTTP specific properties such as
 a common base path to all HTTP requests.

### `Method` Expression

The service methods are described using `Method`. This function defines the
method payload (input) and result (output) types. It may also list an arbitrary
number of error return values. An error return value has a name and optionally a
type. Omitting the payload or result type has the same effect as using the
built-in type `Empty` which maps to an empty body in HTTP and to the `Empty`
message in gRPC.

```go
    Method("update", func() {
        Description("Change account name")
        Payload(UpdateAccount)
        Result(Empty)
        Error(ErrNotFound)
        Error(ErrBadRequest, ErrorResult)
```

The payload, result and error types define the input and output *independently
of the transport*.

The `HTTP` function defines the mapping of the payload and result type
attributes to the HTTP request path and query string values as well as the HTTP
request and response and bodies. The `HTTP` function also defines other HTTP
specific properties such as the request path, the response HTTP status codes
etc.

```go
        HTTP(func() {
            PUT("/{accountID}")    // "accountID" request attribute
            Body(func() {
                Attribute("name")  // "name" request attribute
                Required("name")
            })
            Response(NoContent)
            Error(ErrNotFound, NotFound)
            Error(ErrBadRequest, BadRequest, ErrorResult)
        })
```

### Method Payload Type

In the example above the `accountID` HTTP request path parameter is defined by
the attribute of the `UpdateAccount` type with the same name and so is the body
attribute `name`.

Any attribute that is no explicitly mapped by the `HTTP` function is implicitly
mapped to request body attributes. This makes is simple to define mappings where
only one of the fields for the payload type is mapped to a HTTP header and all
other fields are mapped to the HTTP request body.

The body attributes may also be listed explicitly using the `Body` function.
This function accepts either a DSL listing the body attributes or the name of a
request type attribute whose type defines the body as a whole. The latter makes
it possible to use any arbitrary type to describe request body and not just
object, for example the attribute (and thus the body) could be an array.

Implicit request body definition:

```go
        HTTP(func() {
            PUT("/{accountID}")    // "accountID" request attribute
            Response(NoContent)
            Error(ErrNotFound, NotFound)
            Error(ErrBadRequest, BadRequest, ErrorResult)
        })
```

Array body definition:

```go
        HTTP(func() {
            PUT("/")
            Body("names") // Assumes request type has attribute "names"
            Response(NoContent)
            Error(ErrNotFound, NotFound)
            Error(ErrBadRequest, BadRequest, ErrorResult)
        })
```

### Method Result Type

While a service may only define one result type the `HTTP` function may list
multiple responses. Each response defines the HTTP status code, response body
shape (if any) and may also list HTTP headers. The `Tag` DSL function makes it
possible to define an attribute of the result type that is used to determine
which HTTP response to send. The function specifies the name of a result type
attribute and the value the attribute must have for the response in which the
tag is defined to be used to write the HTTP response.

By default the shape of the body of responses with HTTP status code 200 is
described by the method result type. The `HTTP` function may optionnally use
result type attributes to define response headers. Any attribute of the result
type that is not explicitly used to define a response header defines a field of
the response body implicitly. This alleviates the need to repeat all the result
type attributes to define the body since in most cases only a few would map to
headers.

The response body may also be explicitly described using the function `Body`.
The function works identically as when used to describe the request body: it may
be given a list of result type attributes in which case the body shape is an
object or the name of a specific attribute in which case the response body shape
is dictated by the type of the attribute.

```go
    Method("index", func() {
        Description("Index all accounts")
        Payload(ListAccounts)
        Result(func() {
            Attribute("marker", String, "Pagination marker")
            Attribute("accounts", CollectionOf(Account), "list of accounts")
        })
        HTTP(func() {
            GET("")
            Response(OK, func() {
                Header("marker")
                Body("accounts")
            })
        })
    })
```

The example above produces response bodies of the form
`[{"name"="foo"},{"name"="bar"}]` assuming the type `Account` only has a `name`
attribute. The same example but with the line defining the response body
(`Body("accounts")`) removed produces response bodies of the form:
`{"accounts":[{"name"="foo"},{"name"="bar"}]` since `accounts` isn't used to
define headers.

### Data Types

Like in v1, the types supported in the DSL are primitive types, array, map and
object types (note the change of nomenclature and DSL from `hash` to `map`).

The list of primitive types in v2 is:

* `Boolean`
* `Int`, `Int32`, `Int64`, `UInt`, `UInt32`, `UInt64`
* `Float32`, `Float64`
* `String`, `Bytes`
* `Any` (maps to any type, primitive or not)

Like in v1 arrays can be declared in one of two ways:

* `ArrayOf()` which accepts any type or result type and returns a type
* `CollectionOf()` which accepts result types only and returns a result type

The result type returned by `CollectionOf` contains the same views as the result
type given as argument. Each view simply renders an array where each element has
been projected using the corresponding element view.

Like in v1 the goa DSL makes it possible to define both user and result types
(called media types in v1). Result types are user types that also define views.
The DSL for defining user types and result types is the same as in v1 (using
`Type` and `ResultType` respectively).

### Payload to HTTP request mapping

The payload types describe the shape of the data given as argument to the
service methods. The HTTP transport specific DSL defines how the data is built
from the incoming HTTP request state.

The HTTP request state comprises four different parts:

- The URL path parameters (for example the route `/bottle/{id}` defines the `id` path parameter)
- The URL query string parameters
- The HTTP headers
- And finally the HTTP request body
 
The HTTP expressions drive how the generated code decodes the request into the
payload type:

* The `Param` expression defines values loaded from path or query string
  parameters.
* The `Header` expression defines values loaded from HTTP headers.
* The `Body` expression defines values loaded from the request body.


The next two sections describe the expressions in more details. 

Note that the generated code provides a default decoder implementation that
ought to be sufficient in most cases however it also makes it possible to plug a
user provided decoder in the (hopefully rare) cases when that's needed.
 
#### Mapping payload with non-object types

When the payload type is a primitive, an array or a map then the value is loaded from:

- the first URL path parameter if any
- otherwise the first query string parameter if any
- otherwise the first header if any
- otherwise the body

with the following restrictions:

- only primitive or array types may be used to define path parameters or headers
- only primitive, array and map types may be used to define query string parameters
- array and map types used to define path parameters, query string parameters or
  headers must use primitive types to define their elements

Arrays in paths and headers are represented using comma separated values.

Examples:

* simple "get by identifier" where identifiers are integers:

```go
Method("show", func() {
    Payload(Int)
    HTTP(func() {
        GET("/{id}")
    })
})
```

| Generated method | Example request | Corresponding call |
| ---------------- | --------------- | ------------------ |
| Show(int)        | GET /1          | Show(1)            |

* bulk "delete by identifiers" where identifiers are strings:

```go
Method("delete", func() {
    Payload(ArrayOf(String))
    HTTP(func() {
        DELETE("/{ids}")
    })
})
```

| Generated method   | Example request | Corresponding call         |
| ------------------ | --------------- | -------------------------- |
| Delete([]string)   | DELETE /a,b     | Delete([]string{"a", "b"}) |


> Note that in both the previous examples the name of the parameter path is
> irrelevant.

* list with filters:

```go
Method("list", func() {
    Payload(ArrayOf(String))
    HTTP(func() {
        GET("")
        Param("filter")
    })
})
```

| Generated method | Example request         | Corresponding call       |
| ---------------- | ----------------------- | ------------------------ |
| List([]string)   | GET /?filter=a&filter=b | List([]string{"a", "b"}) |

list with version:

```go
Method("list", func() {
    Payload(Float32)
    HTTP(func() {
        GET("")
        Header("version")
    })
})
```

| Generated method | Example request     | Corresponding call |
| ---------------- | ------------------- | ------------------ |
| List(float32)    | GET / [version=1.0] | List(1.0)          |

creation:

```go
Method("create", func() {
    Payload(MapOf(String, Int))
    HTTP(func() {
        POST("")
    })
})
```

| Generated method       | Example request         | Corresponding call                     |
| ---------------------- | ----------------------- | -------------------------------------- |
| Create(map[string]int) | POST / {"a": 1, "b": 2} | Create(map[string]int{"a": 1, "b": 2}) |

#### Mapping payload with object types

The HTTP expressions describe how the payload object attributes are loaded from
the HTTP request state. Different attributes may be loaded from different parts
of the request: some attributes may be loaded from the request path, some from
the query string parameters and others from the body for example. The same type
restrictions apply to the path, query string and header attributes (attributes
describing path and headers must be primitives or arrays of primitives and
attributes describing query string parameters must be primitives, arrays or maps
of primitives).

The `Body` expression makes it possible to define the payload type attribute
that describes the request body. Alternatively if the `Body` expression is
omitted then all attributes that make up the payload type and that are not used
to define a path parameter, a query string parameter or a header implicitly
describe the body.

For example, given the payload:

```go
Method("create", func() {
    Payload(func() {
        Attribute("id", Int)
        Attribute("name", String)
        Attribute("age", Int)
    })
})
```

The following HTTP expression causes the `id` attribute to get loaded from the
path parameter while `name` and `age` are loaded from the request body:

```go 
Method("create", func() {
    Payload(func() {
        Attribute("id", Int)
        Attribute("name", String)
        Attribute("age", Int)
    })
    HTTP(func() {
        POST("/{id}")
    })
})
```

| Generated method       | Example request                 | Corresponding call                               |
| ---------------------- | ------------------------------- | ------------------------------------------------ |
| Create(*CreatePayload) | POST /1 {"name": "a", "age": 2} | Create(&CreatePayload{ID: 1, Name: "a", Age: 2}) |

`Body` makes it possible to describe request bodies that are not objects such as
arrays or maps.

Consider the following payload:

```go 
Method("rate", func() {
    Payload(func() {
        Attribute("id", Int)
        Attribute("rates", MapOf(String, Float64))
    })
})
```

Using the following HTTP expression the rates are loaded from the body:

```go 
Method("rate", func() {
    Payload(func() {
        Attribute("id", Int)
        Attribute("rates", MapOf(String, Float64))
    })
    HTTP(func() {
        PUT("/{id}")
        Body("rates")
    })
})
```

| Generated method   | Example request             | Corresponding call                                                       |
| ------------------ | --------------------------- | ------------------------------------------------------------------------ |
| Rate(*RatePayload) | PUT /1 {"a": 0.5, "b": 1.0} | Rate(&RatePayload{ID: 1, Rates: map[string]float64{"a": 0.5, "b": 1.0}}) |

Without `Body` the request body shape would be an object with one key `rates`.

#### Mapping HTTP element names to attribute names

The expressions used to describe the HTTP request elements `Param`, `Header` and
`Body` may provide a mapping between the names of the elements (query string
key, header name or body field name) and the corresponding payload attribute
name. The mapping is defined using the syntax `"attribute name:element name"`,
for example:

```go 
Header("version:X-Api-Version")
```

causes the `version` attribute value to get loaded from the `X-Api-Version` HTTP
header.

The `Body` expression supports an alternative syntax where the attributes that
make up the body can be explicitly listed. This syntax allows for specifying a
mapping between the incoming data field names and the payload attribute names,
for example:

```go 
Method("create", func() {
    Payload(func() {
        Attribute("name", String)
        Attribute("age", Int)
    })
    HTTP(func() {
        POST("")
        Body(func() {
        	Attribute("name:n")
        	Attribute("age:a")
        })
    })
})
```

| Generated method       | Example request            | Corresponding call                               |
| ---------------------- | -------------------------- | ------------------------------------------------ |
| Create(*CreatePayload) | POST /1 {"n": "a", "a": 2} | Create(&CreatePayload{ID: 1, Name: "a", Age: 2}) |
