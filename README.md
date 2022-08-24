# gocommas

<!-- markdownlint-disable no-hard-tabs -->

⚠️ In development ⚠️

A simple tool to find missing trailing commas in Go source code to make it valid.

In several cases, Go requires a trailing comma on the final line before a closing
brace, including:

* composite literals
* function calls
* function signature

Because of this, missing commas are considered invalid syntax, and prevent
most auto-formatting tools from working.

## Example

`example.go`:

```go
func main() {
	x := []string{
		"1",
		"2",
		"3"
	}
	fmt.Println(
		x
	)
}
```

```console
$ gofmt -w example.go
example.go:7:6: missing ',' before newline in composite literal
example.go:10:4: missing ',' before newline in argument list
```

The output file is never written because `gofmt` never got to the point of
formatting it.

## Solution

There [have been proposals](https://github.com/golang/go/issues/18939) to fix this
automatically, but they were rejected because the tools do not have a precedent
for converting invalid syntax into valid syntax.

Well, we can still make external tools that do it! So this is exactly what
`gocommas` does: it reads a potentially invalid Go source file and adds commas
where it thinks they are needed. In several cases, this is enough for the file
to become valid, and it can be passed on to other formatters or the Go compiler.

```console
$ ./gocommas -w  example.go
Missing comma: example.go:7:6
Missing comma: example.go:10:4
```

`example.go`:

```go
package main

func main() {
	x := []string{
		"1",
		"2",
		"3",
	}
	fmt.Println(
		x,
	)
}
```

🎉
Now the file can be built and formatted using your favorite tool!

