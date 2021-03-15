package main

import "go-sdl2/sdl"

type tri = [3][3]float64
type vec3 = [3]float64
type vec4 = [4]float64

type material struct {
	color sdl.Color
}

type object struct {
	dat    []tri
	light  []uint8
	mat    material
	pos    vec3
	update func(*object)
	draw   func(*sdl.Renderer, *object)
	id     int64
}
