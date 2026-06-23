package main

import (
	"traffic-guarder/internal/infrastructure/app"
	"traffic-guarder/internal/router"
)

func main() {

	r := router.NewRouter()
	a := app.New(r)
	a.Start()

}
