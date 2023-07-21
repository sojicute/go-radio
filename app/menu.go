package main

import "github.com/maxence-charriere/go-app/v7/pkg/app"

type menu struct {
	app.Compo
}

func (m *menu) Render() app.UI {
	return app.Nav().Class("menu").Body(
		app.H1().Text("rock radio"),
	)

}
