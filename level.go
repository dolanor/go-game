package main

import "mathlib"

func loadLevel(name string) []object {
	entities := []object{}
	//dat := ioutil.ReadFile("/level/" + name + ".txt")
	cube := initCube()
	cube.update = cubeUpdate
	entities = append(entities, initCube())
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
	rotMat = mathlib.RotationMat(rotation, axis)
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
