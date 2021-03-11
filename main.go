package main

import (
	//	"fmt"
	"frametimer"
	"go-sdl2/sdl"
	"math"
	"mathlib"
	"strconv"
)

var projMat [4][4]float64
var rotMat [3][3]float64
var rotation float64
var axis [3]float64

func init() {
	projMat = mathlib.PerspectiveMat(math.Pi/2, WINW/WINH, 0.1, 100) // correct
	axis = vec3{0.5, 0.5, 0.5}
	rotation = 0.005
	rotMat = mathlib.RotationMat(0.01, axis)
}

const (
	// WINW = Window Width
	WINW = 800
	// WINH = Window Height
	WINH = 600
)

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

func main() {
	win, rdr, surf, cleanup := initSdl()
	defer cleanup()

	surf.FillRect(nil, 0)
	// init objects
	cube := initCube()
	timer := frametimer.Timer{}
	win.UpdateSurface()

	running := true
	timer.RecordTime()
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				if event.GetType() == sdl.KEYDOWN {
					// 79 right 80 left 81 down 82 up
					switch e.Keysym.Scancode {
					case sdl.GetScancodeFromKey(sdl.K_RIGHT):
					case sdl.GetScancodeFromKey(sdl.K_LEFT):
					case sdl.GetScancodeFromKey(sdl.K_DOWN):
					case sdl.GetScancodeFromKey(sdl.K_UP):
					}
				}
			}
		}
		// update
		update(&cube)
		// draw
		rect := sdl.Rect{X: 0, Y: 0, W: WINW, H: WINH}
		rdr.SetDrawColor(10, 200, 10, 100)
		rdr.FillRect(&rect)
		draw(rdr, cube)
		rdr.Present()
		delayMicro := timer.GetElapsedSinceLast()
		delayToLimitFramerate(delayMicro)
		_ = timer.RecordTime()
		if timer.TotalFrames%50 == 0 {
			fps := int(timer.CalcFPS())
			win.SetTitle(strconv.Itoa(fps))
		}
	}
}

func update(o *object) {
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

func draw(rdr *sdl.Renderer, o object) {
	r, g, b, a, err := rdr.GetDrawColor() // get previous draw color
	if err != nil {
		panic(err)
	}
	rdr.SetDrawColor(100, 100, 100, 100)
	oProjected := calcObjectProjection(&o)
	for _, screenTri := range oProjected {
		// scale
		for i := 0; i < 3; i++ {
			screenTri[i][0] += 0.4
			screenTri[i][1] += 0.4
			screenTri[i][0] *= 0.5 * WINW
			screenTri[i][1] *= 0.5 * WINH
		}
		rdr.DrawLine(int32(math.Round(screenTri[0][0])), int32(math.Round(screenTri[0][1])),
			int32(math.Round(screenTri[1][0])), int32(math.Round(screenTri[1][1])))
		rdr.DrawLine(int32(math.Round(screenTri[1][0])), int32(math.Round(screenTri[1][1])),
			int32(math.Round(screenTri[2][0])), int32(math.Round(screenTri[2][1])))
		rdr.DrawLine(int32(math.Round(screenTri[2][0])), int32(math.Round(screenTri[2][1])),
			int32(math.Round(screenTri[0][0])), int32(math.Round(screenTri[0][1])))
	}
	// change the color back
	rdr.SetDrawColor(r, g, b, a)
}

// @TODO: Might want to clean up these temporaries for speed's sake
func calcObjectProjection(o *object) [][3][2]float64 {
	// o.dat is []tri
	var projected [][3][2]float64
	for _, tr := range o.dat {
		// create three 1x4s so we can multipy by the projMat
		var projectedTri [3][2]float64
		for vertex := range tr {
			var projTmp vec4
			projTmp = vec4{tr[vertex][0], tr[vertex][1], tr[vertex][2], 1.0}
			// multiply each 1x4 by the projMat and store the result
			projTmp = mathlib.MultiplyMatByVec4(projMat, projTmp)
			projectedTri[vertex][0], projectedTri[vertex][1] = projTmp[0], projTmp[1]
		}
		//fmt.Printf("%10.6f %10.6f\n", projectedTri[2][0], projectedTri[2][1])
		projected = append(projected, projectedTri)
	}
	if len(projected) != len(o.dat) {
		panic("projected != o.dat")
	}
	return projected
}
