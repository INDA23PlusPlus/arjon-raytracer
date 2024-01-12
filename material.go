package main

type Material[T Float] interface {
	scatter(rIn *Ray[T], rec *HitRecord[T], attenuation *Vec3[T], scattered *Ray[T]) bool
}

// Lambertian
type lambertian[T Float] struct {
	albedo Vec3[T]
}

func (l *lambertian[T]) scatter(rIn *Ray[T], rec *HitRecord[T], attenuation *Vec3[T], scattered *Ray[T]) bool {
	scatterDirection := rec.normal.Add(randomUnitVec3[T]())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.normal
	}

	*scattered = Ray[T]{rec.p, scatterDirection}
	*attenuation = l.albedo
	return true
}

// End Lambertian

// Metal
type metal[T Float] struct {
	albedo Vec3[T]
}

func (m *metal[T]) scatter(rIn *Ray[T], rec *HitRecord[T], attenuation *Vec3[T], scattered *Ray[T]) bool {
	reflected := rIn.direction.UnitVector().Reflect(rec.normal)
	*scattered = Ray[T]{rec.p, reflected}
	*attenuation = m.albedo
	return true
}

// End Metal
