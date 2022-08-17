package foo

type Struct struct {
	i, j int
	k bool
}

func main() {
	x := []string{
		"1",
		"2",
		"3"
	}

	bar := Struct{
		1,
		2,
		false
	}

	baz := Struct{
		1,
		2,
		false
	}

	// weird edge cases: https://github.com/golang/go/issues/18939
	_ = []int{1, 2
		+3, 4}

}

func Foo(
	int i
) {}

func(Struct) Bar(
	int j
) {}

