package main

import (
	"fmt"
	"frametimer"
	"go-sdl2/sdl"
	"math"
	"mathlib"
	"strconv"
)

var projMat [4][4]float64
var rotMat [3][3]float64
var axis [3]float64
var ent []object

func init() {
	projMat = mathlib.PerspectiveMat(math.Pi/2, WINW/WINH, 0.1, 100) // correct
	axis = vec3{0.5, 0.5, 0.5}
	rotMat = mathlib.RotationMat(0.01, axis)
	ent = loadLevel("test")
	fmt.Println("entities:", ent)
}

const (
	// WINW = Window Width
	WINW = 700
	// WINH = Window Height
	WINH = 700
)

func main() {
	win, rdr, surf, cleanup := initSdl()
	defer cleanup()

	surf.FillRect(nil, 0)
	// init objects

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
					case sdl.GetScancodeFromKey(sdl.K_ESCAPE):
						running = false
					}
				}
			}
		}
		// update
		for i := range ent {
			ent[i].update(&ent[i])
		}
		// draw
		rect := sdl.Rect{X: 0, Y: 0, W: WINW, H: WINH}
		rdr.SetDrawColor(10, 200, 200, 100)
		rdr.FillRect(&rect)
		for i := range ent {
			ent[i].draw(rdr, &ent[i])
		}
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
			projTmp = mathlib.MultiplyMatVec4(projMat, projTmp)
			projectedTri[vertex][0], projectedTri[vertex][1] = projTmp[0], projTmp[1]
		}
		projected = append(projected, projectedTri)
	}
	if len(projected) != len(o.dat) {
		panic("projected != o.dat")
	}
	return projected
}
