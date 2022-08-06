package dbrepo

import (
	"myapp2/pkg/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {

	return true

}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(restrict models.RoomRestriction) error {

	return nil
}

func (m *testDBRepo) RoomsAvailibilityByDates(start, end time.Time, roomID int) (bool, error) {

	return true, nil

}

func (m *testDBRepo) RoomsAvaialible(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {

	var u models.User

	return u, nil

}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	return 0, "", nil

}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil

}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil

}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var i models.Reservation

	return i, nil

}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil

}

func (m *testDBRepo) DeleteReservation(id int) error {

	return nil

}

func (m *testDBRepo) UpdateReservationProcess(id, process int) error {

	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction

	return restrictions, nil

}

func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	return nil
}

func (m *testDBRepo) DeleteBlockForRoom(id int) error {
	return nil
}
