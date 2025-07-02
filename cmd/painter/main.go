package main

import (
	"image"
	"net/http"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
	"github.com/roman-mazur/architecture-lab-3/ui"
)

func main() {
	var (
		pv     ui.Visualizer
		opLoop painter.Loop
		parser lang.Parser
	)

	pv.Title = "Simple painter"
	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	pv.OnMove = func(p image.Point) {
		opLoop.Post(painter.Move{NewPos: p})
		opLoop.Post(painter.UpdateOp)
	}

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
