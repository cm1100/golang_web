package driver

import (
	"database/sql"
	"time"

	
	_ "github.com/jackc/pgx/v4"

	_ "github.com/jackc/pgconn"

	_ "github.com/jackc/pgx/v4/stdlib"

)

// DB holds the database connection
type DB struct{
	SQL *sql.DB


}


var DbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute


// connectSql creates database pools for postgress
func ConnectSql(dsn string) (*DB,error){

	d,err := NewDatabase(dsn)
	if err!=nil{
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)
	d.SetConnMaxIdleTime(maxIdleDbConn)

	DbConn.SQL = d

	if err = TestDb(d);err!=nil{
		panic(err)
	}

	return DbConn,err


}

// Tries to ping the database
func TestDb(d *sql.DB) error{
	err := d.Ping()
	if err != nil{
		return err
	}
	return nil
}


// Created a New Database for the applications
func NewDatabase(dsn string) (*sql.DB,error){

	db,err := sql.Open("pgx",dsn)
	if err!=nil{
		return nil,err
	}

	if err = db.Ping();err!=nil{
		return nil,err
	}
	return db, nil
}