package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"myapp2/pkg/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"rooms", "/rooms/room1", "GET", []postData{}, http.StatusOK},
	{"rooms", "/rooms/room2", "GET", []postData{}, http.StatusOK},
	{"book-get", "/book", "GET", []postData{}, http.StatusOK},
	{"book", "/book", "POST", []postData{
		{key: "start", value: "2021-01-01"},
		{key: "end", value: "2021-05-09"},
	}, http.StatusOK}, /*
		{"book-json","/book-json","POST",[]postData{
			{key: "start" , value: "2021-01-01"},
			{key: "end" , value: "2021-05-09"},
		}, http.StatusOK},*/
	//{"choose-room","/choose_room/1","GET",[]postData{},http.StatusOK},
	//{"get-make-reservation","/reservation","GET",[]postData{},http.StatusOK},
	/*{"make-reservation","/reservation","POST",[]postData{
		{key: "first_name" , value: "Chaitanya"},
		{key: "last_name" , value: "Malik"},
		{key: "email" , value: "cm110@g.com"},
		{key: "phone" , value: "9711844678"},

	}, http.StatusOK},*/
}

func TestHandlers(t *testing.T) {
	fmt.Println("starting testing")

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, response got %d but wanted %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {

			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				//t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, response got %d but wanted %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}
		fmt.Println(e)
	}

}

func TestRepository_Reservation(t *testing.T) {

	room_display := 1

	req, _ := http.NewRequest("GET", "/book-json", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "room_id_display", room_display)

	handler := http.HandlerFunc(AvailibilityJson)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("test failed")

	}

	reservation := models.Reservation{
		ID:        0,
		FirstName: "",
		LastName:  "",
		Email:     "",
		Phone:     "",
		StartDate: time.Time{},
		EndDate:   time.Time{},
		RoomID:    0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Room: models.Room{
			ID:        0,
			RoomName:  "",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	session.Put(ctx, "reservation", reservation)

	req, _ = http.NewRequest("POST", "/reservation", nil)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(MakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("Didnt work for reservation")
	}

	session.Put(ctx, "reservation", reservation)

	req, _ = http.NewRequest("POST", "/reservation", nil)

	urlsValues := []postData{
		{key: "first_name", value: "chaitanya"},
		{key: "last_name", value: "malik"},
	}

	values := url.Values{}
	for _, x := range urlsValues {
		values.Add(x.key, x.value)
	}

	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.PostForm = values
	handler = http.HandlerFunc(PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Error("Didnt work for reservation")
	}

	fmt.Println("You are awesome")

}

func TestRepository_PostReservation(t *testing.T) {

	reservation := models.Reservation{
		ID:        0,
		FirstName: "",
		LastName:  "",
		Email:     "",
		Phone:     "",
		StartDate: time.Time{},
		EndDate:   time.Time{},
		RoomID:    0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Room: models.Room{
			ID:        0,
			RoomName:  "",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=cm12")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abc@c.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=90909000")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)

	session.Put(ctx, "reservation", reservation)

	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(PostMakeReservations)

	handler.ServeHTTP(rr, req)

}

func TestRepository_AvailibilityJson(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-01")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//create request

	req, _ := http.NewRequest("POST", "/book-json", strings.NewReader(reqBody))

	ctx := getCtx(req)

	room_display := 1

	session.Put(ctx, "room_id_display", room_display)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "x-www-form-urlencoded")

	handler := http.HandlerFunc(AvailibilityJson)

	//make request

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Error("failed to parse json")
	}
	fmt.Println(j)
	if !j.Ok {
		t.Error("failed the json test")
	}

}

func getCtx(req *http.Request) context.Context {

	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx

}
