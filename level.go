package main

import (
	"math"
	"mathlib"

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
	camera = vec3{0, 0, 0.5}
	light = vec3{0, 0, +1}
	cube := initCube()
	cube.update = cubeUpdate
	cube.draw = cubeDraw
	entities := []object{}
	entities = append(entities, cube)
	return entities
}

func initCube() (o object) {
	o.dat = append(o.dat, tri{{0, 0, 0}, {0, 1, 0}, {1, 1, 0}}) // Sout
	o.dat = append(o.dat, tri{{0, 0, 0}, {1, 1, 0}, {1, 0, 0}})
	o.dat = append(o.dat, tri{{1, 0, 0}, {1, 1, 0}, {1, 1, 1}}) // East)
	o.dat = append(o.dat, tri{{1, 0, 0}, {1, 1, 1}, {1, 0, 1}})
	o.dat = append(o.dat, tri{{1, 0, 1}, {1, 1, 1}, {0, 1, 1}}) // North
	o.dat = append(o.dat, tri{{1, 0, 1}, {0, 1, 1}, {0, 0, 1}})
	o.dat = append(o.dat, tri{{0, 0, 1}, {0, 1, 1}, {0, 1, 0}}) // West
	o.dat = append(o.dat, tri{{0, 0, 1}, {0, 1, 0}, {0, 0, 0}})
	o.dat = append(o.dat, tri{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}}) // Top)
	o.dat = append(o.dat, tri{{0, 1, 0}, {1, 1, 1}, {1, 1, 0}})
	o.dat = append(o.dat, tri{{1, 0, 1}, {0, 0, 1}, {0, 0, 0}}) // Bottom)
	o.dat = append(o.dat, tri{{1, 0, 1}, {0, 0, 0}, {1, 0, 0}})
	return
}

func cubeUpdate(o *object) {
	// apply a rotation to each point
	// tri is [3][3]float64
	rotMat = mathlib.RotationMat(0.01, axis)
	beforeCopy := o.dat[1][2][2]
	for triIdx, tr := range o.dat {
		for vertex := range tr {
			result := mathlib.MultiplyMatVec3(rotMat, tr[vertex])
			o.dat[triIdx][vertex] = result
		}
	}
	afterCopy := o.dat[1][2][2]
	if beforeCopy == afterCopy {
		panic("rotation not applied")
	}
}

func cubeDraw(rdr *sdl.Renderer, o *object) {
	numTris := len(o.dat)
	if len(o.light) < numTris {
		o.light = make([]uint8, numTris)
	}
	if len(o.visible) < numTris {
		o.visible = make([]bool, numTris)
	}
	var projected [][3][2]float64
	{
		// o.dat is []tri
		for trIdx, tr := range o.dat {
			// calculate normal
			normal := mathlib.CrossProductVec3(mathlib.SubtrVec3(tr[1], tr[0]),
				mathlib.SubtrVec3(tr[2], tr[0]))
			normal = mathlib.NormalizeVec3(normal)
			camera = mathlib.NormalizeVec3(camera)
			light = mathlib.NormalizeVec3(light)
			// @TODO: check these normals by drawing them

			similarityToCamera := mathlib.DotProductVec3(normal, camera)
			if similarityToCamera < 0 {
				o.visible[trIdx] = true
			} else {
				o.visible[trIdx] = false
			}
			// calculate whether normal <= 90 degrees with camera
			similarityToLight := mathlib.DotProductVec3(normal, light)
			if similarityToLight < 0 {
				o.light[trIdx] = uint8(-similarityToLight * 255)
			} else {
				o.light[trIdx] = 0
			}
			// calculate similarity of normal to light

			var projectedTri [3][2]float64
			for vertex := range tr {
				var projTmp vec4
				tr[vertex][2] += 2.0
				// create 1x4 so we can multipy by the projMat
				projTmp = vec4{tr[vertex][0], tr[vertex][1], tr[vertex][2], 1.0}
				// multiply each 1x4 by the projMat
				projTmp = mathlib.MultiplyMatVec4(projMat, projTmp)
				// scale by z depth
				projectedTri[vertex][0] = projTmp[0] / tr[vertex][2]
				projectedTri[vertex][1] = projTmp[1] / tr[vertex][2]
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
		if o.visible[i] {
			drawTriangle(rdr, &screenTri, &sdl.Color{R: o.light[i], G: 100, B: 100, A: 100})
		}
	}
}
