package models

// TravelResponse is the object that gets serialized as the json response
type TravelResponse struct {
	CurrentLocation                CurrentGeo  `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool        `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool        `json:"traveFromCurrentGeoSuspicious"`
	PreviousLocation               PreviousGeo `json:"precedingIpAccess"`
	FutureLocation                 FutureGeo   `json:"subsequentIpAccess"`
}
