package painter

import (
	"golang.org/x/exp/shiny/imageutil"
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Operation interface {
	Do(t screen.Texture) bool
}

type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (updateOp) Do(t screen.Texture) bool { return true }

type OperationFunc func(screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

type FillBackground struct {
	Color color.RGBA
}

func (op FillBackground) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), op.Color, screen.Src)
	return false
}

type BgRect struct {
	Rect image.Rectangle
}

func (op BgRect) Do(t screen.Texture) bool {
	t.Fill(op.Rect, color.Black, screen.Src)
	return false
}

type DrawT180 struct {
	PosX, PosY int
	Size       int
	Color      color.RGBA
}

func (op DrawT180) Do(t screen.Texture) bool {
	b := t.Bounds()

	cx := op.PosX
	cy := op.PosY
	if cx == 0 || cy == 0 {
		cx = b.Dx() / 2
		cy = b.Dy() / 2
	}

	thickness := op.Size / 5
	half := op.Size / 2

	horzRect := image.Rect(cx-half, cy-thickness/2, cx+half, cy+thickness/2)
	t.Fill(horzRect, op.Color, screen.Src)

	vertRect := image.Rect(cx-thickness/2, cy-half, cx+thickness/2, cy+thickness/2)
	t.Fill(vertRect, op.Color, screen.Src)

	return false
}

type Border struct {
	Thickness int
	Color     color.Color
}

func (op Border) Do(t screen.Texture) bool {
	bounds := t.Bounds()
	borders := imageutil.Border(bounds, op.Thickness)
	for _, r := range borders {
		t.Fill(r, op.Color, screen.Src)
	}
	return false
}

type Reset struct{}

func (Reset) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.Black, screen.Src)
	return false
}

type Move struct {
	NewPos image.Point
}

func (Move) Do(t screen.Texture) bool {
	return false
}

var (
	WhiteFill = OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color.White, screen.Src)
	})
	GreenFill = OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color.RGBA{0, 128, 0, 255}, screen.Src)
	})
)
