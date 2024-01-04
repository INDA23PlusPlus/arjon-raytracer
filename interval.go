package main

import "math"

type Interval[T Float] struct {
	min, max T
}

func CreateInterval[T Float](min, max T) Interval[T] {
	return Interval[T]{min, max}
}

func CreateIntervalDefault[T Float]() Interval[T] {
	return (*Interval[T])(nil).Empty()
}

func (Interval[T]) Universe() Interval[T] {
	return Interval[T]{T(math.Inf(-1)), T(math.Inf(1))}
}

func (Interval[T]) Empty() Interval[T] {
	return Interval[T]{T(math.Inf(1)), T(math.Inf(-1))}
}

func (i Interval[T]) Contains(t T) bool {
	return t >= i.min && t <= i.max
}

func (i Interval[T]) Surrounds(t T) bool {
	return t > i.min && t < i.max
}

func (i Interval[T]) Clamp(t T) T {
	if t < i.min {
		return i.min
	}
	if t > i.max {
		return i.max
	}
	return t
}
