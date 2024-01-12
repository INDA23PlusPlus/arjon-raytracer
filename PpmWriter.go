package main

import (
	"fmt"
	"math"
	"os"
)

func IsEmpty[T any](array [][]T) bool {
	if len(array) == 0 {
		return true
	}
	for _, row := range array {
		if len(row) != 0 {
			return false
		}
	}
	return true
}

func IsRectangular[T any](array [][]T) bool {
	if len(array) == 0 {
		return true
	}

	length := len(array[0])
	for _, row := range array {
		if len(row) != length {
			return false
		}
	}
	return true
}

func PpmWriter[T Float](filename string, image [][]Vec3[T]) {
	if IsEmpty(image) {
		panic("Image is empty")
	}
	if !IsRectangular(image) {
		panic("Image is not rectangular")
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Close file when function returns
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.WriteString("P3\n")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(fmt.Sprintf("%d %d\n", len(image[0]), len(image)))
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString("255\n")
	if err != nil {
		panic(err)
	}

	// Values should already be in the range [0, 1], but clamp them just in case
	interval := Interval[T]{0, 1}

	for _, row := range image {
		for _, pixel := range row {
			r := int(255.999 * interval.Clamp(linearToGamma(pixel.x)))
			g := int(255.999 * interval.Clamp(linearToGamma(pixel.y)))
			b := int(255.999 * interval.Clamp(linearToGamma(pixel.z)))
			_, err = file.WriteString(fmt.Sprintf("%d %d %d\n", r, g, b))
			if err != nil {
				panic(err)
			}
		}
	}
}

func linearToGamma[T Float](x T) T {
	return T(math.Sqrt(float64(x)))
}
