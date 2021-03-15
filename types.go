package main

import "go-sdl2/sdl"

type vec3 = [3]float64
type vec4 = [4]float64

type tri struct {
	vert    [3][3]float64
	shade   uint8
	color   sdl.Color
	visible bool
}

type object struct {
	dat    []tri
	update func(*object)
	draw   func(*sdl.Renderer, *object)
	id     int64
}
