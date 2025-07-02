package painter

import (
	"image"
	"image/color"
	"log"
	"sync"
	"time"

	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	next screen.Texture
	prev screen.Texture

	mq messageQueue

	stopReq bool
	stopped chan struct{}

	mu sync.Mutex

	bgColor color.RGBA
	bgRect  *image.Rectangle
	border  *Border
	figures []DrawT180
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.bgColor = color.RGBA{0, 128, 0, 255}

	l.figures = append(l.figures, DrawT180{
		PosX:  size.X / 2,
		PosY:  size.Y / 2,
		Size:  200,
		Color: color.RGBA{255, 255, 0, 255},
	})

	l.stopped = make(chan struct{})

	go func() {
		defer close(l.stopped)
		for {
			if l.stopReq && l.mq.empty() {
				return
			}

			op := l.mq.pull()
			if op == nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			l.handleOp(op)
		}
	}()
}

func (l *Loop) handleOp(op Operation) {
	l.mu.Lock()
	defer l.mu.Unlock()

	log.Printf("Handling operation: %T %+v", op, op)

	switch op := op.(type) {
	case OperationList:
		for _, o := range op {
			l.handleOp(o)
		}

	case FillBackground:
		l.bgColor = op.Color
		log.Printf("Background color changed to: %#v", l.bgColor)

	case BgRect:
		l.bgRect = &op.Rect

	case Reset:
		l.bgColor = color.RGBA{0, 128, 0, 255}
		l.bgRect = nil
		l.figures = nil
		l.border = nil

	case DrawT180:
		l.figures = append(l.figures, op)

	case Move:
		for i := range l.figures {
			l.figures[i].PosX = op.NewPos.X
			l.figures[i].PosY = op.NewPos.Y
		}

	case Border:
		l.border = &op

	case updateOp:
		l.next.Fill(l.next.Bounds(), l.bgColor, screen.Src)

		if l.bgRect != nil {
			l.next.Fill(*l.bgRect, color.Black, screen.Src)
		}

		for _, f := range l.figures {
			drawT180(l.next, f.PosX, f.PosY, f.Size, f.Color)
		}

		if l.border != nil {
			borders := imageutil.Border(l.next.Bounds(), l.border.Thickness)
			for _, r := range borders {
				l.next.Fill(r, l.border.Color, screen.Src)
			}
		}

		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next

	default:
		op.Do(l.next)
	}
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.stopReq = true
	<-l.stopped
}

type messageQueue struct {
	mu  sync.Mutex
	ops []Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.ops = append(mq.ops, op)
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if len(mq.ops) == 0 {
		return nil
	}
	op := mq.ops[0]
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.ops) == 0
}

func drawT180(t screen.Texture, cx, cy, size int, col color.RGBA) {
	thickness := size / 5
	half := size / 2

	horzRect := image.Rect(cx-half, cy-thickness/2, cx+half, cy+thickness/2)
	t.Fill(horzRect, col, screen.Src)

	vertRect := image.Rect(cx-thickness/2, cy-half, cx+thickness/2, cy+thickness/2)
	t.Fill(vertRect, col, screen.Src)
}
