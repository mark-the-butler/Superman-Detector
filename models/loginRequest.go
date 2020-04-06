package models

// LoginRequest is used to map the json request payload to
type LoginRequest struct {
	Username      string `json:"username"`
	UnixTimestamp int    `json:"unix_timestamp"`
	EventUUID     string `json:"event_uuid"`
	IPAddress     string `json:"ip_address"`
}
