package main

import (
	"math"
	"math/rand"
)

func degreesToRadians[T Float](degrees T) T {
	return degrees * (math.Pi / 180.0)
}

// If two arguments are provided, returns a random float in the interval [min, max)
// where min is the first argument and max is the second argument.
// If no arguments are provided, returns a random float in the interval [0, 1).
func random[T Float](interval ...T) T {
	if len(interval) == 2 {
		intervalMin, intervalMax := interval[0], interval[1]
		return intervalMin + (intervalMax-intervalMin)*T(rand.Float64())
	}
	if len(interval) == 0 {
		return T(rand.Float64())
	}
	panic("Expected 0 or 2 arguments")
}

type RayTracer[T Float] struct {
	aspectRatio              T       // Width / Height
	imageWidth, imageHeight  uint    // Width and height of the image in pixels
	viewPortHeight           T       // Height of the viewport in world units
	viewPortWidth            T       // Width of the viewport in world units
	cameraCenter             Vec3[T] // Location of the camera
	pixel00Loc               Vec3[T] // Location of the pixel at (0, 0)
	pixelDeltaU, pixelDeltaV Vec3[T] // Pixel delta vectors
	focalLength              T       // Distance from the camera to the viewport
	samplesPerPixel          uint    // Number of samples per pixel
	maxDepth                 uint    // Maximum number of bounces
	vfov                     T       // Vertical field of view in degrees

	lookFrom, lookAt, vUp Vec3[T] // Camera parameters
}

func makeRayTracer[T Float](imageWidth, imageHeight uint, vfov T, samplesPerPixel uint, lookFrom, lookAt, vUp Vec3[T]) RayTracer[T] {
	if imageWidth <= 0 || imageHeight <= 0 {
		panic("Image width and height must be positive")
	}

	cameraCenter := lookFrom

	// Determine the viewport size
	focalLength := lookFrom.Sub(lookAt).Length()
	theta := degreesToRadians[T](vfov)
	h := T(math.Tan(float64(theta / 2.0)))
	aspectRatio := T(imageWidth) / T(imageHeight)
	viewPortHeight := T(2.0) * h * focalLength
	viewPortWidth := aspectRatio * viewPortHeight

	w := lookFrom.Sub(lookAt).UnitVector()
	u := vUp.Cross(w).UnitVector()
	v := w.Cross(u)

	// Calculate the vectors for the viewport
	viewportU, viewportV := u.Mul(viewPortWidth), v.Mul(-viewPortHeight)

	// Calculate the vectors for the pixels
	pixelDeltaU, pixelDeltaV := viewportU.Div(T(imageWidth)), viewportV.Div(T(imageHeight))

	// Calculate the location of the pixel at (0, 0)
	viewportUpperLeftCorner := cameraCenter.Sub(viewportU.Div(2)).Sub(viewportV.Div(2)).Sub(w.Mul(focalLength))
	pixel00Loc := viewportUpperLeftCorner.Add(pixelDeltaU.Add(pixelDeltaV).Div(2))

	maxDepth := uint(10)

	return RayTracer[T]{aspectRatio, imageWidth, imageHeight, viewPortHeight,
		viewPortWidth, cameraCenter, pixel00Loc, pixelDeltaU, pixelDeltaV,
		focalLength, samplesPerPixel, maxDepth, vfov, lookFrom, lookAt, vUp}
}

func (rt *RayTracer[T]) hitSphere(center Vec3[T], radius T, r *Ray[T]) T {
	oc := r.origin.Sub(center)
	a := r.direction.Dot(r.direction)
	b := 2.0 * oc.Dot(r.direction)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	if discriminant < 0.0 {
		return -1.0
	}
	return (-b - T(math.Sqrt(float64(discriminant)))) / (2.0 * a)
}

func (rt *RayTracer[T]) color(r *Ray[T], hittable Hittable[T], depth uint) Vec3[T] {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth == 0 {
		return Vec3[T]{0.0, 0.0, 0.0}
	}

	var hitRecord HitRecord[T]
	if hittable.hit(r, Interval[T]{0.0001, T(math.Inf(1))}, &hitRecord) {
		var scattered Ray[T]
		var attenuation Vec3[T]
		if hitRecord.mat.scatter(r, &hitRecord, &attenuation, &scattered) {
			return attenuation.ElementWiseMul(rt.color(&scattered, hittable, depth-1))
		}
		return Vec3[T]{0.0, 0.0, 0.0}
	}

	unitDirection := r.direction.UnitVector()
	a := 0.5 * (unitDirection.y + 1.0)
	return Vec3[T]{1.0, 1.0, 1.0}.Mul(1.0 - a).Add(Vec3[T]{0.5, 0.7, 1.0}.Mul(a))
}

func (rt *RayTracer[T]) pixelSampleSquare() Vec3[T] {
	deltaPixel := T(-0.5) + random[T]() // Random float in the interval [-0.5, 0.5)
	return (rt.pixelDeltaU.Mul(deltaPixel)).Add(rt.pixelDeltaV.Mul(deltaPixel))
}

func (rt *RayTracer[T]) getRay(i, j uint) Ray[T] {
	u := T(j) + random[T](-0.5, 0.5)
	v := T(i) + random[T](-0.5, 0.5)

	pixelCenter := rt.pixel00Loc.Add(rt.pixelDeltaU.Mul(u)).Add(rt.pixelDeltaV.Mul(v))
	r := Ray[T]{rt.cameraCenter, pixelCenter.Sub(rt.cameraCenter)}
	return r
}

func (rt *RayTracer[T]) traceImage(world HittableList[T]) [][]Vec3[T] {
	var image [][]Vec3[T]

	for i := uint(0); i < rt.imageHeight; i++ {
		image = append(image, make([]Vec3[T], rt.imageWidth))
		for j := uint(0); j < rt.imageWidth; j++ {
			color := Vec3[T]{0, 0, 0}

			for s := uint(0); s < rt.samplesPerPixel; s++ {
				r := rt.getRay(i, j)
				color = color.Add(rt.color(&r, &world, rt.maxDepth))
			}

			// Divide the color by the number of samples to get the average color
			color = color.Div(T(rt.samplesPerPixel))

			image[i][j] = color
		}
	}

	return image
}

func main() {

	type T = float32

	materialGround := lambertian[T]{Vec3[T]{0.8, 0.8, 0.0}}
	materialCenter := lambertian[T]{Vec3[T]{0.1, 0.2, 0.5}}
	materialLeft := dielectric[T]{1.5}
	materialRight := metal[T]{Vec3[T]{0.8, 0.6, 0.2}, 1.0}

	var world HittableList[T]

	world.add(Sphere[T]{Vec3[T]{0.0, -100.5, -1.0}, 100.0, &materialGround})
	world.add(Sphere[T]{Vec3[T]{0.0, 0.0, -1.0}, 0.5, &materialCenter})
	world.add(Sphere[T]{Vec3[T]{-1.0, 0.0, -1.0}, 0.5, &materialLeft})
	world.add(Sphere[T]{Vec3[T]{-1.0, 0.0, -1.0}, -0.4, &materialLeft})
	world.add(Sphere[T]{Vec3[T]{1.0, 0.0, -1.0}, 0.5, &materialRight})

	const width, height = 800, 400
	rt := makeRayTracer[T](width, height, 20.0, 100, Vec3[T]{-2, 2, 1}, Vec3[T]{0, 0, -1}, Vec3[T]{0, 1, 0})
	image := rt.traceImage(world)
	PpmWriter("test.ppm", image)
}
