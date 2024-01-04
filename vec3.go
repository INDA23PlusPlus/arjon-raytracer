package main

import (
	"fmt"
	"math"
)

type Float = interface{ float64 | float32 }

type Vec3[T Float] struct {
	x, y, z T
}

type Vec3f = Vec3[float32]
type Vec3d = Vec3[float64]

type Point3f = Vec3f
type Point3d = Vec3d

func (v Vec3[T]) Add(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v.x + v2.x, v.y + v2.y, v.z + v2.z}
}

func (v Vec3[T]) Sub(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v.x - v2.x, v.y - v2.y, v.z - v2.z}
}

func (v Vec3[T]) Mul(t T) Vec3[T] {
	return Vec3[T]{v.x * t, v.y * t, v.z * t}
}

func (v Vec3[T]) Div(t T) Vec3[T] {
	if t == 0 {
		panic("Divide by zero")
	}
	return Vec3[T]{v.x / t, v.y / t, v.z / t}
}

func (v Vec3[T]) Dot(v2 Vec3[T]) T {
	return v.x*v2.x + v.y*v2.y + v.z*v2.z
}

func (v Vec3[T]) Cross(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v.y*v2.z - v.z*v2.y, v.z*v2.x - v.x*v2.z, v.x*v2.y - v.y*v2.x}
}

func (v Vec3[T]) SquaredLength() T {
	return v.Dot(v)
}

func (v Vec3[T]) Length() T {
	return T(math.Sqrt(float64(v.SquaredLength())))
}

func (v Vec3[T]) UnitVector() Vec3[T] {
	return v.Div(v.Length())
}

func (v Vec3[T]) PpmString() string {
	const factor = 255.99
	r, g, b := int(factor*v.x), int(factor*v.y), int(factor*v.z)
	return fmt.Sprintf("%d %d %d\n", r, g, b)
}
