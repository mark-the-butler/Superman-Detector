package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mysteryboy73/Superman-Detector/models"
)

// DataRepo is an interface for retrieving data from a db
type DataRepo interface {
	getLocation(ipAddress string) (models.GeoLocation, error)
	saveLogin(login models.LoginRequest) bool
	getPreviousAndFutureIPAdress(username string, currentIP string, currentTimeStamp int) (previousLogin, futureLogin models.LoginRequest)
}

type loginAttempt struct {
	RowNum        int    `db:"row_num"`
	Username      string `db:"user_name"`
	UnixTimestamp int    `db:"time_stamp"`
	EventUUID     string `db:"event_uuid"`
	IPAddress     string `db:"ip_address"`
}

// GeoRepository implements the DataRepo interface for retrieving data from a db
type GeoRepository struct {
}

func (gr *GeoRepository) saveLogin(login models.LoginRequest) bool {
	database, err := sqlx.Open("sqlite3", "./db/geolite2.db")
	checkErr(err)

	statement, err := database.Prepare("INSERT INTO logins(user_name, time_stamp, event_uuid, ip_address) values(?,?,?,?)")
	checkErr(err)

	res, err := statement.Exec(login.Username, login.UnixTimestamp, login.EventUUID, login.IPAddress)
	checkErr(err)

	rows, err := res.RowsAffected()

	// Need to figure out what rows is and not saved record to db if it exist
	fmt.Printf(string(rows))

	if err != nil {
		return false
	}

	return true
}

// GetCurrentLocation retrieves users current location from db
func (gr *GeoRepository) getLocation(ipAddress string) (models.GeoLocation, error) {
	var currentLocation models.GeoLocation

	database, err := sqlx.Open("sqlite3", "./db/geolite2.db")
	checkErr(err)

	statement, err := database.Prepare("SELECT latitude, longitude, accuracy_radius FROM blocks WHERE network =?")
	checkErr(err)

	rows, err := statement.Query(ipAddress)
	checkErr(err)

	if rows.Next() {
		err = rows.Scan(&currentLocation.Lat, &currentLocation.Lon, &currentLocation.Radius)
		checkErr(err)
	}

	rows.Close()
	statement.Close()
	return currentLocation, nil
}

func (gr *GeoRepository) getPreviousAndFutureIPAdress(username string, currentIP string, currentTimeStamp int) (previousLogin, futureLogin models.LoginRequest) {
	database, err := sqlx.Open("sqlite3", "./db/geolite2.db")
	checkErr(err)

	loginAttempts := []loginAttempt{}
	statement, err := database.Preparex(`with cte as (select row_number() over(order by time_stamp desc) row_num,* from logins where user_name = $1),
	current as (select row_num from cte where ip_address = $2)
	select cte.* from cte, current where abs(cte.row_num - current.row_num) <= 1 order by cte.time_stamp desc;`)

	if queryError := statement.Select(&loginAttempts, username, currentIP); queryError != nil {
		fmt.Printf(queryError.Error())
		panic(queryError)
	}

	for _, attempt := range loginAttempts {
		if attempt.IPAddress != currentIP && attempt.UnixTimestamp < currentTimeStamp {
			previousLogin.Username = attempt.Username
			previousLogin.UnixTimestamp = attempt.UnixTimestamp
			previousLogin.EventUUID = attempt.EventUUID
			previousLogin.IPAddress = attempt.IPAddress
		}

		if attempt.IPAddress != currentIP && attempt.UnixTimestamp > currentTimeStamp {
			futureLogin.IPAddress = attempt.IPAddress
			futureLogin.UnixTimestamp = attempt.UnixTimestamp
			futureLogin.EventUUID = attempt.EventUUID
			futureLogin.IPAddress = attempt.IPAddress
		}
	}

	statement.Close()

	return previousLogin, futureLogin
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
