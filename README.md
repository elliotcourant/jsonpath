# jsonpath

[![PkgGoDev](https://pkg.go.dev/badge/github.com/elliotcourant/jsonpath)](https://pkg.go.dev/github.com/elliotcourant/jsonpath)
[![Build Status](https://travis-ci.com/elliotcourant/jsonpath.svg?branch=main)](https://travis-ci.com/elliotcourant/jsonpath)
[![codecov](https://codecov.io/gh/elliotcourant/jsonpath/branch/main/graph/badge.svg)](https://codecov.io/gh/elliotcourant/jsonpath)

Go implementation of jsonpath

This is still a work in progress, updates to come.


## Example

```go
const jsonString = `{
  "firstName": "John",
  "lastName" : "doe",
  "age"      : 26,
  "address"  : {
    "streetAddress": "naist street",
    "city"         : "Nara",
    "postalCode"   : "630-0192"
  },
  "phoneNumbers": [
    {
      "type"  : "iPhone",
      "number": "0123-4567-8888"
    },
    {
      "type"  : "home",
      "number": "0123-4567-8910"
    },
    {
      "type"  : "mobile",
      "number": "0913-8532-8492"
    }
  ]
}`

result, err := Jsonpath([]byte(jsonString), "$.phoneNumbers[*].type")
if err != nil {
    log.Fatal(err)
}

for _, item := range result {
    fmt.Println(item)
}
// Output:
// iPhone
// home
// mobile
```

## Supported operations

There are still a few operations which this library does not support but the
goal is to support all of these path operations.

Operation | Supported | Description
---|---|---
$ | Yes | The root object/element.
@ | No | The current object/element. (To be added).
. or [] | Yes | Child operator. Within [] single quotes or double quotes can be used.
.. | Yes | Recursive decent.
* | Yes | Wildcard. All objects/elements regardless of their names.
[] | Yes | subscript operator. XPath uses it to iterate over element collections and for predicates. In Javascript and JSON it is the native array operator. 
[,] | Yes | Union operator in XPath results in a combination of node sets. JSONPath allows alternate names or array indices as a set.
[start:end:step] | No | Array slice operator borrowed from ES4. (To be added).
?() | No | Applies a filter (script) expression. (To be added).
() | No | Script expression, using the underlying script engine. (To be added).