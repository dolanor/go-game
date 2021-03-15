package main

import (
	"go-sdl2/sdl"
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
