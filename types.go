package main

type tri = [3][3]float64
type vec3 = [3]float64
type vec4 = [4]float64

type material struct {
}

type object struct {
	dat    []tri
	mat    material
	update func(*object)
	id     int64
}
