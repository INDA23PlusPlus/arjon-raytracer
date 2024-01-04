package main

type Ray[T Float] struct {
	origin, direction Vec3[T]
}

func (r Ray[T]) PointAtParameter(t T) Vec3[T] {
	return r.origin.Add(r.direction.Mul(t))
}
