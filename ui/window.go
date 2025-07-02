package ui

import (
	"image"
	"image/color"
	"log"
	"sync"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz size.Event
	mu sync.Mutex

	bgColor    color.Color
	figurePos  image.Point
	figureSize int

	currentTexture screen.Texture

	OnMove func(p image.Point)
}

func (v *Visualizer) Update(t screen.Texture) {
	log.Println("Visualizer: received texture update")
	v.mu.Lock()
	v.currentTexture = t
	v.mu.Unlock()

	v.tx <- t
}

func (v *Visualizer) Main() {
	v.tx = make(chan screen.Texture)
	v.done = make(chan struct{})
	driver.Main(v.run)
}

func (v *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  v.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	defer func() {
		w.Release()
		close(v.done)
	}()

	v.w = w
	v.bgColor = color.RGBA{0, 128, 0, 255}
	v.figureSize = 200
	v.figurePos = image.Point{X: 400, Y: 400}
	if v.OnScreenReady != nil {
		v.OnScreenReady(s)
	}

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if v.Debug {
				log.Printf("event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			v.handleEvent(e, t)

		case t = <-v.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		return e.To == lifecycle.StageDead
	case key.Event:
		return e.Code == key.CodeEscape
	}
	return false
}

func (v *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	case size.Event:
		v.sz = e

	case error:
		log.Printf("ERROR: %v", e)

	case mouse.Event:
		if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
			if v.OnMove != nil {
				v.OnMove(image.Point{X: int(e.X), Y: int(e.Y)})
			}
			v.mu.Lock()
			v.figurePos = image.Point{X: int(e.X), Y: int(e.Y)}
			v.mu.Unlock()
			v.w.Send(paint.Event{})
		}

	case paint.Event:
		v.mu.Lock()
		defer v.mu.Unlock()

		if v.currentTexture == nil {
			v.drawScene()
		} else {
			v.w.Scale(v.sz.Bounds(), v.currentTexture, v.currentTexture.Bounds(), draw.Src, nil)
		}
		v.w.Publish()
	}
}

func (v *Visualizer) drawScene() {
	v.w.Fill(v.sz.Bounds(), v.bgColor, draw.Src)

	v.drawT180(v.figurePos, v.figureSize)

	for _, br := range imageutil.Border(v.sz.Bounds(), 10) {
		v.w.Fill(br, color.White, draw.Src)
	}
}

func (v *Visualizer) drawT180(center image.Point, size int) {
	thickness := size / 5
	half := size / 2

	horzRect := image.Rect(
		center.X-half,
		center.Y-thickness/2,
		center.X+half,
		center.Y+thickness/2,
	)
	vertRect := image.Rect(
		center.X-thickness/2,
		center.Y-thickness/2,
		center.X+thickness/2,
		center.Y-half,
	)

	v.w.Fill(horzRect, color.RGBA{255, 255, 0, 255}, draw.Src)
	v.w.Fill(vertRect, color.RGBA{255, 255, 0, 255}, draw.Src)
}
