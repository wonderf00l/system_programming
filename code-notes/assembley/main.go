package main

const (
	a = 2
	b
)

var (
	c int
	d = 3
	//m = make(map[string]string, 3)
)

func main() {
	e := a + b
	_ = e

	mm := new(int)
	_ = *mm
}
