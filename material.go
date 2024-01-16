package main

import "math"

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
	fuzz   T
}

func (m *metal[T]) scatter(rIn *Ray[T], rec *HitRecord[T], attenuation *Vec3[T], scattered *Ray[T]) bool {
	reflected := rIn.direction.UnitVector().Reflect(rec.normal)
	*scattered = Ray[T]{rec.p, reflected.Add(randomUnitVec3[T]().Mul(m.fuzz))}
	*attenuation = m.albedo
	return scattered.direction.Dot(rec.normal) > 0
}

// End Metal

// Dielectric

type dielectric[T Float] struct {
	ir T // Index of Refraction
}

func (d *dielectric[T]) scatter(rIn *Ray[T], rec *HitRecord[T], attenuation *Vec3[T], scattered *Ray[T]) bool {
	*attenuation = Vec3[T]{1, 1, 1}
	var refractionRatio T
	if rec.frontFace {
		refractionRatio = 1 / d.ir
	} else {
		refractionRatio = d.ir
	}

	unitDirection := rIn.direction.UnitVector()
	cosTheta := T(math.Min(float64(unitDirection.Mul(-1).Dot(rec.normal)), 1.0))
	sinTheta := T(math.Sqrt(float64(1.0 - cosTheta*cosTheta)))
	cannotRefract := refractionRatio*sinTheta > 1.0

	var direction Vec3[T]

	if cannotRefract || reflectance(cosTheta, refractionRatio) > random[T]() {
		direction = unitDirection.Reflect(rec.normal)
	} else {
		direction = unitDirection.Refract(rec.normal, refractionRatio)
	}

	*scattered = Ray[T]{rec.p, direction}
	return true
}

// End Dielectric

func reflectance[T Float](cosine, refIdx T) T {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*T(math.Pow(float64(1-cosine), 5))
}
