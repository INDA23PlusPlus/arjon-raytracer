package main

type HitRecord[T Float] struct {
	p, normal Vec3[T]
	t         T
	frontFace bool
	mat       Material[T]
}

func (rec *HitRecord[T]) setFaceNormal(r *Ray[T], outwardNormal Vec3[T]) {
	rec.frontFace = r.direction.Dot(outwardNormal) < 0
	if rec.frontFace {
		rec.normal = outwardNormal
	} else {
		rec.normal = outwardNormal.Mul(-1)
	}
}
