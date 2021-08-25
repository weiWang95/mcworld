package main

import (
	"github.com/weiWang95/mcworld/app"
	_ "github.com/weiWang95/mcworld/level/demo"
	_ "github.com/weiWang95/mcworld/level/world"
)

func main() {
	app.Create().Run()
}
