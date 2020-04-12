package models

// LoginRequest holds details about the login request that was just made
type LoginRequest struct {
	Username      string `json:"username"`
	UnixTimestamp int    `json:"unix_timestamp"`
	EventUUID     string `json:"event_uuid"`
	IPAddress     string `json:"ip_address"`
}
