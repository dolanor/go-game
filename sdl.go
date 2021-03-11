package main

import (
	"go-sdl2/sdl"
)

func initSdl() (win *sdl.Window, rdr *sdl.Renderer, surf *sdl.Surface, cleanupFunc func()) {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	win, rdr, err = sdl.CreateWindowAndRenderer(WINW, WINH, sdl.WINDOW_SHOWN)

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

func delayToLimitFramerate(dt int64) {
	pauseTime := float64(16_333 - dt)
	if pauseTime > 0 {
		delayMilli := uint32(pauseTime / 1_000)
		sdl.Delay(delayMilli)
	}
}
