package models

//when requesting NetBox API for a list of all device types, you will get a paginated response
//SitePage struct represents a single page of that anwer
//when "next" has a value of null, it means that you've reached the last page
//when "previous" has a value of null, it means, that you are curently viewing the first page

type DeviceTypePage struct {
	Count    int          `json:"count"` //count represents a count of ALL device types across ALL pages of the response
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []DeviceType `json:"results"`
}

//all of the necessary structs to unmarshal the DeviceType object
// \/\/\/

type DeviceType struct {
	Id          int    `json:"id"`
	Model       string `json:"slug"`
	UnitHeight  int    `json:"u_height"`
	IsFullDepth bool   `json:"is_full_depth"`
	DeviceCount int    `json:"device_count"`
}

// /\/\/\

//site repo interface

type DeviceTypeRepository interface {
	GetAllDeviceTypes() ([]DeviceType, error)
}
