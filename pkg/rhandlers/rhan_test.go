package rhandlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)


func TestNoSurf(t *testing.T){

	var myH myHandler

	h:= NoSurf(&myH)

	switch v:= h.(type){
	case http.Handler:


	default: t.Errorf(fmt.Sprintf("test failed by nosurf type is not http.Handler %T",v))
	}
}


func TestSessionLoad( t *testing.T){


	var myH myHandler

	h:= SessionLoad(&myH)

	switch v:= h.(type){
	case http.Handler:


	default: t.Errorf(fmt.Sprintf("test failed by nosurf type is not http.Handler %T",v))
	}

}


func TestRoutes(t *testing.T){

	mux := Routes2(app)

	switch mux.(type){
	case *chi.Mux:

	default:
		t.Errorf("not of type http handler routes test")
	}

}



