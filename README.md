
# 2019-09-27 Introduction - Part 1

## Syntax basics

* `package` declaration
* packages `import`
* `fmt` package basic usage
* package exported names
* `var` declarations
* basic types: `bool`, `byte`, `int`, `float64`, `string`, `error` and more
* `type` declarations: creating aliases for basic types / functions
* `func` declarations
* functions as values (as arguments, returned values, variables), closures
* type conversion with `T(v)`
* `const` declarations
* pointers: `&` and `*` operators, passing arguments by value vs by reference
* `if else` statements
* `:=` statement: shadowing "gotcha"
* `for` loop: typical usage / replacement for while / "forever" loop
* `switch` statement: typical usage / replacement for long if-else
* declaring arrays, arrays literals
* `struct` anonymously and with `type` alias, accessing fields
* using fields names in `struct` literal

## Generic slices: []T

* what is slice? slice underlying array
* `len` and `cap`
* slice zero value declaration: nil as ready-to-use value
* `append` builtin
* creating slice from array and other slice: underlying array reusage
* slicing default bounds
* slice allocation with `make`: limiting number of allocations

# 2019-10-03 Introduction - Part 2

## Generic maps: map[K]V

* map value is a pointer: zero value (nil) can not be used
* allocating map: 2 equivalent ways
* map literals: possible to omit type name from values literals
* inserting/updating elements
* retrieving elements: zero value for missing elements
* removing element: `delete(m, key)`
* checking for presence: two-value assignment
* key can be of any comparable type: structs as keys

## Range

* range over slices
* range over maps: random iteration order

## Methods

* adding method to a type: only for types declared in package
* can have value or pointer receiver
* just like normal function but can be called differently
* can be called by accessing method on receiver
* can use "pointer indirection"
* methods set can implement interface

## Interfaces

* set of method signatures
* implemented implicitly
* `interface{}` implemented by all types
* type assertion `v.(T)`: with and without `ok`
* type switches `switch v := v.(type)`
* `error` is an interface with single method: `Error() string`

## 2019-10-09 `testing` package

* implementing simple tests for div
* runnig tests on package: verbose flag, selecting tests to run
* failure messages: can be many
* runnig subtests
* -cover, -coverprofile
* implementing simple benchmarks
* getting cpu profile out of benchmark
* assertions library
* workshop

## Useful URLs

* https://tour.golang.org
* https://golang.org/doc/
* https://golang.org/doc/editors.html
* https://golang.org/doc/effective_go.html
