package rhandlers

import (
	"myapp2/pkg/config"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	var app = config.AppConfig{}

	SetConfig(&app)

	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
