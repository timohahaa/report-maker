package models

//when requesting NetBox API for a list of all devices, you will get a paginated response
//DevicePage struct represents a single page of that anwer
//when "next" has a value of null, it means that you've reached the last page
//when "previous" has a value of null, it means, that you are curently viewing the first page

type DevicePage struct {
	Count    int      `json:"count"` //count represents a count of ALL devices across ALL pages of the response
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Device `json:"results"`
}

type Device struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	//note: device type here is more concise, because NetBox sends it that way, when requesting for devices
	//if there is a need to get a full device type description for a particular device
	//(witch there is almost no need to do so)
	//request for it using the device type id
	DeviceType struct {
		Id    int    `json:"id"`
		Model string `json:"slug"`
	} `json:"device_type"`
	//same as for DeviceType
	Site struct {
		Id   int    `json:"id"`
		Name string `json:"slug"`
	}
}

//site repo interface

type DeviceRepository interface {
	GetAllDevices() ([]Device, error)
}
