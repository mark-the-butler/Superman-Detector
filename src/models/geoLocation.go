package models

// GeoLocation holds details location details of login attempts
type GeoLocation struct {
	IP        string  `json:"ip,omitempty"`
	Speed     int     `json:"speed,omitempty"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Radius    int     `json:"radius"`
	TimeStamp int     `json:"timestamp,omitempty"`
}
