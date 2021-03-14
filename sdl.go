package main

import (
	"go-sdl2/sdl"
	"math"
)

func initSdl(width, height int32) (win *sdl.Window, rdr *sdl.Renderer, surf *sdl.Surface, cleanupFunc func()) {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	win, rdr, err = sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	if surf, err = win.GetSurface(); err != nil {
		panic(err)
	}
	cleanup := func() {
		defer sdl.Quit()
		defer win.Destroy()
	}
	return win, rdr, surf, cleanup
}

func delayToLimitFramerate(frametime int64) {
	pauseTime := float64(16_333 - frametime)
	if pauseTime > 0 {
		delayMilli := uint32(pauseTime / 1_000)
		sdl.Delay(delayMilli)
	}
}

func drawTriangle(rdr *sdl.Renderer, tri2d *[3][2]float64, color *sdl.Color) {
	r, g, b, a, err := rdr.GetDrawColor() // get previous draw color
	if err != nil {
		panic(err)
	}
	rdr.SetDrawColor(color.R, color.G, color.B, color.A)
	rdr.DrawLine(int32(math.Round(tri2d[0][0])), int32(math.Round(tri2d[0][1])),
		int32(math.Round(tri2d[1][0])), int32(math.Round(tri2d[1][1])))
	rdr.DrawLine(int32(math.Round(tri2d[1][0])), int32(math.Round(tri2d[1][1])),
		int32(math.Round(tri2d[2][0])), int32(math.Round(tri2d[2][1])))
	rdr.DrawLine(int32(math.Round(tri2d[2][0])), int32(math.Round(tri2d[2][1])),
		int32(math.Round(tri2d[0][0])), int32(math.Round(tri2d[0][1])))
	rdr.SetDrawColor(r, g, b, a)
}
