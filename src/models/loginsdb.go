package models

type Logins struct {
	RowNum        int    `db:"row_num"`
	Username      string `db:"user_name"`
	UnixTimestamp int    `db:"time_stamp"`
	EventUUID     string `db:"event_uuid"`
	IPAddress     string `db:"ip_address"`
}
