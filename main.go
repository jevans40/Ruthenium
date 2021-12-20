package main

import (
	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/game"
)

//A simple example of how to create a game with this library.
func main() {
	newGame := game.NewGameECS()
	x := component.NewDenseStorage[component.BaseComponent]()
	_ = x
	newGame.Init()

}
