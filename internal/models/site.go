package models

//when requesting NetBox API for a list of all sites, you will get a paginated response
//SitePage struct represents a single page of that anwer
//when "next" has a value of null, it means that you've reached the last page
//when "previous" has a value of null, it means, that you are curently viewing the first page

type SitePage struct {
	Count    int    `json:"count"` //count represents a count of ALL sites across ALL pages of the response
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Site `json:"results"`
}

type Site struct {
	Id   int    `json:"id"`
	Name string `json:"slug"`
	//note this is enough information to know about a region, and because it is only two fields - nested struct is used
	Region struct {
		Id   int    `json:"id"`
		Name string `json:"slug"`
	} `json:"region"`

	Facility            string `json:"facility"`
	Description         string `json:"description"`
	PhysicalAddress     string `json:"physical_address"`
	ShippingAddress     string `json:"shipping_address"`
	CircuitCount        int    `json:"circuit_count"`
	DeviceCount         int    `json:"device_count"`
	PrefixCount         int    `json:"prefix_count"`
	RackCount           int    `json:"rack_count"`
	VirtualMachineCount int    `json:"virtualmachine_count"`
	VlanCount           int    `json:"vlan_count"`
}

//site repo interface

type SiteRepository interface {
	GetAllSites() ([]Site, error)
}
