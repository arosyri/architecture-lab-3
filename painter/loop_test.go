package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	var testOps []string

	l.Start(mockScreen{})

	l.Post(logOp(t, "do white fill", FillBackground{Color: color.RGBA{255, 255, 255, 255}}))

	l.Post(logOp(t, "do green fill", FillBackground{Color: color.RGBA{0, 128, 0, 255}}))

	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", FillBackground{Color: color.RGBA{0, 128, 0, 255}}))
	}

	l.Post(OperationFunc(func(tx screen.Texture) {
		testOps = append(testOps, "op 1")

		func(tx screen.Texture) {
			testOps = append(testOps, "op 2")
		}(tx)
	}))

	l.Post(OperationFunc(func(tx screen.Texture) {
		testOps = append(testOps, "op 3")
	}))

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type")
	}
	if !reflect.DeepEqual(mt.Colors[0], color.RGBA{255, 255, 255, 255}) {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) < 2 {
		t.Error("Not enough Fill operations:", mt.Colors)
	}

	want := []string{"op 1", "op 2", "op 3"}
	if !reflect.DeepEqual(testOps, want) {
		t.Errorf("Bad order of nested ops: got %v, want %v", testOps, want)
	}
}

func TestLoop_Figure(t *testing.T) {
	var l Loop
	l.Receiver = &testReceiver{}

	l.Start(mockScreen{})

	l.Post(Reset{})

	l.Post(UpdateOp)

	l.Post(DrawT180{
		PosX:  100,
		PosY:  150,
		Size:  50,
		Color: color.RGBA{255, 0, 0, 255},
	})

	l.Post(Move{
		NewPos: image.Pt(300, 400),
	})

	l.Post(UpdateOp)

	l.StopAndWait()

	if len(l.figures) != 1 {
		t.Fatalf("expected 1 figure, got %d", len(l.figures))
	}

	fig := l.figures[0]
	if fig.PosX != 300 || fig.PosY != 400 {
		t.Errorf("figure not moved correctly, got PosX=%d PosY=%d, expected 300 400", fig.PosX, fig.PosY)
	}
}

func logOp(t *testing.T, msg string, op Operation) Operation {
	return OperationFunc(func(tx screen.Texture) {
		t.Log(msg)
		op.Do(tx)
	})
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("not implemented")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("not implemented")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point {
	return image.Pt(800, 800)
}

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	if rgba, ok := src.(color.RGBA); ok {
		m.Colors = append(m.Colors, rgba)
	} else {
		r, g, b, a := src.RGBA()
		m.Colors = append(m.Colors, color.RGBA{
			uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8),
		})
	}
}
