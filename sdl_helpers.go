package main

import (
	"fmt"
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
