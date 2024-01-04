package main

import "math"

// Sphere[T Float] has the interface of Hittable[T]
type Sphere[T Float] struct {
	center Vec3[T]
	radius T
}

func (s Sphere[T]) hit(r *Ray[T], rayT Interval[T], rec *HitRecord[T]) bool {
	oc := r.origin.Sub(s.center)
	a := r.direction.SquaredLength()
	b := oc.Dot(r.direction)
	c := oc.SquaredLength() - s.radius*s.radius
	discriminant := b*b - a*c
	if discriminant < 0 {
		return false
	}

	temp := (-b - T(math.Sqrt(float64(discriminant)))) / a
	if !rayT.Surrounds(temp) {
		temp = (-b + T(math.Sqrt(float64(discriminant)))) / a
		if !rayT.Surrounds(temp) {
			return false
		}
	}

	rec.t = temp
	rec.p = r.PointAtParameter(rec.t)
	outwardNormal := rec.p.Sub(s.center).Div(s.radius)
	rec.setFaceNormal(r, outwardNormal)

	return true
}
