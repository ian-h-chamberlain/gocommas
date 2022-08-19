package foo

type Struct struct {
	i, j int
	k bool
}

type Intf interface {
	DoSomething()
}

func main() {
	x := []string{
		"1",
		"2",
		"3",
	}

	bar := Struct{
		1,
		2,
		false,
	}

	baz := Struct{
		1,
		2,
		false,
	}

	// weird edge cases: https://github.com/golang/go/issues/18939
	_ = []int{1, 2
		+3, 4}

	foo(
		1,
		2,
		3,
	)

}

func Foo(
	int i,
) (
	bool,
) {}

func Baz[
	X any,
	Y any,
](
	int i,
) (
	bool,
) {}

func(
	Struct,
) Bar(
	int j,
) (
	bool,
) {}

