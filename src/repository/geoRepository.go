package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Driver to connect to sqlite db
	"github.com/mysteryboy73/Superman-Detector/models"
)

// DataRepo is an interface for retrieving data from a db
type DataRepo interface {
	GetLocation(ipAddress string) (models.GeoLocation, error)
	SaveLogin(login models.LoginRequest) bool
	GetPreviousAndFutureIPAdress(username string, currentIP string, currentTimeStamp int) (previousLogin, futureLogin models.LoginRequest)
}

// GeoRepository implements the DataRepo interface for retrieving data from a db
type GeoRepository struct {
	database sqlx.DB
}

// NewGeoRepository creates a geo repository with database connection
func NewGeoRepository() *GeoRepository {
	database, err := sqlx.Open("sqlite3", "../db/geolite2.db")
	checkErr(err)

	geoRepository := GeoRepository{database: *database}
	return &geoRepository
}

// SaveLogin saves login attempts to the database if they have not been saved before
func (gr *GeoRepository) SaveLogin(login models.LoginRequest) bool {
	statement, prepErr := gr.database.Prepare("INSERT OR IGNORE INTO logins(user_name, time_stamp, event_uuid, ip_address) values(?,?,?,?)")
	checkErr(prepErr)

	_, err := statement.Exec(login.Username, login.UnixTimestamp, login.EventUUID, login.IPAddress)
	checkErr(err)
	defer statement.Close()

	if err != nil {
		return false
	}

	return true
}

// GetLocation retrieves users current location from db
func (gr *GeoRepository) GetLocation(ipAddress string) (models.GeoLocation, error) {
	var currentLocation models.GeoLocation

	statement, err := gr.database.Prepare("SELECT latitude, longitude, accuracy_radius FROM blocks WHERE network =?")
	checkErr(err)

	rows, err := statement.Query(ipAddress)
	checkErr(err)
	if err != nil {
		return models.GeoLocation{}, err
	}
	defer statement.Close()

	if rows.Next() {
		err = rows.Scan(&currentLocation.Lat, &currentLocation.Lon, &currentLocation.Radius)
		checkErr(err)
	}
	defer rows.Close()

	return currentLocation, nil
}

// GetPreviousAndFutureIPAdress gets the ip addresses of the previous login attempt and the subsequent login attempt from the current login attempt
func (gr *GeoRepository) GetPreviousAndFutureIPAdress(username string, currentIP string, currentTimeStamp int) (previousLogin, futureLogin models.LoginRequest) {
	loginAttempts := []models.Logins{}
	statement, err := gr.database.Preparex(`with cte as (select row_number() over(order by time_stamp desc) row_num,* from logins where user_name = $1),
	current as (select row_num from cte where ip_address = $2)
	select cte.* from cte, current where abs(cte.row_num - current.row_num) <= 1 order by cte.time_stamp desc;`)
	checkErr(err)
	defer statement.Close()

	queryError := statement.Select(&loginAttempts, username, currentIP)
	checkErr(queryError)

	for _, attempt := range loginAttempts {
		if attempt.IPAddress != currentIP && attempt.UnixTimestamp < currentTimeStamp {
			previousLogin.Username = attempt.Username
			previousLogin.UnixTimestamp = attempt.UnixTimestamp
			previousLogin.EventUUID = attempt.EventUUID
			previousLogin.IPAddress = attempt.IPAddress
		}

		if attempt.IPAddress != currentIP && attempt.UnixTimestamp > currentTimeStamp {
			futureLogin.Username = attempt.Username
			futureLogin.UnixTimestamp = attempt.UnixTimestamp
			futureLogin.EventUUID = attempt.EventUUID
			futureLogin.IPAddress = attempt.IPAddress
		}
	}

	return previousLogin, futureLogin
}

func checkErr(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}
