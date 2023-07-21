package main

import (
	"log"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", &app.Handler{
		Author:      "soji💎",
		Title:       "24/7 indie / pop / rock radio",
		Name:        "Radio",
		Description: "We're broadcasting a mix of indie, pop and rock music provided and handpicked by new, up-and-coming artists worldwide. This playlist is a collection of the freshest songs for you to enjoy.",
		Icon: app.Icon{
			Default: "/web/favicon.png",
		},
		ThemeColor:      "#000000",
		BackgroundColor: "#000000",
		Styles: []string{
			"/web/home.css",
			"/web/range.css",
		},
		RawHeaders: []string{
			`
		<meta property="og:type" content="website">
		<meta property="og:url" content="https://secure-dawn-84089.herokuapp.com/">
		<meta property="og:title" content="24/7 indie / pop / rock radio">
		<meta property="og:description" content="We're broadcasting a mix of indie, pop and rock music provided and handpicked by new, up-and-coming artists worldwide. This playlist is a collection of the freshest songs for you to enjoy.">
		<meta property="og:image" content="https://images-na.ssl-images-amazon.com/images/I/91Z418spOWL._AC_SL1500_.jpg">
		`,
		},
	})
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
