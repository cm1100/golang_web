package dbrepo

import (
	"context"
	"errors"
	"myapp2/pkg/models"

	//"hash"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDb) AllUsers() bool {

	return true

}

func (m *postgresDb) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,end_date
		,room_id,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName, res.LastName,
		res.Email, res.Phone, res.StartDate, res.EndDate,
		res.RoomID, time.Now(), time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDb) InsertRoomRestriction(restrict models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date,end_date,room_id,reservation_id,created_at,updated_at,
		restriction_id) values($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt, restrict.StartDate, restrict.EndDate, restrict.RoomID,
		restrict.ReservationID, time.Now(), time.Now(), restrict.RestrictionID)
	if err != nil {
		return err
	}
	//fmt.Println("added without error")
	return nil
}

func (m *postgresDb) RoomsAvailibilityByDates(start, end time.Time, roomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

		select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date
	`

	var numRows int

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, nil
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil

}

func (m *postgresDb) RoomsAvaialible(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select r.id,r.room_name from rooms r where id  not in  (select rr.room_id from room_restrictions rr 
												where 
												$1>rr.start_date 
												 and 
												 $2<rr.end_date)
					
	`
	rows, err := m.DB.QueryContext(ctx, query, end, start)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (m *postgresDb) GetRoomByID(id int) (models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select * from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

//Returns User by Id
func (m *postgresDb) GetUserByID(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

		select id , first_name,last_name,email,password ,acesss_level, created_at, updated_at
		from users where id = $1

	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.AccessLevel, &u.CreatedAt,
		&u.UpdatedAt)
	if err != nil {
		log.Println(err)
		return u, err

	}

	return u, nil

}

func (m *postgresDb) UpdateUser(u models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

			update users first_name=$1,last_name=$2,email=$3,access_level=$4,
			updated_at = $5
	`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName,
		u.LastName, u.Email, u.AccessLevel, time.Now())

	if err != nil {
		return err
	}

	return nil

}

// Authentcates a User
func (m *postgresDb) Authenticate(email, testPassword string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashPassword string

	row := m.DB.QueryRowContext(ctx, "select id , password from users where email = $1", email)

	err := row.Scan(&id, &hashPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {

		return 0, "", errors.New("Incorrect Password")

	} else if err != nil {

		return 0, "", err

	}

	return id, hashPassword, nil

}

//Returns the slice of all reservations
func (m *postgresDb) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `

		select r.id, r.first_name , r.last_name , r.email,
		r.phone,r.start_date,r.end_date,r.room_id,r.created_at,
		r.updated_at,r.processed,rm.id,rm.room_name

		from reservations r
		left join rooms rm on r.room_id = rm.id
		order by r.start_date asc

	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt, &i.UpdatedAt, &i.Processed,
			&i.Room.ID, &i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {

		return reservations, err

	}

	return reservations, nil

}

//Returns the slice of new reservations
func (m *postgresDb) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `

		select r.id, r.first_name , r.last_name , r.email,
		r.phone,r.start_date,r.end_date,r.room_id,r.created_at,
		r.updated_at,r.processed,rm.id,rm.room_name

		from reservations r
		left join rooms rm on r.room_id = rm.id
		where r.processed=0
		order by r.start_date asc

	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt, &i.UpdatedAt, &i.Processed,
			&i.Room.ID, &i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {

		return reservations, err

	}

	return reservations, nil

}

// Returns One Reservation By id
func (m *postgresDb) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var i models.Reservation

	query := `
		select r.id, r.first_name , r.last_name , r.email,
		r.phone,r.start_date,r.end_date,r.room_id,r.created_at,
		r.updated_at,r.processed,rm.id,rm.room_name

		from reservations r
		left join rooms rm on r.room_id = rm.id
		where r.id = $1
		order by r.start_date asc
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.StartDate,
		&i.EndDate,
		&i.RoomID,
		&i.CreatedAt, &i.UpdatedAt, &i.Processed,
		&i.Room.ID, &i.Room.RoomName,
	)
	if err != nil {
		return i, err
	}

	return i, nil

}

func (m *postgresDb) UpdateReservation(u models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

			update reservations set first_name=$1,last_name=$2,email=$3,phone=$4,
			updated_at=$5 where id=$6
	`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName,
		u.LastName, u.Email, u.Phone, time.Now(), u.ID)

	if err != nil {
		return err
	}

	return nil

}

func (m *postgresDb) DeleteReservation(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

		delete from reservations where id = $1

	`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil

}

//updated process for the reservation by ID
func (m *postgresDb) UpdateReservationProcess(id, process int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `

		update reservations set processed=$1 where id = $2

	`

	_, err := m.DB.ExecContext(ctx, query, process, id)

	if err != nil {
		return err
	}

	return nil

}

func (m *postgresDb) AllRooms() ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `

		select id,room_name,created_at,updated_at from rooms order by room_name
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, nil
	}
	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)
		if err != nil {
			return rooms, nil
		}

		rooms = append(rooms, rm)
	}

	err = rows.Err()
	if err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDb) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `

	select id , coalesce(reservation_id,0) , restriction_id , room_id , start_date , end_date
	from room_restrictions where $1<end_date and $2>=start_date and room_id = $3

	`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, err
	}

	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, r)
	}

	if err := rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil

}

// Inserts a Room Restriction
func (m *postgresDb) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		insert into room_restrictions (start_date,end_date,room_id, restriction_id,created_at,updated_at)
		values ($1,$2,$3,$4,$5,$6)
	`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate, id, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil

}

// Deletes a Room Restrictions
func (m *postgresDb) DeleteBlockForRoom(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		delete from room_restrictions where id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

}
