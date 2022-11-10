package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/callbacks"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/joypad"
	"github.com/gabe565/gones/internal/ppu"
	"image/color"
	"os"
	"path/filepath"
)

func main() {
	if err := NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func run(path string, callback callbacks.CallbackHandler) error {
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, 3*ppu.Width, 3*ppu.Height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return err
	}

	c, err := console.New(path, func(ppu *ppu.PPU, joypad1 *joypad.Joypad) {
		for button, key := range joypad.Keymap {
			if win.JustPressed(button) {
				joypad1.Set(key, true)
			} else if win.JustReleased(button) {
				joypad1.Set(key, false)
			}
		}

		win.Clear(color.Black)
		pic := ppu.Render()
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))
		win.Update()

		if win.Closed() {
			os.Exit(0)
		}
	})
	if err != nil {
		return err
	}
	c.Reset()

	win.SetTitle(filepath.Base(path) + " | GoNES")

	if callback != nil {
		c.Callback = callback(win)
	}

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
