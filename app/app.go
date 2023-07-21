package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

type home struct {
	app.Compo
}

func (h *home) Render() app.UI {
	return app.Div().Body(
		&player{},
		app.Shell().
			Menu(&menu{}).
			Content(&detail{}),
	)
}

func main() {
	app.Route("/", &home{})
	app.Run()
}
