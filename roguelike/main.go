package main

import (
	"log"
	"os"
	"roguelike/internal/game"
)

func main() {
	builder := game.NewDefaultGameBuilder()

	// Вычитываем конфигурацию локации из файла
	err := builder.SetLocation("resources/location_1.json")
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	// Строим экран для рендеринга игры
	err = builder.BuildScreen()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	// Строим движок для игры по данным указанным в конфигурации локации
	builder.BuildEngine()

	game := builder.GetResult()
	defer game.Stop()
	game.Run()
}
