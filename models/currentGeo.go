package models

// CurrentGeo holds details of the current login attempts for the response
type CurrentGeo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius int     `json:"radius"`
}
