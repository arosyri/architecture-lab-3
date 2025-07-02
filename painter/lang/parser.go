package lang

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fmt.Println(">> received command:", line)
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		cmd := strings.ToLower(fields[0])

		switch cmd {
		case "white":
			res = append(res, painter.FillBackground{Color: color.RGBA{255, 255, 255, 255}})

		case "green":
			res = append(res, painter.FillBackground{Color: color.RGBA{0, 128, 0, 255}})

		case "update":
			res = append(res, painter.UpdateOp)

		case "bgrect":
			if len(fields) != 5 {
				return nil, fmt.Errorf("bgrect requires 4 arguments")
			}
			x1, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			y1, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, err
			}
			x2, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, err
			}
			y2, err := strconv.ParseFloat(fields[4], 64)
			if err != nil {
				return nil, err
			}

			r := image.Rect(int(x1*800), int(y1*800), int(x2*800), int(y2*800))
			res = append(res, painter.BgRect{Rect: r})

		case "figure":
			if len(fields) != 3 {
				return nil, fmt.Errorf("figure requires 2 arguments")
			}
			x, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(fields[2])
			if err != nil {
				return nil, err
			}
			res = append(res, painter.DrawT180{
				PosX:  x,
				PosY:  y,
				Size:  100,
				Color: color.RGBA{255, 255, 0, 255},
			})

		case "move":
			if len(fields) != 3 {
				return nil, fmt.Errorf("move requires 2 arguments")
			}
			x, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(fields[2])
			if err != nil {
				return nil, err
			}
			res = append(res, painter.Move{NewPos: image.Point{X: x, Y: y}})

		case "border":
			var c color.Color = color.Black
			if len(fields) >= 2 {
				switch fields[1] {
				case "green":
					c = color.RGBA{0, 255, 0, 255}
				case "white":
					c = color.White
				case "red":
					c = color.RGBA{255, 0, 0, 255}
				}
			}
			res = append(res, painter.Border{Thickness: 10, Color: c})

		case "reset":
			res = append(res, painter.Reset{})

		default:
			return nil, fmt.Errorf("unknown command: %s", cmd)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
