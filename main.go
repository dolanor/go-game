package main

import (
	"frametimer"
	"strconv"

	"go-sdl2/sdl"
)

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

var (
	// screenRect is the screenRectangle
	screenRect sdl.Rect
	// ent is the entity vector
	ent []object
)

func main() {
	// init objects
	win, rdr, _, cleanup := initSdl(WINW, WINH)
	defer cleanup()
	timer := frametimer.Timer{}
	screenRect = sdl.Rect{0, 0, WINW, WINH}
	clearScreen(rdr)
	rdr.Present()

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
		clearScreen(rdr)
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
