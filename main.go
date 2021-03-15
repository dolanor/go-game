package main

import (
	"frametimer"
	"strconv"

	"go-sdl2/sdl"
)

var ent []object

//var levelList []string

func init() {
	ent = loadLevel("test")
}

const (
	// WINW = Window Width
	WINW = 600
	// WINH = Window Height
	WINH = 600
)

func main() {
	win, rdr, surf, cleanup := initSdl(WINW, WINH)
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
		_ = timer.RecordTime()
		if timer.TotalFrames%50 == 0 {
			fps := int(timer.CalcFPS())
			win.SetTitle(strconv.Itoa(fps))
		}
	}
}
