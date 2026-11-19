package app

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mygo-music/internal/music"
)

// smokeOf returns the bar's smoke for a track whose cover carries tint, with
// the hue and value it came out at.
func smokeOf(t *testing.T, cover string, tint ui.Color, p palette, tr *music.Track) (ui.Color, float64, float64) {
	t.Helper()
	a := &app{art: newArtCache(func(f func()) { f() })}
	a.art.add(sized(cover, 120), &artwork{tint: tint})
	got := a.barSmokeColor(tr, p)
	h, _, v := hsv(float64(got.R)/255, float64(got.G)/255, float64(got.B)/255)
	return got, h, v
}

// The bar shows on every page, so its smoke takes the playing cover's color, as
// Now Playing's does: a blue cover gives blue smoke, rather than the same red
// for every song.
func TestBarSmokeTakesTheCoversColor(t *testing.T) {
	cover := "https://p2.music.126.net/x/1.jpg?param=120y120"
	tr := &music.Track{ID: "blue", Album: music.Album{ImageURL: cover}}

	_, h, _ := smokeOf(t, cover, ui.RGB(40, 90, 220), palette{dark: true}, tr)
	if h < 200 || h > 240 {
		t.Fatalf("hue %.0f, want the cover's blue", h)
	}
}

// The depth is the bar's own, not the cover's: a pale cover would otherwise
// leave white smoke on the light theme's white bar.
func TestBarSmokeKeepsADepthOfItsOwn(t *testing.T) {
	cover := "https://p2.music.126.net/x/2.jpg?param=120y120"
	tr := &music.Track{ID: "pale", Album: music.Album{ImageURL: cover}}

	_, _, dark := smokeOf(t, cover, ui.RGB(235, 240, 250), palette{dark: true}, tr)
	_, _, light := smokeOf(t, cover, ui.RGB(235, 240, 250), palette{}, tr)
	if light > 0.75 {
		t.Fatalf("value %.2f in light mode, too pale to see", light)
	}
	if dark < light {
		t.Fatalf("dark mode smoke is deeper than light mode's: %.2f against %.2f", dark, light)
	}
}

// A cover loads a moment after its song does. Until it has, the last color
// stays: the bar would otherwise flash the accent and back on every song.
func TestBarSmokeHoldsItsColorWhileTheCoverLoads(t *testing.T) {
	cover := "https://p2.music.126.net/x/3.jpg?param=120y120"
	p := palette{dark: true, accent: ui.Oklch(0.66, 0.2, 22)}
	a := &app{art: newArtCache(func(f func()) { f() })}
	first := &music.Track{ID: "first", Album: music.Album{ImageURL: cover}}
	second := &music.Track{ID: "second", Album: music.Album{ImageURL: "https://p2.music.126.net/x/4.jpg?param=120y120"}}

	a.art.add(sized(cover, 120), &artwork{tint: ui.RGB(40, 90, 220)})
	blue := a.barSmokeColor(first, p)

	// The next song's cover isn't in yet.
	if got := a.barSmokeColor(second, p); got != blue {
		t.Fatalf("the bar went to %v while the cover loaded, want the last color %v", got, blue)
	}
}
