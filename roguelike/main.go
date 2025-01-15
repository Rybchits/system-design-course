package main

import (
	"log"
	"os"
	"roguelike/components"

	"github.com/gdamore/tcell/v2"
)

func drawString(screen tcell.Screen, x, y int, msg string) {
	for i, char := range msg {
		screen.SetContent(x+i, y, char, nil, tcell.StyleDefault)
	}
}

func canMove(newPosition *components.Position, width, height int) bool {
	return newPosition.X < width && newPosition.Y < height && newPosition.X >= 0 && newPosition.Y >= 0
}

func main() {
	width, height := 10, 10

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
		os.Exit(1)
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
		os.Exit(1)
	}
	playerPosition := components.NewPosition().WithX(0).WithY(0)

	running := true
	for running {
		screen.Clear()
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				screen.SetContent(int(x), int(y), '.', nil, tcell.StyleDefault)
			}
		}
		screen.SetContent(int(playerPosition.X), int(playerPosition.Y), '@', nil, tcell.StyleDefault)
		drawString(screen, 0, height, "Use WASD + q")

		screen.Show()

		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			newPosition := playerPosition.Clone()
			switch event.Rune() {
			case 'q':
				running = false
			case 'w', 'W':
				newPosition.WithY(newPosition.Y - 1)
			case 'a', 'A':
				newPosition.WithX(newPosition.X - 1)
			case 's', 'S':
				newPosition.WithY(newPosition.Y + 1)
			case 'd', 'D':
				newPosition.WithX(newPosition.X + 1)
			}

			if canMove(newPosition, width, height) {
				playerPosition = newPosition
			}
		}
	}
}
