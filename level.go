package main

import (
	//	"fmt"
	"math"
	"mathlib"

	"go-sdl2/gfx"
	"go-sdl2/sdl"
)

var projMat [4][4]float64
var rotMat [3][3]float64
var axis [3]float64
var camera vec3
var light vec3

func loadLevel(name string) []object {
	//dat := ioutil.ReadFile("/level/" + name + ".txt")
	projMat = mathlib.PerspectiveMat(math.Pi/2, WINW/WINH, 0.1, 100)
	axis = vec3{0.5, 0.5, 0.5}
	camera = vec3{0, 0, +1}
	light = vec3{0, 0, +1}
	cube := initCube()
	cube.update = cubeUpdate
	cube.draw = cubeDraw
	entities := []object{}
	entities = append(entities, cube)
	return entities
}

func initCube() (o object) {
	o.dat = append(o.dat, tri{vert: [3][3]float64{{0, 0, 0}, {0, 1, 0}, {1, 1, 0}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{0, 0, 0}, {1, 1, 0}, {1, 0, 0}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{1, 0, 0}, {1, 1, 0}, {1, 1, 1}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{1, 0, 0}, {1, 1, 1}, {1, 0, 1}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{1, 0, 1}, {1, 1, 1}, {0, 1, 1}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{1, 0, 1}, {0, 1, 1}, {0, 0, 1}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{0, 0, 1}, {0, 1, 1}, {0, 1, 0}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{0, 0, 1}, {0, 1, 0}, {0, 0, 0}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{0, 1, 0}, {1, 1, 1}, {1, 1, 0}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{1, 0, 1}, {0, 0, 1}, {0, 0, 0}}})
	o.dat = append(o.dat, tri{vert: [3][3]float64{{1, 0, 1}, {0, 0, 0}, {1, 0, 0}}})
	return
}

func cubeUpdate(o *object) {
	// apply a rotation to each point
	// tri is [3][3]float64
	rotMat = mathlib.RotationMat(0.01, axis)
	beforeCopy := o.dat[1].vert[2][2]
	for triIdx, tr := range o.dat {
		for vertex := range tr.vert {
			result := mathlib.MultiplyMatVec3(rotMat, tr.vert[vertex])
			o.dat[triIdx].vert[vertex] = result
		}
	}
	afterCopy := o.dat[1].vert[2][2]
	if beforeCopy == afterCopy {
		panic("rotation not applied")
	}
}

func cubeDraw(rdr *sdl.Renderer, o *object) {
	var projected [][3][2]float64
	{
		// sort o.dat ( []tri ) by distance of midpoint to camera
		// calculate triangle midpoint
		// calc distance from camera to midpoint
		// draw farthest triangles first
		// o.dat is []tri
		for trIdx, tr := range o.dat {
			// calculate normal
			normal := mathlib.CrossProductVec3(mathlib.SubtrVec3(tr.vert[1], tr.vert[0]),
				mathlib.SubtrVec3(tr.vert[2], tr.vert[0]))
			normal = mathlib.NormalizeVec3(normal)
			camera = mathlib.NormalizeVec3(camera)
			light = mathlib.NormalizeVec3(light)
			//fmt.Println("normal, camera, light", normal, camera, light)

			// @TODO: check these normals by drawing them
			similarityToCamera := mathlib.DotProductVec3(normal, camera)
			if similarityToCamera < 0 {
				o.dat[trIdx].visible = true
			} else {
				o.dat[trIdx].visible = false
			}
			// calculate whether normal <= 90 degrees with camera
			similarityToLight := mathlib.DotProductVec3(normal, light)
			if similarityToLight < 0 {
				o.dat[trIdx].shade = uint8(-similarityToLight * 255)
			} else {
				o.dat[trIdx].shade = 0
			}

			var projectedTri [3][2]float64
			for vertex := range tr.vert {
				var projTmp vec4
				tr.vert[vertex][2] += 2.0
				// create 1x4 so we can multipy by the projMat
				projTmp = vec4{tr.vert[vertex][0], tr.vert[vertex][1], tr.vert[vertex][2], 1.0}
				// multiply each 1x4 by the projMat
				projTmp = mathlib.MultiplyMatVec4(projMat, projTmp)
				// scale by z depth
				projectedTri[vertex][0] = projTmp[0] / tr.vert[vertex][2]
				projectedTri[vertex][1] = projTmp[1] / tr.vert[vertex][2]
			}
			projected = append(projected, projectedTri)
		}
	}
	for i, screenTri := range projected {
		// scale
		for i := 0; i < 3; i++ {
			screenTri[i][0] += 0.5
			screenTri[i][1] += 0.5
			screenTri[i][0] *= 0.5 * WINW
			screenTri[i][1] *= 0.5 * WINH
		}
		if o.dat[i].visible {
			drawTriangle(rdr, &screenTri, &sdl.Color{R: o.dat[i].shade, G: 100, B: 100, A: 255})
		}
	}
}

func drawTriangle(rdr *sdl.Renderer, tri2d *[3][2]float64, color *sdl.Color) {
	r, g, b, a, err := rdr.GetDrawColor() // get previous draw color
	if err != nil {
		panic(err)
	}
	// convert coordinates to rounded integers
	var tri2dInt [3][2]int32
	tri2dInt[0] = [2]int32{int32(math.Round(tri2d[0][0])), int32(math.Round(tri2d[0][1]))}
	tri2dInt[1] = [2]int32{int32(math.Round(tri2d[1][0])), int32(math.Round(tri2d[1][1]))}
	tri2dInt[2] = [2]int32{int32(math.Round(tri2d[2][0])), int32(math.Round(tri2d[2][1]))}
	// draw wireframe
	rdr.SetDrawColor(100, 100, 100, 255)
	rdr.DrawLine(tri2dInt[0][0], tri2dInt[0][1], tri2dInt[1][0], tri2dInt[1][1])
	rdr.DrawLine(tri2dInt[1][0], tri2dInt[1][1], tri2dInt[2][0], tri2dInt[2][1])
	rdr.DrawLine(tri2dInt[2][0], tri2dInt[2][1], tri2dInt[0][0], tri2dInt[0][1])

	rdr.SetDrawColor(color.R, color.G, color.B, color.A)
	gfx.FilledTrigonRGBA(rdr, tri2dInt[0][0], tri2dInt[0][1],
		tri2dInt[1][0], tri2dInt[1][1],
		tri2dInt[2][0], tri2dInt[2][1], color.R, color.G, color.B, color.A)
	// shade triangle faces
	// change back to previous color
	rdr.SetDrawColor(r, g, b, a)
}
