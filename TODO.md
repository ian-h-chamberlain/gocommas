# Stuff to do before this is actually usable

## Use cases

* Function definition (parameter / return value list)
* Function calls
* ???

## Features / usability

* Proper CLI interface
* Run on package / module? recursive find files kinda thing
* `golines` integration? run `golines` automatically if present or something?
* `gofmt` automatically after running
* Unit / acceptance testing rather than plain executable
* CI pipeline of some sort

## Test cases

> NOTE: investigate <https://pkg.go.dev/golang.org/x/tools@v0.1.12/go/expect>,
> but a more classic diff testing might be a better approach in this case since
> some test code is probably gonna be unparsable

* Read through <https://github.com/golang/go/issues/18939> and see if there's
  anything that stands out there.

  e.g. unclosed braces in composite literals, cases like below, etc.

  ```go
  []int{1, 2
      +3, 4}
  ```

* Unparseable code. See how resilient the "parser.ParseFile"
* UTF-8 code (string literals, var names)
* Comments after thing that needs trailing comma
