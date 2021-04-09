package main

import (
	"flag"
	"fmt"
	"frametimer"
	"log"
	"mathlib"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"

	"go-sdl2/sdl"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

//var levelList []string

func init() {
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
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// init objects
	win, rdr, cleanup := initSdl(WINW, WINH)
	defer cleanup()
	printRenderInfo(rdr)
	ent = loadLevel("test")
	timer := frametimer.Timer{}
	screenRect = sdl.Rect{X: 0, Y: 0, W: WINW, H: WINH}
	clearScreen(rdr)
	rdr.Present()
	running := true
	timer.RecordTime()
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
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
			if ent[i].update != nil {
				ent[i].update(&ent[i])
			}
		}
		// draw
		clearScreen(rdr)
		// sort ent by z value before drawing
		// TODO: Verify that this is really working and not just a comedy of errors
		// where the midpt and this are both broken
		sort.Slice(ent, func(i, j int) bool {
			imid := ent[i].midpoint()
			jmid := ent[j].midpoint()
			idistToCamera := mathlib.DistVec3(imid, camera)
			jdistToCamera := mathlib.DistVec3(jmid, camera)
			return idistToCamera < jdistToCamera
		})
		for i := range ent {
			// sort o.mesh ( []tri ) by distance of midpoint to camera
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

func initSdl(width, height int32) (win *sdl.Window, rdr *sdl.Renderer, cleanupFunc func()) {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	win, rdr, err = sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		defer sdl.Quit()
		defer win.Destroy()
	}
	return win, rdr, cleanup
}

func clearScreen(rdr *sdl.Renderer) {
	rdr.SetDrawColor(10, 200, 200, 255)
	rdr.FillRect(&screenRect)
}

func printRenderInfo(rdr *sdl.Renderer) {
	renderInfo, err := rdr.GetInfo()
	if err == nil {
		fmt.Printf("render info: %v\n", renderInfo.Flags)
		i := 0x00000002 | 0x00000008
		fmt.Printf("flags target texture and hardware accel: %v\n", i)
	}
}
