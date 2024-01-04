package main

type Hittable[T Float] interface {
	hit(r *Ray[T], rayT Interval[T], rec *HitRecord[T]) bool
}
