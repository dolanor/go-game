package main

import (
	"go-sdl2/sdl"
	"mathlib"
)

func loadLevel(name string) []object {
	entities := []object{}
	//dat := ioutil.ReadFile("/level/" + name + ".txt")
	cube := initCube()
	cube.update = cubeUpdate
	cube.draw = cubeDraw
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
	var projected [][3][2]float64
	{
		// o.dat is []tri
		for _, tr := range o.dat {
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
	for _, screenTri := range projected {
		// scale
		for i := 0; i < 3; i++ {
			screenTri[i][0] += 0.5
			screenTri[i][1] += 0.5
			screenTri[i][0] *= 0.5 * WINW
			screenTri[i][1] *= 0.5 * WINH
		}
		drawTriangle(rdr, &screenTri, &sdl.Color{100, 100, 100, 100})
	}
}
