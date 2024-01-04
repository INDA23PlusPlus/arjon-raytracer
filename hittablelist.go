package main

// HittableList[T Float] has the interface of Hittable[T]
type HittableList[T Float] struct {
	objects []Hittable[T]
}

func (hl *HittableList[T]) add(object Hittable[T]) {
	hl.objects = append(hl.objects, object)
}

func (hl *HittableList[T]) clear() {
	hl.objects = nil
}

func (hl *HittableList[T]) hit(r *Ray[T], rayT Interval[T], rec *HitRecord[T]) bool {
	var tempRec HitRecord[T]
	hitAnything := false
	closestSoFar := rayT.max

	for _, object := range hl.objects {
		if object.hit(r, Interval[T]{rayT.min, closestSoFar}, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t
			*rec = tempRec
		}
	}

	return hitAnything
}
