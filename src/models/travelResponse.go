package models

// TravelResponse is the object that gets serialized as the json response
type TravelResponse struct {
	CurrentLocation                GeoLocation  `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool         `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool         `json:"traveFromCurrentGeoSuspicious"`
	PreviousLocation               *GeoLocation `json:"precedingIpAccess,omitempty"`
	FutureLocation                 *GeoLocation `json:"subsequentIpAccess,omitempty"`
}
