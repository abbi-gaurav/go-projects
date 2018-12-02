package calculator

type Random interface {
	Random(limit int) int
}

type calc struct {
	rnd Random
}

func NewCalculator(rnd Random) calc {
	return calc{
		rnd: rnd,
	}
}

func (c calc) Add(x, y int) int {
	return x + y
}

func (c calc) Random() int {
	return c.rnd.Random(100)
}
