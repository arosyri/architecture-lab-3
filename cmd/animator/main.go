package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func sendCommand(cmd string) error {
	resp, err := http.Get("http://localhost:17000/?cmd=" + url.QueryEscape(cmd))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	return nil
}

func main() {
	fmt.Println("Starting animation client")

	x, y := 100, 100
	dx, dy := 10, 15

	for {
		cmdMove := fmt.Sprintf("move %d %d", x, y)
		if err := sendCommand(cmdMove); err != nil {
			fmt.Println("Error sending move:", err)
		} else {
			fmt.Println("Sent:", cmdMove)
		}

		if err := sendCommand("update"); err != nil {
			fmt.Println("Error sending update:", err)
		} else {
			fmt.Println("Sent: update")
		}

		x += dx
		y += dy

		if x < 0 || x > 800 {
			dx = -dx
		}
		if y < 0 || y > 800 {
			dy = -dy
		}

		time.Sleep(1 * time.Second)
	}
}
