package models

// TravelResponse is the object that gets serialized as the json response
type TravelResponse struct {
	CurrentLocation CurrentGeo `json:"currentGeo"`
}
