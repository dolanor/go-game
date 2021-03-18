package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"mathlib"
	"os"
	"sort"
	"strconv"
	"strings"

	"go-sdl2/gfx"
	"go-sdl2/sdl"
)

var projMat [4][4]float64
var rotMat [3][3]float64
var axis [3]float64
var camera vec3
var light vec3
var wireframeFlag bool

func loadLevel(name string) []object {
	//dat := ioutil.ReadFile("/level/" + name + ".txt")
	projMat = mathlib.PerspectiveMat(math.Pi/2, WINW/WINH, 0.1, 100)
	axis = vec3{0.5, 0.5, 0.5}
	camera = vec3{0, 0.0, +1}
	light = vec3{0, 0, 1}
	//cube := initCube()
	obj, err := newObjectFromFile("icosphere")
	if err != nil {
		panic(err)
	}
	obj.update = cubeUpdate
	obj.draw = drawObject
	entities := []object{}
	entities = append(entities, obj)
	return entities
}

// newObjectFromFile loads the triangle vertices of a .obj file
// # is a comment
// v is a vertex position
// vn is a normal
// vt is a texture coordinate
func newObjectFromFile(name string) (o object, e error) {
	// verify that file exists
	fname := `.\objects\` + name + ".obj"
	if _, e = os.Stat(fname); os.IsNotExist(e) {
		fmt.Printf("file %v does not exist\n", fname)
		return o, e
	}
	dat, e := ioutil.ReadFile(fname)
	fmt.Printf("read in %v bytes from %v\n", len(dat), fname)
	if e != nil {
		panic(e)
	}
	// load in vertices
	// load in faces
	// do error checks
	var vertices [][3]float64
	var faces [][3]int
	lines := strings.Split(string(dat), "\n")
	for i := range lines {
		fields := strings.Split(lines[i], " ")
		switch fields[0] {
		case "v":
			// reading vertex
			if len(fields) != 4 {
				return o, errors.New("vertex field count != 4")
			}
			f1, err1 := strconv.ParseFloat(fields[1], 64)
			f2, err2 := strconv.ParseFloat(fields[2], 64)
			f3, err3 := strconv.ParseFloat(fields[3], 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return o, errors.New(fmt.Sprintf("error parsing float from vertex on line %v\n", i))
			}
			vertices = append(vertices, [3]float64{f1, f2, f3})
		case "f":
			// reading face
			//	i, err := strconv.ParseInt("-42", 10, 64)
			if len(fields) != 4 {
				return object{}, errors.New("face field count != 4")
			}
			i1, err1 := strconv.ParseInt(fields[1], 10, 64)
			i2, err2 := strconv.ParseInt(fields[2], 10, 64)
			i3, err3 := strconv.ParseInt(fields[3], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return o, errors.New(fmt.Sprintf("error parsing int from face on line %v\n", i))
			}
			faces = append(faces, [3]int{int(i1), int(i2), int(i3)})
		}
	}
	fmt.Printf("%v vertices. %v faces. vertices / faces= %v\n",
		len(vertices), len(faces), float64(len(faces))/float64(len(vertices)))
	// create triangles from face #s
	var tmpTri tri
	for i := range faces {
		vIdx1, vIdx2, vIdx3 := faces[i][0]-1, faces[i][1]-1, faces[i][2]-1
		tmpTri.vert[0] = [3]float64{vertices[vIdx1][0], vertices[vIdx1][1], vertices[vIdx1][2]}
		tmpTri.vert[1] = [3]float64{vertices[vIdx2][0], vertices[vIdx2][1], vertices[vIdx2][2]}
		tmpTri.vert[2] = [3]float64{vertices[vIdx3][0], vertices[vIdx3][1], vertices[vIdx3][2]}
		o.dat = append(o.dat, tmpTri)
	}
	return
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
	rotMat = mathlib.RotationMat(0.0095, axis)

	for triIdx, tr := range o.dat {
		for vertex := range tr.vert {
			result := mathlib.MultiplyMatVec3(rotMat, tr.vert[vertex])
			o.dat[triIdx].vert[vertex] = result
		}
	}
}

func drawObject(rdr *sdl.Renderer, o *object) {
	var projectedTriangles [][3][2]float64
	{
		if len(o.dat) == 0 {
			fmt.Println("drawing object with 0 triangles")
		}
		// sort o.dat ( []tri ) by distance of midpoint to camera
		sort.SliceStable(o.dat, func(i, j int) bool {
			imid := mathlib.MidpointTri(o.dat[i].vert)
			jmid := mathlib.MidpointTri(o.dat[j].vert)
			idistToCamera := mathlib.DistVec3(imid, camera)
			jdistToCamera := mathlib.DistVec3(jmid, camera)
			return idistToCamera > jdistToCamera
		})

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
				// translate backward
				tr.vert[vertex][2] += 0.75
				// create 1x4 so we can multipy by the projMat
				projTmp = vec4{tr.vert[vertex][0], tr.vert[vertex][1], tr.vert[vertex][2], 1.0}
				// multiply each 1x4 by the projMat
				projTmp = mathlib.MultiplyMatVec4(projMat, projTmp)
				// scale by z depth
				projectedTri[vertex][0] = projTmp[0] / tr.vert[vertex][2]
				projectedTri[vertex][1] = projTmp[1] / tr.vert[vertex][2]
			}
			projectedTriangles = append(projectedTriangles, projectedTri)
		}
	}
	for i, screenTri := range projectedTriangles {
		// scale
		for i := 0; i < 3; i++ {
			screenTri[i][0] += 0.5
			screenTri[i][1] += 0.5
			screenTri[i][0] *= 0.5 * WINW
			screenTri[i][1] *= 0.5 * WINH
		}
		if o.dat[i].visible {
			// shade triangle faces
			RenderProjectedTri(rdr, &screenTri, &sdl.Color{R: o.dat[i].shade, G: 100, B: 100, A: 255})
		}
	}
}

// RenderProjectedTri draws a 2d triangle
func RenderProjectedTri(rdr *sdl.Renderer, tri2d *[3][2]float64, color *sdl.Color) {
	r, g, b, a, err := rdr.GetDrawColor() // get previous draw color
	if err != nil {
		panic(err)
	}
	// convert coordinates to rounded integers
	var tri2dInt [3][2]int32
	tri2dInt[0] = [2]int32{int32(math.Round(tri2d[0][0])), int32(math.Round(tri2d[0][1]))}
	tri2dInt[1] = [2]int32{int32(math.Round(tri2d[1][0])), int32(math.Round(tri2d[1][1]))}
	tri2dInt[2] = [2]int32{int32(math.Round(tri2d[2][0])), int32(math.Round(tri2d[2][1]))}

	if wireframeFlag {
		// draw wireframe
		rdr.SetDrawColor(255, 255, 255, 255)
		rdr.DrawLine(tri2dInt[0][0], tri2dInt[0][1], tri2dInt[1][0], tri2dInt[1][1])
		rdr.DrawLine(tri2dInt[1][0], tri2dInt[1][1], tri2dInt[2][0], tri2dInt[2][1])
		rdr.DrawLine(tri2dInt[2][0], tri2dInt[2][1], tri2dInt[0][0], tri2dInt[0][1])
	}
	// draw filled triangle
	rdr.SetDrawColor(color.R, color.G, color.B, color.A)
	gfx.FilledTrigonRGBA(rdr, tri2dInt[0][0], tri2dInt[0][1],
		tri2dInt[1][0], tri2dInt[1][1],
		tri2dInt[2][0], tri2dInt[2][1], color.R, color.G, color.B, color.A)
	// change back to previous color
	rdr.SetDrawColor(r, g, b, a)
}
