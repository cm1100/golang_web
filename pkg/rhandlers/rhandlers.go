package rhandlers

import (
	//"fmt"
	"github.com/bmizerany/pat"
	"myapp2/pkg/config"
	"myapp2/pkg/handlers"
	"net/http"
)

func Routes(app *config.AppConfig) http.Handler {

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))

	mux.Get("/about", http.HandlerFunc(handlers.About))
	mux.Get("/index", http.HandlerFunc(handlers.Index))
	return mux

}
