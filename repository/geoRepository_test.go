package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mysteryboy73/Superman-Detector/models"
)

func TestGeoRepository(t *testing.T) {
	loginRequest := models.LoginRequest{
		Username:      "Bob",
		UnixTimestamp: 1611692790,
		EventUUID:     "763 - 283",
		IPAddress:     "128.92.3.1/24",
	}

	mockdb, mock, _ := sqlmock.New()
	defer mockdb.Close()

	sqlxDB := sqlx.NewDb(mockdb, "sqlmock")
	geoRepository := GeoRepository{database: *sqlxDB}

	t.Run("SaveLogin returns true when login event saved", func(t *testing.T) {
		mock.ExpectPrepare("INSERT OR IGNORE INTO logins")
		mock.ExpectExec("INSERT OR IGNORE INTO logins").WillReturnResult(sqlmock.NewResult(1, 1))

		loginSaved := geoRepository.SaveLogin(loginRequest)

		if loginSaved != true {
			t.Errorf("login did not get saved")
		}
	})

	t.Run("SaveLogin returns false when login event cannot be saved", func(t *testing.T) {
		expectedErr := errors.New("could not save login")
		mock.ExpectPrepare("INSERT OR IGNORE INTO logins")
		mock.ExpectExec("INSERT OR IGNORE INTO logins").WillReturnError(expectedErr)

		loginSaved := geoRepository.SaveLogin(loginRequest)

		if loginSaved != false {
			t.Errorf("login did get saved")
		}
	})

	t.Run("GetLocation returns location details", func(t *testing.T) {
		ipAddress := "1.2.3.4/5"
		expectedLocation := models.GeoLocation{
			Lat:    34.2,
			Lon:    -83.2,
			Radius: 20,
		}

		rows := sqlmock.NewRows([]string{"latitude", "longitude", "accuracy_radius"}).
			AddRow(34.2, -83.2, 20)
		mock.ExpectPrepare("SELECT (.+) FROM blocks WHERE").ExpectQuery().WithArgs(ipAddress).WillReturnRows(rows)

		actualLocation, _ := geoRepository.GetLocation(ipAddress)

		if reflect.DeepEqual(actualLocation, expectedLocation) != true {
			t.Errorf("GetLocation returned incorrect location: got %+v\n want %+v\n", actualLocation, expectedLocation)
		}
	})

	t.Run("GetLocation returns error when it can't retrieve location", func(t *testing.T) {
		ipAddress := "1.2.3.4/5"
		expectedErr := errors.New("could not retrieve location")

		mock.ExpectPrepare("SELECT (.+) FROM blocks WHERE").ExpectQuery().WithArgs(ipAddress).WillReturnError(expectedErr)

		_, err := geoRepository.GetLocation(ipAddress)

		if err.Error() != expectedErr.Error() {
			t.Errorf("GetLocation returned incorrect error: got %v want %v", err.Error(), expectedErr.Error())
		}
	})

	t.Run("GetPreviousAndFutureIPAdress returns the ip address of the previous and subsequent login attempts", func(t *testing.T) {
		username := "Bob"
		currentIP := "1.2.3.4/5"
		timeStamp := 1611692790
		expectedPreviousLogin := models.LoginRequest{
			IPAddress:     "5.4.3.2/1",
			UnixTimestamp: 1611689190,
			EventUUID:     "5CC1A437-F1AE-409A-22CB-7E3769858939",
			Username:      "Bob",
		}
		expectedFutureLogin := models.LoginRequest{
			IPAddress:     "6.7.8.9/0",
			UnixTimestamp: 1611696390,
			EventUUID:     "775CB1A3-FF0F-14E6-15D9-3FDB98B5C530",
			Username:      "Bob",
		}

		rows := sqlmock.NewRows([]string{"user_name", "time_stamp", "event_uuid", "ip_address"}).
			AddRow("Bob", 1611696390, "775CB1A3-FF0F-14E6-15D9-3FDB98B5C530", "6.7.8.9/0").
			AddRow("Bob", 1611692790, "926C4CDA-A5F0-19AB-C097-A61431CB4BFA", "1.2.3.4/5").
			AddRow("Bob", 1611689190, "5CC1A437-F1AE-409A-22CB-7E3769858939", "5.4.3.2/1")
		mock.ExpectPrepare("with cte").ExpectQuery().WillReturnRows(rows)

		previousLogin, futureLogin := geoRepository.GetPreviousAndFutureIPAdress(username, currentIP, timeStamp)

		if reflect.DeepEqual(previousLogin, expectedPreviousLogin) != true {

			t.Errorf("returned incorrect previous login: got %+v\n want %+v\n", previousLogin, expectedPreviousLogin)
		}

		if reflect.DeepEqual(futureLogin, expectedFutureLogin) != true {
			t.Errorf("returned incorrect future login: got %+v\n want %+v\n", futureLogin, expectedFutureLogin)
		}
	})
}
