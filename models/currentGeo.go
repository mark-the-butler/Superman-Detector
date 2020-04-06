package models

// CurrentGeo is used to map the current location data from the db to
type CurrentGeo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius int     `json:"radius"`
}
