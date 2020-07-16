package calculator_test

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/use-testify/calculator"
)

type randomMock struct {
	mock.Mock
}

func (o randomMock) Random(limit int) int {
	args := o.Called(limit)
	return args.Int(0)
}

func TestAdd(t *testing.T) {
	calc := calculator.NewCalculator(nil)
	assert.Equal(t, 9, calc.Add(5, 4))
}

func TestRandom(t *testing.T) {
	rnd := new(randomMock)
	rnd.On("Random", 100).Return(7)
	calc := calculator.NewCalculator(rnd)
	assert.Equal(t, 7, calc.Random())
}
