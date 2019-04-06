package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	l := bastion.LoggerFromCtx(r.Context())
	l.Info().Msg("handler")

	render.NewJSON().Send(w, res)
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	app.Logger.Info().Str("app", "test").Msg("main")
	fmt.Fprintln(os.Stderr, app.Serve())
}
