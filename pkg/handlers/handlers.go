package handlers

import (
	//"html/template"
	//"fmt"
	"encoding/json"
	"errors"
	"myapp2/pkg/config"
	"myapp2/pkg/driver"
	"myapp2/pkg/forms"
	"myapp2/pkg/helpers"
	"myapp2/pkg/models"
	"myapp2/pkg/renderers"
	"myapp2/pkg/repository"
	"myapp2/pkg/repository/db_repo"

	"strings"

	//"html/template"
	"strconv"
	"time"

	//"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	//"github.com.backup/ian-kent/go-log/layout"
	"github.com/go-chi/chi"
)

var app *config.AppConfig
var db repository.DatabaseRepo

func SetConfig(config *config.AppConfig) {

	app = config

}

func SetRepository(db1 *driver.DB) {
	db = dbrepo.NewPostgresRepo(db1.SQL, app)
}

func SetTestRepositery() {
	db = dbrepo.NewTestDBRepo(app)
}

func Home(w http.ResponseWriter, r *http.Request) {

	//fmt.Println(db.AllUsers())

	remoteIp := r.RemoteAddr
	app.Session.Put(r.Context(), "remote_ip", remoteIp)

	renderers.Render("/home.page.tmpl", w, &models.TemplateData{}, r)
	return

}

func Index(w http.ResponseWriter, r *http.Request) {
	renderers.Render("/index.page.tmpl", w, &models.TemplateData{}, r)
	return
}

func About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again"
	remoteIp := (*app).Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	fmt.Println(app.Session.Cookie)
	renderers.Render("/about.page.tmpl", w, &models.TemplateData{
		StringMap: stringMap,
	}, r)
	return
}

func Rooms(w http.ResponseWriter, r *http.Request) {

	a := filepath.Base((r.RequestURI))
	stringMap := make(map[string]string)

	roomID, err := strconv.Atoi(string(a[len(a)-1]))
	if err != nil {
		helpers.ServerError(w, err)
	}

	app.Session.Put(r.Context(), "room_id_display", roomID)

	stringMap["room"] = a + ".png"
	fmt.Println(stringMap)
	renderers.Render("/room1.page.tmpl", w, &models.TemplateData{
		StringMap: stringMap,
	}, r)
}

func Reservation(w http.ResponseWriter, r *http.Request) {

	renderers.Render("/reservation.page.tmpl", w, &models.TemplateData{}, r)
}

func MakeReservations(w http.ResponseWriter, r *http.Request) {
	layout := "2006-01-02"
	res, ok := app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot parse or get reservation"))
		return
	}

	room, err := db.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room = room

	app.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format(layout)
	ed := res.EndDate.Format(layout)

	data := make(map[string]interface{})
	data["reservation"] = res

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	renderers.Render("/make_reservation.page.tmpl", w, &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	}, r)
}

func PostMakeReservations(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in reservation")
	layout := "2006-01-02"
	reservation, ok := app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("not able to typecast or cant get from session line 143"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)
	//form.Has("first_name",r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation

		stringMap := make(map[string]string)
		stringMap["start_date"] = form.Get("start_date")
		stringMap["end_date"] = form.Get("end_date")

		//http.Error(w, , http.StatusSeeOther)
		renderers.Render("/make_reservation.page.tmpl", w, &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		}, r)

		return

	}
	newReservation, err := db.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservation,
		RestrictionID: 1,
	}

	err = db.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
	}

	app.Session.Put(r.Context(), "reservation", reservation)

	//send notifications - first to guest

	htmlMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation </strong><br>

		Dear %s,<br>

		This is to confirm your Reservation from %s to %s

	`, reservation.FirstName, reservation.StartDate.Format(layout), reservation.EndDate.Format(layout))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "m@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	app.MailChan <- msg

	http.Redirect(w, r, "reservation-sum", http.StatusSeeOther)

}

func PostReservation(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		return
	}

	rooms, err := db.RoomsAvaialible(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	if len(rooms) == 0 {
		// no rooms availibility
		app.Session.Put(r.Context(), "error", "No Availibility")
		http.Redirect(w, r, "/book", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{

		StartDate: startDate,
		EndDate:   endDate,
	}

	app.Session.Put(r.Context(), "reservation", res)
	renderers.Render("/choose_room.page.tmpl", w, &models.TemplateData{
		Data: data,
	}, r)
	return
	//w.Write([]byte("posted to search availibility"+start+end))
}

type jsonResponse struct {
	Ok  bool   `json:"ok"`
	MSG string `json:"message"`
}

func AvailibilityJson(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		resp := jsonResponse{
			Ok:  false,
			MSG: "internal server form error",
		}

		out, _ := json.MarshalIndent(resp, "", " ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	layout := "2006-01-02"
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	//roomID , _ := strconv.Atoi(r.Form.Get("room_id"))

	roomID, ok := app.Session.Get(r.Context(), "room_id_display").(int)
	if !ok {
		resp := jsonResponse{
			Ok:  false,
			MSG: "internal server session error",
		}

		out, _ := json.MarshalIndent(resp, "", " ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	availible, err := db.RoomsAvailibilityByDates(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			Ok:  false,
			MSG: "internal server db error",
		}

		out, _ := json.MarshalIndent(resp, "", " ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		Ok:  availible,
		MSG: "",
	}

	out, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		resp := jsonResponse{
			Ok:  false,
			MSG: "internal server  indent error",
		}

		out, _ := json.MarshalIndent(resp, "", " ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		app.ErrorLog.Println("cant get reservation from session")
		log.Println("cannot get reservation from session")
		app.Session.Put(r.Context(), "error", "cant get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.Session.Get(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	layout := "2006-01-02"
	sd := reservation.StartDate.Format(layout)
	ed := reservation.EndDate.Format(layout)

	stringMap := make(map[string]string)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	renderers.Render("/reservation-summary.page.tmpl", w, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	}, r)
	return

}

func ChooseRoom(w http.ResponseWriter, r *http.Request) {

	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("type assertion failed"))
		return
	}
	res.RoomID = roomID
	app.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/reservation", http.StatusSeeOther)
	return

}

func ShowLogin(w http.ResponseWriter, r *http.Request) {

	renderers.Render("/login.page.tmpl", w, &models.TemplateData{
		Form: forms.New(nil),
	}, r)

	return

}

func PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = app.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	var email, password string

	email = r.Form.Get("email")
	password = r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		fmt.Println("in here")
		http.Error(w, "error", http.StatusForbidden)
		renderers.Render("/login.page.tmpl", w, &models.TemplateData{
			Form: form,
		}, r)
		return
	}

	id, _, err := db.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		app.Session.Put(r.Context(), "error", "invalid login details")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "user_id", id)
	app.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return

}

func Logout(w http.ResponseWriter, r *http.Request) {

	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {

	renderers.RenderAdmin("/admin.page.tmpl", w, &models.TemplateData{}, r)
}

func AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := db.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	renderers.RenderAdmin("/admin_new_reservations.page.tmpl", w, &models.TemplateData{
		Data: data,
	}, r)

}

func AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := db.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	renderers.RenderAdmin("/admin_all_reservations.page.tmpl", w, &models.TemplateData{
		Data: data,
	}, r)

}

func AdminShowReservation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	res, err := db.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	renderers.RenderAdmin("/admin_reservations_show.page.tmpl", w, &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	}, r)

}

func AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := db.GetReservationByID(id)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	fmt.Println(res.Phone)

	err = db.UpdateReservation(res)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")
	app.Session.Put(r.Context(), "flash", "Changes Saved")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

func AdminReservationsCalander(w http.ResponseWriter, r *http.Request) {

	// assume there is no month/year specified

	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())

	}

	data := make(map[string]interface{})

	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear

	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()

	currentLocation := now.Location()

	firstMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	lastOfMonth := firstMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)

	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := db.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		// create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-02")] = 0
			blockMap[d.Format("2006-01-02")] = 0

		}
		// get all the restrictions for the current room

		restrictions, err := db.GetRestrictionsForRoomByDate(x.ID, firstMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
		}
		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// its a reservation

				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ID
				}

			} else {
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					blockMap[d.Format("2006-01-2")] = y.ID
				}

			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		app.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)

	}

	renderers.RenderAdmin("/admin_reservations_calender.page.tmpl", w, &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	}, r)

}

// marks a reservation as processed
func AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	src := chi.URLParam(r, "src")

	err := db.UpdateReservationProcess(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	app.Session.Put(r.Context(), "flash", "reservation marked as processed")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?m=%s&y=%s", month, year), http.StatusSeeOther)
	}

}

func AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	src := chi.URLParam(r, "src")

	err := db.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	app.Session.Put(r.Context(), "flash", "deleted reservation")

	if year == "" {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?m=%s&y=%s", month, year), http.StatusSeeOther)

	}

}

// Handles post of reservation calender
func AdminPostReservationsCalander(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//process blocks

	rooms, err := db.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	for _, x := range rooms {
		// Get the Blockmap from the session. Loop through entire map , if we have entry in map
		// that does not exist in posted data , and if the restriction_id is 0 , than it is a block we need to remove

		currMap := app.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range currMap {
			// ok will be false if the value is not in the map
			if val, ok := currMap[name]; ok {
				// only pay attention to values >0, and that are not in the form post
				// the rest are just placeholder
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name), r) {
						//delete the restriction by id
						err := db.DeleteBlockForRoom(value)
						if err != nil {
							helpers.ServerError(w, err)
							return
						}
					}
				}
			}

		}

	}

	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, err := strconv.Atoi(exploded[2])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

			startDate, err := time.Parse("2006-01-2", exploded[3])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

			// insert a new block
			err = db.InsertBlockForRoom(roomID, startDate)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
		}
	}

	app.Session.Put(r.Context(), "flash", "changes saved")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)

	return

}
