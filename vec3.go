package main

import (
	"fmt"
	"math"
	"math/rand"
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

// RandomVec3[T Float] returns a random vector in the interval [0, 1) if no interval is specified
// Otherwise, it returns a random vector in the interval [min, max)
func randomVec3[T Float](interval ...T) Vec3[T] {
	x, y, z := random[T](interval...), random[T](interval...), random[T](interval...)
	return Vec3[T]{x, y, z}
}

/*
func randomInUnitSphere[T Float]() Vec3[T] {
	for {
		p := randomVec3[T](-1, 1)
		length := p.Length()
		if length < 1 && length != 0 {
			return p
		}
	}
}

func randomUnitVec3[T Float]() Vec3[T] {
	return randomInUnitSphere[T]().UnitVector()
}
*/

// RandomUnitVec3[T Float] returns a random and evenly distributed unit vector
func randomUnitVec3[T Float]() Vec3[T] {
	// Generates a random unit vector in the unit sphere
	// See https://mathworld.wolfram.com/SpherePointPicking.html
	for {
		x := rand.NormFloat64()
		y := rand.NormFloat64()
		z := rand.NormFloat64()
		v := Vec3[T]{T(x), T(y), T(z)}

		// Reject the vector if it is zero (extremely unlikely), since we can't normalize it
		if v.Length() != 0 {
			return v.UnitVector()
		}
	}
}

func randomInHemisphere[T Float](normal Vec3[T]) Vec3[T] {
	inUnitSphere := randomUnitVec3[T]()
	if inUnitSphere.Dot(normal) > 0.0 {
		return inUnitSphere
	} else {
		return inUnitSphere.Mul(-1)
	}
}

func (v Vec3[T]) NearZero() bool {
	const s = 1e-8
	return math.Abs(float64(v.x)) < s && math.Abs(float64(v.y)) < s && math.Abs(float64(v.z)) < s
}

func (v Vec3[T]) Reflect(normal Vec3[T]) Vec3[T] {
	return v.Sub(normal.Mul(2 * v.Dot(normal)))
}

func (v Vec3[T]) ElementWiseMul(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v.x * v2.x, v.y * v2.y, v.z * v2.z}
}
