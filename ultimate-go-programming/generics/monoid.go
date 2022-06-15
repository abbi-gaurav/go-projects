package main

type monoid[T any] interface {
	op(x, y T) T
	zero() T
}
